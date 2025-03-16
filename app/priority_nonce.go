package app

import (
	"context"
	"sync"

	"fmt"
	"math"

	"github.com/huandu/skiplist"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

const MAX_TXS_PRE_SENDER_IN_MEMPOOL int = 48

var (
	_ mempool.Mempool  = (*PriorityNonceMempool)(nil)
	_ mempool.Iterator = (*PriorityNonceIterator)(nil)

	errMempoolTxGasPriceTooLow = errors.New("gas price is too low")
	errMempoolTooManyTxs       = errors.New("tx sender has too many txs in mempool")
	errMempoolIsFull           = errors.New("mempool is full")
	errTxInMempool             = errors.New("tx already in mempool")
)

// PriorityNonceMempool is a mempool implementation that stores txs
// in a partially ordered set by 2 dimensions: priority, and sender-nonce
// (sequence number). Internally it uses one priority ordered skip list and one
// skip list per sender ordered by sender-nonce (sequence number). When there
// are multiple txs from the same sender, they are not always comparable by
// priority to other sender txs and must be partially ordered by both sender-nonce
// and priority.
type PriorityNonceMempool struct {
	mtx            sync.Mutex
	priorityIndex  *skiplist.SkipList
	priorityCounts map[int64]int
	senderIndices  map[string]*skiplist.SkipList
	scores         map[txMeta]txMeta
	onRead         func(tx sdk.Tx)
	txReplacement  func(op, np int64, oTx, nTx sdk.Tx) bool
	maxTx          int

	senderTxCntLock sync.RWMutex
	counterBySender map[string]int
	txRecord        map[txMeta]struct{}

	txReplacedCallback func(ctx context.Context, oldTx, newTx *TxInfo)
}

type PriorityNonceIterator struct {
	senderCursors map[string]*skiplist.Element
	nextPriority  int64
	sender        string
	priorityNode  *skiplist.Element
	mempool       *PriorityNonceMempool
}

// txMeta stores transaction metadata used in indices
type txMeta struct {
	// nonce is the sender's sequence number
	nonce uint64
	// priority is the transaction's priority
	priority int64
	// sender is the transaction's sender
	sender string
	// weight is the transaction's weight, used as a tiebreaker for transactions with the same priority
	weight int64
	// senderElement is a pointer to the transaction's element in the sender index
	senderElement *skiplist.Element
}

// txMetaLess is a comparator for txKeys that first compares priority, then weight,
// then sender, then nonce, uniquely identifying a transaction.
//
// Note, txMetaLess is used as the comparator in the priority index.
func txMetaLess(a, b any) int {
	keyA := a.(txMeta)
	keyB := b.(txMeta)
	res := skiplist.Int64.Compare(keyA.priority, keyB.priority)
	if res != 0 {
		return res
	}

	// Weight is used as a tiebreaker for transactions with the same priority.
	// Weight is calculated in a single pass in .Select(...) and so will be 0
	// on .Insert(...).
	res = skiplist.Int64.Compare(keyA.weight, keyB.weight)
	if res != 0 {
		return res
	}

	// Because weight will be 0 on .Insert(...), we must also compare sender and
	// nonce to resolve priority collisions. If we didn't then transactions with
	// the same priority would overwrite each other in the priority index.
	res = skiplist.String.Compare(keyA.sender, keyB.sender)
	if res != 0 {
		return res
	}

	return skiplist.Uint64.Compare(keyA.nonce, keyB.nonce)
}

type PriorityNonceMempoolOption func(*PriorityNonceMempool)

// PriorityNonceWithOnRead sets a callback to be called when a tx is read from
// the mempool.
func PriorityNonceWithOnRead(onRead func(tx sdk.Tx)) PriorityNonceMempoolOption {
	return func(mp *PriorityNonceMempool) {
		mp.onRead = onRead
	}
}

// PriorityNonceWithTxReplacement sets a callback to be called when duplicated
// transaction nonce detected during mempool insert. An application can define a
// transaction replacement rule based on tx priority or certain transaction fields.
func PriorityNonceWithTxReplacement(txReplacementRule func(op, np int64, oTx, nTx sdk.Tx) bool) PriorityNonceMempoolOption {
	return func(mp *PriorityNonceMempool) {
		mp.txReplacement = txReplacementRule
	}
}

// PriorityNonceWithMaxTx sets the maximum number of transactions allowed in the
// mempool with the semantics:
//
// <0: disabled, `Insert` is a no-op
// 0: unlimited
// >0: maximum number of transactions allowed
func PriorityNonceWithMaxTx(maxTx int) PriorityNonceMempoolOption {
	return func(mp *PriorityNonceMempool) {
		mp.maxTx = maxTx
	}
}

func PriorityNonceWithTxReplacedCallback(cb func(ctx context.Context, oldTx, newTx *TxInfo)) PriorityNonceMempoolOption {
	return func(mp *PriorityNonceMempool) {
		mp.txReplacedCallback = cb
	}
}

// DefaultPriorityMempool returns a priorityNonceMempool with no options.
func DefaultPriorityMempool() mempool.Mempool {
	return NewPriorityMempool()
}

// NewPriorityMempool returns the SDK's default mempool implementation which
// returns txs in a partial order by 2 dimensions; priority, and sender-nonce.
func NewPriorityMempool(opts ...PriorityNonceMempoolOption) *PriorityNonceMempool {
	mp := &PriorityNonceMempool{
		priorityIndex:   skiplist.New(skiplist.LessThanFunc(txMetaLess)),
		priorityCounts:  make(map[int64]int),
		senderIndices:   make(map[string]*skiplist.SkipList),
		scores:          make(map[txMeta]txMeta),
		counterBySender: make(map[string]int),
		txRecord:        make(map[txMeta]struct{}),
	}

	for _, opt := range opts {
		opt(mp)
	}

	return mp
}

// NextSenderTx returns the next transaction for a given sender by nonce order,
// i.e. the next valid transaction for the sender. If no such transaction exists,
// nil will be returned.
func (mp *PriorityNonceMempool) NextSenderTx(sender string) sdk.Tx {
	senderIndex, ok := mp.senderIndices[sender]
	if !ok {
		return nil
	}

	cursor := senderIndex.Front()
	return cursor.Value.(sdk.Tx)
}

// Insert attempts to insert a Tx into the app-side mempool in O(log n) time,
// returning an error if unsuccessful. Sender and nonce are derived from the
// transaction's first signature.
//
// Transactions are unique by sender and nonce. Inserting a duplicate tx is an
// O(log n) no-op.
//
// Inserting a duplicate tx with a different priority overwrites the existing tx,
// changing the total order of the mempool.
func (mp *PriorityNonceMempool) Insert(ctx context.Context, tx sdk.Tx) error {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	// if mp.maxTx > 0 && mp.CountTx() >= mp.maxTx {
	// 	return mempool.ErrMempoolTxMaxCapacity
	// } else
	if mp.maxTx < 0 {
		return nil
	}

	sdkContext := sdk.UnwrapSDKContext(ctx)
	priority := sdkContext.Priority()

	txInfo, err := extractTxInfo(tx)
	if err != nil {
		return err
	}

	if !mp.canInsert(txInfo.Sender) {
		return errors.Wrapf(errMempoolTooManyTxs, "[%d@%s]sender has too many txs in mempool", txInfo.Nonce, txInfo.Sender)
	}

	// init sender index if not exists
	senderIndex, ok := mp.senderIndices[txInfo.Sender]
	if !ok {
		senderIndex = skiplist.New(skiplist.LessThanFunc(func(a, b any) int {
			return skiplist.Uint64.Compare(b.(txMeta).nonce, a.(txMeta).nonce)
		}))

		// initialize sender index if not found
		mp.senderIndices[txInfo.Sender] = senderIndex
	}

	newKey := txMeta{nonce: txInfo.Nonce, priority: priority, sender: txInfo.Sender}

	// Since mp.priorityIndex is scored by priority, then sender, then nonce, a
	// changed priority will create a new key, so we must remove the old key and
	// re-insert it to avoid having the same tx with different priorityIndex indexed
	// twice in the mempool.
	//
	// This O(log n) remove operation is rare and only happens when a tx's priority
	// changes.

	sk := txMeta{nonce: txInfo.Nonce, sender: txInfo.Sender}
	if oldScore, txExists := mp.scores[sk]; txExists {
		if oldScore.priority < priority {
			oldTx := senderIndex.Get(newKey).Value.(sdk.Tx)
			return mp.doTxReplace(ctx, newKey, oldScore, oldTx, tx)
		}
		return errors.Wrapf(errTxInMempool, "[%d@%s] tx already in mempool", txInfo.Nonce, txInfo.Sender)
	} else {
		mempoolSize := mp.priorityIndex.Len()
		if mempoolSize >= mp.maxTx {
			lowestPriority := mp.getLowestPriority()
			// find one to replace
			if lowestPriority > 0 && priority <= lowestPriority {
				return errors.Wrapf(errMempoolTxGasPriceTooLow, "[%d@%s]tx with priority %d is too low, current lowest priority is %d", newKey.nonce, newKey.sender, priority, lowestPriority)
			}

			var maxIndexSize int
			var lowerPriority int64 = math.MaxInt64
			var selectedElement *skiplist.Element
			for sender, index := range mp.senderIndices {
				indexSize := index.Len()
				if sender == txInfo.Sender {
					continue
				}

				if indexSize > 1 {
					tail := index.Back()
					if tail != nil {
						tailKey := tail.Key().(txMeta)
						if tailKey.priority < lowerPriority {
							lowerPriority = tailKey.priority
							maxIndexSize = indexSize
							selectedElement = tail
						} else if tailKey.priority == lowerPriority {
							if indexSize > maxIndexSize {
								maxIndexSize = indexSize
								selectedElement = tail
							}
						}
					}
				}
			}

			if selectedElement != nil {
				key := selectedElement.Key().(txMeta)
				replacedTx, _ := mp.doRemove(key, true)

				// insert new tx
				mp.doInsert(newKey, tx, true)

				if mp.txReplacedCallback != nil && replacedTx != nil {
					sdkContext.Logger().Debug("txn replaced caused by full of mempool", "old", fmt.Sprintf("%d@%s", key.nonce, key.sender), "new", fmt.Sprintf("%d@%s", newKey.nonce, newKey.sender), "mempoolSize", mempoolSize)
					mp.txReplacedCallback(ctx,
						&TxInfo{Sender: key.sender, Nonce: key.nonce, Tx: replacedTx},
						&TxInfo{Sender: newKey.sender, Nonce: newKey.nonce, Tx: tx},
					)
				}
			} else {
				// not found any index more than 1 except sender's index
				// We do not replace the sender's only tx in the mempool
				return errors.Wrapf(errMempoolIsFull, "%d@%s with priority%d", newKey.nonce, newKey.sender, newKey.priority)
			}
		} else {
			mp.doInsert(newKey, tx, true)
		}
		return nil
	}
}

func (mp *PriorityNonceMempool) doInsert(newKey txMeta, tx sdk.Tx, incrCnt bool) {
	senderIndex, ok := mp.senderIndices[newKey.sender]
	if !ok {
		senderIndex = skiplist.New(skiplist.LessThanFunc(func(a, b any) int {
			return skiplist.Uint64.Compare(b.(txMeta).nonce, a.(txMeta).nonce)
		}))

		// initialize sender index if not found
		mp.senderIndices[newKey.sender] = senderIndex
	}

	mp.priorityCounts[newKey.priority]++
	newKey.senderElement = senderIndex.Set(newKey, tx)

	mp.scores[txMeta{nonce: newKey.nonce, sender: newKey.sender}] = txMeta{priority: newKey.priority}
	mp.priorityIndex.Set(newKey, tx)

	if incrCnt {
		mp.incrSenderTxCnt(newKey.sender, newKey.nonce)
	}
}

func (mp *PriorityNonceMempool) doRemove(oldKey txMeta, decrCnt bool) (sdk.Tx, error) {
	scoreKey := txMeta{nonce: oldKey.nonce, sender: oldKey.sender}
	score, ok := mp.scores[scoreKey]
	if !ok {
		return nil, errors.Wrapf(mempool.ErrTxNotFound, "%d@%s not found", oldKey.nonce, oldKey.sender)
	}
	tk := txMeta{nonce: oldKey.nonce, priority: score.priority, sender: oldKey.sender, weight: score.weight}

	senderTxs, ok := mp.senderIndices[oldKey.sender]
	if !ok {
		return nil, fmt.Errorf("%d@%s not found", oldKey.nonce, oldKey.sender)
	}

	mp.priorityIndex.Remove(tk)
	removedElem := senderTxs.Remove(tk)
	delete(mp.scores, scoreKey)
	mp.priorityCounts[score.priority]--

	if decrCnt {
		mp.decrSenderTxCnt(oldKey.sender, oldKey.nonce)
	}

	if removedElem == nil {
		return nil, mempool.ErrTxNotFound
	}

	return removedElem.Value.(sdk.Tx), nil
}

func (mp *PriorityNonceMempool) doTxReplace(ctx context.Context, newMate, oldMate txMeta, oldTx, newTx sdk.Tx) error {
	if mp.txReplacement != nil && !mp.txReplacement(oldMate.priority, newMate.priority, oldTx, newTx) {
		return fmt.Errorf(
			"tx doesn't fit the replacement rule, oldPriority: %v, newPriority: %v, oldTx: %v, newTx: %v",
			oldMate.priority,
			newMate.priority,
			oldTx,
			newTx,
		)
	}

	e := mp.priorityIndex.Remove(txMeta{
		nonce:    newMate.nonce,
		sender:   newMate.sender,
		priority: oldMate.priority,
		weight:   oldMate.weight,
	})
	replacedTx := e.Value.(sdk.Tx)
	mp.priorityCounts[oldMate.priority]--

	mp.doInsert(newMate, newTx, false)

	if mp.txReplacedCallback != nil && replacedTx != nil {
		sdkContext := sdk.UnwrapSDKContext(ctx)
		sdkContext.Logger().Debug("txn update", "txn", fmt.Sprintf("%d@%s", newMate.nonce, newMate.sender), "oldPriority", oldMate.priority, "newPriority", newMate.priority)
		mp.txReplacedCallback(ctx,
			&TxInfo{Sender: newMate.sender, Nonce: newMate.nonce, Tx: replacedTx},
			&TxInfo{Sender: newMate.sender, Nonce: newMate.nonce, Tx: newTx},
		)
	}

	return nil
}

func (i *PriorityNonceIterator) iteratePriority() mempool.Iterator {
	// beginning of priority iteration
	if i.priorityNode == nil {
		i.priorityNode = i.mempool.priorityIndex.Front()
	} else {
		i.priorityNode = i.priorityNode.Next()
	}

	// end of priority iteration
	if i.priorityNode == nil {
		return nil
	}

	i.sender = i.priorityNode.Key().(txMeta).sender

	nextPriorityNode := i.priorityNode.Next()
	if nextPriorityNode != nil {
		i.nextPriority = nextPriorityNode.Key().(txMeta).priority
	} else {
		i.nextPriority = math.MinInt64
	}

	return i.Next()
}

func (i *PriorityNonceIterator) Next() mempool.Iterator {
	if i.priorityNode == nil {
		return nil
	}

	cursor, ok := i.senderCursors[i.sender]
	if !ok {
		// beginning of sender iteration
		cursor = i.mempool.senderIndices[i.sender].Front()
	} else {
		// middle of sender iteration
		cursor = cursor.Next()
	}

	// end of sender iteration
	if cursor == nil {
		return i.iteratePriority()
	}

	key := cursor.Key().(txMeta)

	// We've reached a transaction with a priority lower than the next highest
	// priority in the pool.
	if key.priority < i.nextPriority {
		return i.iteratePriority()
	} else if key.priority == i.nextPriority && i.priorityNode.Next() != nil {
		// Weight is incorporated into the priority index key only (not sender index)
		// so we must fetch it here from the scores map.
		weight := i.mempool.scores[txMeta{nonce: key.nonce, sender: key.sender}].weight
		if weight < i.priorityNode.Next().Key().(txMeta).weight {
			return i.iteratePriority()
		}
	}

	i.senderCursors[i.sender] = cursor
	return i
}

func (i *PriorityNonceIterator) Tx() sdk.Tx {
	return i.senderCursors[i.sender].Value.(sdk.Tx)
}

// Select returns a set of transactions from the mempool, ordered by priority
// and sender-nonce in O(n) time. The passed in list of transactions are ignored.
// This is a readonly operation, the mempool is not modified.
//
// The maxBytes parameter defines the maximum number of bytes of transactions to
// return.
//
// NOTE: It is not safe to use this iterator while removing transactions from
// the underlying mempool.
func (mp *PriorityNonceMempool) Select(ctx context.Context, txs [][]byte) mempool.Iterator {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	return mp.doSelect(ctx, txs)
}

func (mp *PriorityNonceMempool) SelectBy(ctx context.Context, txs [][]byte, callback func(sdk.Tx) bool) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	iter := mp.doSelect(ctx, txs)
	for iter != nil && callback(iter.Tx()) {
		iter = iter.Next()
	}
}

func (mp *PriorityNonceMempool) doSelect(_ context.Context, _ [][]byte) mempool.Iterator {
	if mp.priorityIndex.Len() == 0 {
		return nil
	}

	mp.reorderPriorityTies()

	iterator := &PriorityNonceIterator{
		mempool:       mp,
		senderCursors: make(map[string]*skiplist.Element),
	}

	return iterator.iteratePriority()
}

func (mp *PriorityNonceMempool) GetSenderUncommittedTxnCount(ctx context.Context, sender string) int {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	if _, exists := mp.counterBySender[sender]; exists {
		return mp.counterBySender[sender]
	}
	return 0
}

type reorderKey struct {
	deleteKey txMeta
	insertKey txMeta
	tx        sdk.Tx
}

func (mp *PriorityNonceMempool) reorderPriorityTies() {
	node := mp.priorityIndex.Front()

	var reordering []reorderKey
	for node != nil {
		key := node.Key().(txMeta)
		if mp.priorityCounts[key.priority] > 1 {
			newKey := key
			newKey.weight = senderWeight(key.senderElement)
			reordering = append(reordering, reorderKey{deleteKey: key, insertKey: newKey, tx: node.Value.(sdk.Tx)})
		}

		node = node.Next()
	}

	for _, k := range reordering {
		mp.priorityIndex.Remove(k.deleteKey)
		delete(mp.scores, txMeta{nonce: k.deleteKey.nonce, sender: k.deleteKey.sender})
		mp.priorityIndex.Set(k.insertKey, k.tx)
		mp.scores[txMeta{nonce: k.insertKey.nonce, sender: k.insertKey.sender}] = k.insertKey
	}
}

// senderWeight returns the weight of a given tx (t) at senderCursor. Weight is
// defined as the first (nonce-wise) same sender tx with a priority not equal to
// t. It is used to resolve priority collisions, that is when 2 or more txs from
// different senders have the same priority.
func senderWeight(senderCursor *skiplist.Element) int64 {
	if senderCursor == nil {
		return 0
	}

	weight := senderCursor.Key().(txMeta).priority
	senderCursor = senderCursor.Next()
	for senderCursor != nil {
		p := senderCursor.Key().(txMeta).priority
		if p != weight {
			weight = p
		}

		senderCursor = senderCursor.Next()
	}

	return weight
}

// CountTx returns the number of transactions in the mempool.
func (mp *PriorityNonceMempool) CountTx() int {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	return mp.priorityIndex.Len()
}

// Remove removes a transaction from the mempool in O(log n) time, returning an
// error if unsuccessful.
func (mp *PriorityNonceMempool) Remove(tx sdk.Tx) error {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	txInfo, err := extractTxInfo(tx)
	if err != nil {
		return err
	}

	mp.decrSenderTxCnt(txInfo.Sender, txInfo.Nonce)

	scoreKey := txMeta{nonce: txInfo.Nonce, sender: txInfo.Sender}
	score, ok := mp.scores[scoreKey]
	if !ok {
		return mempool.ErrTxNotFound
	}
	tk := txMeta{nonce: txInfo.Nonce, priority: score.priority, sender: txInfo.Sender, weight: score.weight}

	senderTxs, ok := mp.senderIndices[txInfo.Sender]
	if !ok {
		return fmt.Errorf("sender %s not found", txInfo.Sender)
	}

	mp.priorityIndex.Remove(tk)
	senderTxs.Remove(tk)
	delete(mp.scores, scoreKey)
	mp.priorityCounts[score.priority]--

	return nil
}

func (mp *PriorityNonceMempool) getLowestPriority() int64 {
	if mp.priorityIndex.Len() == 0 {
		return 0
	}

	min := int64(math.MaxInt64)
	for priority, count := range mp.priorityCounts {
		if count > 0 {
			if priority < min {
				min = priority
			}
		}
	}

	return min
}

func (mp *PriorityNonceMempool) canInsert(sender string) bool {
	mp.senderTxCntLock.RLock()
	defer mp.senderTxCntLock.RUnlock()

	if _, exists := mp.counterBySender[sender]; exists {
		return mp.counterBySender[sender] < MAX_TXS_PRE_SENDER_IN_MEMPOOL
	}

	return true
}

func (mp *PriorityNonceMempool) incrSenderTxCnt(sender string, nonce uint64) error {
	mp.senderTxCntLock.Lock()
	defer mp.senderTxCntLock.Unlock()

	existsKey := txMeta{nonce: nonce, sender: sender}
	if _, exists := mp.txRecord[existsKey]; !exists {
		mp.txRecord[existsKey] = struct{}{}

		if _, exists := mp.counterBySender[sender]; !exists {
			mp.counterBySender[sender] = 1
		} else {
			if mp.counterBySender[sender] < MAX_TXS_PRE_SENDER_IN_MEMPOOL {
				mp.counterBySender[sender] += 1
			} else {
				return fmt.Errorf("tx sender has too many txs in mempool")
			}
		}
	}
	return nil
}

func (mp *PriorityNonceMempool) decrSenderTxCnt(sender string, nonce uint64) {
	mp.senderTxCntLock.Lock()
	defer mp.senderTxCntLock.Unlock()

	existsKey := txMeta{nonce: nonce, sender: sender}
	if _, exists := mp.txRecord[existsKey]; exists {
		delete(mp.txRecord, existsKey)

		if _, exists := mp.counterBySender[sender]; exists {
			if mp.counterBySender[sender] > 1 {
				mp.counterBySender[sender] -= 1
			} else {
				delete(mp.counterBySender, sender)
			}
		}
	}
}

func IsEmpty(mempool mempool.Mempool) error {
	mp := mempool.(*PriorityNonceMempool)
	if mp.priorityIndex.Len() != 0 {
		return fmt.Errorf("priorityIndex not empty")
	}

	var countKeys []int64
	for k := range mp.priorityCounts {
		countKeys = append(countKeys, k)
	}

	for _, k := range countKeys {
		if mp.priorityCounts[k] != 0 {
			return fmt.Errorf("priorityCounts not zero at %v, got %v", k, mp.priorityCounts[k])
		}
	}

	var senderKeys []string
	for k := range mp.senderIndices {
		senderKeys = append(senderKeys, k)
	}

	for _, k := range senderKeys {
		if mp.senderIndices[k].Len() != 0 {
			return fmt.Errorf("senderIndex not empty for sender %v", k)
		}
	}

	return nil
}

type TxInfo struct {
	Sender string
	Nonce  uint64
	Tx     sdk.Tx
}

func extractTxInfo(tx sdk.Tx) (*TxInfo, error) {
	var sender string
	var nonce uint64

	sigs, err := tx.(signing.SigVerifiableTx).GetSignaturesV2()
	if err != nil {
		return nil, err
	}

	if len(sigs) == 0 {
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return nil, fmt.Errorf("tx must have at least one signer")
		}
		msgEthTx, ok := msgs[0].(*evmtypes.MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("tx must have at least one signer")
		}
		ethTx := msgEthTx.AsTransaction()
		signer := gethtypes.NewEIP2930Signer(ethTx.ChainId())
		ethSender, err := signer.Sender(ethTx)
		if err != nil {
			return nil, fmt.Errorf("tx must have at least one signer")
		}
		sender = sdk.AccAddress(ethSender.Bytes()).String()
		nonce = ethTx.Nonce()
	} else {
		sig := sigs[0]
		sender = sdk.AccAddress(sig.PubKey.Address()).String()
		nonce = sig.Sequence
	}

	return &TxInfo{Sender: sender, Nonce: nonce, Tx: tx}, nil
}

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cmtstore "github.com/cometbft/cometbft/proto/tendermint/store"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbName = "dbQuery"
)

var (
	appDB       dbm.DB
	blockDB     dbm.DB
	stateDB     dbm.DB
	indexDB     dbm.DB
	encodingCfg EncodingConfig
)

func InitDB(backend dbm.BackendType, dir string) error {
	encodingCfg = MakeEncodingConfig()
	var err error
	appDB, err = dbm.NewDB("application", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create application db")
	}
	blockDB, err = dbm.NewDB("blockstore", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create blockstore db")
	}

	stateDB, err = dbm.NewDB("state", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create state db")
	}
	indexDB, err = dbm.NewDB("tx_index", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create tx_index db")
	}

	iterateBlocks(1600000)
	return nil
}

func GetTx(hash string) (sdk.Tx, error) {
	if len(hash) == 0 {
		return nil, errors.New("tx hash cannot be empty")
	}

	txByteSlice, err := queryTx(hash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Errorf(codes.NotFound, "tx not found: %s", hash)
		}

		return nil, errors.Wrapf(err, "failed to get tx")
	}

	// input tx byte slice
	tx, err := encodingCfg.TxConfig.TxDecoder()(txByteSlice)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode tx")
	}

	return tx, nil
}

func queryTx(hashHexStr string) ([]byte, error) {
	r, err := loadTxFromIndex(hashHexStr)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, errors.New("tx not found")
	}

	return r.Tx, nil
}

func loadTxFromIndex(hashHexStr string) (*abci.TxResult, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	rawBytes, err := indexDB.Get(hash)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get tx from indexDB")
	}

	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(abci.TxResult)
	err = proto.Unmarshal(rawBytes, txResult)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal txResult")
	}

	jsonStr, _ := AbciTxResultToJSONString(txResult)

	println("abci.TxResult: ", jsonStr)

	return txResult, nil
}

func calcBlockMetaKey(height int64) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}

func calcBlockPartKey(height int64, partIndex int) []byte {
	return []byte(fmt.Sprintf("P:%v:%v", height, partIndex))
}

func loadBlockMeta(height int64) *cometbfttypes.BlockMeta {
	pbbm := new(cmtproto.BlockMeta)
	bz, err := blockDB.Get(calcBlockMetaKey(height))
	if err != nil {
		panic(err)
	}

	if len(bz) == 0 {
		return nil
	}

	err = proto.Unmarshal(bz, pbbm)
	if err != nil {
		panic(fmt.Errorf("unmarshal to cmtproto.BlockMeta: %w", err))
	}

	blockMeta, err := types.BlockMetaFromTrustedProto(pbbm)
	if err != nil {
		panic(fmt.Errorf("error from proto blockMeta: %w", err))
	}

	return blockMeta
}

func loadBlockPart(height int64, index int) *types.Part {
	pbpart := new(cmtproto.Part)

	bz, err := blockDB.Get(calcBlockPartKey(height, index))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil
	}

	err = proto.Unmarshal(bz, pbpart)
	if err != nil {
		panic(fmt.Errorf("unmarshal to cmtproto.Part failed: %w", err))
	}
	part, err := types.PartFromProto(pbpart)
	if err != nil {
		panic(fmt.Sprintf("Error reading block part: %v", err))
	}

	return part
}

func loadBlock(height int64) *cometbfttypes.Block {
	blockMeta := loadBlockMeta(height)
	if blockMeta == nil {
		return nil
	}
	pbb := new(cmtproto.Block)
	buf := []byte{}
	for i := 0; i < int(blockMeta.BlockID.PartSetHeader.Total); i++ {
		part := loadBlockPart(height, i)
		// If the part is missing (e.g. since it has been deleted after we
		// loaded the block meta) we consider the whole block to be missing.
		if part == nil {
			return nil
		}
		buf = append(buf, part.Bytes...)
	}
	err := proto.Unmarshal(buf, pbb)
	if err != nil {
		// NOTE: The existence of meta should imply the existence of the
		// block. So, make sure meta is only saved after blocks are saved.
		panic(fmt.Sprintf("Error reading block: %v", err))
	}

	block, err := types.BlockFromProto(pbb)
	if err != nil {
		panic(fmt.Errorf("error from proto block: %w", err))
	}

	return block
}

func iterateBlocks(start int64) {
	h := start
	for {
		block := loadBlock(h)
		if block == nil {
			println("last block height: ", h)
			break
		}

		printBlockTxsAsJSON(block)
		h++
	}
}

func printLastTx(dir string) {
	// 打开 blockstore 数据库
	db, err := leveldb.OpenFile(path.Join(dir, "blockstore.db"), nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 遍历区块数据库查找最后一个区块
	iter := db.NewIterator(nil, nil)
	var lastBlockKey, lastBlockValue []byte
	for iter.Next() {
		lastBlockKey = iter.Key()
		lastBlockValue = iter.Value()

		var bsj cmtstore.BlockStoreState
		if err := json.Unmarshal(lastBlockValue, &bsj); err != nil {
			log.Printf("Failed to unmarshal block data: %v", err)
		} else {
			fmt.Printf("Block Height: %d\n", bsj.Height)
		}
	}
	iter.Release()

	// 检查迭代器错误
	if err := iter.Error(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}

	// 假设交易数据存储在最后一个区块数据中，解析并输出
	var blockData map[string]interface{}
	if err := json.Unmarshal(lastBlockValue, &blockData); err != nil {
		log.Fatalf("Failed to unmarshal block data: %v", err)
	}

	// 输出区块信息，获取交易数据
	fmt.Printf("Last block key: %s\n", lastBlockKey)
	fmt.Printf("Last block data: %+v\n", blockData)

	// 获取并解析交易数据
	// 假设交易数据位于 blockData["txs"]
	txs, ok := blockData["txs"].([]interface{})
	if !ok {
		log.Fatalf("Transaction data not found in last block.")
	}
	lastTx := txs[len(txs)-1]

	// 输出最后一个交易
	fmt.Printf("Last transaction: %+v\n", lastTx)

}

func printBlockTxsAsJSON(block *cometbfttypes.Block) {
	var txsJSON []string

	for _, tx := range block.Data.Txs {
		txJSON, err := json.Marshal(tx)
		if err != nil {
			log.Printf("Failed to marshal tx to JSON: %v", err)
			continue
		}
		txsJSON = append(txsJSON, string(txJSON))
	}

	blockData := map[string]interface{}{
		"height": block.Header.Height,
		"txs":    txsJSON,
	}

	blockJSON, err := json.MarshalIndent(blockData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal block data to JSON: %v", err)
	}

	fmt.Println(string(blockJSON))
}

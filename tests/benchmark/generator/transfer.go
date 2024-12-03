package generator

// 1, prepare accounts
// 2, make transfer from faucet account to other accounts
// 3, make transfer between accounts

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/0glabs/0g-chain/tests/benchmark/account"
	"github.com/0glabs/0g-chain/tests/benchmark/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type transferGeneratorImlp struct {
	client     *ethclient.Client
	chainId    *big.Int
	signer     types.Signer
	accountMgr account.AccountManager
	taskPool   chan *task
	txPool     chan *types.Transaction
	poolSize   uint32
	state      sync.Once
	context    context.Context
	cancelFunc context.CancelFunc
}

func NewTransferGenerator(poolSize uint32, faucetPrivateKey string, ethClient *ethclient.Client, accountMgr account.AccountManager) (Generator, error) {
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &transferGeneratorImlp{
		client:     ethClient,
		chainId:    chainID,
		signer:     types.NewEIP155Signer(chainID),
		accountMgr: accountMgr,
		poolSize:   poolSize,
		taskPool:   make(chan *task, 64),
		txPool:     make(chan *types.Transaction, poolSize),
		context:    ctx,
		cancelFunc: cancel,
	}, nil
}

func (g *transferGeneratorImlp) WarmUp() error {
	// make transfer from faucet account to other accounts
	taskList := make([]*task, 0, g.accountMgr.GetAccountCount())
	for i := 0; i < int(g.accountMgr.GetAccountCount()); i++ {
		targetAccount := g.accountMgr.GetAccount(uint32(i))

		balance, err := targetAccount.GetBalance(g.client)
		if err != nil {
			return errors.Wrap(err, "failed to get account balance")
		}
		if balance.Cmp(big.NewInt(defaultTransferVal)) < 0 {
			taskList = append(taskList, &task{
				fromAccount: g.accountMgr.GetFaucetAccount(),
				toAccout:    targetAccount,
				value:       utils.ToBigInt(initialTransferVal),
			})
		} else {
			println("account:", targetAccount.Address.Hex(), "balance:", balance.String())
		}
	}

	for i := range taskList {
		tx, err := g.generateTransaction(taskList[i])
		if err != nil {
			return err
		}
		g.txPool <- tx
	}
	return nil
}

func (g *transferGeneratorImlp) generateTransaction(t *task) (*types.Transaction, error) {
	ctx := context.Background()

	nonce := t.fromAccount.GetAndIncrementNonce()

	gasLimit := defaultTransferGasLimit
	gasPrice, err := g.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// tx := types.NewTransaction(nonce, t.toAccout.Address, t.value, gasLimit, big.NewInt(0), nil)
	tx := types.NewTransaction(nonce, t.toAccout.Address, t.value, gasLimit, gasPrice, nil)
	t.fromAccount.ReqChan <- &account.TxSignRequest{
		Nonce: nonce,
		Tx:    tx,
	}

	res := <-t.fromAccount.ResChan

	return res.SignedTx, nil
}

func (g *transferGeneratorImlp) GenerateTransfer() <-chan *types.Transaction {
	g.state.Do(func() {
		go func(ctx context.Context) {
			for {
				acctIdxList := make([]uint32, 0, g.accountMgr.GetAccountCount())

				for i := 0; i < int(g.accountMgr.GetAccountCount()); i++ {
					acctIdxList = append(acctIdxList, uint32(i))
				}

				acctIdxList = utils.Shuffle(acctIdxList)

				usedFrom := make(map[uint32]bool)
				usedTo := make(map[uint32]bool)

				for i := 0; i < len(acctIdxList); i++ {
					for j := 0; j < len(acctIdxList); j++ {
						if acctIdxList[i] != acctIdxList[j] && !usedFrom[acctIdxList[i]] && !usedTo[acctIdxList[j]] {
							fromAcct := g.accountMgr.GetAccount(acctIdxList[i])
							toAcct := g.accountMgr.GetAccount(acctIdxList[j])

							if !fromAcct.IsAvailable() {
								balance, err := fromAcct.GetBalance(g.client)
								if err == nil {
									if balance.Cmp(big.NewInt(defaultTransferVal)) < 0 {
										g.taskPool <- &task{
											fromAccount: g.accountMgr.GetFaucetAccount(),
											toAccout:    fromAcct,
											value:       utils.ToBigInt(initialTransferVal),
										}
									} else {
										fromAcct.MarkAsAvailable()
									}
								}
							}

							if !toAcct.IsAvailable() {
								balance, err := toAcct.GetBalance(g.client)
								if err == nil {
									if balance.Cmp(big.NewInt(defaultTransferVal)) < 0 {
										g.taskPool <- &task{
											fromAccount: g.accountMgr.GetFaucetAccount(),
											toAccout:    toAcct,
											value:       utils.ToBigInt(initialTransferVal),
										}
									} else {
										toAcct.MarkAsAvailable()
									}
								}
							}

							if fromAcct.IsAvailable() && toAcct.IsAvailable() {
								usedFrom[acctIdxList[i]] = true
								usedTo[acctIdxList[j]] = true

								g.taskPool <- &task{
									fromAccount: fromAcct,
									toAccout:    toAcct,
									value:       utils.ToBigInt(defaultTransferVal),
								}
								break
							}
						}
					}
				}

				select {
				case <-ctx.Done():
					close(g.taskPool)
					return
				default:
					continue
				}
			}
		}(g.context)

		go func(ctx context.Context) {
			for {
				t := <-g.taskPool
				tx, err := g.generateTransaction(t)
				if err != nil {
					log.Fatal("generate transaction error: ", err.Error())
				}
				g.txPool <- tx

				select {
				case <-ctx.Done():
					close(g.txPool)
					return
				default:
					continue
				}
			}
		}(g.context)
	})

	return g.txPool
}

func (g *transferGeneratorImlp) TearDown() {
	g.cancelFunc()
}

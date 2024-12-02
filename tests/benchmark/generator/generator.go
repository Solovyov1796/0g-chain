package generator

import (
	"math/big"

	"github.com/0glabs/0g-chain/tests/benchmark/account"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultTransferGasLimit      = uint64(21000)
	defaultErc20TransferGasLimit = uint64(3000000)
	initialTransferVal           = 50 * 1e18 // 50 a0gi
	defaultTransferVal           = 1e10      // 0.01 ua0gi
)

type task struct {
	fromAccount *account.Account
	toAccout    *account.Account
	value       *big.Int
}

type Generator interface {
	WarmUp() error
	GenerateTransfer() <-chan *types.Transaction
	TearDown()
}

type DeployedErc20 struct {
	Address string
}

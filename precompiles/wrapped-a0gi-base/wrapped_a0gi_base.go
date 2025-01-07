package wrappeda0gibase

import (
	"strings"

	precompiles_common "github.com/0glabs/0g-chain/precompiles/common"
	wrappeda0gibasekeeper "github.com/0glabs/0g-chain/x/wrapped-a0gi-base/keeper"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	PrecompileAddress = "0x0000000000000000000000000000000000001002"
)

var _ vm.PrecompiledContract = &WrappedA0giBasePrecompile{}
var _ precompiles_common.PrecompileCommon = &WrappedA0giBasePrecompile{}

type WrappedA0giBasePrecompile struct {
	abi                   abi.ABI
	wrappeda0gibaseKeeper wrappeda0gibasekeeper.Keeper
}

// Abi implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) Abi() *abi.ABI {
	panic("unimplemented")
}

// IsTx implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) IsTx(string) bool {
	panic("unimplemented")
}

// KVGasConfig implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) KVGasConfig() types.GasConfig {
	panic("unimplemented")
}

// Address implements vm.PrecompiledContract.
func (w *WrappedA0giBasePrecompile) Address() common.Address {
	return common.HexToAddress(PrecompileAddress)
}

// RequiredGas implements vm.PrecompiledContract.
func (w *WrappedA0giBasePrecompile) RequiredGas(input []byte) uint64 {
	panic("unimplemented")
}

func NewWrappedA0giBasePrecompile(wrappeda0gibaseKeeper wrappeda0gibasekeeper.Keeper) (*WrappedA0giBasePrecompile, error) {
	abi, err := abi.JSON(strings.NewReader(Wrappeda0gibaseABI))
	if err != nil {
		return nil, err
	}
	return &WrappedA0giBasePrecompile{
		abi:                   abi,
		wrappeda0gibaseKeeper: wrappeda0gibaseKeeper,
	}, nil
}

// Run implements vm.PrecompiledContract.
func (w *WrappedA0giBasePrecompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	panic("unimplemented")
}

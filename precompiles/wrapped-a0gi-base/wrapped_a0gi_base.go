package wrappeda0gibase

import (
	"strings"

	precompiles_common "github.com/0glabs/0g-chain/precompiles/common"
	wrappeda0gibasekeeper "github.com/0glabs/0g-chain/x/wrapped-a0gi-base/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	PrecompileAddress = "0x0000000000000000000000000000000000001002"

	// txs
	WrappedA0GIBaseFunctionMint = "mint"
	WrappedA0GIBaseFunctionBurn = "burn"
	// queries
	WrappedA0GIBaseFunctionGetWA0GI     = "getWA0GI"
	WrappedA0GIBaseFunctionMinterSupply = "minterSupply"
)

var _ vm.PrecompiledContract = &WrappedA0giBasePrecompile{}
var _ precompiles_common.PrecompileCommon = &WrappedA0giBasePrecompile{}

type WrappedA0giBasePrecompile struct {
	abi                   abi.ABI
	wrappeda0gibaseKeeper wrappeda0gibasekeeper.Keeper
}

// Abi implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) Abi() *abi.ABI {
	return &w.abi
}

// IsTx implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) IsTx(method string) bool {
	switch method {
	case WrappedA0GIBaseFunctionMint,
		WrappedA0GIBaseFunctionBurn:
		return true
	default:
		return false
	}
}

// KVGasConfig implements common.PrecompileCommon.
func (w *WrappedA0giBasePrecompile) KVGasConfig() storetypes.GasConfig {
	return storetypes.KVGasConfig()
}

// Address implements vm.PrecompiledContract.
func (w *WrappedA0giBasePrecompile) Address() common.Address {
	return common.HexToAddress(PrecompileAddress)
}

// RequiredGas implements vm.PrecompiledContract.
func (w *WrappedA0giBasePrecompile) RequiredGas(input []byte) uint64 {
	return 0
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
	ctx, stateDB, method, initialGas, args, err := precompiles_common.InitializePrecompileCall(w, evm, contract, readonly)
	if err != nil {
		return nil, err
	}

	var bz []byte
	switch method.Name {
	// queries
	case WrappedA0GIBaseFunctionGetWA0GI:
		bz, err = w.GetW0GI(ctx, evm, method, args)
	case WrappedA0GIBaseFunctionMinterSupply:
		bz, err = w.MinterSupply(ctx, evm, method, args)
	// txs
	case WrappedA0GIBaseFunctionMint:
		bz, err = w.Mint(ctx, evm, stateDB, contract, method, args)
	case WrappedA0GIBaseFunctionBurn:
		bz, err = w.Burn(ctx, evm, stateDB, contract, method, args)
	}

	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}
	return bz, nil
}

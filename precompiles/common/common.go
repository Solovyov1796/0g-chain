package common

import (
	"errors"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

type PrecompileCommon interface {
	Abi() *abi.ABI
	IsTx(string) bool
	KVGasConfig() storetypes.GasConfig
}

func InitializePrecompileCall(
	p PrecompileCommon,
	evm *vm.EVM,
	contract *vm.Contract,
	readonly bool,
) (
	ctx sdk.Context,
	stateDB *statedb.StateDB,
	method *abi.Method,
	initialGas storetypes.Gas,
	args []interface{},
	err error,
) {
	// parse input
	if len(contract.Input) < 4 {
		return sdk.Context{}, nil, nil, uint64(0), nil, vm.ErrExecutionReverted
	}
	method, err = p.Abi().MethodById(contract.Input[:4])
	if err != nil {
		return sdk.Context{}, nil, nil, uint64(0), nil, vm.ErrExecutionReverted
	}
	args, err = method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return sdk.Context{}, nil, nil, uint64(0), nil, err
	}
	// readonly check
	if readonly && p.IsTx(method.Name) {
		return sdk.Context{}, nil, nil, uint64(0), nil, errors.New(ErrWriteOnReadOnly)
	}
	// get state db and context
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return sdk.Context{}, nil, nil, uint64(0), nil, errors.New(ErrGetStateDB)
	}
	ctx, err = stateDB.GetCachedContextForPrecompile()
	if err != nil {
		return sdk.Context{}, nil, nil, uint64(0), nil, err
	}
	// initial gas
	initialGas = ctx.GasMeter().GasConsumed()
	ctx = ctx.WithKVGasConfig(p.KVGasConfig())

	return ctx, stateDB, method, initialGas, args, nil
}

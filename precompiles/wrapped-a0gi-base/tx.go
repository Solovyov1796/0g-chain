package wrappeda0gibase

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/statedb"
)

func (w *WrappedA0giBasePrecompile) Mint(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgMint(args)
	if err != nil {
		return nil, err
	}
	// validation
	wa0gi := common.BytesToAddress(w.wrappeda0gibaseKeeper.GetWA0GIAddress(ctx))
	if contract.CallerAddress != wa0gi {
		return nil, errors.New(ErrSenderNotWA0GI)
	}
	// execute
	_, err = w.wrappeda0gibaseKeeper.Mint(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack()
}

func (w *WrappedA0giBasePrecompile) Burn(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgBurn(args)
	if err != nil {
		return nil, err
	}
	// validation
	wa0gi := common.BytesToAddress(w.wrappeda0gibaseKeeper.GetWA0GIAddress(ctx))
	if contract.CallerAddress != wa0gi {
		return nil, errors.New(ErrSenderNotWA0GI)
	}
	// execute
	_, err = w.wrappeda0gibaseKeeper.Burn(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack()
}

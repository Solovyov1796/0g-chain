package wrappeda0gibase

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func (w *WrappedA0giBasePrecompile) GetW0GI(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewGetW0GIRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := w.wrappeda0gibaseKeeper.GetWA0GI(ctx, req)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(common.BytesToAddress(response.Address))
}

func (w *WrappedA0giBasePrecompile) MinterSupply(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewMinterSupplyRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := w.wrappeda0gibaseKeeper.MinterSupply(ctx, req)
	if err != nil {
		return nil, err
	}
	supply := Supply{
		Cap:   new(big.Int).SetBytes(response.Cap),
		Total: new(big.Int).SetBytes(response.Supply),
	}
	return method.Outputs.Pack(supply)
}

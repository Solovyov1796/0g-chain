package wrappeda0gibase

import (
	"fmt"
	"math/big"

	precompiles_common "github.com/0glabs/0g-chain/precompiles/common"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	"github.com/ethereum/go-ethereum/common"
)

type Supply = struct {
	Cap   *big.Int "json:\"cap\""
	Total *big.Int "json:\"total\""
}

func NewGetW0GIRequest(args []interface{}) (*types.GetWA0GIRequest, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 0, len(args))
	}
	return &types.GetWA0GIRequest{}, nil
}

func NewMinterSupplyRequest(args []interface{}) (*types.MinterSupplyRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 1, len(args))
	}
	return &types.MinterSupplyRequest{
		Address: args[0].(common.Address).Bytes(),
	}, nil
}

func NewMsgMint(args []interface{}) (*types.MsgMint, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 1, len(args))
	}
	return &types.MsgMint{
		Minter: args[0].(common.Address).Bytes(),
		Amount: args[1].(*big.Int).Bytes(),
	}, nil
}

func NewMsgBurn(args []interface{}) (*types.MsgBurn, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 1, len(args))
	}
	return &types.MsgBurn{
		Minter: args[0].(common.Address).Bytes(),
		Amount: args[1].(*big.Int).Bytes(),
	}, nil
}

package keeper

import (
	"context"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ types.QueryServer = Keeper{}

// MinterSupply implements types.QueryServer.
func (k Keeper) MinterSupply(c context.Context, request *types.MinterSupplyRequest) (*types.MinterSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	account := common.BytesToAddress(request.Address)
	cap, err := k.getMinterCap(ctx, account)
	if err != nil {
		return nil, err
	}
	supply, err := k.getMinterSupply(ctx, account)
	if err != nil {
		return nil, err
	}
	return &types.MinterSupplyResponse{
		Cap:    cap.Bytes(),
		Supply: supply.Bytes(),
	}, nil
}

// GetWA0GI implements types.QueryServer.
func (k Keeper) GetWA0GI(c context.Context, request *types.GetWA0GIRequest) (*types.GetWA0GIResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.GetWA0GIResponse{
		Address: k.GetWA0GIAddress(ctx),
	}, nil
}

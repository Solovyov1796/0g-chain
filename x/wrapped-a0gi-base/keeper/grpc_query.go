package keeper

import (
	"context"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
)

var _ types.QueryServer = Keeper{}

// MinterSupply implements types.QueryServer.
func (k Keeper) MinterSupply(c context.Context, request *types.MinterSupplyRequest) (*types.MinterSupplyResponse, error) {
	panic("unimplemented")
}

package keeper

import (
	"context"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
)

/*
	errorsmod "cosmossdk.io/errors"
	"github.com/0glabs/0g-chain/crypto/bn254util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	etherminttypes "github.com/evmos/ethermint/types"
*/

var _ types.MsgServer = &Keeper{}

// Burn implements types.MsgServer.
func (k Keeper) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	panic("unimplemented")
}

// Mint implements types.MsgServer.
func (k Keeper) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	panic("unimplemented")
}

// SetMinterCap implements types.MsgServer.
func (k Keeper) SetMinterCap(goCtx context.Context, msg *types.MsgSetMintCap) (*types.MsgSetMintCapResponse, error) {
	panic("unimplemented")
}

// SetWA0GI implements types.MsgServer.
func (k Keeper) SetWA0GI(goCtx context.Context, msg *types.MsgSetWA0GI) (*types.MsgSetWA0GIResponse, error) {
	panic("unimplemented")
}

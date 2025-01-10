package keeper

import (
	"context"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	precisebanktypes "github.com/0glabs/0g-chain/x/precisebank/types"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ types.MsgServer = &Keeper{}

// Burn implements types.MsgServer.
func (k Keeper) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	minter := common.BytesToAddress(msg.Minter)
	supply, err := k.getMinterSupply(ctx, minter)
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).SetBytes(msg.Amount)
	// check & update mint supply
	supply.Sub(supply, amount)
	if supply.Cmp(big.NewInt(0)) < 0 {
		return nil, types.ErrInsufficientMintSupply
	}
	// burn
	c := sdk.NewCoin(precisebanktypes.ExtendedCoinDenom, sdk.NewIntFromBigInt(amount))
	if err = k.pbkeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	if err = k.setMinterSupply(ctx, minter, supply); err != nil {
		return nil, err
	}
	return &types.MsgBurnResponse{}, nil
}

// Mint implements types.MsgServer.
func (k Keeper) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	minter := common.BytesToAddress(msg.Minter)
	cap, err := k.getMinterCap(ctx, minter)
	if err != nil {
		return nil, err
	}
	supply, err := k.getMinterSupply(ctx, minter)
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).SetBytes(msg.Amount)
	// check & update mint supply
	supply.Add(supply, amount)
	if supply.Cmp(cap) > 0 {
		return nil, types.ErrInsufficientMintCap
	}
	// mint
	c := sdk.NewCoin(precisebanktypes.ExtendedCoinDenom, sdk.NewIntFromBigInt(amount))
	if err = k.pbkeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	if err = k.setMinterSupply(ctx, minter, supply); err != nil {
		return nil, err
	}
	return &types.MsgMintResponse{}, nil
}

// SetMinterCap implements types.MsgServer.
func (k Keeper) SetMinterCap(goCtx context.Context, msg *types.MsgSetMintCap) (*types.MsgSetMintCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	minter := common.BytesToAddress(msg.Minter)
	// validate authority
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}
	// update minter cap
	if err := k.setMinterCap(ctx, minter, new(big.Int).SetBytes(msg.Cap)); err != nil {
		return nil, err
	}
	return &types.MsgSetMintCapResponse{}, nil
}

// SetWA0GI implements types.MsgServer.
func (k Keeper) SetWA0GI(goCtx context.Context, msg *types.MsgSetWA0GI) (*types.MsgSetWA0GIResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// validate authority
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}
	// save new address
	k.SetWA0GIAddress(ctx, common.BytesToAddress(msg.Address))
	return &types.MsgSetWA0GIResponse{}, nil
}

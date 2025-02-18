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
	s, err := k.getMinterSupply(ctx, minter)
	if err != nil {
		return nil, err
	}
	supply := new(big.Int).SetBytes(s.Supply)

	amount := new(big.Int).SetBytes(msg.Amount)
	// check & update mint supply
	supply.Sub(supply, amount)
	if supply.Cmp(big.NewInt(0)) < 0 {
		return nil, types.ErrInsufficientMintSupply
	}
	// transfer from wa0gi contract address & burn
	c := sdk.NewCoin(precisebanktypes.ExtendedCoinDenom, sdk.NewIntFromBigInt(amount))
	wa0gi := sdk.AccAddress(k.GetWA0GIAddress(ctx))
	if err = k.pbkeeper.SendCoinsFromAccountToModule(ctx, wa0gi, types.ModuleName, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	if err = k.pbkeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	// update supply
	s.Supply = supply.Bytes()
	if err = k.setMinterSupply(ctx, minter, s); err != nil {
		return nil, err
	}
	return &types.MsgBurnResponse{}, nil
}

// Mint implements types.MsgServer.
func (k Keeper) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	minter := common.BytesToAddress(msg.Minter)
	s, err := k.getMinterSupply(ctx, minter)
	if err != nil {
		return nil, err
	}
	supply := new(big.Int).SetBytes(s.Supply)
	cap := new(big.Int).SetBytes(s.Cap)

	amount := new(big.Int).SetBytes(msg.Amount)
	// check & update mint supply
	supply.Add(supply, amount)
	if supply.Cmp(cap) > 0 {
		return nil, types.ErrInsufficientMintCap
	}
	// mint & transfer to wa0gi contract address
	c := sdk.NewCoin(precisebanktypes.ExtendedCoinDenom, sdk.NewIntFromBigInt(amount))
	if err = k.pbkeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	wa0gi := sdk.AccAddress(k.GetWA0GIAddress(ctx))
	if err = k.pbkeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, wa0gi, sdk.NewCoins(c)); err != nil {
		return nil, err
	}
	// update supply
	s.Supply = supply.Bytes()
	if err = k.setMinterSupply(ctx, minter, s); err != nil {
		return nil, err
	}
	return &types.MsgMintResponse{}, nil
}

// SetMinterCap implements types.MsgServer.
func (k Keeper) SetMinterCap(goCtx context.Context, msg *types.MsgSetMinterCap) (*types.MsgSetMinterCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	minter := common.BytesToAddress(msg.Minter)
	// validate authority
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}
	// get previous minter supply
	s, err := k.getMinterSupply(ctx, minter)
	if err != nil {
		return nil, err
	}
	supply := new(big.Int).SetBytes(s.Supply)
	currentInitialSupply := new(big.Int).SetBytes(s.InitialSupply)
	newInitialSupply := new(big.Int).SetBytes(msg.InitialSupply)
	difference := new(big.Int).Sub(newInitialSupply, currentInitialSupply)
	supply.Add(supply, difference)
	if supply.Cmp(big.NewInt(0)) < 0 {
		supply = big.NewInt(0)
	}

	s.Cap = msg.Cap
	s.InitialSupply = msg.InitialSupply
	s.Supply = supply.Bytes()

	// update minter supply
	if err := k.setMinterSupply(ctx, minter, s); err != nil {
		return nil, err
	}
	return &types.MsgSetMinterCapResponse{}, nil
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

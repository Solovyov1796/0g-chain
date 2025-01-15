package keeper

import (
	"math/big"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"

	precisebankkeeper "github.com/0glabs/0g-chain/x/precisebank/keeper"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey  storetypes.StoreKey
	cdc       codec.BinaryCodec
	pbkeeper  precisebankkeeper.Keeper
	authority string // the address capable of changing signers params. Should be the gov module account
}

// NewKeeper creates a new wrapped a0gi base keeper instance
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	pbkeeper precisebankkeeper.Keeper,
	authority string,
) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		pbkeeper:  pbkeeper,
		authority: authority,
	}
}

func (k Keeper) SetWA0GIAddress(ctx sdk.Context, addr common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.WA0GIKey, addr.Bytes())
}

func (k Keeper) GetWA0GIAddress(ctx sdk.Context) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.WA0GIKey)
	return bz
}

func (k Keeper) setMinterCap(ctx sdk.Context, account common.Address, cap *big.Int) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MinterCapKeyPrefix)
	store.Set(account.Bytes(), cap.Bytes())
	return nil
}

func (k Keeper) getMinterCap(ctx sdk.Context, account common.Address) (*big.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MinterCapKeyPrefix)
	bz := store.Get(account.Bytes())
	return new(big.Int).SetBytes(bz), nil
}

func (k Keeper) setMinterSupply(ctx sdk.Context, account common.Address, supply *big.Int) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MinterSupplyKeyPrefix)
	store.Set(account.Bytes(), supply.Bytes())
	return nil
}

func (k Keeper) getMinterSupply(ctx sdk.Context, account common.Address) (*big.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MinterSupplyKeyPrefix)
	bz := store.Get(account.Bytes())
	return new(big.Int).SetBytes(bz), nil
}

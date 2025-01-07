package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	storeKey  storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string // the address capable of changing signers params. Should be the gov module account
}

// NewKeeper creates a new wrapped a0gi base keeper instance
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	authority string,
) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		authority: authority,
	}
}

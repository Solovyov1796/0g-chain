package wrappeda0gibase

import (
	"fmt"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/keeper"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the store state from a genesis state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, gs types.GenesisState) {
	if err := gs.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", types.ModuleName, err))
	}
	panic("unimplemented")
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	panic("unimplemented")
}

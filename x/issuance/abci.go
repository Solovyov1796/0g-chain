package issuance

import (
	"time"

	"github.com/0glabs/0g-chain/x/issuance/keeper"
	"github.com/0glabs/0g-chain/x/issuance/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker iterates over each asset and seizes coins from blocked addresses by returning them to the asset owner
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	err := k.SeizeCoinsForBlockableAssets(ctx)
	if err != nil {
		panic(err)
	}
	k.SynchronizeBlockList(ctx)
	k.UpdateTimeBasedSupplyLimits(ctx)
}

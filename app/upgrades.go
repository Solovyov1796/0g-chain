package app

import (
	"fmt"

	wrappeda0gibasetypes "github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName_Testnet = "v0.5.0"
)

// RegisterUpgradeHandlers registers the upgrade handlers for the app.
func (app App) RegisterUpgradeHandlers() {
	app.upgradeKeeper.SetUpgradeHandler(
		UpgradeName_Testnet,
		upgradeHandler(app, UpgradeName_Testnet),
	)

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	doUpgrade := upgradeInfo.Name == UpgradeName_Testnet

	if doUpgrade && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				wrappeda0gibasetypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// upgradeHandler returns an UpgradeHandler for the given upgrade parameters.
func upgradeHandler(
	app App,
	name string,
) upgradetypes.UpgradeHandler {
	return func(
		ctx sdk.Context,
		plan upgradetypes.Plan,
		fromVM module.VersionMap,
	) (module.VersionMap, error) {
		logger := app.Logger()
		logger.Info(fmt.Sprintf("running %s upgrade handler", name))

		// Run migrations for all modules and return new consensus version map.
		versionMap, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			return nil, err
		}

		app.mintKeeper.InitGenesis(ctx, app.accountKeeper, &minttypes.GenesisState{
			BondDenom: "ua0gi",
		})

		logger.Info("completed store migrations")

		return versionMap, nil
	}
}

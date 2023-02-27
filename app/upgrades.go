package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type UpgradeCreator struct {
	name    string
	handler func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler
}

var DummyUpgradeHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

var ExampleUpgradeHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		// Update some params in the state as an example
		params := app.ValidatorKeeper.GetParams(ctx)
		params.MaxEntries++
		app.ValidatorKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

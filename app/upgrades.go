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

var FixSendUpgradeHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		adrSender := sdk.MustAccAddressFromBech32("d01k2mj3dkq3llcr9f3pkg2vzla8n24dq3uwjrrsw")
		adrReceiver := sdk.MustAccAddressFromBech32("d01dn8uqlcpjxtvvdzt7zwlcaw9hxvzx0vfk7y67r")
		coins := app.BankKeeper.GetAllBalances(ctx, adrSender)
		app.BankKeeper.SendCoins(ctx, adrSender, adrReceiver, coins)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

var MigrationUpgradeHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		adrSender := sdk.MustAccAddressFromBech32("d01k2mj3dkq3llcr9f3pkg2vzla8n24dq3uwjrrsw")
		adrReceiver := sdk.MustAccAddressFromBech32("d01dn8uqlcpjxtvvdzt7zwlcaw9hxvzx0vfk7y67r")
		coins := app.BankKeeper.GetAllBalances(ctx, adrSender)
		app.BankKeeper.SendCoins(ctx, adrSender, adrReceiver, coins)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

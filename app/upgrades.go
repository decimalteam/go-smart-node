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
		adrReceive := sdk.MustAccAddressFromBech32("d01x5gemuf2lhuy3dl0y5m27fq34v9xz8z9fnd9zm")
		// send coin reserve
		coinPool := app.AccountKeeper.GetModuleAccount(ctx, "coin").GetAddress()
		coins := app.BankKeeper.GetAllBalances(ctx, coinPool)
		err := app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "coin", adrReceive, coins)
		if err != nil {
			return nil, err
		}

		nbPool := app.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress()
		coins = app.BankKeeper.GetAllBalances(ctx, nbPool)
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "not_bonded_tokens_pool", adrReceive, coins)
		if err != nil {
			return nil, err
		}
		bPool := app.ValidatorKeeper.GetBondedPool(ctx).GetAddress()
		coins = app.BankKeeper.GetAllBalances(ctx, bPool)
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "bonded_tokens_pool", adrReceive, coins)
		if err != nil {
			return nil, err
		}

		nftPool := app.AccountKeeper.GetModuleAccount(ctx, "reserved_pool").GetAddress()
		coins = app.BankKeeper.GetAllBalances(ctx, nftPool)
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "reserved_pool", adrReceive, coins)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

var UpdateRewardAndMaxVars = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		//params := app.ValidatorKeeper.GetParams(ctx)
		//params.MaxEntries = 10
		//
		//app.ValidatorKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// var FixValidatorDuplicates = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
// 	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

// 		validators := app.GetStakingKeeper().GetAllValidators(ctx)

// 			for _, validator := range validators {
				
// 				// store := ctx.KVStore(app.)
// 				store := ctx.KVStore(app.GetKey(stakingtypes.StoreKey))
				
// 				deleted := false

// 				iterator := sdk.KVStorePrefixIterator(store, stakingtypes.ValidatorsByPowerIndexKey)
// 				defer iterator.Close()

// 				for ; iterator.Valid(); iterator.Next() {
// 					valAddr := stakingtypes.ParseValidatorPowerRankKey(iterator.Key())
// 					if bytes.Equal(valAddr, validator.GetOperator()) {
// 						if deleted {
// 							panic("found duplicate power index key")
// 						} else {
// 							deleted = true
// 						}

// 						store.Delete(iterator.Key())
// 					}
// 				}

// 				app.GetStakingKeeper().SetValidatorByPowerIndex(ctx, validator)
// 				_, err := app.GetStakingKeeper().ApplyAndReturnValidatorSetUpdates(ctx)

// 				if err != nil {
// 					panic(err)
// 				}
// 			}

// 			return mm.RunMigrations(ctx, configurator, fromVM)
// 	}
// }
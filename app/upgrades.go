package app

import (
	"bytes"
	"fmt"
	"time"

	dsctypes "bitbucket.org/decimalteam/go-smart-node/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
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

var ValidatorDuplicatesHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", "v2.2.3")

		validators := app.GetStakingKeeper().GetAllValidators(ctx)

		logger.Info("Start changing validators.")

		for _, validator := range validators {

			store := ctx.KVStore(app.GetKey(validatortypes.StoreKey))

			deleted := false

			iterator := sdk.KVStorePrefixIterator(store, validatortypes.GetValidatorsByPowerIndexKey())
			defer iterator.Close()

			for ; iterator.Valid(); iterator.Next() {
				valAddr := validatortypes.ParseValidatorPowerKey(iterator.Key())
				val := sdk.ValAddress(valAddr).String()

				if bytes.Equal(valAddr, validator.GetOperator()) {
					if deleted {
						logger.Info("Duplicate validator address is: " + val)
					} else {
						deleted = true
					}
					store.Delete(iterator.Key())
				}
			}

			app.GetStakingKeeper().SetValidatorByPowerIndex(ctx, validator)
			_, err := app.GetStakingKeeper().ApplyAndReturnValidatorSetUpdates(ctx)
			if err != nil {
				panic(err)
			}

		}

		logger.Info("Updated all validator successfully.")

		return app.mm.RunMigrations(ctx, configurator, fromVM)
	}
}

var TransferDaoAndVals = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		var oldDaoAccount = sdk.MustAccAddressFromBech32("d01pk2rurh73er88p032qrd6kq5xmu53thjqc22mu")
		var oldDevelopAccount = sdk.MustAccAddressFromBech32("d01gsa4w0cuyjqwt9j7qtc32m6n0lkyxfan9s2ghh")

		var newDaoAccount = sdk.MustAccAddressFromBech32("d01zafwcqd3vwcjmtcfgwnt37r02ta38mr9w0da3k")
		var newDevelopAccount = sdk.MustAccAddressFromBech32("d01hv3zxnm2x4sgnyaap7luwt783c04xxjfdlnt9u")

		if err := app.BankKeeper.SendCoins(
			ctx, oldDevelopAccount, newDevelopAccount, app.BankKeeper.GetAllBalances(ctx, oldDevelopAccount),
		); err != nil {
			panic(err)
		}

		if err := app.BankKeeper.SendCoins(
			ctx, oldDaoAccount, newDaoAccount, app.BankKeeper.GetAllBalances(ctx, oldDaoAccount),
		); err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

var AutoUnbondMigrationHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)

		// Set OfflineSince = ctx.BlockTime() for all currently offline validators.
		// This gives them the full auto-unbond timeout window from the upgrade moment.
		validators := app.ValidatorKeeper.GetAllValidators(ctx)
		count := 0
		for _, val := range validators {
			if !val.Online {
				valAddr := val.GetOperator()
				if _, found := app.ValidatorKeeper.GetValidatorOfflineSince(ctx, valAddr); !found {
					app.ValidatorKeeper.SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())
					count++
				}
			}
		}
		logger.Info(fmt.Sprintf("auto-unbond migration: initialized OfflineSince for %d offline validators", count))

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// StakeMigration defines a single old->new address pair for stake migration.
type StakeMigration struct {
	OldHex string
	NewHex string
}

// NewMigrateStakesHandler builds an upgrade handler that migrates all delegations
// from old addresses to new addresses.
// Approach: set UndelegationTime=0, undelegate from old, complete instantly, delegate to new, restore time.
// This preserves validator power and CustomCoinStaked (subtract then add = net zero).
// NOTE: This only migrates Cosmos-side state. EVM-side delegation contract state must be updated
// separately via admin migrateStakes() call, since CallEVM does not trigger PostTxProcessing hooks.
func NewMigrateStakesHandler(migrations []StakeMigration) func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
		return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			logger := ctx.Logger().With("upgrade", plan.Name)

			// 1. Save original undelegation time and set to 0 for instant completion
			// Use Subspace.Set directly to bypass SetParams validation (which rejects 0)
			validatorSubspace := app.GetSubspace(validatortypes.ModuleName)
			var originalUndelegationTime time.Duration
			validatorSubspace.Get(ctx, validatortypes.KeyUndelegationTime, &originalUndelegationTime)
			validatorSubspace.Set(ctx, validatortypes.KeyUndelegationTime, time.Duration(0))

			for _, m := range migrations {
				oldAddr, err := dsctypes.GetDecimalAddressFromHex(m.OldHex)
				if err != nil {
					panic(fmt.Errorf("invalid old address %s: %w", m.OldHex, err))
				}
				newAddr, err := dsctypes.GetDecimalAddressFromHex(m.NewHex)
				if err != nil {
					panic(fmt.Errorf("invalid new address %s: %w", m.NewHex, err))
				}

				// Get all delegations from old delegator
				delegations := app.ValidatorKeeper.GetDelegatorDelegations(ctx, oldAddr, 65535)
				logger.Info(fmt.Sprintf("Migrating %d delegations from %s to %s", len(delegations), m.OldHex, m.NewHex))

				for _, del := range delegations {
					valAddr := del.GetValidator()
					stake := del.Stake

					// Create zero remain stake for full unbond
					var remainStake validatortypes.Stake
					switch stake.Type {
					case validatortypes.StakeType_Coin:
						remainStake = validatortypes.NewStakeCoin(sdk.Coin{Denom: stake.Stake.Denom, Amount: sdk.ZeroInt()})
					case validatortypes.StakeType_NFT:
						remainStake = validatortypes.NewStakeNFT(stake.ID, nil, sdk.Coin{Denom: stake.Stake.Denom, Amount: sdk.ZeroInt()})
					}

					// 2a. Undelegate from old address (instant with period=0)
					_, err = app.ValidatorKeeper.Undelegate(ctx, oldAddr, valAddr, stake, remainStake, nil)
					if err != nil {
						panic(fmt.Errorf("failed to undelegate %s from validator %s: %w", stake.ID, valAddr, err))
					}

					// 2b. Complete unbonding immediately (period=0 means entry is already mature)
					err = app.ValidatorKeeper.CompleteUnbonding(ctx, oldAddr, valAddr)
					if err != nil {
						panic(fmt.Errorf("failed to complete unbonding for %s: %w", stake.ID, err))
					}

					// 2c. Delegate to new address (handles merging if delegation already exists)
					validator, found := app.ValidatorKeeper.GetValidator(ctx, valAddr)
					if !found {
						panic(fmt.Errorf("validator not found: %s", valAddr))
					}

					err = app.ValidatorKeeper.Delegate(ctx, newAddr, validator, stake)
					if err != nil {
						panic(fmt.Errorf("failed to delegate %s to new address: %w", stake.ID, err))
					}

					logger.Info(fmt.Sprintf("  Migrated stake %s (%s) on validator %s", stake.ID, stake.Stake.Amount, valAddr))
				}
			}

			// 3. Restore original undelegation time
			validatorSubspace.Set(ctx, validatortypes.KeyUndelegationTime, originalUndelegationTime)

			logger.Info("Stakes migration completed successfully")

			return mm.RunMigrations(ctx, configurator, fromVM)
		}
	}
}

// Mainnet migration addresses
var MainnetStakeMigrations = []StakeMigration{
	{"0x35679b820c6318a159e4700b94645b90819b5bce", "0x50bd8f9af4c26bc8083cea3db84730dc2ac7412a"},
	{"0x227b033e84d038d6f6148d0ab92d2beb70e92ddc", "0x091dfd363d8721918e1285baecf16dafae38a70d"},
	{"0x31cb3bdccf4b3fceed1a84c2d745f59e64e9e6ae", "0xb23c6ba9dd63924d1d94b041a73f5d34c9467aff"},
	{"0x972ee514def1004f99413f954a07a34b2db9f128", "0x10cb963ff74a1134a81dfb10586236cf61ef22c0"},
}

// Testnet migration addresses
var TestnetStakeMigrations = []StakeMigration{
	{"0x0b94c822596dbb629be24dab4c5e0d6aa79a2980", "0xb48eb2301ca6a5af54b11e1769fae5a5b6d63c8b"},
	{"0xe2d9b3af39c9a9b74aea41b7bc85e2aec7e975d9", "0xa13b78f8aa52c2f0ffd13326a1ba4eab977e70c0"},
	{"0x1d1db8d3264b667759d737a95a39d66852e43549", "0xd5602db45634282a6f545675e7dce76d2a40ad1d"},
	{"0x8d76bea7bb18f9b92157e0d45414c9c8a3793386", "0xa1bd770ddf1bbeda438938bbf90fe6caefdca8eb"},
}

// MigrateStakesHandlerCreator is the mainnet handler (kept for backwards compatibility with upgradeslist.go)
var MigrateStakesHandlerCreator = NewMigrateStakesHandler(MainnetStakeMigrations)

// CombinedTestnetUpgradeHandlerCreator runs testnet migrations in a single upgrade:
// 1. Migrate stakes from old delegator addresses to new ones
// 2. Initialize auto-unbond tracking for offline validators
var CombinedTestnetUpgradeHandlerCreator = func(app *DSC, mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)

		// --- Step 1: Migrate stakes ---
		logger.Info("step 1: migrating stakes")
		{
			validatorSubspace := app.GetSubspace(validatortypes.ModuleName)
			var originalUndelegationTime time.Duration
			validatorSubspace.Get(ctx, validatortypes.KeyUndelegationTime, &originalUndelegationTime)
			validatorSubspace.Set(ctx, validatortypes.KeyUndelegationTime, time.Duration(0))

			for _, m := range TestnetStakeMigrations {
				oldAddr, err := dsctypes.GetDecimalAddressFromHex(m.OldHex)
				if err != nil {
					panic(fmt.Errorf("invalid old address %s: %w", m.OldHex, err))
				}
				newAddr, err := dsctypes.GetDecimalAddressFromHex(m.NewHex)
				if err != nil {
					panic(fmt.Errorf("invalid new address %s: %w", m.NewHex, err))
				}

				delegations := app.ValidatorKeeper.GetDelegatorDelegations(ctx, oldAddr, 65535)
				logger.Info(fmt.Sprintf("migrating %d delegations from %s to %s", len(delegations), m.OldHex, m.NewHex))

				for _, del := range delegations {
					valAddr := del.GetValidator()
					stake := del.Stake

					var remainStake validatortypes.Stake
					switch stake.Type {
					case validatortypes.StakeType_Coin:
						remainStake = validatortypes.NewStakeCoin(sdk.Coin{Denom: stake.Stake.Denom, Amount: sdk.ZeroInt()})
					case validatortypes.StakeType_NFT:
						remainStake = validatortypes.NewStakeNFT(stake.ID, nil, sdk.Coin{Denom: stake.Stake.Denom, Amount: sdk.ZeroInt()})
					}

					_, err = app.ValidatorKeeper.Undelegate(ctx, oldAddr, valAddr, stake, remainStake, nil)
					if err != nil {
						panic(fmt.Errorf("failed to undelegate %s from validator %s: %w", stake.ID, valAddr, err))
					}

					err = app.ValidatorKeeper.CompleteUnbonding(ctx, oldAddr, valAddr)
					if err != nil {
						panic(fmt.Errorf("failed to complete unbonding for %s: %w", stake.ID, err))
					}

					validator, found := app.ValidatorKeeper.GetValidator(ctx, valAddr)
					if !found {
						panic(fmt.Errorf("validator not found: %s", valAddr))
					}

					err = app.ValidatorKeeper.Delegate(ctx, newAddr, validator, stake)
					if err != nil {
						panic(fmt.Errorf("failed to delegate %s to new address: %w", stake.ID, err))
					}

					logger.Info(fmt.Sprintf("  migrated stake %s (%s) on validator %s", stake.ID, stake.Stake.Amount, valAddr))
				}
			}

			validatorSubspace.Set(ctx, validatortypes.KeyUndelegationTime, originalUndelegationTime)
			logger.Info("step 1 complete: stakes migrated")
		}

		// --- Step 2: Initialize auto-unbond tracking ---
		logger.Info("step 2: initializing auto-unbond tracking for offline validators")
		{
			validators := app.ValidatorKeeper.GetAllValidators(ctx)
			count := 0
			for _, val := range validators {
				if !val.Online {
					valAddr := val.GetOperator()
					if _, found := app.ValidatorKeeper.GetValidatorOfflineSince(ctx, valAddr); !found {
						app.ValidatorKeeper.SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())
						count++
					}
				}
			}
			logger.Info(fmt.Sprintf("step 2 complete: initialized OfflineSince for %d offline validators", count))
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

package app

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	dsctypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var testPKs = simapp.CreateTestPubKeys(500)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}


func TestMigrateStakesHandler(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	valK := dsc.ValidatorKeeper
	defaultCoins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000000000))))

	// Use the first address pair from the handler's hardcoded list
	oldHex := "0x35679b820c6318a159e4700b94645b90819b5bce"
	newHex := "0x50bd8f9af4c26bc8083cea3db84730dc2ac7412a"

	oldAddr, err := dsctypes.GetDecimalAddressFromHex(oldHex)
	require.NoError(t, err)
	newAddr, err := dsctypes.GetDecimalAddressFromHex(newHex)
	require.NoError(t, err)

	// Fund the old address
	initAccountWithCoins(dsc, ctx, oldAddr, defaultCoins)

	// Create a validator using a separate test address
	valAccounts := AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	valAddrs := ConvertAddrsToValAddrs(valAccounts)

	val, err := validatortypes.NewValidator(valAddrs[0], valAccounts[0], testPKs[0], validatortypes.Description{Moniker: "test-val"}, sdk.ZeroDec())
	require.NoError(t, err)
	val.Status = validatortypes.BondStatus_Bonded
	val.Online = true
	val.Stake = 10
	valK.CreateValidator(ctx, val)

	// Delegate from the hardcoded old address
	stakeAmount := keeper.TokensFromConsensusPower(5)
	stake := validatortypes.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, stakeAmount))
	err = valK.Delegate(ctx, oldAddr, val, stake)
	require.NoError(t, err)

	// Verify pre-conditions
	oldDels := valK.GetDelegatorDelegations(ctx, oldAddr, 100)
	require.Len(t, oldDels, 1)
	require.Equal(t, stakeAmount, oldDels[0].Stake.Stake.Amount)

	newDels := valK.GetDelegatorDelegations(ctx, newAddr, 100)
	require.Len(t, newDels, 0)

	rsBefore, err := valK.GetValidatorRS(ctx, valAddrs[0])
	require.NoError(t, err)
	powerBefore := rsBefore.Stake

	originalUndelegationTime := valK.UndelegationTime(ctx)

	// Call the actual handler
	handler := MigrateStakesHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err = handler(ctx, upgradetypes.Plan{Name: "test-migrate"}, fromVM)
	require.NoError(t, err)

	// Verify: old address has no delegations
	oldDels = valK.GetDelegatorDelegations(ctx, oldAddr, 100)
	require.Len(t, oldDels, 0, "old address should have no delegations after handler")

	// Verify: new address has the delegation
	newDels = valK.GetDelegatorDelegations(ctx, newAddr, 100)
	require.Len(t, newDels, 1, "new address should have 1 delegation after handler")
	require.Equal(t, stakeAmount, newDels[0].Stake.Stake.Amount, "delegation amount preserved")
	require.Equal(t, val.GetOperator(), newDels[0].GetValidator(), "same validator")

	// Verify: validator power preserved
	rsAfter, err := valK.GetValidatorRS(ctx, valAddrs[0])
	require.NoError(t, err)
	require.Equal(t, powerBefore, rsAfter.Stake, "validator power preserved")

	// Verify: undelegation time restored
	require.Equal(t, originalUndelegationTime, valK.UndelegationTime(ctx), "undelegation time restored")
}

func TestMigrateStakesHandlerMultipleAddresses(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	valK := dsc.ValidatorKeeper
	defaultCoins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000000000))))

	// All 4 address pairs from the handler
	type addrPair struct {
		oldHex string
		newHex string
	}
	migrations := []addrPair{
		{"0x35679b820c6318a159e4700b94645b90819b5bce", "0x50bd8f9af4c26bc8083cea3db84730dc2ac7412a"},
		{"0x227b033e84d038d6f6148d0ab92d2beb70e92ddc", "0x091dfd363d8721918e1285baecf16dafae38a70d"},
		{"0x31cb3bdccf4b3fceed1a84c2d745f59e64e9e6ae", "0xb23c6ba9dd63924d1d94b041a73f5d34c9467aff"},
		{"0x972ee514def1004f99413f954a07a34b2db9f128", "0x10cb963ff74a1134a81dfb10586236cf61ef22c0"},
	}

	// Create 2 validators
	valAccounts := AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	valAddrs := ConvertAddrsToValAddrs(valAccounts)

	var validators [2]validatortypes.Validator
	for i := 0; i < 2; i++ {
		validators[i], _ = validatortypes.NewValidator(valAddrs[i], valAccounts[i], testPKs[i], validatortypes.Description{Moniker: "val"}, sdk.ZeroDec())
		validators[i].Status = validatortypes.BondStatus_Bonded
		validators[i].Online = true
		validators[i].Stake = 10
		valK.CreateValidator(ctx, validators[i])
	}

	// Set up delegations from all 4 old addresses to different validators
	stakeAmounts := []int64{3, 5, 7, 11}
	for i, m := range migrations {
		oldAddr, err := dsctypes.GetDecimalAddressFromHex(m.oldHex)
		require.NoError(t, err)
		initAccountWithCoins(dsc, ctx, oldAddr, defaultCoins)

		valIdx := i % 2
		stake := validatortypes.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(stakeAmounts[i])))
		err = valK.Delegate(ctx, oldAddr, validators[valIdx], stake)
		require.NoError(t, err)
	}

	// Record validator powers before
	rs0Before, _ := valK.GetValidatorRS(ctx, validators[0].GetOperator())
	rs1Before, _ := valK.GetValidatorRS(ctx, validators[1].GetOperator())

	// Call the actual handler
	handler := MigrateStakesHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err := handler(ctx, upgradetypes.Plan{Name: "test-migrate"}, fromVM)
	require.NoError(t, err)

	// Verify all migrations
	for i, m := range migrations {
		oldAddr, _ := dsctypes.GetDecimalAddressFromHex(m.oldHex)
		newAddr, _ := dsctypes.GetDecimalAddressFromHex(m.newHex)

		// Old address empty
		oldDels := valK.GetDelegatorDelegations(ctx, oldAddr, 100)
		require.Len(t, oldDels, 0, "old address %d should have no delegations", i)

		// New address has delegation with correct amount
		newDels := valK.GetDelegatorDelegations(ctx, newAddr, 100)
		require.Len(t, newDels, 1, "new address %d should have 1 delegation", i)

		expectedAmount := keeper.TokensFromConsensusPower(stakeAmounts[i])
		require.Equal(t, expectedAmount, newDels[0].Stake.Stake.Amount, "delegation %d amount preserved", i)
	}

	// Verify validator powers preserved
	rs0After, _ := valK.GetValidatorRS(ctx, validators[0].GetOperator())
	rs1After, _ := valK.GetValidatorRS(ctx, validators[1].GetOperator())
	require.Equal(t, rs0Before.Stake, rs0After.Stake, "validator[0] power preserved")
	require.Equal(t, rs1Before.Stake, rs1After.Stake, "validator[1] power preserved")
}

func TestMigrateStakesHandlerNoOp(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	// Call handler on clean chain - no delegations to migrate, should succeed
	handler := MigrateStakesHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err := handler(ctx, upgradetypes.Plan{Name: "test-noop"}, fromVM)
	require.NoError(t, err, "handler should succeed with no delegations to migrate")

	// Undelegation time should be unchanged
	originalTime := validatortypes.DefaultUndelegationTime
	require.Equal(t, originalTime, dsc.ValidatorKeeper.UndelegationTime(ctx))
}

func TestDenomWhitelistValidation(t *testing.T) {
	// All whitelisted denoms must pass ValidateDenom
	for _, denom := range invalidDelegatedDenoms {
		require.NoError(t, sdk.ValidateDenom(denom), "whitelisted denom %q should pass ValidateDenom", denom)
	}
}

func TestDenomWhitelistNewCoin(t *testing.T) {
	// sdk.NewCoin must not panic for whitelisted denoms
	for _, denom := range invalidDelegatedDenoms {
		require.NotPanics(t, func() {
			sdk.NewCoin(denom, sdk.OneInt())
		}, "sdk.NewCoin should not panic for whitelisted denom %q", denom)
	}
}

func TestCombinedTestnetUpgradeHandler_NoOp(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Empty TestnetStakeMigrations, no offline validators → should succeed cleanly
	handler := CombinedTestnetUpgradeHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err := handler(ctx, upgradetypes.Plan{Name: "test-combined-noop"}, fromVM)
	require.NoError(t, err)

	// Undelegation time unchanged
	require.Equal(t, validatortypes.DefaultUndelegationTime, dsc.ValidatorKeeper.UndelegationTime(ctx))
}

func TestCombinedTestnetUpgradeHandler_StakeMigration(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	valK := dsc.ValidatorKeeper
	defaultCoins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000000000))))

	oldHex := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	newHex := "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	oldAddr, err := dsctypes.GetDecimalAddressFromHex(oldHex)
	require.NoError(t, err)
	newAddr, err := dsctypes.GetDecimalAddressFromHex(newHex)
	require.NoError(t, err)

	initAccountWithCoins(dsc, ctx, oldAddr, defaultCoins)

	// Create validator
	valAccounts := AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	valAddrs := ConvertAddrsToValAddrs(valAccounts)
	val, err := validatortypes.NewValidator(valAddrs[0], valAccounts[0], testPKs[0], validatortypes.Description{Moniker: "test-val"}, sdk.ZeroDec())
	require.NoError(t, err)
	val.Status = validatortypes.BondStatus_Bonded
	val.Online = true
	val.Stake = 10
	valK.CreateValidator(ctx, val)

	// Delegate from old address
	stakeAmount := keeper.TokensFromConsensusPower(5)
	stake := validatortypes.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, stakeAmount))
	err = valK.Delegate(ctx, oldAddr, val, stake)
	require.NoError(t, err)

	// Temporarily set TestnetStakeMigrations
	origMigrations := TestnetStakeMigrations
	TestnetStakeMigrations = []StakeMigration{{oldHex, newHex}}
	defer func() { TestnetStakeMigrations = origMigrations }()

	rsBefore, err := valK.GetValidatorRS(ctx, valAddrs[0])
	require.NoError(t, err)
	originalUndelegationTime := valK.UndelegationTime(ctx)

	handler := CombinedTestnetUpgradeHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err = handler(ctx, upgradetypes.Plan{Name: "test-combined-migrate"}, fromVM)
	require.NoError(t, err)

	// Old address empty
	oldDels := valK.GetDelegatorDelegations(ctx, oldAddr, 100)
	require.Len(t, oldDels, 0)

	// New address has delegation
	newDels := valK.GetDelegatorDelegations(ctx, newAddr, 100)
	require.Len(t, newDels, 1)
	require.Equal(t, stakeAmount, newDels[0].Stake.Stake.Amount)
	require.Equal(t, val.GetOperator(), newDels[0].GetValidator())

	// Validator power preserved
	rsAfter, err := valK.GetValidatorRS(ctx, valAddrs[0])
	require.NoError(t, err)
	require.Equal(t, rsBefore.Stake, rsAfter.Stake)

	// Undelegation time restored
	require.Equal(t, originalUndelegationTime, valK.UndelegationTime(ctx))
}

func TestCombinedTestnetUpgradeHandler_OfflineTracking(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	blockTime := time.Now()
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{
		Time: blockTime,
	})

	valK := dsc.ValidatorKeeper
	defaultCoins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000000000))))

	// Create 3 validators: 1 online, 2 offline
	valAccounts := AddTestAddrsIncremental(dsc, ctx, 3, defaultCoins)
	valAddrs := ConvertAddrsToValAddrs(valAccounts)

	for i := 0; i < 3; i++ {
		val, err := validatortypes.NewValidator(valAddrs[i], valAccounts[i], testPKs[i], validatortypes.Description{Moniker: "val"}, sdk.ZeroDec())
		require.NoError(t, err)
		val.Status = validatortypes.BondStatus_Bonded
		val.Online = i == 0 // only first is online
		val.Stake = 10
		valK.CreateValidator(ctx, val)
	}

	// Verify no offline tracking exists yet
	for i := 1; i < 3; i++ {
		_, found := valK.GetValidatorOfflineSince(ctx, valAddrs[i])
		require.False(t, found)
	}

	handler := CombinedTestnetUpgradeHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err := handler(ctx, upgradetypes.Plan{Name: "test-combined-offline"}, fromVM)
	require.NoError(t, err)

	// Online validator should NOT have offline tracking
	_, found := valK.GetValidatorOfflineSince(ctx, valAddrs[0])
	require.False(t, found, "online validator should not have offline tracking")

	// Offline validators should have tracking set to block time
	for i := 1; i < 3; i++ {
		ts, found := valK.GetValidatorOfflineSince(ctx, valAddrs[i])
		require.True(t, found, "offline validator %d should have tracking", i)
		require.Equal(t, blockTime.UTC(), ts.UTC(), "offline-since should equal block time")
	}
}

func TestCombinedTestnetUpgradeHandler_SkipsAlreadyTrackedOffline(t *testing.T) {
	dsc := Setup(t, false, feemarkettypes.DefaultGenesisState())
	blockTime := time.Now()
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{
		Time: blockTime,
	})

	valK := dsc.ValidatorKeeper
	defaultCoins := sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000000000))))

	valAccounts := AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	valAddrs := ConvertAddrsToValAddrs(valAccounts)

	val, err := validatortypes.NewValidator(valAddrs[0], valAccounts[0], testPKs[0], validatortypes.Description{Moniker: "val"}, sdk.ZeroDec())
	require.NoError(t, err)
	val.Status = validatortypes.BondStatus_Bonded
	val.Online = false
	val.Stake = 10
	valK.CreateValidator(ctx, val)

	// Pre-set offline-since to an earlier time
	earlierTime := blockTime.Add(-24 * time.Hour)
	valK.SetValidatorOfflineSince(ctx, valAddrs[0], earlierTime)

	handler := CombinedTestnetUpgradeHandlerCreator(dsc, dsc.mm, dsc.configurator)
	fromVM := dsc.mm.GetVersionMap()
	_, err = handler(ctx, upgradetypes.Plan{Name: "test-combined-skip"}, fromVM)
	require.NoError(t, err)

	// Should NOT overwrite existing tracking
	ts, found := valK.GetValidatorOfflineSince(ctx, valAddrs[0])
	require.True(t, found)
	require.Equal(t, earlierTime.UTC(), ts.UTC(), "pre-existing offline-since should not be overwritten")
}

func TestDenomWhitelistDoesNotBreakStandardDenoms(t *testing.T) {
	// Standard denoms must still pass
	for _, denom := range []string{"del", "btc", "eth", "usdt", "atom"} {
		require.NoError(t, sdk.ValidateDenom(denom), "standard denom %q should still pass", denom)
	}

	// Invalid denoms not in the whitelist must still fail
	require.Error(t, sdk.ValidateDenom(""), "empty denom should fail")
	require.Error(t, sdk.ValidateDenom("a"), "too short denom should fail")
	require.Error(t, sdk.ValidateDenom("1abc"), "denom starting with digit should fail")
}

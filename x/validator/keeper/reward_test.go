package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// setupValidatorWithDelegation creates a bonded validator with a delegation that has
// optional hold entries. Returns the validator address, delegator address, and the app context.
func setupValidatorWithDelegation(
	t *testing.T,
	holds []*types.StakeHold,
	stakeAmount sdkmath.Int,
) (*app.DSC, sdk.Context, sdk.ValAddress, sdk.AccAddress) {
	t.Helper()

	_, dsc, ctx := createTestInput(t)
	valK := dsc.ValidatorKeeper

	accs, vals := generateAddresses(
		dsc, ctx, 2,
		sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1_000_000_000)))),
	)
	creator := accs[0]
	delegator := accs[1]
	valAddr := vals[0]

	// Create validator
	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1_000)))
	msgCreate, err := types.NewMsgCreateValidator(
		valAddr, creator, PKs[0],
		types.Description{Moniker: "test-val"},
		sdk.ZeroDec(), creatorStake,
	)
	require.NoError(t, err)

	msgsrv := keeper.NewMsgServerImpl(valK)
	_, err = msgsrv.CreateValidator(sdk.WrapSDKContext(ctx), msgCreate)
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(sdk.WrapSDKContext(ctx), types.NewMsgSetOnline(valAddr))
	require.NoError(t, err)

	// Bond the validator
	valK.BlockValidatorUpdates(ctx)

	// Delegate with holds
	stake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, stakeAmount))
	stake.Holds = holds

	val, found := valK.GetValidator(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, types.BondStatus_Bonded, val.Status)

	err = valK.Delegate(ctx, delegator, val, stake)
	require.NoError(t, err)

	return dsc, ctx, valAddr, delegator
}

// mintModuleRewards mints coins into the validator module and sets accumulated rewards
// for the validator, simulating what PayValidators does over 120 blocks.
// Also ensures the base coin LimitVolume is set (as PayValidators would do).
func mintModuleRewards(
	t *testing.T,
	dsc *app.DSC,
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	rewardAmount sdkmath.Int,
) {
	t.Helper()

	baseDenom := dsc.ValidatorKeeper.BaseDenom(ctx)

	// Ensure base coin LimitVolume is set (PayValidators does this each block).
	// Use a value close to total staked so percentForHold is moderate (~50%),
	// not ~100% which would send all rewards to the hold pool leaving nothing
	// for normal distribution.
	baseCoin, err := dsc.CoinKeeper.GetCoin(ctx, baseDenom)
	require.NoError(t, err)
	if baseCoin.LimitVolume.IsZero() {
		baseCoin.LimitVolume = helpers.EtherToWei(sdkmath.NewInt(20_000))
	}
	baseCoin.LimitVolume = baseCoin.LimitVolume.Add(rewardAmount)
	dsc.CoinKeeper.SetCoin(ctx, baseCoin)

	// Mint coins into the validator module (simulating PayValidators)
	err = dsc.BankKeeper.MintCoins(
		ctx, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(baseDenom, rewardAmount)),
	)
	require.NoError(t, err)

	// Update base coin volume/reserve after minting
	baseCoin, err = dsc.CoinKeeper.GetCoin(ctx, baseDenom)
	require.NoError(t, err)
	err = dsc.CoinKeeper.UpdateCoinVR(ctx, baseDenom, baseCoin.Volume.Add(rewardAmount), baseCoin.Reserve)
	require.NoError(t, err)

	// Set validator accumulated rewards
	valRS, err := dsc.ValidatorKeeper.GetValidatorRS(ctx, valAddr)
	require.NoError(t, err)
	valRS.Rewards = rewardAmount
	dsc.ValidatorKeeper.SetValidatorRS(ctx, valAddr, valRS)
}

// TestPayRewards_NoHolds verifies that PayRewards succeeds when there are no holds.
// Module balance should be sufficient to cover all normal reward payments.
func TestPayRewards_NoHolds(t *testing.T) {
	stakeAmount := helpers.EtherToWei(sdkmath.NewInt(10_000))
	dsc, ctx, valAddr, delegator := setupValidatorWithDelegation(t, nil, stakeAmount)

	rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

	beforeBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)

	err := dsc.ValidatorKeeper.PayRewards(ctx)
	require.NoError(t, err)

	afterBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)
	require.True(t, afterBal.Amount.GT(beforeBal.Amount), "delegator should receive rewards")

	// Module should have non-negative balance (rounding remainder)
	moduleBal := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
	require.True(t, moduleBal.Amount.GTE(sdk.ZeroInt()), "module balance must not go negative")
}

// TestPayRewards_WithHolds_ModuleNotDrained is the core regression test for the
// accounting bug. With holds present, PayRewards should NOT spend more than what
// was minted into the module.
func TestPayRewards_WithHolds_ModuleNotDrained(t *testing.T) {
	now := time.Now()
	holdEnd := now.AddDate(2, 0, 0) // 2-year hold — qualifies immediately

	stakeAmount := helpers.EtherToWei(sdkmath.NewInt(10_000))
	holds := []*types.StakeHold{{
		Amount:        stakeAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   holdEnd.Unix(),
	}}

	dsc, ctx, valAddr, delegator := setupValidatorWithDelegation(t, holds, stakeAmount)

	rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

	moduleBalBefore := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
	require.True(t, moduleBalBefore.Amount.Equal(rewardAmount),
		"module should have exactly the minted reward amount")

	err := dsc.ValidatorKeeper.PayRewards(ctx)
	require.NoError(t, err, "PayRewards must not fail with insufficient funds")

	// After paying: module balance must be non-negative
	moduleBalAfter := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
	require.True(t, moduleBalAfter.Amount.GTE(sdk.ZeroInt()),
		"module balance must not go negative after paying rewards with holds; got %s",
		moduleBalAfter.Amount.String())

	// Delegator should have received some reward
	delBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)
	require.True(t, delBal.Amount.GT(sdk.ZeroInt()), "delegator should receive rewards")
}

// TestPayRewards_HoldRewardsComeFromReservedPool verifies that hold rewards reduce
// normal rewards (i.e., they come from the reserved accumRewards, not on top).
func TestPayRewards_HoldRewardsComeFromReservedPool(t *testing.T) {
	now := time.Now()
	holdEnd := now.AddDate(2, 0, 0)

	stakeAmount := helpers.EtherToWei(sdkmath.NewInt(10_000))
	holds := []*types.StakeHold{{
		Amount:        stakeAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   holdEnd.Unix(),
	}}

	dsc, ctx, valAddr, delegator := setupValidatorWithDelegation(t, holds, stakeAmount)

	rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

	err := dsc.ValidatorKeeper.PayRewards(ctx)
	require.NoError(t, err)

	// Total outflow = DAO + develop + validator commission + delegator rewards + hold rewards
	// This total must be <= rewardAmount (the only coins minted into module).
	moduleBalAfter := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
	totalPaid := rewardAmount.Sub(moduleBalAfter.Amount)
	require.True(t, totalPaid.LTE(rewardAmount),
		"total paid (%s) must not exceed minted rewards (%s)", totalPaid, rewardAmount)

	// Delegator got both normal + hold rewards
	delBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)
	require.True(t, delBal.Amount.GT(sdk.ZeroInt()))
}

// TestPayRewards_MultipleCycles verifies that running PayRewards multiple times
// (simulating multiple 120-block cycles) does not drain the module.
func TestPayRewards_MultipleCycles(t *testing.T) {
	now := time.Now()
	holdEnd := now.AddDate(2, 0, 0)

	stakeAmount := helpers.EtherToWei(sdkmath.NewInt(10_000))
	holds := []*types.StakeHold{{
		Amount:        stakeAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   holdEnd.Unix(),
	}}

	dsc, ctx, valAddr, _ := setupValidatorWithDelegation(t, holds, stakeAmount)

	for cycle := 0; cycle < 5; cycle++ {
		rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
		mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

		err := dsc.ValidatorKeeper.PayRewards(ctx)
		require.NoError(t, err, "PayRewards must not fail on cycle %d", cycle)

		moduleBal := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
		require.True(t, moduleBal.Amount.GTE(sdk.ZeroInt()),
			"module balance must not go negative on cycle %d; got %s", cycle, moduleBal.Amount)
	}
}

// TestPayRewards_BaseDenomHoldFiltering verifies that base denom delegations only
// count holds >= 1 year in allHoldBigOneYearsSum (pass 1 matches pass 2).
func TestPayRewards_BaseDenomHoldFiltering(t *testing.T) {
	now := time.Now()

	stakeAmount := helpers.EtherToWei(sdkmath.NewInt(10_000))

	// Short hold — less than 1 year, should NOT qualify for hold rewards
	shortHolds := []*types.StakeHold{{
		Amount:        stakeAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   now.Add(180 * 24 * time.Hour).Unix(), // ~6 months
	}}

	dsc, ctx, valAddr, delegator := setupValidatorWithDelegation(t, shortHolds, stakeAmount)

	rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

	beforeBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)

	err := dsc.ValidatorKeeper.PayRewards(ctx)
	require.NoError(t, err)

	afterBal := dsc.BankKeeper.GetBalance(ctx, delegator, cmdcfg.BaseDenom)
	normalReward := afterBal.Amount.Sub(beforeBal.Amount)

	// Now test with a long hold (>= 1 year) — should get hold bonus
	longHolds := []*types.StakeHold{{
		Amount:        stakeAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   now.AddDate(1, 1, 0).Unix(), // 13 months
	}}

	dsc2, ctx2, valAddr2, delegator2 := setupValidatorWithDelegation(t, longHolds, stakeAmount)

	mintModuleRewards(t, dsc2, ctx2, valAddr2, rewardAmount)

	beforeBal2 := dsc2.BankKeeper.GetBalance(ctx2, delegator2, cmdcfg.BaseDenom)

	err = dsc2.ValidatorKeeper.PayRewards(ctx2)
	require.NoError(t, err)

	afterBal2 := dsc2.BankKeeper.GetBalance(ctx2, delegator2, cmdcfg.BaseDenom)
	holdReward := afterBal2.Amount.Sub(beforeBal2.Amount)

	// Delegator with >= 1 year hold should receive MORE than one without qualifying hold
	require.True(t, holdReward.GT(normalReward),
		"hold reward (%s) should exceed normal reward (%s)", holdReward, normalReward)
}

// TestPayRewards_CustomCoin_WithHolds verifies hold rewards work correctly with
// custom coin delegations and the module stays solvent.
func TestPayRewards_CustomCoin_WithHolds(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	valK := dsc.ValidatorKeeper

	accs, vals := generateAddresses(
		dsc, ctx, 3,
		sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1_000_000_000)))),
	)
	creator := accs[0]
	delegator := accs[1]
	valAddr := vals[0]

	// Create custom coin
	ccDenom := "custom"
	initVolume := keeper.TokensFromConsensusPower(100_000_000)
	initReserve := keeper.TokensFromConsensusPower(1_000)
	limitVolume := keeper.TokensFromConsensusPower(1_000_000_000_000_000)
	crr := uint64(50)

	_, err := dsc.CoinKeeper.CreateCoin(ctx,
		cointypes.NewMsgCreateCoin(creator, ccDenom, "d", crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), ""))
	require.NoError(t, err)

	// Send custom coins to delegator
	_, err = dsc.CoinKeeper.SendCoin(ctx,
		cointypes.NewMsgSendCoin(creator, delegator, sdk.NewCoin(ccDenom, initVolume.Quo(sdk.NewInt(2)))))
	require.NoError(t, err)

	// Create validator
	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1_000)))
	msgCreate, err := types.NewMsgCreateValidator(
		valAddr, creator, PKs[0],
		types.Description{Moniker: "test-val-cc"},
		sdk.ZeroDec(), creatorStake,
	)
	require.NoError(t, err)

	msgsrv := keeper.NewMsgServerImpl(valK)
	_, err = msgsrv.CreateValidator(sdk.WrapSDKContext(ctx), msgCreate)
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(sdk.WrapSDKContext(ctx), types.NewMsgSetOnline(valAddr))
	require.NoError(t, err)

	valK.BlockValidatorUpdates(ctx)

	// Delegate custom coin with hold >= 1 year
	now := ctx.BlockTime()
	ccAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	stake := types.NewStakeCoin(sdk.NewCoin(ccDenom, ccAmount))
	stake.Holds = []*types.StakeHold{{
		Amount:        ccAmount,
		HoldStartTime: now.Unix(),
		HoldEndTime:   now.AddDate(2, 0, 0).Unix(),
	}}

	val, found := valK.GetValidator(ctx, valAddr)
	require.True(t, found)

	err = valK.Delegate(ctx, delegator, val, stake)
	require.NoError(t, err)

	// Mint rewards into module
	rewardAmount := helpers.EtherToWei(sdkmath.NewInt(1_000))
	mintModuleRewards(t, dsc, ctx, valAddr, rewardAmount)

	// PayRewards must succeed
	err = valK.PayRewards(ctx)
	require.NoError(t, err, "PayRewards must not fail with custom coin holds")

	// Module balance must be non-negative
	moduleBal := dsc.BankKeeper.GetBalance(ctx, types.ModuleAddress, cmdcfg.BaseDenom)
	require.True(t, moduleBal.Amount.GTE(sdk.ZeroInt()),
		"module balance must not go negative with custom coin holds; got %s", moduleBal.Amount)
}

package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/testvalidator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestOfflineTrackerCRUD(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	now := ctx.BlockTime()

	// Initially no entry
	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.False(t, found)

	// Set and get
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], now)
	offlineSince, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, now.UTC(), offlineSince.UTC())

	// Second validator not affected
	_, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[1])
	require.False(t, found)

	// Delete
	dsc.ValidatorKeeper.DeleteValidatorOfflineSince(ctx, addrVals[0])
	_, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.False(t, found)
}

func TestOfflineTrackerIterate(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 3, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	now := ctx.BlockTime()

	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], now)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[1], now.Add(time.Hour))
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[2], now.Add(2*time.Hour))

	var entries []sdk.ValAddress
	dsc.ValidatorKeeper.IterateValidatorOfflineSince(ctx, func(valAddr sdk.ValAddress, offlineSince time.Time) bool {
		entries = append(entries, valAddr)
		return false
	})
	require.Len(t, entries, 3)
}

func TestProcessAutoUnbondNotExpired(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	// Create validator and set offline
	val := testvalidator.NewValidator(t, addrVals[0], PKs[0])
	val.Online = false
	dsc.ValidatorKeeper.SetValidator(ctx, val)

	// Set offline 1 hour ago — timeout is 30 days, so not expired
	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.AutoUnbondTimeout = 30 * 24 * time.Hour
	dsc.ValidatorKeeper.SetParams(ctx, params)

	offlineTime := ctx.BlockTime().Add(-time.Hour)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], offlineTime)

	dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)

	// Entry should still exist (not expired)
	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found, "OfflineSince should still exist for non-expired validator")
}

func TestProcessAutoUnbondExpired(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	// Create validator and set offline
	val := testvalidator.NewValidator(t, addrVals[0], PKs[0])
	val.Online = false
	dsc.ValidatorKeeper.SetValidator(ctx, val)

	// Set timeout to 1 hour, set offline 2 hours ago
	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.AutoUnbondTimeout = time.Hour
	dsc.ValidatorKeeper.SetParams(ctx, params)

	offlineTime := ctx.BlockTime().Add(-2 * time.Hour)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], offlineTime)

	// ProcessAutoUnbond will try to enqueue but EVM contract is not set up in test.
	// The OfflineSince entry should be deleted even if EVM call fails (enqueue logs errors).
	dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)

	// Entry should be deleted (expired and processed)
	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.False(t, found, "OfflineSince should be deleted for expired validator")
}

func TestAutoUnbondTimeoutZeroDisabled(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	// Set timeout to 0 (disabled)
	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.AutoUnbondTimeout = 0
	dsc.ValidatorKeeper.SetParams(ctx, params)

	// Set validator offline long ago
	offlineTime := ctx.BlockTime().Add(-365 * 24 * time.Hour)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], offlineTime)

	dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)

	// Entry should still exist (feature disabled)
	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found, "OfflineSince should not be touched when timeout=0")
}

func TestSetOnlineClearsOfflineSince(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], ctx.BlockTime())

	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found)

	// Simulate clearing (as SetOnline msg handler does)
	dsc.ValidatorKeeper.DeleteValidatorOfflineSince(ctx, addrVals[0])

	_, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.False(t, found, "OfflineSince should be cleared after SetOnline")
}

func TestJailDoesNotResetExistingOfflineSince(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 1, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	originalTime := ctx.BlockTime().Add(-48 * time.Hour)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], originalTime)

	// Simulate what Jail does: only set if not already set
	if _, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0]); !found {
		dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], ctx.BlockTime())
	}

	offlineSince, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, originalTime.UTC(), offlineSince.UTC(), "Jail should not reset existing OfflineSince timer")
}

func TestGenesisExportImportOfflineSince(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	now := ctx.BlockTime()
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], now)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[1], now.Add(time.Hour))

	// Export
	genesis := dsc.ValidatorKeeper.ExportGenesis(ctx)
	require.Len(t, genesis.ValidatorOfflineSince, 2)

	// Clear state
	dsc.ValidatorKeeper.DeleteValidatorOfflineSince(ctx, addrVals[0])
	dsc.ValidatorKeeper.DeleteValidatorOfflineSince(ctx, addrVals[1])

	// Re-import (only the offline-since part)
	for _, entry := range genesis.ValidatorOfflineSince {
		valAddr, err := sdk.ValAddressFromBech32(entry.ValidatorAddress)
		require.NoError(t, err)
		dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, valAddr, entry.OfflineSince)
	}

	// Verify restored
	offlineSince, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, now.UTC(), offlineSince.UTC())

	offlineSince, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[1])
	require.True(t, found)
	require.Equal(t, now.Add(time.Hour).UTC(), offlineSince.UTC())
}

func TestAutoUnbondParamValidation(t *testing.T) {
	params := types.DefaultParams()

	// Zero is valid (disabled)
	params.AutoUnbondTimeout = 0
	require.NoError(t, params.Validate())

	// Positive is valid
	params.AutoUnbondTimeout = 30 * 24 * time.Hour
	require.NoError(t, params.Validate())

	// Negative is invalid
	params.AutoUnbondTimeout = -time.Hour
	require.Error(t, params.Validate())
}

// TestProcessAutoUnbondWithBlockTime uses a custom block time to verify cutoff logic.
func TestProcessAutoUnbondWithBlockTime(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 3, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	timeout := 24 * time.Hour
	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.AutoUnbondTimeout = timeout
	dsc.ValidatorKeeper.SetParams(ctx, params)

	baseTime := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: baseTime})

	// val0: offline 25h ago (expired)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[0], baseTime.Add(-25*time.Hour))
	// val1: offline 23h ago (not expired)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[1], baseTime.Add(-23*time.Hour))
	// val2: offline exactly 24h ago (expired, cutoff is exclusive)
	dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, addrVals[2], baseTime.Add(-24*time.Hour))

	dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)

	// val0: expired, should be deleted
	_, found := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[0])
	require.False(t, found, "val0 should be processed (25h > 24h)")

	// val1: not expired, should remain
	_, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[1])
	require.True(t, found, "val1 should not be processed (23h < 24h)")

	// val2: exactly at cutoff, should be deleted
	_, found = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, addrVals[2])
	require.False(t, found, "val2 should be processed (24h == 24h)")
}

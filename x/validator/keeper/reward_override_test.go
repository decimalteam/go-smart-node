package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestRewardPerBlockOverride_StoreRoundTrip(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	k := dsc.ValidatorKeeper

	// Absent by default.
	_, ok := k.GetRewardPerBlockOverride(ctx)
	require.False(t, ok)

	// Set a non-zero value.
	want := sdkmath.NewInt(1234567890)
	k.SetRewardPerBlockOverride(ctx, want)
	got, ok := k.GetRewardPerBlockOverride(ctx)
	require.True(t, ok)
	require.True(t, want.Equal(got), "want %s got %s", want, got)

	// Zero is a valid, explicit override (distinct from cleared).
	k.SetRewardPerBlockOverride(ctx, sdkmath.ZeroInt())
	got, ok = k.GetRewardPerBlockOverride(ctx)
	require.True(t, ok)
	require.True(t, got.IsZero())

	// Clearing removes it; falls back to "absent".
	k.ClearRewardPerBlockOverride(ctx)
	_, ok = k.GetRewardPerBlockOverride(ctx)
	require.False(t, ok)
}

func TestGetBlockReward_ReduceOnlyClamp(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	k := dsc.ValidatorKeeper

	const height = uint64(1)
	scheduled := types.GetRewardForBlock(height)
	require.True(t, scheduled.IsPositive(), "precondition: schedule should be positive at height 1")

	// No override -> schedule.
	require.True(t, k.GetBlockReward(ctx, height).Equal(scheduled))

	// Override below schedule -> override wins (emission lowered).
	below := scheduled.Sub(sdkmath.OneInt())
	k.SetRewardPerBlockOverride(ctx, below)
	require.True(t, k.GetBlockReward(ctx, height).Equal(below))

	// Override equal to schedule -> schedule (not strictly below).
	k.SetRewardPerBlockOverride(ctx, scheduled)
	require.True(t, k.GetBlockReward(ctx, height).Equal(scheduled))

	// Override above schedule -> clamped to schedule (cannot inflate).
	k.SetRewardPerBlockOverride(ctx, scheduled.Add(sdkmath.NewInt(1_000_000)))
	require.True(t, k.GetBlockReward(ctx, height).Equal(scheduled))

	// Explicit zero override (below positive schedule) -> zero emission.
	k.SetRewardPerBlockOverride(ctx, sdkmath.ZeroInt())
	require.True(t, k.GetBlockReward(ctx, height).IsZero())

	// Clearing restores the schedule.
	k.ClearRewardPerBlockOverride(ctx)
	require.True(t, k.GetBlockReward(ctx, height).Equal(scheduled))
}

func TestGetBlockReward_AfterEmissionEnd(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	k := dsc.ValidatorKeeper

	// Past the emission schedule end, scheduled reward is zero. The clamp must
	// keep it at zero even with a large override (never mint past the cap).
	const endedHeight = uint64(46_656_000)
	require.True(t, types.GetRewardForBlock(endedHeight).IsZero(), "precondition: schedule ended")

	k.SetRewardPerBlockOverride(ctx, sdkmath.NewInt(1_000_000_000))
	require.True(t, k.GetBlockReward(ctx, endedHeight).IsZero())
}

// TestRewardPerBlockUpdated_UnpackLog validates the ABI fragment + generated
// struct that the EVM hook relies on to decode the contract event.
func TestRewardPerBlockUpdated_UnpackLog(t *testing.T) {
	abiV, err := validator.ValidatorMetaData.GetAbi()
	require.NoError(t, err)

	ev, ok := abiV.Events["RewardPerBlockUpdated"]
	require.True(t, ok, "RewardPerBlockUpdated must be present in the binding ABI")

	// Both args are non-indexed -> all live in log.Data; the only topic is the event id.
	data, err := ev.Inputs.NonIndexed().Pack(big.NewInt(4242), true)
	require.NoError(t, err)
	log := &ethtypes.Log{Topics: []common.Hash{ev.ID}, Data: data}

	var out validator.ValidatorRewardPerBlockUpdated
	require.NoError(t, contracts.UnpackLog(abiV, &out, "RewardPerBlockUpdated", log))
	require.Equal(t, int64(4242), out.RewardPerBlock.Int64())
	require.True(t, out.Enabled)
}

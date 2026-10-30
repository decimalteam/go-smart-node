package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func hold(amount, start, end int64) *types.StakeHold {
	return &types.StakeHold{Amount: sdkmath.NewInt(amount), HoldStartTime: start, HoldEndTime: end}
}

func sumAt(holds []*types.StakeHold, end int64) int64 {
	var s int64
	for _, h := range holds {
		if h.HoldEndTime == end {
			s += h.Amount.Int64()
		}
	}
	return s
}

func TestRebuildHolds(t *testing.T) {
	const now = 1000
	t.Run("unchanged", func(t *testing.T) {
		in := []*types.StakeHold{hold(50, 10, 2000)}
		out, changed := RebuildHolds(in, map[int64]sdkmath.Int{2000: sdkmath.NewInt(50)}, now)
		require.False(t, changed)
		require.Equal(t, in, out)
	})
	t.Run("grown bucket gets a new segment starting a year before its end, old start kept", func(t *testing.T) {
		out, changed := RebuildHolds([]*types.StakeHold{hold(50, 10, 2000)}, map[int64]sdkmath.Int{2000: sdkmath.NewInt(80)}, now)
		require.True(t, changed)
		require.Equal(t, int64(80), sumAt(out, 2000))
		require.Len(t, out, 2)
		require.Equal(t, int64(10), out[0].HoldStartTime)
		require.Equal(t, int64(2000-holdYear), out[1].HoldStartTime)
		require.Equal(t, int64(30), out[1].Amount.Int64())
	})
	t.Run("new bucket", func(t *testing.T) {
		out, changed := RebuildHolds(nil, map[int64]sdkmath.Int{3000: sdkmath.NewInt(5)}, now)
		require.True(t, changed)
		require.Equal(t, []*types.StakeHold{hold(5, 3000-holdYear, 3000)}, out)
	})
	t.Run("restored start: one year before end, never in the future", func(t *testing.T) {
		const n = int64(1_790_000_000)
		require.Equal(t, n+100-holdYear, restoredHoldStart(n+100, n), "ends soon: counts as a full-year hold")
		require.Equal(t, n, restoredHoldStart(n+2*holdYear, n), "ends in 2 years: starts now")
		require.Equal(t, n, restoredHoldStart(n+holdYear, n), "exactly a year away: starts now")
	})
	t.Run("shrunk bucket", func(t *testing.T) {
		out, changed := RebuildHolds([]*types.StakeHold{hold(50, 10, 2000), hold(20, 20, 2000)}, map[int64]sdkmath.Int{2000: sdkmath.NewInt(30)}, now)
		require.True(t, changed)
		require.Equal(t, int64(30), sumAt(out, 2000))
	})
	t.Run("bucket absent from contract is removed", func(t *testing.T) {
		out, changed := RebuildHolds([]*types.StakeHold{hold(50, 10, 2000), hold(9, 10, 2500)}, map[int64]sdkmath.Int{2000: sdkmath.NewInt(50)}, now)
		require.True(t, changed)
		require.Equal(t, int64(0), sumAt(out, 2500))
		require.Equal(t, int64(50), sumAt(out, 2000))
	})
	t.Run("matured buckets are left to DeleteHoldMature", func(t *testing.T) {
		in := []*types.StakeHold{hold(50, 10, 900)}
		out, changed := RebuildHolds(in, map[int64]sdkmath.Int{}, now)
		require.False(t, changed)
		require.Equal(t, in, out)
		_, changed = RebuildHolds(nil, map[int64]sdkmath.Int{500: sdkmath.NewInt(7)}, now)
		require.False(t, changed)
	})
	t.Run("HoldsMatch does not mutate", func(t *testing.T) {
		in := []*types.StakeHold{hold(50, 10, 2000)}
		require.False(t, HoldsMatch(in, map[int64]sdkmath.Int{2000: sdkmath.NewInt(10)}, now))
		require.Equal(t, int64(50), in[0].Amount.Int64())
	})
}

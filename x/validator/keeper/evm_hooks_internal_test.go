package keeper

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	validatorType "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func mkHold(amount, start, end int64) *validatorType.StakeHold {
	return &validatorType.StakeHold{
		Amount:        math.NewInt(amount),
		HoldStartTime: start,
		HoldEndTime:   end,
	}
}

func sumHolds(holds []*validatorType.StakeHold) int64 {
	total := int64(0)
	for _, h := range holds {
		total += h.Amount.Int64()
	}
	return total
}

// TestApplyTransferredHold locks in the redelegation hold-accounting fixes:
//   - the moved amount is subtracted distributively across same-end-time holds and clamped
//     at zero, so hold amounts never go negative (regression for the over-subtraction bug);
//   - the moved amount is split into per-source-start sub-holds, each carrying the REAL start
//     of the source hold it came from, so a young same-end hold cannot inherit an older
//     hold's >=1yr reward window (regression for the start-piggyback bug);
//   - on desync the destination is credited only what was sourced (sum(movedHolds) =
//     moved - leftover), never the full moved amount.
func TestApplyTransferredHold(t *testing.T) {
	const (
		end       = int64(2000)
		otherEnd  = int64(9000)
		origStart = int64(100)
		nowStart  = int64(1500) // placeholder "reset to now" start on the moved hold
	)

	t.Run("single hold: partial transfer keeps source start, reduces amount", func(t *testing.T) {
		remain := []*validatorType.StakeHold{mkHold(100, origStart, end)}
		moved := mkHold(30, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		require.Len(t, kept, 1)
		require.Equal(t, int64(70), kept[0].Amount.Int64())
		require.Len(t, movedHolds, 1)
		require.Equal(t, int64(30), movedHolds[0].Amount.Int64())
		require.Equal(t, origStart, movedHolds[0].HoldStartTime, "moved hold carries source start")
		require.Equal(t, int64(30), moved.Amount.Int64(), "moved.Amount must be untouched")
	})

	t.Run("single hold: full transfer drops the consumed hold", func(t *testing.T) {
		remain := []*validatorType.StakeHold{mkHold(100, origStart, end)}
		moved := mkHold(100, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		require.Len(t, kept, 0)
		require.Len(t, movedHolds, 1)
		require.Equal(t, int64(100), movedHolds[0].Amount.Int64())
		require.Equal(t, origStart, movedHolds[0].HoldStartTime)
	})

	t.Run("duplicate same-end holds: distribute, never over-subtract", func(t *testing.T) {
		// Old over-subtraction bug: subtracted 500 from BOTH holds -> both 0 (1000 lost).
		remain := []*validatorType.StakeHold{mkHold(500, origStart, end), mkHold(500, 200, end)}
		moved := mkHold(500, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		for _, h := range kept {
			require.False(t, h.Amount.IsNegative(), "hold amount must never be negative")
		}
		require.Equal(t, int64(500), sumHolds(kept), "remaining held = 1000-500")
		require.Equal(t, int64(500), sumHolds(movedHolds), "moved held credit = 500")
	})

	t.Run("duplicate same-end holds: large transfer spans entries with per-start segments", func(t *testing.T) {
		// Old over-subtraction bug: subtracted 800 from BOTH -> -300 each (negatives).
		remain := []*validatorType.StakeHold{mkHold(500, origStart, end), mkHold(500, 200, end)}
		moved := mkHold(800, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		require.Equal(t, int64(200), sumHolds(kept), "remaining held = 1000-800")
		// 500 drawn from the first hold (start origStart), 300 from the second (start 200).
		require.Len(t, movedHolds, 2)
		require.Equal(t, int64(500), movedHolds[0].Amount.Int64())
		require.Equal(t, origStart, movedHolds[0].HoldStartTime)
		require.Equal(t, int64(300), movedHolds[1].Amount.Int64())
		require.Equal(t, int64(200), movedHolds[1].HoldStartTime)
		require.Equal(t, int64(800), sumHolds(movedHolds))
	})

	t.Run("anti-piggyback: young same-end coins keep their own young start", func(t *testing.T) {
		// Old start-piggyback bug: the whole moved amount inherited the FIRST (oldest) start,
		// so 999_999 young coins claimed the >=1yr window of a 1-token dust hold.
		const oldStart = int64(10)   // qualifies: end-oldStart large
		const youngStart = int64(1900) // does not qualify: end-youngStart small
		remain := []*validatorType.StakeHold{mkHold(1, oldStart, end), mkHold(1_000_000, youngStart, end)}
		moved := mkHold(1_000_000, nowStart, end)

		_, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		require.Len(t, movedHolds, 2)
		// Only the 1-token dust segment carries the old start; the 999_999 keep the young start.
		require.Equal(t, int64(1), movedHolds[0].Amount.Int64())
		require.Equal(t, oldStart, movedHolds[0].HoldStartTime)
		require.Equal(t, int64(999_999), movedHolds[1].Amount.Int64())
		require.Equal(t, youngStart, movedHolds[1].HoldStartTime, "young coins must NOT inherit the old start")
	})

	t.Run("holds with other end-times are untouched", func(t *testing.T) {
		remain := []*validatorType.StakeHold{mkHold(100, origStart, end), mkHold(40, 50, otherEnd)}
		moved := mkHold(100, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		require.Len(t, kept, 1)
		require.Equal(t, otherEnd, kept[0].HoldEndTime)
		require.Equal(t, int64(40), kept[0].Amount.Int64())
		require.Equal(t, int64(100), sumHolds(movedHolds))
	})

	t.Run("desync: moved exceeds available holds -> destination credited only sourced amount", func(t *testing.T) {
		remain := []*validatorType.StakeHold{mkHold(100, origStart, end)}
		moved := mkHold(150, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.Equal(t, int64(50), leftover.Int64(), "leftover = moved - sourced")
		require.Len(t, kept, 0) // the 100 hold fully consumed and dropped
		require.Equal(t, int64(100), sumHolds(movedHolds), "destination credited only the sourced 100, not 150")
	})

	t.Run("pre-existing non-positive source holds are dropped, not consumed into moved", func(t *testing.T) {
		remain := []*validatorType.StakeHold{mkHold(0, origStart, end), mkHold(-5, 200, end), mkHold(60, 300, end)}
		moved := mkHold(40, nowStart, end)

		kept, movedHolds, leftover := applyTransferredHold(remain, moved)

		require.True(t, leftover.IsZero())
		// Only the positive 60 hold is touched: 40 drawn, 20 remains; zero/negative dropped.
		require.Len(t, kept, 1)
		require.Equal(t, int64(20), kept[0].Amount.Int64())
		require.Equal(t, int64(300), kept[0].HoldStartTime)
		require.Len(t, movedHolds, 1)
		require.Equal(t, int64(40), movedHolds[0].Amount.Int64())
		require.Equal(t, int64(300), movedHolds[0].HoldStartTime)
	})
}

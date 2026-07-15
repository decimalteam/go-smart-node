package keeper

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResolveHoldStartTime(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	contractStart := int64(1_699_000_000)

	// New bucket + non-zero event start -> use the event start.
	require.Equal(t, contractStart, resolveHoldStartTime(now, big.NewInt(contractStart), true))
	// New bucket + zero event start (pre-backfill) -> fall back to block time.
	require.Equal(t, now.Unix(), resolveHoldStartTime(now, big.NewInt(0), true))
	// Top-up (not a new bucket) -> always block time, ignore event start.
	require.Equal(t, now.Unix(), resolveHoldStartTime(now, big.NewInt(contractStart), false))
}

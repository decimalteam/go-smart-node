package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// ReliabilityFixActivationHeights is a release-time, per-chain activation table.
// Keep empty until a coordinated upgrade height has been agreed. It must never
// be derived from local configuration, wall-clock time, or environment variables.
// Treat it as immutable after application startup.
var ReliabilityFixActivationHeights = map[string]int64{}

func (k Keeper) reliabilityFixesEnabled(ctx sdk.Context) bool {
	height := ReliabilityFixActivationHeights[ctx.ChainID()]
	return height > 0 && ctx.BlockHeight() >= height
}

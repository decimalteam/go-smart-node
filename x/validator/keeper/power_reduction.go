package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/decimalteam/ethermint/types"
)

// TokensToConsensusPower converts input tokens to potential consensus-engine power
func TokensToConsensusPower(tokens sdkmath.Int) int64 {
	return sdk.TokensToConsensusPower(tokens, ethtypes.PowerReduction)
}

// TokensFromConsensusPower - convert input power to tokens
func TokensFromConsensusPower(power int64) sdkmath.Int {
	return sdk.TokensFromConsensusPower(power, ethtypes.PowerReduction)
}

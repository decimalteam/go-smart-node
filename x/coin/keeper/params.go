package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// GetParams returns the total set of the module parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	// Effective optimizations to reduce retrieving param values
	k.baseDenom = params.BaseSymbol

	k.ps.SetParamSet(ctx, &params)
}

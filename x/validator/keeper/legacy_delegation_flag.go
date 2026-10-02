package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EnableLegacyDelegationEvents turns on decoding of the legacy (7-field Stake)
// delegation contract events in PostTxProcessing. It is set once by the mainnet
// delegation-desync repair upgrade so blocks before that upgrade replay exactly
// as they were executed.
func (k Keeper) EnableLegacyDelegationEvents(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Set(types.GetLegacyDelegationEventsKey(), []byte{1})
}

// LegacyDelegationEventsEnabled reports whether legacy delegation events are decoded.
func (k Keeper) LegacyDelegationEventsEnabled(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetLegacyDelegationEventsKey())
}

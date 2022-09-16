package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// RegisterInvariants registers all the module's invariants.
func RegisterInvariants(registry sdk.InvariantRegistry, k Keeper) {
	registry.RegisterRoute(types.ModuleName, "supply", SupplyInvariant(k))
}

// AllInvariants runs all invariants of the module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return SupplyInvariant(k)(ctx)
	}
}

// SupplyInvariant checks that the total amount of NFT tokens in the collections matches the total amount NFT tokens in the KVStore.
func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		collections := k.GetCollections(ctx)
		msg, invariant := types.SupplyInvariantCheck(collections)
		return sdk.FormatInvariant(types.ModuleName, "supply", msg), invariant
	}
}

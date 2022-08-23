package keeper

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(
		types.ModuleName, "supply",
		SupplyInvariant(k),
	)
}

// AllInvariants runs all invariants of the nfts module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return SupplyInvariant(k)(ctx)
	}
}

// SupplyInvariant checks that the total amount of nfts on collections matches the total amount nfts in store
func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		collections := k.GetCollections(ctx)
		nfts, err := k.GetNFTs(ctx)
		if err != nil {
			panic(err)
		}

		msg, invariant := types.SupplyInvariantCheck(collections, nfts)

		return sdk.FormatInvariant(types.ModuleName, "supply", msg), invariant
	}
}

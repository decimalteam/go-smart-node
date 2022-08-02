package keeper

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(
		types.ModuleName, "supply",
		BaseDenomPriceInvariant(k),
	)
}

// AllInvariants runs all invariants of the nfts module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return BaseDenomPriceInvariant(k)(ctx)
	}
}

// BaseDenomPriceInvariant checks that the total amount of nfts on collections matches the total amount nfts in store
func BaseDenomPriceInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		price, found := k.GetBaseDenomPrice(ctx)
		if found && price <= 0 {
			msg := fmt.Sprintf("wrong price: %f", price)

			return sdk.FormatInvariant(types.ModuleName, "price", msg), true
		}

		return "", false
	}
}

package nft_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	reserve sdk.Int = types.NewMinReserve2
)

// CheckInvariants checks the invariants
//func CheckInvariants(k keeper.Keeper, ctx sdk.Context) bool {
//	collectionsSupply := make(map[string]int)
//	ownersCollectionsSupply := make(map[string]int)
//
//	k.IterateCollections(ctx, func(collection types.Collection) bool {
//		collectionsSupply[collection.Denom] = collection.Supply()
//		return false
//	})
//
//	owners := k.GetOwners(ctx)
//	for _, owner := range owners {
//		for _, idCollection := range owner.IDCollections {
//			ownersCollectionsSupply[idCollection.Denom] += idCollection.Supply()
//		}
//	}
//
//	for denom, supply := range collectionsSupply {
//		if supply != ownersCollectionsSupply[denom] {
//			fmt.Printf("denom is %s, supply is %d, ownerSupply is %d", denom, supply, ownersCollectionsSupply[denom])
//			return false
//		}
//	}
//	return true
//}

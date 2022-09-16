package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	// Initialize params
	k.SetParams(ctx, gs.Params)

	// TODO: Initialize NFT collections
	// TODO: Initialize NFT tokens
	// TODO: Initialize NFT sub-tokens

	// nftDenoms := make(map[string]string)
	// for _, collection := range gs.Collections {
	// 	for _, id := range collection.NFTs {
	// 		nftDenoms[id] = collection.Denom
	// 	}

	// 	k.SetCollection(ctx, collection.Denom, collection)
	// }

	// for nftID, subTokens := range gs.SubTokens {
	// 	for _, subToken := range subTokens.SubTokens {
	// 		k.setSubToken(ctx, nftID, subToken)
	// 	}
	// }

	// for _, nft := range gs.NFTs {
	// 	denom := nftDenoms[nft.GetID()]

	// 	for _, owner := range nft.GetOwners() {
	// 		err := k.MintNFTAndCollection(ctx, denom, nft.GetID(), nft.GetReserve(), nft.GetCreator(), owner.GetAddress(), nft.GetTokenURI(), nft.GetAllowMint(), owner.SubTokenIDs)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:      k.GetParams(ctx),
		Collections: k.GetCollections(ctx),
	}
}

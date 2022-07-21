package nft

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets nft information for genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	nftDenoms := make(map[string]string)
	for _, collection := range data.Collections {
		for _, id := range collection.NFTs {
			nftDenoms[id] = collection.Denom
		}

		k.SetCollection(ctx, collection.Denom, collection)
	}

	for nftID, subTokens := range data.SubTokens {
		for _, subToken := range subTokens.SubTokens {
			k.SetSubToken(ctx, nftID, subToken)
		}
	}

	for _, nft := range data.Nfts {
		denom := nftDenoms[nft.GetID()]

		for _, owner := range nft.GetOwners() {
			err := k.MintNFTAndCollection(ctx, denom, nft.GetID(), nft.GetReserve(), nft.GetCreator(), owner.GetAddress(), nft.GetTokenURI(), nft.GetAllowMint(), owner.SubTokenIDs)
			if err != nil {
				panic(err)
			}
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	collections := k.GetCollections(ctx)

	nfts := k.GetNFTs(ctx)

	subTokens := make(map[string]types.SubTokens)
	for _, nft := range nfts {
		nftSubTokens := k.GetSubTokens(ctx, nft.GetID())
		subTokens[nft.GetID()] = types.SubTokens{
			SubTokens: make([]types.SubToken, len(nftSubTokens)),
		}
	}

	return types.NewGenesisState(collections, nfts, subTokens)
}

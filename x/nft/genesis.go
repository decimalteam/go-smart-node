package nft

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets nft information for genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetOwners(ctx, data.Owners)

	for _, c := range data.Collections {
		sortedCollection := types.NewCollection(c.Denom, c.NFTs.Sort())
		k.SetCollection(ctx, c.Denom, sortedCollection)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(k.GetOwners(ctx), k.GetCollections(ctx))
}

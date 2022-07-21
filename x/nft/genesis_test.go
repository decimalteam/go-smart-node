package nft_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 2)
	genesisState := types.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Owners))
	require.Equal(t, 0, len(genesisState.Collections))

	ids := []string{ID1, ID2, ID3}
	idCollection := types.NewIDCollection(Denom1, ids)
	idCollection2 := types.NewIDCollection(Denom2, ids)
	owner := types.NewOwner(addrs[0].String(), idCollection)

	owner2 := types.NewOwner(addrs[1].String(), idCollection2)

	owners := []types.Owner{owner, owner2}

	reserve := sdk.NewInt(100)
	subTokenIds := []int64{}

	nft1 := types.NewBaseNFT(ID1, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, subTokenIds, true)
	nft2 := types.NewBaseNFT(ID2, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, subTokenIds, true)
	nft3 := types.NewBaseNFT(ID3, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, subTokenIds, true)
	nfts := types.NewNFTs(nft1, nft2, nft3)
	collection := types.NewCollection(Denom1, nfts)

	nftx := types.NewBaseNFT(ID1, addrs[1].String(), addrs[1].String(), TokenURI1, reserve, subTokenIds, true)
	nft2x := types.NewBaseNFT(ID2, addrs[1].String(), addrs[1].String(), TokenURI1, reserve, subTokenIds, true)
	nft3x := types.NewBaseNFT(ID3, addrs[1].String(), addrs[1].String(), TokenURI1, reserve, subTokenIds, true)
	nftsx := types.NewNFTs(nftx, nft2x, nft3x)
	collection2 := types.NewCollection(Denom2, nftsx)

	collections := types.NewCollections(collection, collection2)

	genesisState = types.NewGenesisState(owners, collections)

	nft.InitGenesis(ctx, dsc.NftKeeper, *genesisState)

	returnedOwners := dsc.NftKeeper.GetOwners(ctx)
	require.Equal(t, 2, len(owners))
	require.Equal(t, returnedOwners[0].String(), owners[0].String())
	require.Equal(t, returnedOwners[1].String(), owners[1].String())

	returnedCollections := dsc.NftKeeper.GetCollections(ctx)
	require.Equal(t, 2, len(returnedCollections))
	require.Equal(t, returnedCollections[0].String(), collections[0].String())
	require.Equal(t, returnedCollections[1].String(), collections[1].String())

	exportedGenesisState := nft.ExportGenesis(ctx, dsc.NftKeeper)
	require.Equal(t, len(genesisState.Owners), len(exportedGenesisState.Owners))
	require.Equal(t, genesisState.Owners[0].String(), exportedGenesisState.Owners[0].String())
	require.Equal(t, genesisState.Owners[1].String(), exportedGenesisState.Owners[1].String())

	require.Equal(t, len(genesisState.Collections), len(exportedGenesisState.Collections))
	require.Equal(t, genesisState.Collections[0].String(), exportedGenesisState.Collections[0].String())
	require.Equal(t, genesisState.Collections[1].String(), exportedGenesisState.Collections[1].String())
}

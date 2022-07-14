package keeper

//import (
//	"strconv"
//	"strings"
//	"testing"
//
//	"bitbucket.org/decimalteam/go-node/x/nft/exported"
//	nftTypes "bitbucket.org/decimalteam/go-node/x/nft/internal/types"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//
//	"github.com/stretchr/testify/require"
//
//	abci "github.com/tendermint/tendermint/abci/types"
//)
//
//func TestNewQuerier(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	querier := NewQuerier(NFTKeeper)
//	query1 := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//	_, err := querier(ctx, []string{"foo", "bar"}, query1)
//	require.Error(t, err)
//}
//
//func TestQuerySupply(t *testing.T) {
//	ctx, cdc, NFTKeeper := createTestApp(t, false)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := nftTypes.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//
//	require.NoError(t, err)
//
//	querier := NewQuerier(NFTKeeper)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//
//	query.Path = "/custom/nft/supply"
//	query.Data = []byte("?")
//
//	res, err := querier(ctx, []string{"supply"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	queryCollectionParams := nftTypes.NewQueryCollectionParams(Denom2)
//	bz, errRes := cdc.MarshalJSON(queryCollectionParams)
//	require.Nil(t, errRes)
//	query.Data = bz
//	res, err = querier(ctx, []string{"supply"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	queryCollectionParams = nftTypes.NewQueryCollectionParams(Denom1)
//	bz, errRes = cdc.MarshalJSON(queryCollectionParams)
//	require.Nil(t, errRes)
//	query.Data = bz
//
//	res, err = querier(ctx, []string{"supply"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	supplyResp := strings.Trim(string(res), "\"")
//	supply, _ := strconv.Atoi(supplyResp)
//	require.Equal(t, 1, supply)
//}
//
//func TestQueryCollection(t *testing.T) {
//	ctx, cdc, NFTKeeper := createTestApp(t, false)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := nftTypes.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	querier := NewQuerier(NFTKeeper)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//
//	query.Path = "/custom/nft/collection"
//
//	query.Data = []byte("?")
//	res, err := querier(ctx, []string{"collection"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	queryCollectionParams := nftTypes.NewQueryCollectionParams(Denom2)
//	bz, errRes := cdc.MarshalJSON(queryCollectionParams)
//	require.Nil(t, errRes)
//
//	query.Data = bz
//	res, err = querier(ctx, []string{"collection"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	queryCollectionParams = nftTypes.NewQueryCollectionParams(Denom1)
//	bz, errRes = cdc.MarshalJSON(queryCollectionParams)
//	require.Nil(t, errRes)
//
//	query.Data = bz
//	res, err = querier(ctx, []string{"collection"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	var collections nftTypes.Collections
//	nftTypes.ModuleCdc.MustUnmarshalJSON(res, &collections)
//	require.Len(t, collections, 1)
//	require.Len(t, collections[0].NFTs, 1)
//}
//
//func TestQueryOwner(t *testing.T) {
//	ctx, cdc, NFTKeeper := createTestApp(t, false)
//
//	// MintNFT shouldn't fail when collection does not exist
//
//	nft := nftTypes.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	_, err = NFTKeeper.MintNFT(ctx, Denom2, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	querier := NewQuerier(NFTKeeper)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//	query.Path = "/custom/nft/ownerByDenom"
//
//	query.Data = []byte("?")
//	res, err := querier(ctx, []string{"ownerByDenom"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	// query the balance using the first denom
//	params := nftTypes.NewQueryBalanceParams(Addrs[0], Denom1)
//	bz, err2 := cdc.MarshalJSON(params)
//	require.Nil(t, err2)
//	query.Data = bz
//
//	res, err = querier(ctx, []string{"ownerByDenom"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	var out nftTypes.Owner
//	cdc.MustUnmarshalJSON(res, &out)
//
//	// build the owner using only the first denom
//	idCollection1 := nftTypes.NewIDCollection(Denom1, []string{ID1})
//	owner := nftTypes.NewOwner(Addrs[0], idCollection1)
//
//	require.Equal(t, out.String(), owner.String())
//
//	// query the balance using no denom so that all denoms will be returns
//	params = nftTypes.NewQueryBalanceParams(Addrs[0], "")
//	bz, err2 = cdc.MarshalJSON(params)
//	require.Nil(t, err2)
//
//	query.Path = "/custom/nft/owner"
//	query.Data = []byte("?")
//	_, err = querier(ctx, []string{"owner"}, query)
//	require.Error(t, err)
//
//	query.Data = bz
//	res, err = querier(ctx, []string{"owner"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	cdc.MustUnmarshalJSON(res, &out)
//
//	// build the owner using both denoms
//	idCollection2 := nftTypes.NewIDCollection(Denom2, []string{ID1})
//	owner = nftTypes.NewOwner(Addrs[0], idCollection2, idCollection1)
//
//	out.IDCollections = out.IDCollections.Sort()
//	owner.IDCollections = owner.IDCollections.Sort()
//
//	require.Equal(t, out.String(), owner.String())
//}
//
//func TestQueryNFT(t *testing.T) {
//	ctx, cdc, NFTKeeper := createTestApp(t, false)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := nftTypes.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{1},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	querier := NewQuerier(NFTKeeper)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//	query.Path = "/custom/nft/nft"
//	var res []byte
//
//	query.Data = []byte("?")
//	res, err = querier(ctx, []string{"nft"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	params := nftTypes.NewQueryNFTParams(Denom2, ID2)
//	bz, err2 := cdc.MarshalJSON(params)
//	require.Nil(t, err2)
//
//	query.Data = bz
//	res, err = querier(ctx, []string{"nft"}, query)
//	require.Error(t, err)
//	require.Nil(t, res)
//
//	params = nftTypes.NewQueryNFTParams(Denom1, ID1)
//	bz, err2 = cdc.MarshalJSON(params)
//	require.Nil(t, err2)
//
//	query.Data = bz
//	res, err = querier(ctx, []string{"nft"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	var out exported.NFT
//	cdc.MustUnmarshalJSON(res, &out)
//
//	require.Equal(t, out.String(), nft.String())
//}
//
//func TestQueryDenoms(t *testing.T) {
//	ctx, cdc, NFTKeeper := createTestApp(t, false)
//
//	// MintNFT shouldn't fail when collection does not exist
//	nft := nftTypes.NewBaseNFT(
//		ID1,
//		Addrs[0],
//		Addrs[0],
//		TokenURI1,
//		sdk.NewInt(100),
//		[]int64{},
//		true,
//	)
//	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	_, err = NFTKeeper.MintNFT(ctx, Denom2, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
//	require.NoError(t, err)
//
//	querier := NewQuerier(NFTKeeper)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//	var res []byte
//
//	query.Path = "/custom/nft/denoms"
//
//	res, err = querier(ctx, []string{"denoms"}, query)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	denoms := []string{Denom2, Denom1}
//
//	var out []string
//	cdc.MustUnmarshalJSON(res, &out)
//
//	for key, denomInQuestion := range out {
//		require.Equal(t, denomInQuestion, denoms[key])
//	}
//}

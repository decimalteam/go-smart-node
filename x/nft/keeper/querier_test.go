package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"strconv"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())
	query1 := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(ctx, []string{"foo", "bar"}, query1)
	require.Error(t, err)
}

func TestQuerySupply(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())

	query := abci.RequestQuery{
		Path: "/custom/nft/supply",
		Data: []byte{},
	}

	query.Data = []byte("?")

	res, err := querier(ctx, []string{"supply"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	queryCollectionParams := types.NewQueryCollectionParams(Denom2)
	bz, errRes := dsc.LegacyAmino().MarshalJSON(queryCollectionParams)
	require.Nil(t, errRes)
	query.Data = bz
	res, err = querier(ctx, []string{"supply"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	queryCollectionParams = types.NewQueryCollectionParams(Denom1)
	bz, errRes = dsc.LegacyAmino().MarshalJSON(queryCollectionParams)
	require.Nil(t, errRes)
	query.Data = bz

	res, err = querier(ctx, []string{"supply"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	supplyResp := strings.Trim(string(res), "\"")
	supply, _ := strconv.Atoi(supplyResp)
	require.Equal(t, 1, supply)
}

func TestQueryCollection(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	// MintNFT shouldn't fail when collection does not exist
	addrs := getAddrs(dsc, ctx, 1)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = "/custom/nft/collection"

	query.Data = []byte("?")
	res, err := querier(ctx, []string{"collection"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	queryCollectionParams := types.NewQueryCollectionParams(Denom2)
	bz, errRes := dsc.LegacyAmino().MarshalJSON(queryCollectionParams)
	require.Nil(t, errRes)

	query.Data = bz
	res, err = querier(ctx, []string{"collection"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	queryCollectionParams = types.NewQueryCollectionParams(Denom1)
	bz, errRes = dsc.LegacyAmino().MarshalJSON(queryCollectionParams)
	require.Nil(t, errRes)

	query.Data = bz
	res, err = querier(ctx, []string{"collection"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	var collections types.Collections
	dsc.LegacyAmino().MustUnmarshalJSON(res, &collections)
	require.Len(t, collections, 1)
	require.Len(t, collections[0].NFTs, 1)
}

func TestQueryOwner(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	msg = types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom2,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = "/custom/nft/ownerByDenom"

	query.Data = []byte("?")
	res, err := querier(ctx, []string{"ownerByDenom"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	// query the balance using the first denom
	params := types.NewQueryBalanceParams(addrs[0], Denom1)
	bz, err2 := dsc.LegacyAmino().MarshalJSON(params)
	require.Nil(t, err2)
	query.Data = bz

	res, err = querier(ctx, []string{"ownerByDenom"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	var out types.Owner
	dsc.LegacyAmino().MustUnmarshalJSON(res, &out)

	// build the owner using only the first denom
	idCollection1 := types.NewIDCollection(Denom1, []string{ID1})
	owner := types.NewOwner(addrs[0].String(), idCollection1)

	require.Equal(t, out.String(), owner.String())

	// query the balance using no denom so that all denoms will be returns
	params = types.NewQueryBalanceParams(addrs[0], "")
	bz, err2 = dsc.LegacyAmino().MarshalJSON(params)
	require.Nil(t, err2)

	query.Path = "/custom/nft/owner"
	query.Data = []byte("?")
	_, err = querier(ctx, []string{"owner"}, query)
	require.Error(t, err)

	query.Data = bz
	res, err = querier(ctx, []string{"owner"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	dsc.LegacyAmino().MustUnmarshalJSON(res, &out)

	// build the owner using both denoms
	idCollection2 := types.NewIDCollection(Denom2, []string{ID1})
	owner = types.NewOwner(addrs[0].String(), idCollection2, idCollection1)

	out.IDCollections = out.IDCollections.Sort()
	owner.IDCollections = owner.IDCollections.Sort()

	require.Equal(t, out.String(), owner.String())
}

func TestQueryNFT(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		addrs[0].String(),
		addrs[0].String(),
		TokenURI1,
		sdk.NewInt(100),
		[]int64{1},
		true,
	)

	_, err := dsc.NftKeeper.Mint(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), nft.GetCreator(), nft.GetTokenURI(), nft.GetAllowMint())
	require.NoError(t, err)

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = "/custom/nft/nft"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(ctx, []string{"nft"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params := types.NewQueryNFTParams(Denom2, ID2)
	bz, err2 := dsc.LegacyAmino().MarshalJSON(params)
	require.Nil(t, err2)

	query.Data = bz
	res, err = querier(ctx, []string{"nft"}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params = types.NewQueryNFTParams(Denom1, ID1)
	bz, err2 = dsc.LegacyAmino().MarshalJSON(params)
	require.Nil(t, err2)

	query.Data = bz
	res, err = querier(ctx, []string{"nft"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	var out types.BaseNFT
	dsc.LegacyAmino().MustUnmarshalJSON(res, &out)

	require.Equal(t, out.String(), nft.String())
}

func TestQueryDenoms(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)

	// MintNFT shouldn't fail when collection does not exist
	msg := types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom1,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err := dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	msg = types.MsgMintNFT{
		Sender:    addrs[0].String(),
		Recipient: addrs[0].String(),
		ID:        ID1,
		Denom:     Denom2,
		Quantity:  sdk.NewInt(50),
		TokenURI:  TokenURI1,
		Reserve:   types.NewMinReserve2,
		AllowMint: true,
	}

	_, err = dsc.NftKeeper.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	require.NoError(t, err)

	querier := keeper.NewQuerier(dsc.NftKeeper, dsc.LegacyAmino())

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	var res []byte

	query.Path = "/custom/nft/denoms"

	res, err = querier(ctx, []string{"denoms"}, query)
	require.NoError(t, err)
	require.NotNil(t, res)

	denoms := []string{Denom2, Denom1}

	var out []string
	dsc.LegacyAmino().MustUnmarshalJSON(res, &out)

	for key, denomInQuestion := range out {
		require.Equal(t, denomInQuestion, denoms[key])
	}
}

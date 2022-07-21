package nft_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	moduleKey            = "module"
	denom                = "denom"
	nftID                = "nft_id"
	sender               = "sender"
	recipient            = "recipient"
	tokenURI             = "token_uri"
	amount               = "amount"
	subTokenIdStartRange = "sub_token_id_start_range"
)

func TestInvalidMsg(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()
	h := nft.NewHandler(dsc.NftKeeper)
	_, err := h(ctx, testdata.NewTestMsg())

	require.Error(t, err)
}

func TestTransferNFTMsg(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()
	h := nft.NewHandler(dsc.NftKeeper)

	addrs := getAddrs(dsc, ctx, 1)
	// An NFT to be transferred
	reserve := sdk.NewInt(100)
	basenft := types.NewBaseNFT(ID1, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, []int64{}, true)

	// Define MsgTransferNft
	transferNftMsg := types.NewMsgTransferNFT(addrs[0], addrs[1], Denom1, ID1, []int64{})

	// handle should fail trying to transfer NFT that doesn't exist
	res, err := h(ctx, transferNftMsg)
	require.Error(t, err)

	// Create token (collection and owner)
	_, err = dsc.NftKeeper.Mint(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(1), basenft.GetCreator(), addrs[0].String(), basenft.GetTokenURI(), basenft.GetAllowMint())
	require.Nil(t, err)
	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case tokenURI:
				require.Equal(t, value, TokenURI1)
			case sender:
				require.Equal(t, value, addrs[0])
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// nft should have been transferred as a result of the message
	nftAfterwards, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, nftAfterwards.GetOwners().GetOwners()[1].GetAddress(), addrs[1])

	transferNftMsg = types.NewMsgTransferNFT(addrs[1], addrs[2], Denom1, ID1, []int64{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))

	// Create token (collection and owner)
	_, err = dsc.NftKeeper.Mint(ctx,
		Denom2, basenft.GetID(),
		basenft.GetReserve(), sdk.NewInt(100),
		basenft.GetCreator(),
		addrs[1].String(),
		basenft.GetTokenURI(), basenft.GetAllowMint(),
	)
	require.Nil(t, err)
	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))

	transferNftMsg = types.NewMsgTransferNFT(addrs[1], addrs[2], Denom2, ID1, []int64{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))
}

func TestEditNFTMetadataMsg(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()
	h := nft.NewHandler(dsc.NftKeeper)

	addrs := getAddrs(dsc, ctx, 1)
	reserve := sdk.NewInt(101)

	// An NFT to be edited
	basenft := types.NewBaseNFT(ID1, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, []int64{}, true)

	// Create token (collection and address)
	_, err := dsc.NftKeeper.Mint(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(1), basenft.GetCreator(), addrs[0].String(), basenft.GetTokenURI(), basenft.GetAllowMint())

	require.Nil(t, err)

	// Define MsgTransferNft
	failingEditNFTMetadata := types.NewMsgEditNFTMetadata(addrs[0], ID1, Denom2, TokenURI2)

	res, err := h(ctx, failingEditNFTMetadata)
	require.Error(t, err)

	// Define MsgTransferNft
	editNFTMetadata := types.NewMsgEditNFTMetadata(addrs[0], ID1, Denom1, TokenURI2)

	res, err = h(ctx, editNFTMetadata)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, addrs[0].String())
			case tokenURI:
				require.Equal(t, value, TokenURI2)
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, TokenURI2, nftAfterwards.GetTokenURI())
}

func TestMintNFTMsg(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()
	h := nft.NewHandler(dsc.NftKeeper)

	// Define MsgMintNFT
	addrs := getAddrs(dsc, ctx, 1)
	reserve := sdk.NewInt(101)
	mintNFT := types.NewMsgMintNFT(addrs[0], addrs[0], ID1, Denom1, TokenURI1, sdk.NewInt(1), reserve, false)

	// minting a token should succeed
	res, err := h(ctx, mintNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, addrs[0].String())
			case tokenURI:
				require.Equal(t, value, TokenURI1)
			case subTokenIdStartRange:
				require.Equal(t, value, ID1)
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)

	require.NoError(t, err)
	require.Equal(t, TokenURI1, nftAfterwards.GetTokenURI())

	// minting the same token should fail if allowMint=false
	res, err = h(ctx, mintNFT)
	require.Error(t, err)

	require.True(t, CheckInvariants(dsc.NftKeeper, ctx))
}

func TestBurnNFTMsg(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()
	h := nft.NewHandler(dsc.NftKeeper)

	addrs := getAddrs(dsc, ctx, 1)

	reserve := sdk.NewInt(100)
	// An NFT to be burned
	basenft := types.NewBaseNFT(ID1, addrs[0].String(), addrs[0].String(), TokenURI1, reserve, []int64{1, 2, 3}, true)

	// Create token (collection and address)
	_, err := dsc.NftKeeper.Mint(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(3), basenft.GetCreator(), addrs[0].String(), basenft.GetTokenURI(), basenft.GetAllowMint())
	require.Nil(t, err)

	bnft, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.NotEmpty(t, bnft)

	// burning a non-existent NFT should fail
	failBurnNFT := types.NewMsgBurnNFT(addrs[0], ID2, Denom1, []int64{4})
	res, err := h(ctx, failBurnNFT)
	require.Error(t, err)

	// NFT should still exist
	bnft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.NotEmpty(t, bnft)

	// burning the NFt should succeed
	burnNFT := types.NewMsgBurnNFT(addrs[0], ID1, Denom1, []int64{2})

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			if event.Type != sdk.EventTypeMessage || event.Type != types.EventTypeBurnNFT {
				continue
			}
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, addrs[0].String())
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nft, err := dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 3}, nft.GetOwners().GetOwners()[0].GetSubTokenIDs())

	// the NFT should not exist after burn
	bnft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.NotEmpty(t, bnft)

	ownerReturned := dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.Equal(t, 1, ownerReturned.Supply())

	burnNFT = types.NewMsgBurnNFT(addrs[0], ID1, Denom1, []int64{1})

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)

	nft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, []int64{3}, nft.GetOwners().GetOwners()[0].GetSubTokenIDs())

	// the NFT should not exist after burn
	bnft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.NotEmpty(t, bnft)

	ownerReturned = dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.Equal(t, 1, ownerReturned.Supply())

	burnNFT = types.NewMsgBurnNFT(addrs[0], ID1, Denom1, []int64{3})

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)

	nft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, []int64(nil), nft.GetOwners().GetOwners()[0].GetSubTokenIDs())

	// the NFT should not exist after burn
	bnft, err = dsc.NftKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.NotEmpty(t, bnft)

	ownerReturned = dsc.NftKeeper.GetOwner(ctx, addrs[0])
	require.Equal(t, 1, ownerReturned.Supply())
}

func TestUniqueTokenURI(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	reserve := sdk.NewInt(100)

	const tokenURI1 = "tokenURI1"
	const tokenURI2 = "tokenURI2"

	msg := types.NewMsgMintNFT(addrs[0], addrs[0], "token1", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err := dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(addrs[0], addrs[0], "token1", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(addrs[0], addrs[0], "token2", "denom1", tokenURI2, sdk.NewInt(1), reserve, true)
	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(addrs[0], addrs[0], "token3", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, types.ErrNotUniqueTokenURI(), err)
}

func TestUniqueTokenID(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := getAddrs(dsc, ctx, 1)
	reserve := sdk.NewInt(100)

	const tokenURI1 = "tokenURI1"
	const tokenURI2 = "tokenURI2"

	msg := types.NewMsgMintNFT(addrs[0], addrs[0], "token1", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err := dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(addrs[0], addrs[0], "token1", "denom2", tokenURI2, sdk.NewInt(1), reserve, true)
	_, err = dsc.NftKeeper.MintNFT(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, types.ErrNotUniqueTokenID(), err)
}

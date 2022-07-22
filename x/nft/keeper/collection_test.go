package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	firstID         string = "first_id"
	secondID        string = "second_id"
	firstDenom      string = "first_test_denom"
	secondDenom     string = "second_test_denom"
	firstTokenURI   string = "first_token_uri"
	secondTokenURI  string = "second_token_uri"
	firstAllowMint  bool   = true
	secondAllowMint bool   = true
)

var (
	firstReserve  sdk.Int = types.NewMinReserve2
	secondReserve sdk.Int = types.NewMinReserve2.MulRaw(2)
	firstQuantity sdk.Int = sdk.NewInt(10)
)

func TestSetCollections(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	collectionDenomsToStore := []string{firstDenom, secondDenom}
	for _, denom := range collectionDenomsToStore {
		collection := types.NewCollection(denom, []string{})
		dsc.NFTKeeper.SetCollection(ctx, collection.Denom, collection)
	}

	// Check throw GetCollections method
	storedCollections := dsc.NFTKeeper.GetCollections(ctx)
	require.Len(t, storedCollections, len(collectionDenomsToStore))

	for _, denom := range collectionDenomsToStore {
		var stored bool
		for _, storeCollection := range storedCollections {
			if storeCollection.GetDenom() == denom {
				stored = true
				break
			}
		}

		require.True(t, stored, denom)
	}

	// Check throw GetCollection method
	for _, denom := range collectionDenomsToStore {
		storedCollection, found := dsc.NFTKeeper.GetCollection(ctx, denom)
		require.True(t, found, denom)
		require.Equal(t, denom, storedCollection.GetDenom())
	}
}

func TestGetDenoms(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	collectionDenomsToStore := []string{firstDenom, secondDenom}
	for _, denom := range collectionDenomsToStore {
		collection := types.NewCollection(denom, []string{})
		dsc.NFTKeeper.SetCollection(ctx, collection.Denom, collection)
	}

	// Check throw GetDenoms method
	storedDenoms := dsc.NFTKeeper.GetDenoms(ctx)
	require.Equal(t, collectionDenomsToStore, storedDenoms)
}

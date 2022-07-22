package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetOwnerCollection(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrs := app.GetAddrs(dsc, ctx, 1)

	ownerCollection1 := types.NewOwnerCollection(firstDenom, []string{"firstID", "secondID"})
	ownerCollection2 := types.NewOwnerCollection(secondDenom, []string{"thirdID", "fourthID"})

	ownerCollectionsToStore := []types.OwnerCollection{ownerCollection1, ownerCollection2}

	for _, ownerCollection := range ownerCollectionsToStore {
		dsc.NFTKeeper.SetOwnerCollectionByDenom(ctx, addrs[0], ownerCollection.Denom, ownerCollection)
	}

	// Check throw GetOwnerCollections method
	storedOwnerCollections := dsc.NFTKeeper.GetOwnerCollections(ctx, addrs[0])
	require.Len(t, storedOwnerCollections, len(ownerCollectionsToStore))

	for _, ownerCollectionToStore := range ownerCollectionsToStore {
		var stored bool
		for _, storedOwnerCollection := range storedOwnerCollections {
			if ownerCollectionToStore.Denom == storedOwnerCollection.Denom {
				stored = true
				require.Equal(t, ownerCollectionToStore, storedOwnerCollection)
				break
			}
		}

		require.True(t, stored, ownerCollectionToStore.Denom)
	}

	// Check throw GetOwnerCollectionByDenom method
	for _, ownerCollectionToStore := range ownerCollectionsToStore {
		storedOwnerCollection, found := dsc.NFTKeeper.GetOwnerCollectionByDenom(ctx, addrs[0], ownerCollectionToStore.Denom)
		require.True(t, found)
		require.Equal(t, ownerCollectionToStore, storedOwnerCollection)
	}
}

package keeper_test

import (
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetSubTokens(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper()

	subToken1 := types.NewSubToken(
		1,
		firstReserve,
	)
	subToken2 := types.NewSubToken(
		2,
		secondReserve,
	)
	subToken3 := types.NewSubToken(
		3,
		firstReserve,
	)

	subTokensToStore := []types.SubToken{subToken1, subToken2, subToken3}
	for _, subTokenToStore := range subTokensToStore {
		dsc.NFTKeeper.SetSubToken(ctx, firstID, subTokenToStore)
	}

	// Check throw GetSubTokens method
	storedSubTokens := dsc.NFTKeeper.GetSubTokens(ctx, firstID)
	require.Len(t, subTokensToStore, len(storedSubTokens))

	for _, subTokenToStore := range subTokensToStore {
		var stored bool
		for _, storedSubToken := range storedSubTokens {
			if storedSubToken.ID == subTokenToStore.ID {
				stored = true
				require.Equal(t, subTokenToStore, storedSubToken)
				break
			}
		}

		require.True(t, stored, subTokenToStore.ID)
	}

	// Check throw GetSubToken method
	for _, subTokenToStore := range subTokensToStore {
		storedSubToken, found := dsc.NFTKeeper.GetSubToken(ctx, firstID, subTokenToStore.ID)
		require.True(t, found, subTokenToStore.ID)
		require.Equal(t, subTokenToStore.ID, storedSubToken.ID)
	}
}

func TestUpdateSubToken(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper()

	subToken := types.NewSubToken(
		1,
		firstReserve,
	)

	// set first version of sub token
	dsc.NFTKeeper.SetSubToken(ctx, firstID, subToken)
	storedSubToken, found := dsc.NFTKeeper.GetSubToken(ctx, firstID, subToken.ID)
	require.True(t, found, subToken.ID)
	require.Equal(t, subToken, storedSubToken)

	subToken.Reserve = secondReserve
	// set updated version of sub token
	dsc.NFTKeeper.SetSubToken(ctx, firstID, subToken)
	storedSubToken, found = dsc.NFTKeeper.GetSubToken(ctx, firstID, subToken.ID)
	require.True(t, found, subToken.ID)
	require.Equal(t, subToken, storedSubToken)
}

func TestRemoveSubTokens(t *testing.T) {
	dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper()

	subToken := types.NewSubToken(
		1,
		firstReserve,
	)

	// set sub token to store
	dsc.NFTKeeper.SetSubToken(ctx, firstID, subToken)
	storedSubToken, found := dsc.NFTKeeper.GetSubToken(ctx, firstID, subToken.ID)
	require.True(t, found, subToken.ID)
	require.Equal(t, subToken, storedSubToken)

	subToken.Reserve = secondReserve
	// remove sub token from the store
	dsc.NFTKeeper.RemoveSubToken(ctx, firstID, subToken.ID)

	// check if there is a sub token in the store
	storedSubToken, found = dsc.NFTKeeper.GetSubToken(ctx, firstID, subToken.ID)
	require.False(t, found, subToken.ID)
	require.Equal(t, types.SubToken{}, storedSubToken)
}

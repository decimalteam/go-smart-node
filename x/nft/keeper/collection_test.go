package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestGetAndSetCollection() {
	require := s.Require()
	denom := "Test_Set_Collection"
	pk := ed25519.GenPrivKey().PubKey() // TODO replace this pks
	owner := sdk.AccAddress(pk.Address())

	collection := types.Collection{
		Denom:   denom,
		Creator: owner.String(),
		Supply:  1,
	}

	s.nftKeeper.SetCollection(s.ctx, collection)

	// positive case get collection
	storeCollection, found := s.nftKeeper.GetCollection(s.ctx, owner, denom)
	require.True(found)
	require.Equal(collection.String(), storeCollection.String())
	// positive case get all collections
	storeCollections := s.nftKeeper.GetCollections(s.ctx)
	for _, storeCollection := range storeCollections {
		if storeCollection.Denom == denom {
			require.Equal(collection.String(), storeCollection.String())
		}
	}
}

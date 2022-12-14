package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	firstID        string = "first_id"
	secondID       string = "second_id"
	thirdID        string = "third_id"
	firstDenom     string = "first_test_denom"
	secondDenom    string = "second_test_denom"
	firstTokenURI  string = "first_token_uri"
	secondTokenURI string = "second_token_uri"
	thirdTokenURI  string = "third_token_uri"
)

var (
	firstReserve  = sdk.NewCoin(config.BaseDenom, types.DefaultMinReserveAmount)
	secondReserve = sdk.NewCoin(config.BaseDenom, types.DefaultMinReserveAmount.MulRaw(2))
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
	require.True(collection.Equal(storeCollection))
	// positive case get all collections
	storeCollections := s.nftKeeper.GetCollections(s.ctx)
	for _, storeCollection := range storeCollections {
		if storeCollection.Denom == denom {
			require.True(collection.Equal(storeCollection))
		}
	}
}

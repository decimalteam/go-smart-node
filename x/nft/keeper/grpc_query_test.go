package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (s *KeeperTestSuite) TestGRPCQueryCollections() {
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	denom := "Test_Query_Collections"
	pk := ed25519.GenPrivKey().PubKey()
	owner := sdk.AccAddress(pk.Address())

	collection := types.Collection{
		Denom:   denom,
		Creator: owner.String(),
	}
	keeper.SetCollection(ctx, collection)

	var req *types.QueryCollectionsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCollectionsRequest{
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Collections(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				//require.True(collection.Equal(&res.Validator))
				require.Equal(collection.String(), res.GetCollections()[0].String())
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCollection() {
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	denom := "Test_Query_Collection"
	pk := ed25519.GenPrivKey().PubKey()
	owner := sdk.AccAddress(pk.Address())
	pk = ed25519.GenPrivKey().PubKey()
	invalidOwner := sdk.AccAddress(pk.Address())

	collection := types.Collection{
		Denom:   denom,
		Creator: owner.String(),
	}
	keeper.SetCollection(ctx, collection)

	var req *types.QueryCollectionRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCollectionRequest{
					Creator: owner.String(),
					Denom:   denom,
				}
			},
			true,
		},
		{
			"invalid owner",
			func() {
				req = &types.QueryCollectionRequest{
					Creator: invalidOwner.String(),
					Denom:   denom,
				}
			},
			false,
		},
		{
			"invalid denom",
			func() {
				req = &types.QueryCollectionRequest{
					Creator: invalidOwner.String(),
					Denom:   "Invalid_denom",
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Collection(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				//require.True(collection.Equal(&res.Validator))
				require.Equal(collection.String(), res.Collection.String())
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCollectionsByCreator() {
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	pk := ed25519.GenPrivKey().PubKey()
	owner := sdk.AccAddress(pk.Address())
	pk = ed25519.GenPrivKey().PubKey()
	notOwner := sdk.AccAddress(pk.Address())
	denom1 := "Test_Query_Collections_By_Owner_1"
	denom2 := "Test_Query_Collections_By_Owner_2"
	denom3 := "Test_Query_Collections_By_Owner_3"

	collections := []types.Collection{
		{
			Denom:   denom1,
			Creator: owner.String(),
		},
		{
			Denom:   denom2,
			Creator: owner.String(),
		},
		{
			Denom:   denom3,
			Creator: owner.String(),
		},
	}

	for _, v := range collections {
		keeper.SetCollection(ctx, v)
	}

	var hits map[string]bool
	var req *types.QueryCollectionsByCreatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request with collections",
			func() {
				req = &types.QueryCollectionsByCreatorRequest{
					Creator: owner.String(),
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]bool{
					denom1: false,
					denom2: false,
					denom3: false,
				}
			},
			true,
		},
		{
			"valid request without collections",
			func() {
				req = &types.QueryCollectionsByCreatorRequest{
					Creator: notOwner.String(),
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]bool{}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.CollectionsByCreator(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, v := range res.Collections {
					if _, ok := hits[v.Denom]; ok {
						delete(hits, v.Denom)
					} else {
						s.T().Fatal("collection does not set, but it was included in the resp")
					}
				}
				if len(hits) != 0 {
					s.T().Fatal("not all set collections were returned")
				}
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryToken() {

}
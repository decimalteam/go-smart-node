package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"strconv"
)

func (s *KeeperTestSuite) TestGRPCQueryCollections() {
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	var (
		denom      = "Test_Query_Collections"
		pk         = ed25519.GenPrivKey().PubKey()
		owner      = sdk.AccAddress(pk.Address())
		collection = types.Collection{
			Denom:   denom,
			Creator: owner.String(),
		}
	)

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
				require.True(collection.Equal(res.GetCollections()[0]))
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

	var (
		denom   = "Test_Query_Collection"
		pk      = ed25519.GenPrivKey().PubKey()
		owner   = sdk.AccAddress(pk.Address())
		ID1     = "query_collection_1"
		ID2     = "query_collection_2"
		reserve = defaultCoin
	)
	pk = ed25519.GenPrivKey().PubKey()
	invalidOwner := sdk.AccAddress(pk.Address())

	collection := types.Collection{
		Denom:   denom,
		Creator: owner.String(),
		Tokens: types.Tokens{
			{
				Creator: owner.String(),
				Denom:   denom,
				ID:      ID1,
				URI:     ID1,
				Reserve: reserve,
			},
			{
				Creator: owner.String(),
				Denom:   denom,
				ID:      ID2,
				URI:     ID2,
				Reserve: reserve,
			},
		},
	}

	keeper.SetCollection(ctx, collection)
	for _, v := range collection.Tokens {
		keeper.CreateToken(ctx, collection, *v)
	}

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
				require.True(collection.Equal(res.Collection))
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

	var (
		pk      = ed25519.GenPrivKey().PubKey()
		owner   = sdk.AccAddress(pk.Address())
		denom1  = "Test_Query_Collections_By_Owner_1"
		denom2  = "Test_Query_Collections_By_Owner_2"
		denom3  = "Test_Query_Collections_By_Owner_3"
		ID1     = "query_collection_by_owner_1"
		ID2     = "query_collection_by_owner_2"
		reserve = defaultCoin
	)
	pk = ed25519.GenPrivKey().PubKey()
	notOwner := sdk.AccAddress(pk.Address())

	collections := []types.Collection{
		{
			Denom:   denom1,
			Creator: owner.String(),
			Tokens: types.Tokens{
				{
					Creator: owner.String(),
					Denom:   denom1,
					ID:      ID1,
					URI:     ID1,
					Reserve: reserve,
				},
				{
					Creator: owner.String(),
					Denom:   denom1,
					ID:      ID2,
					URI:     ID2,
					Reserve: reserve,
				},
			},
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

	for _, collection := range collections {
		keeper.SetCollection(ctx, collection)
		if collection.Tokens == nil || len(collection.Tokens) < 1 {
			continue
		}
		for _, token := range collection.Tokens {
			keeper.CreateToken(ctx, collection, *token)
		}
	}

	var hits map[string]types.Collection
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
				hits = map[string]types.Collection{
					denom1: collections[0],
					denom2: collections[1],
					denom3: collections[2],
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
				hits = map[string]types.Collection{}
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
				for _, resCollection := range res.Collections {
					if collection, ok := hits[resCollection.Denom]; ok {
						require.True(collection.Equal(resCollection))
						delete(hits, resCollection.Denom)
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
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	var (
		pk         = ed25519.GenPrivKey().PubKey()
		owner      = sdk.AccAddress(pk.Address())
		denom1     = "Test_Query_Token"
		ID         = "query_token_1"
		collection = types.Collection{
			Creator: owner.String(),
			Denom:   denom1,
		}
	)

	token := types.Token{
		Creator:   owner.String(),
		Denom:     denom1,
		ID:        ID,
		URI:       ID,
		Reserve:   defaultCoin,
		AllowMint: false,
		Minted:    0,
		Burnt:     0,
		SubTokens: types.SubTokens{
			{
				ID:      1,
				Owner:   owner.String(),
				Reserve: &defaultCoin,
			},
			{
				ID:      2,
				Owner:   owner.String(),
				Reserve: &defaultCoin,
			},
			{
				ID:      3,
				Owner:   owner.String(),
				Reserve: &defaultCoin,
			},
		},
	}

	keeper.SetCollection(ctx, collection)
	keeper.CreateToken(ctx, collection, token)

	for _, subtoken := range token.SubTokens {
		keeper.SetSubToken(ctx, token.ID, *subtoken)
	}

	var req *types.QueryTokenRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryTokenRequest{TokenId: ID}
			},
			true,
		},
		{
			"not exist token ID",
			func() {
				req = &types.QueryTokenRequest{TokenId: "not_exist"}
			},
			false,
		},
		{
			"empty token ID",
			func() {
				req = &types.QueryTokenRequest{}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Token(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(token.Equal(res.GetToken()))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQuerySubToken() {
	ctx, keeper, queryClient := s.ctx, s.nftKeeper, s.queryClient
	require := s.Require()

	var (
		pk         = ed25519.GenPrivKey().PubKey()
		owner      = sdk.AccAddress(pk.Address())
		denom1     = "Test_Query_Subtoken"
		ID         = "query_subtoken_1"
		reserve    = defaultCoin
		collection = types.Collection{
			Creator: owner.String(),
			Denom:   denom1,
		}
	)

	token := types.Token{
		Creator:   owner.String(),
		Denom:     denom1,
		ID:        ID,
		URI:       ID,
		Reserve:   reserve,
		AllowMint: false,
		Minted:    0,
		Burnt:     0,
		SubTokens: nil,
	}

	subtoken := types.SubToken{
		ID:      1,
		Owner:   owner.String(),
		Reserve: &reserve,
	}
	keeper.SetCollection(ctx, collection)
	keeper.CreateToken(ctx, collection, token)
	keeper.SetSubToken(ctx, ID, subtoken)

	var req *types.QuerySubTokenRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QuerySubTokenRequest{
					TokenId:    ID,
					SubTokenId: strconv.Itoa(int(subtoken.ID)),
				}
			},
			true,
		},
		{
			"not exist nft sub token",
			func() {
				req = &types.QuerySubTokenRequest{
					TokenId:    ID,
					SubTokenId: "20",
				}
			},
			false,
		},
		{
			"not exist nft",
			func() {
				req = &types.QuerySubTokenRequest{
					TokenId:    "not_exist",
					SubTokenId: strconv.Itoa(int(subtoken.ID)),
				}
			},
			false,
		},
		{
			"invalid subtoken ID",
			func() {
				req = &types.QuerySubTokenRequest{
					TokenId:    ID,
					SubTokenId: "invalid",
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.SubToken(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(subtoken.Equal(res.SubToken))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

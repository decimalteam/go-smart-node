package keeper_test

import (
	gocontext "context"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (s *KeeperTestSuite) TestGRPCQueryCoins() {
	ctx, k, queryClient := s.ctx, s.coinKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.Coin)
	coins := []types.Coin{
		{
			Denom:       "test",
			Title:       "test_query_coins",
			Creator:     addr.String(),
			CRR:         80,
			LimitVolume: sdk.NewInt(100000),
			MinVolume:   sdk.ZeroInt(),
			Volume:      sdk.NewInt(1000),
			Reserve:     sdk.NewInt(100),
		},
		{
			Denom:       "test1",
			Title:       "test_query_coins1",
			Creator:     addr.String(),
			CRR:         80,
			LimitVolume: sdk.NewInt(100000),
			MinVolume:   sdk.ZeroInt(),
			Volume:      sdk.NewInt(1000),
			Reserve:     sdk.NewInt(100),
		},
	}

	for _, v := range coins {
		k.SetCoin(ctx, v)
	}

	var req *types.QueryCoinsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCoinsRequest{
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]types.Coin{
					coins[0].Denom: coins[0],
					coins[1].Denom: coins[1],
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Coins(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resCoin := range res.Coins {
					if coin, ok := hits[resCoin.Denom]; ok {
						require.True(coin.Equal(resCoin))
						delete(hits, resCoin.Denom)
					} else {
						s.T().Fatal("coin does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all coins were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPQueryCoin() {
	ctx, k, queryClient := s.ctx, s.coinKeeper, s.queryClient
	require := s.Require()

	coin := types.Coin{
		Denom:       "test",
		Title:       "test_query_coin",
		Creator:     addr.String(),
		CRR:         80,
		LimitVolume: sdk.NewInt(100000),
		MinVolume:   sdk.ZeroInt(),
		Volume:      sdk.NewInt(1000),
		Reserve:     sdk.NewInt(100),
	}
	k.SetCoin(ctx, coin)

	var req *types.QueryCoinRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCoinRequest{
					Denom: coin.Denom,
				}
			},
			true,
		},
		{
			"not exist coin",
			func() {
				req = &types.QueryCoinRequest{
					Denom: "not exist coin",
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Coin(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(coin.Equal(res.GetCoin()))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPQueryChecks() {
	ctx, k, queryClient := s.ctx, s.coinKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.Check)
	checks := []types.Check{
		{
			ChainID:  "chain-id",
			Coin:     sdk.NewCoin(baseDenom, baseAmount),
			Nonce:    []byte("nonce"),
			DueBlock: 1,
			Lock:     []byte("lock"),
			V:        sdk.NewInt(1),
			R:        sdk.NewInt(1),
			S:        sdk.NewInt(1),
		},
		{
			ChainID:  "chain-id1",
			Coin:     sdk.NewCoin(baseDenom, baseAmount),
			Nonce:    []byte("nonce1"),
			DueBlock: 2,
			Lock:     []byte("lock1"),
			V:        sdk.NewInt(1),
			R:        sdk.NewInt(1),
			S:        sdk.NewInt(1),
		},
	}

	for _, v := range checks {
		k.SetCheck(ctx, &v)
	}

	var req *types.QueryChecksRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"Valid request",
			func() {
				req = &types.QueryChecksRequest{
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]types.Check{
					checks[0].ChainID: checks[0],
					checks[1].ChainID: checks[1],
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Checks(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resCheck := range res.Checks {
					if check, ok := hits[resCheck.ChainID]; ok {
						require.True(check.Equal(resCheck))
						delete(hits, resCheck.ChainID)
					} else {
						s.T().Fatal("check does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all check were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPQueryCheck() {
	ctx, k, queryClient := s.ctx, s.coinKeeper, s.queryClient
	require := s.Require()

	check := types.Check{
		ChainID:  "chain-id",
		Coin:     sdk.NewCoin(baseDenom, baseAmount),
		Nonce:    []byte("nonce"),
		DueBlock: 1,
		Lock:     []byte("lock"),
		V:        sdk.NewInt(1),
		R:        sdk.NewInt(1),
		S:        sdk.NewInt(1),
	}
	k.SetCheck(ctx, &check)
	checkHash := check.HashFull()

	invalidCheck := types.Check{
		ChainID:  "chain-id1",
		Coin:     sdk.NewCoin(baseDenom, baseAmount),
		Nonce:    []byte("nonce1"),
		DueBlock: 1,
		Lock:     []byte("lock1"),
		V:        sdk.NewInt(1),
		R:        sdk.NewInt(1),
		S:        sdk.NewInt(1),
	}
	invalidHash := invalidCheck.HashFull()

	var req *types.QueryCheckRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCheckRequest{
					Hash: checkHash[:],
				}
			},
			true,
		},
		{
			"not exist check",
			func() {
				req = &types.QueryCheckRequest{
					Hash: invalidHash[:],
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Check(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(check.Equal(res.GetCheck()))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryParams() {
	ctx, k, queryClient := s.ctx, s.coinKeeper, s.queryClient
	require := s.Require()

	params := k.GetParams(ctx)

	res, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	require.NoError(err)
	require.NotNil(res)
	require.True(res.Params.Equal(params))
}

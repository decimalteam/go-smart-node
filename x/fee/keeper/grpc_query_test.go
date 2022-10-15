package keeper_test

import (
	gocontext "context"
	"fmt"
	"time"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

func (s *KeeperTestSuite) TestGRPCQueryCoinPrices() {
	ctx, k, queryClient := s.ctx, s.feeKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.CoinPrice)
	coinPrices := []types.CoinPrice{
		{
			Denom:     "query_coin_prices1",
			Quote:     "",
			Price:     sdk.NewDec(1),
			UpdatedAt: time.Now(),
		},
		{
			Denom:     "query_coin_prices2",
			Quote:     "",
			Price:     sdk.NewDec(2),
			UpdatedAt: time.Now(),
		},
	}

	for _, v := range coinPrices {
		require.NoError(k.SavePrice(ctx, v))
	}

	var req *types.QueryCoinPricesRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryCoinPricesRequest{}
				hits = map[string]types.CoinPrice{
					coinPrices[0].Denom: coinPrices[0],
					coinPrices[1].Denom: coinPrices[1],
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.CoinPrices(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resPrice := range res.Prices {
					if price, ok := hits[resPrice.Denom]; ok {
						require.True(price.Equal(resPrice))
						delete(hits, resPrice.Denom)
					} else {
						s.T().Fatal("price does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all prices were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCoinPrice() {
	ctx, k, queryClient := s.ctx, s.feeKeeper, s.queryClient
	require := s.Require()

	coinPrices := []types.CoinPrice{
		{
			Denom:     "query_coin_price1",
			Quote:     "",
			Price:     sdk.NewDec(1),
			UpdatedAt: time.Now(),
		},
		{
			Denom:     "query_coin_price2",
			Quote:     "",
			Price:     sdk.NewDec(2),
			UpdatedAt: time.Now(),
		},
	}

	require.NoError(k.SavePrice(ctx, coinPrices[0]))

	var req *types.QueryCoinPriceRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"Valid request",
			func() {
				req = &types.QueryCoinPriceRequest{
					Denom: coinPrices[0].Denom,
					Quote: "",
				}
			},
			true,
		},
		{
			"not exists coin price",
			func() {
				req = &types.QueryCoinPriceRequest{
					Denom: coinPrices[1].Denom,
					Quote: "",
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.CoinPrice(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(coinPrices[0].Equal(res.GetPrice()))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryModuleParams() {
	ctx, k, queryClient := s.ctx, s.feeKeeper, s.queryClient
	require := s.Require()

	mparams := k.GetModuleParams(ctx)
	res, err := queryClient.ModuleParams(gocontext.Background(), &types.QueryModuleParamsRequest{})

	require.NoError(err)
	require.True(mparams.Equal(res.GetParams()))
}

func (s *KeeperTestSuite) TestGRPCQueryParams() {
	ctx, k, queryClient := s.ctx, s.feeKeeper, s.fmQueryClient
	require := s.Require()
	err := k.SavePrice(ctx, types.CoinPrice{
		Denom: config.BaseDenom,
		Quote: feeconfig.DefaultQuote,
		Price: sdk.OneDec(),
	})
	require.NoError(err)

	params := k.GetParams(ctx)
	res, err := queryClient.Params(gocontext.Background(), &feemarkettypes.QueryParamsRequest{})
	require.NoError(err)
	resParams := res.GetParams()
	require.Equal(params.String(), resParams.String())
}

func (s *KeeperTestSuite) TestGRPCQueryBaseFee() {
	ctx, k, queryClient := s.ctx, s.feeKeeper, s.fmQueryClient
	require := s.Require()

	err := k.SavePrice(ctx, types.CoinPrice{
		Denom:     config.BaseDenom,
		Quote:     feeconfig.DefaultQuote,
		Price:     sdk.NewDec(1),
		UpdatedAt: time.Now(),
	})
	require.NoError(err)

	basefee := k.GetBaseFee(ctx)

	res, err := queryClient.BaseFee(gocontext.Background(), &feemarkettypes.QueryBaseFeeRequest{})
	require.NoError(err)
	require.Equal(basefee.String(), res.BaseFee.String())
}

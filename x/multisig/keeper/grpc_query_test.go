package keeper_test

import (
	gocontext "context"
	"fmt"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func (s *KeeperTestSuite) TestGRPCQueryWallets() {
	ctx, k, queryClient := s.ctx, s.msKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.Wallet)
	k.SetWallet(ctx, *existsWallet)
	k.SetWallet(ctx, *defaultWallet)

	var req *types.QueryWalletsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryWalletsRequest{
					Owner: user1.String(),
					Pagination: &query.PageRequest{
						Offset: 0,
						Limit:  100,
					},
				}
				hits = map[string]types.Wallet{
					defaultWallet.Address: *defaultWallet,
					existsWallet.Address:  *existsWallet,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Wallets(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resWallet := range res.Wallets {
					if wallet, ok := hits[resWallet.Address]; ok {
						require.Equal(wallet.String(), resWallet.String()) // TODO replace after regenerate proto
						delete(hits, resWallet.Address)
					} else {
						s.T().Fatal("wallet does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all wallets were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryWallet() {
	ctx, k, queryClient := s.ctx, s.msKeeper, s.queryClient
	require := s.Require()

	k.SetWallet(ctx, *existsWallet)
	var req *types.QueryWalletRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryWalletRequest{
					Wallet: existsWallet.Address,
				}
			},
			true,
		},
		{
			"not exists wallet",
			func() {
				req = &types.QueryWalletRequest{
					Wallet: defaultWallet.Address,
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Wallet(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				wallet := res.GetWallet()
				require.Equal(existsWallet.String(), wallet.String()) // TODO replace after regenerate proto
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryTransactions() {
	ctx, k, queryClient := s.ctx, s.msKeeper, s.queryClient
	require := s.Require()

	hits := make(map[string]types.Transaction)
	transactions := []types.Transaction{
		{
			Id:     "first",
			Wallet: existsWallet.Address,
		},
		{
			Id:     "Second",
			Wallet: existsWallet.Address,
		},
	}

	for _, v := range transactions {
		k.SetTransaction(ctx, v)
	}

	var req *types.QueryTransactionsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryTransactionsRequest{
					Wallet: existsWallet.Address,
				}
				hits = map[string]types.Transaction{
					transactions[0].Id: transactions[0],
					transactions[1].Id: transactions[1],
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Transactions(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				for _, resTx := range res.Transactions {
					if tx, ok := hits[resTx.Id]; ok {
						// TODO: add more checks
						require.Equal(tx.Id, resTx.Id)
						require.Equal(tx.Wallet, resTx.Wallet)
						delete(hits, resTx.Id)
					} else {
						s.T().Fatal("wallet does not set, but it was included in the resp")
					}
				}
				require.Equal(0, len(hits), "not all wallets were returned")
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryCoinPrice() {
	ctx, k, queryClient := s.ctx, s.msKeeper, s.queryClient
	require := s.Require()

	tx := types.Transaction{
		Id:     "first",
		Wallet: existsWallet.Address,
		Message: *sdkcodec.UnsafePackAny(cointypes.NewMsgSendCoin(
			sdk.MustAccAddressFromBech32(existsWallet.Address),
			user1,
			sdk.NewInt64Coin("del", 1)),
		),
		CreatedAt: 100,
	}
	k.SetTransaction(ctx, tx)

	var req *types.QueryTransactionRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryTransactionRequest{
					Id: "first",
				}
			},
			true,
		},
		{
			"not exist tx",
			func() {
				req = &types.QueryTransactionRequest{
					Id: "not exist",
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Transaction(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.Equal(tx.Id, res.Transaction.Id)
				require.Equal(tx.Wallet, res.Transaction.Wallet)
				require.Equal(tx.Message.Value, res.Transaction.Message.Value)
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

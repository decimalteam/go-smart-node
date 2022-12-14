package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

var (
	defaultDenom = cmdcfg.BaseDenom
	defaultCoins = sdk.NewCoins(sdk.NewCoin(defaultDenom, sdk.NewInt(1)))
)

func (s *KeeperTestSuite) TestMsgCreateWallet() {
	ctx, _, msgServer := s.ctx, s.msKeeper, s.msgServer
	require := s.Require()

	testCases := []struct {
		name   string
		input  *types.MsgCreateWallet
		expErr bool
		ctx    sdk.Context
	}{
		{
			"valid request",
			types.NewMsgCreateWallet(user4, defaultOwners, defaultWeights, defaultThreeshold),
			false,
			ctx,
		},
		{
			"wallet exists",
			types.NewMsgCreateWallet(user4, defaultOwners, defaultWeights, defaultThreeshold),
			true,
			ctx,
		},
		{
			"account address exists",
			types.NewMsgCreateWallet(user4, defaultOwners, defaultWeights, defaultThreeshold),
			true,
			ctx.WithTxBytes([]byte{1}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.CreateWallet(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

/* TODO: need fix with MsgServiceRouter
func (s *KeeperTestSuite) TestMsgCreateTransaction() {
	ctx, k, msgServer := s.ctx, s.msKeeper, s.msgServer
	require := s.Require()

	k.SetWallet(ctx, *defaultWallet)
	createTestCases := []struct {
		do       func()
		name     string
		wallet   sdk.AccAddress
		receiver sdk.AccAddress
		expErr   bool
	}{
		{
			func() {},
			"valid request -- first confirm",
			defaultWalletAddress,
			user1,
			false,
		},
		{
			func() {},
			"wallet not exists",
			existsWalletAddress,
			user4,
			true,
		},
		{
			func() { k.SetWallet(ctx, *existsWallet) },
			"insufuccient funds",
			existsWalletAddress,
			user4,
			true,
		},
	}

	createResults := make([]string, 0)
	for _, tc := range createTestCases {
		tc := tc
		tc.do()
		s.T().Run(tc.name, func(t *testing.T) {
			tx, err := types.NewMsgCreateTransaction(user3, tc.wallet.String(), cointypes.NewMsgSendCoin(tc.wallet, tc.receiver, defaultCoins[0]))
			require.NoError(err)
			res, err := msgServer.CreateTransaction(ctx, tx)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
				createResults = append(createResults, res.ID)
			}
		})
	}

	signTestCases := []struct {
		name   string
		input  *types.MsgSignTransaction
		expErr bool
	}{
		{
			"valid request -- second confirm",
			types.NewMsgSignTransaction(user2, createResults[0]),
			false,
		},
		{
			"third confirm -- enough conformations",
			types.NewMsgSignTransaction(user1, createResults[0]),
			true,
		},
		{
			"not exists tx with current ID",
			types.NewMsgSignTransaction(user2, "2"),
			true,
		},
	}

	for _, tc := range signTestCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.SignTransaction(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
*/

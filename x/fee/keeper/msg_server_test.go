package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"time"
)

func (s *KeeperTestSuite) TestMsgUpdateCoinPrices() {
	ctx, _, msgServer := s.ctx, s.feeKeeper, s.msgServer
	require := s.Require()

	pk := ed25519.GenPrivKey().PubKey()
	invalidOracle := sdk.AccAddress(pk.Address()).String()

	testCases := []struct {
		name   string
		input  *types.MsgUpdateCoinPrices
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgUpdateCoinPrices(oracle, []types.CoinPrice{
				{
					Denom:     baseDenom,
					Quote:     "",
					Price:     sdk.NewDec(10000),
					UpdatedAt: time.Now(),
				},
			}),
			false,
		},
		{
			"unknown oracle",
			types.NewMsgUpdateCoinPrices(invalidOracle, []types.CoinPrice{
				{
					Denom:     baseDenom,
					Quote:     "",
					Price:     sdk.NewDec(10000),
					UpdatedAt: time.Now(),
				},
			}),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.UpdateCoinPrices(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

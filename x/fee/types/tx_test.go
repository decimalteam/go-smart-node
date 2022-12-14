package types

import (
	"testing"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSavePriceValidation(t *testing.T) {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)

	var testCases = []struct {
		tag         string
		sender      string
		prices      []CoinPrice
		expectError bool
	}{
		{
			"valid",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj",
			[]CoinPrice{
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(1),
				},
				{
					Denom: "btc",
					Quote: "usd",
					Price: sdk.NewDec(1),
				},
			},
			false,
		},
		{
			"invalid sender",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj1",
			[]CoinPrice{
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(1),
				},
			},
			true,
		},
		{
			"invalid price 1",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj",
			[]CoinPrice{
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(-1),
				},
			},
			true,
		},
		{
			"invalid price 2",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj",
			[]CoinPrice{
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(0),
				},
			},
			true,
		},
		{
			"invalid price 3",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj",
			[]CoinPrice{},
			true,
		},
		{
			"invalid price 4",
			"d01qql8ag4cluz6r4dz28p3w00dnc9w8ueuak90zj",
			[]CoinPrice{
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(1),
				},
				{
					Denom: cmdcfg.BaseDenom,
					Quote: "usd",
					Price: sdk.NewDec(2),
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		msg := NewMsgUpdateCoinPrices(tc.sender, tc.prices)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

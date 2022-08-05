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
		price       sdk.Dec
		expectError bool
	}{
		{
			"valid",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			sdk.NewDec(1),
			false,
		},
		{
			"invalid sender",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd1",
			sdk.NewDec(1),
			true,
		},
		{
			"invalid price 1",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			sdk.NewDec(-1),
			true,
		},
		{
			"invalid price 2",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			sdk.NewDec(0),
			true,
		},
	}

	for _, tc := range testCases {
		msg := NewMsgSaveBaseDenomPrice(tc.sender, "del", tc.price)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

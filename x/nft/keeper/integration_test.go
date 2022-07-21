package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	Denom1    = "test_denom_1"
	Denom2    = "test_denom_2"
	Denom3    = "test_denom_3"
	ID1       = "1"
	ID2       = "2"
	ID3       = "3"
	TokenURI1 = "https://google.com/token-1.json"
	TokenURI2 = "https://google.com/token-2.json"
)

func getAddrs(dsc *app.DSC, ctx sdk.Context, number int) []sdk.AccAddress {
	addrs := app.AddTestAddrsIncremental(dsc, ctx, number, sdk.Coins{
		{
			Denom:  "del",
			Amount: helpers.EtherToWei(sdk.NewInt(1000000000000)),
		},
	})

	return addrs
}

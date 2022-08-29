package coin_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// create invalid coin
func invalidCoin() sdk.Coin {
	return sdk.Coin{
		Denom:  "invalidDenom",
		Amount: sdk.NewInt(100000000),
	}
}

// create valid coin
func validCoin(denom string, amount int64) sdk.Coin {
	return sdk.Coin{
		Denom:  denom,
		Amount: helpers.EtherToWei(sdk.NewInt(amount)),
	}
}

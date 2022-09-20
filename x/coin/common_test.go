package coin_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// create invalid coin
func invalidCoin() sdk.Coin {
	return sdk.Coin{
		Denom:  "invalidDenom",
		Amount: sdkmath.NewInt(100000000),
	}
}

// create valid coin
func validCoin(denom string, amount int64) sdk.Coin {
	return sdk.Coin{
		Denom:  denom,
		Amount: helpers.EtherToWei(sdkmath.NewInt(amount)),
	}
}

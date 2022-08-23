package types

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const regName = "^[a-zA-Z0-9_-]{1,255}$"

var MinReserve = sdk.NewInt(100)

var NewMinReserve = helpers.BipToPip(sdk.NewInt(100))
var NewMinReserve2 = helpers.BipToPip(sdk.NewInt(1))

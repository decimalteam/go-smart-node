package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

const maxCoinNameBytes = 64
const allowedCoinSymbols = "^[a-zA-Z][a-zA-Z0-9]{2,9}$"

var minCoinSupply = sdk.NewInt(1)
var maxCoinSupply = helpers.EtherToWei(sdk.NewInt(1000000000000000))
var MinCoinReserve = helpers.EtherToWei(sdk.NewInt(1000))

// []byte{'c', 'o', 'i', 'n'} is encoded as 'ds1vdhkjmsygul4t' or 'dc1vdhkjms35u4c0'...
var StubCoinAddress = sdk.AccAddress([]byte{'c', 'o', 'i', 'n'})

// pool for legacy balances
var LegacyCoinPool = "legacy_coin_pool"

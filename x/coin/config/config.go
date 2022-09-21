package config

import (
	"regexp"

	sdkmath "cosmossdk.io/math"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

var (
	CoinDenomRegExp    = "^[a-z][a-z0-9]{2,9}$"
	CoinDenomValidator = regexp.MustCompile(CoinDenomRegExp)
	MaxCoinTitleLength = 64
	MinCoinSupply      = helpers.EtherToWei(sdkmath.NewInt(1))
	MaxCoinSupply      = helpers.EtherToWei(sdkmath.NewInt(1_000_000_000_000_000))
	MinCoinReserve     = helpers.EtherToWei(sdkmath.NewInt(1000))
)

package keeper

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDecreasingFactor(t *testing.T) {
	var testCases = []struct {
		name              string
		crr               uint32
		volume            sdkmath.Int
		reserve           sdkmath.Int
		amountInCollector sdkmath.Int
		amountToBurn      sdkmath.Int
		expectFactor      sdk.Dec
	}{
		{
			name:              "minimal",
			crr:               100,
			volume:            config.MinCoinSupply,
			reserve:           config.MinCoinReserve,
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(10)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(10)),
			expectFactor:      sdk.ZeroDec(),
		},
		{
			name:              "normal without collector",
			crr:               100,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(0)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(10)),
			expectFactor:      sdk.OneDec(),
		},
		{
			name:              "normal with collector",
			crr:               100,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(20)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(20)),
			expectFactor:      sdk.OneDec(),
		},
		{
			name:              "normal burn all",
			crr:               100,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(500)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(500)),
			expectFactor:      sdk.MustNewDecFromStr("0.001"),
		},
		{
			name:              "normal without collector 50 crr",
			crr:               50,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(0)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(100)),
			expectFactor:      sdk.OneDec(),
		},
		{
			name:              "normal without collector 25 crr",
			crr:               25,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(0)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(100)),
			expectFactor:      sdk.OneDec(),
		},
		{
			name:              "normal without collector 10 crr, big burn",
			crr:               10,
			volume:            config.MinCoinSupply.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			reserve:           config.MinCoinReserve.Add(helpers.EtherToWei(sdkmath.NewInt(1000))),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(0)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(500)),
			expectFactor:      sdk.MustNewDecFromStr("0.134067950943311560"),
		},
	}
	for _, tc := range testCases {
		coinInfo := types.Coin{
			CRR:     tc.crr,
			Volume:  tc.volume,
			Reserve: tc.reserve,
		}
		factor := CalculateDecreasingFactor(coinInfo, tc.amountInCollector, tc.amountToBurn)
		require.True(t, tc.expectFactor.Equal(factor), tc.name)
	}
}

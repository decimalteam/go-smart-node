package keeper

import (
	"fmt"
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
		{
			name:              "normal, result 0.5",
			crr:               100,
			volume:            helpers.EtherToWei(sdkmath.NewInt(1000)),
			reserve:           helpers.EtherToWei(sdkmath.NewInt(2000)),
			amountInCollector: helpers.EtherToWei(sdkmath.NewInt(400)),
			amountToBurn:      helpers.EtherToWei(sdkmath.NewInt(200)),
			expectFactor:      sdk.MustNewDecFromStr("0.5"),
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

// This test need for finding coin parameters for specific CalculateDecreasingFactor
// Normaly must skip
func TestDecreasingParams(t *testing.T) {
	t.Skip()
	CRR := []uint32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	VOLUME := []int64{1000, 2000}
	RESERVE := []int64{1000, 2000}
	ACOL := []int64{100, 200, 300, 400, 500}
	ABURN := []int64{100, 200, 300, 400, 500}

	for _, crr := range CRR {
		for _, vol := range VOLUME {
			for _, res := range RESERVE {
				coinInfo := types.Coin{
					CRR:     crr,
					Volume:  helpers.EtherToWei(sdk.NewInt(vol)),
					Reserve: helpers.EtherToWei(sdk.NewInt(res)),
				}
				for _, col := range ACOL {
					for _, burn := range ABURN {
						factor := CalculateDecreasingFactor(coinInfo, helpers.EtherToWei(sdk.NewInt(col)), helpers.EtherToWei(sdk.NewInt(burn)))
						fmt.Printf("crr=%d, vol=%d, res=%d, coll=%d, burn=%d, factor=%s\n", crr, vol, res, col, burn, factor)
					}
				}
			}
		}
	}
}

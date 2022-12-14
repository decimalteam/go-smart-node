package types

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinconfig "bitbucket.org/decimalteam/go-smart-node/x/coin/config"
)

// Init global cosmos sdk cmdcfg
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(cmdcfg.Bech32PrefixAccAddr, cmdcfg.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(cmdcfg.Bech32PrefixValAddr, cmdcfg.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(cmdcfg.Bech32PrefixConsAddr, cmdcfg.Bech32PrefixConsPub)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func TestCreateCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag            string
		sender         string
		denom          string
		title          string
		crr            uint32
		initialVolume  sdkmath.Int
		initialReserve sdkmath.Int
		limitVolume    sdkmath.Int
		identity       string
		expectError    bool
	}{
		{
			tag:            "valid coin",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    false,
		},
		{
			tag:            "invalid sender",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid title",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          RandomString(65),
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid denom 1",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          "de",
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid denom 2",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          "del45678901",
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "low crr",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            9,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "high crr",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            101,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "low volume",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.FinneyToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "high volume",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  coinconfig.MaxCoinSupply.Add(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid reserve",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(999)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "initial > limit",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(2)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdkmath.NewInt(1)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid limit volume",
			sender:         "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:          cmdcfg.BaseDenom,
			title:          "some coin",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdkmath.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdkmath.NewInt(2000)),
			limitVolume:    coinconfig.MaxCoinSupply.Add(sdkmath.NewInt(1)),
			identity:       "some coin",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		msg := MsgCreateCoin{
			Sender:         tc.sender,
			Denom:          tc.denom,
			Title:          tc.title,
			CRR:            tc.crr,
			InitialVolume:  tc.initialVolume,
			InitialReserve: tc.initialReserve,
			LimitVolume:    tc.limitVolume,
			Identity:       tc.identity,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestUpdateCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag         string
		sender      string
		denom       string
		limitVolume sdkmath.Int
		identity    string
		expectError bool
	}{
		{
			tag:         "valid coin update",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:       cmdcfg.BaseDenom,
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			denom:       cmdcfg.BaseDenom,
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid denom 1",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:       "de",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid denom 2",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:       "del45678901",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid limit",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			denom:       cmdcfg.BaseDenom,
			limitVolume: coinconfig.MaxCoinSupply.Add(sdkmath.NewInt(1)),
			identity:    "some coin",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgUpdateCoin{
			Sender:      tc.sender,
			Denom:       tc.denom,
			LimitVolume: tc.limitVolume,
			Identity:    tc.identity,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}

}

func TestSendCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag         string
		sender      string
		recipient   string
		coin        sdk.Coin
		expectError bool
	}{
		{
			tag:         "valid send",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			recipient:   "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			recipient:   "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid recipient",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			recipient:   "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m0",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid recipient (=sender)",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			recipient:   "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid amount",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			recipient:   "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(0)),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgSendCoin{
			Sender:    tc.sender,
			Recipient: tc.recipient,
			Coin:      tc.coin,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}

}

func TestMultiSendCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag         string
		sender      string
		sends       []MultiSendEntry
		expectError bool
	}{
		{
			tag:    "valid send",
			sender: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			sends: []MultiSendEntry{
				{Recipient: "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m", Coin: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))},
				{Recipient: "d018c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg4x0y6k", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: false,
		},
		{
			tag:    "invalid sender",
			sender: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			sends: []MultiSendEntry{
				{Recipient: "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m", Coin: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))},
				{Recipient: "d018c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg4x0y6k", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid recipient",
			sender: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			sends: []MultiSendEntry{
				{Recipient: "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m", Coin: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))},
				{Recipient: "d018c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg4x0y6k0", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid recipient=sender",
			sender: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			sends: []MultiSendEntry{
				{Recipient: "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m", Coin: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))},
				{Recipient: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid amount",
			sender: "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			sends: []MultiSendEntry{
				{Recipient: "d01w98j4vk6dkpyndjnv5dn2eemesq6a2c2kzwv2m", Coin: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))},
				{Recipient: "d018c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg4x0y6k", Coin: sdk.NewCoin("btc", sdkmath.NewInt(0))},
			},
			expectError: true,
		},
		{
			tag:         "no sends",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			sends:       []MultiSendEntry{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgMultiSendCoin{
			Sender: tc.sender,
			Sends:  tc.sends,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}

}

func TestBuyCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag           string
		sender        string
		coinToBuy     sdk.Coin
		maxCoinToSell sdk.Coin
		expectError   bool
	}{
		{
			tag:           "valid buy",
			sender:        "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToBuy:     sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   false,
		},
		{
			tag:           "invalid sender",
			sender:        "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			coinToBuy:     sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid buy amount",
			sender:        "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToBuy:     sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(0)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid sell amount",
			sender:        "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToBuy:     sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(0)),
			expectError:   true,
		},
		{
			tag:           "same coin",
			sender:        "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToBuy:     sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		msg := MsgBuyCoin{
			Sender:        tc.sender,
			CoinToBuy:     tc.coinToBuy,
			MaxCoinToSell: tc.maxCoinToSell,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestSellCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag          string
		sender       string
		coinToSell   sdk.Coin
		minCoinToBuy sdk.Coin
		expectError  bool
	}{
		{
			tag:          "valid sell",
			sender:       "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToSell:   sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  false,
		},
		{
			tag:          "invalid sender",
			sender:       "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			coinToSell:   sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "invalid sell amount",
			sender:       "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToSell:   sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(0)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "same coin",
			sender:       "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinToSell:   sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		msg := MsgSellCoin{
			Sender:       tc.sender,
			CoinToSell:   tc.coinToSell,
			MinCoinToBuy: tc.minCoinToBuy,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestSellAllCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag             string
		sender          string
		coinDenomToSell string
		minCoinToBuy    sdk.Coin
		expectError     bool
	}{
		{
			tag:             "valid sell",
			sender:          "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinDenomToSell: cmdcfg.BaseDenom,
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:     false,
		},
		{
			tag:             "invalid sender",
			sender:          "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			coinDenomToSell: cmdcfg.BaseDenom,
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:     true,
		},
		{
			tag:             "invalid buy amount",
			sender:          "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinDenomToSell: cmdcfg.BaseDenom,
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(0)),
			expectError:     true,
		},
		{
			tag:             "same coin",
			sender:          "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coinDenomToSell: cmdcfg.BaseDenom,
			minCoinToBuy:    sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		msg := MsgSellAllCoin{
			Sender:          tc.sender,
			CoinDenomToSell: tc.coinDenomToSell,
			MinCoinToBuy:    tc.minCoinToBuy,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestBurnCoin(t *testing.T) {
	initConfig()

	var testCases = []struct {
		tag         string
		sender      string
		coin        sdk.Coin
		expectError bool
	}{
		{
			tag:         "valid burn",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv0",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid burn amount",
			sender:      "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			coin:        sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(0)),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgBurnCoin{
			Sender: tc.sender,
			Coin:   tc.coin,
		}
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

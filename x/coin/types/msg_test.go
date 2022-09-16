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
		crr            uint64
		initialVolume  sdkmath.Int
		initialReserve sdkmath.Int
		limitVolume    sdkmath.Int
		identity       string
		expectError    bool
	}{
		{
			tag:            "valid coin",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:          "del",
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
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:       "del",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			denom:       "del",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid denom 1",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:       "de",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid denom 2",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:       "del45678901",
			limitVolume: helpers.EtherToWei(sdkmath.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid limit",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			denom:       "del",
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
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			recipient:   "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			recipient:   "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid recipient",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			recipient:   "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy0",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid recipient (=sender)",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			recipient:   "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid amount",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			recipient:   "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(0)),
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
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []MultiSendEntry{
				{Recipient: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy", Coin: sdk.NewCoin("del", sdkmath.NewInt(1))},
				{Recipient: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: false,
		},
		{
			tag:    "invalid sender",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			sends: []MultiSendEntry{
				{Recipient: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy", Coin: sdk.NewCoin("del", sdkmath.NewInt(1))},
				{Recipient: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid recipient",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []MultiSendEntry{
				{Recipient: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy", Coin: sdk.NewCoin("del", sdkmath.NewInt(1))},
				{Recipient: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f0", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid recipient=sender",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []MultiSendEntry{
				{Recipient: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy", Coin: sdk.NewCoin("del", sdkmath.NewInt(1))},
				{Recipient: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn", Coin: sdk.NewCoin("btc", sdkmath.NewInt(2))},
			},
			expectError: true,
		},
		{
			tag:    "invalid amount",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []MultiSendEntry{
				{Recipient: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy", Coin: sdk.NewCoin("del", sdkmath.NewInt(1))},
				{Recipient: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f", Coin: sdk.NewCoin("btc", sdkmath.NewInt(0))},
			},
			expectError: true,
		},
		{
			tag:         "no sends",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
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
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   false,
		},
		{
			tag:           "invalid sender",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinToBuy:     sdk.NewCoin("del", sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid buy amount",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdkmath.NewInt(0)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid sell amount",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdkmath.NewInt(0)),
			expectError:   true,
		},
		{
			tag:           "same coin",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdkmath.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("del", sdkmath.NewInt(1)),
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
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  false,
		},
		{
			tag:          "invalid sender",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinToSell:   sdk.NewCoin("del", sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "invalid sell amount",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdkmath.NewInt(0)),
			minCoinToBuy: sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "same coin",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdkmath.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("del", sdkmath.NewInt(1)),
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
			sender:          "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinDenomToSell: "del",
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:     false,
		},
		{
			tag:             "invalid sender",
			sender:          "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinDenomToSell: "del",
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(1)),
			expectError:     true,
		},
		{
			tag:             "invalid buy amount",
			sender:          "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinDenomToSell: "del",
			minCoinToBuy:    sdk.NewCoin("btc", sdkmath.NewInt(0)),
			expectError:     true,
		},
		{
			tag:             "same coin",
			sender:          "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinDenomToSell: "del",
			minCoinToBuy:    sdk.NewCoin("del", sdkmath.NewInt(1)),
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
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid burn amount",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coin:        sdk.NewCoin("del", sdkmath.NewInt(0)),
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

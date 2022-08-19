package types

import (
	"math/rand"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Init global cosmos sdk config
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.Bech32PrefixValAddr, config.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.Bech32PrefixConsAddr, config.Bech32PrefixConsPub)
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
		title          string
		symbol         string
		crr            uint64
		initialVolume  sdk.Int
		initialReserve sdk.Int
		limitVolume    sdk.Int
		identity       string
		expectError    bool
	}{
		{
			tag:            "valid coin",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    false,
		},
		{
			tag:            "invalid sender",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid title",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          RandomString(65),
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid symbol 1",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "de",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid symbol 2",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del45678901",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "low crr",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            9,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "high crr",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            101,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "low volume",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.FinneyToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "high volume",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  maxCoinSupply.Add(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid reserve",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(999)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(100000)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "initial > limit",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(2)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    helpers.EtherToWei(sdk.NewInt(1)),
			identity:       "some coin",
			expectError:    true,
		},
		{
			tag:            "invalid limit volume",
			sender:         "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			title:          "some coin",
			symbol:         "del",
			crr:            20,
			initialVolume:  helpers.EtherToWei(sdk.NewInt(1)),
			initialReserve: helpers.EtherToWei(sdk.NewInt(2000)),
			limitVolume:    maxCoinSupply.Add(sdk.NewInt(1)),
			identity:       "some coin",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		msg := MsgCreateCoin{
			Sender:         tc.sender,
			Title:          tc.title,
			Symbol:         tc.symbol,
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
		symbol      string
		limitVolume sdk.Int
		identity    string
		expectError bool
	}{
		{
			tag:         "valid coin update",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			symbol:      "del",
			limitVolume: helpers.EtherToWei(sdk.NewInt(100000)),
			identity:    "some coin",
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			symbol:      "del",
			limitVolume: helpers.EtherToWei(sdk.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid symbol 1",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			symbol:      "de",
			limitVolume: helpers.EtherToWei(sdk.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid symbol 2",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			symbol:      "del45678901",
			limitVolume: helpers.EtherToWei(sdk.NewInt(100000)),
			identity:    "some coin",
			expectError: true,
		},
		{
			tag:         "invalid limit",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			symbol:      "del",
			limitVolume: maxCoinSupply.Add(sdk.NewInt(1)),
			identity:    "some coin",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgUpdateCoin{
			Sender:      tc.sender,
			Symbol:      tc.symbol,
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
		receiver    string
		coin        sdk.Coin
		expectError bool
	}{
		{
			tag:         "valid send",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			receiver:    "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			receiver:    "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid receiver",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			receiver:    "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy0",
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid receiver (=sender)",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			receiver:    "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid amount",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			receiver:    "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			coin:        sdk.NewCoin("del", sdk.NewInt(0)),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := MsgSendCoin{
			Sender:   tc.sender,
			Receiver: tc.receiver,
			Coin:     tc.coin,
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
		sends       []Send
		expectError bool
	}{
		{
			tag:    "valid send",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []Send{
				{Coin: sdk.NewCoin("del", sdk.NewInt(1)), Receiver: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"},
				{Coin: sdk.NewCoin("btc", sdk.NewInt(2)), Receiver: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f"},
			},
			expectError: false,
		},
		{
			tag:    "invalid sender",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			sends: []Send{
				{Coin: sdk.NewCoin("del", sdk.NewInt(1)), Receiver: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"},
				{Coin: sdk.NewCoin("btc", sdk.NewInt(2)), Receiver: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f"},
			},
			expectError: true,
		},
		{
			tag:    "invalid receiver",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []Send{
				{Coin: sdk.NewCoin("del", sdk.NewInt(1)), Receiver: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"},
				{Coin: sdk.NewCoin("btc", sdk.NewInt(2)), Receiver: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f0"},
			},
			expectError: true,
		},
		{
			tag:    "invalid receiver=sender",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []Send{
				{Coin: sdk.NewCoin("del", sdk.NewInt(1)), Receiver: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"},
				{Coin: sdk.NewCoin("btc", sdk.NewInt(2)), Receiver: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn"},
			},
			expectError: true,
		},
		{
			tag:    "invalid amount",
			sender: "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends: []Send{
				{Coin: sdk.NewCoin("del", sdk.NewInt(1)), Receiver: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"},
				{Coin: sdk.NewCoin("btc", sdk.NewInt(0)), Receiver: "dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f"},
			},
			expectError: true,
		},
		{
			tag:         "no sends",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sends:       []Send{},
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
			coinToBuy:     sdk.NewCoin("del", sdk.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:   false,
		},
		{
			tag:           "invalid sender",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinToBuy:     sdk.NewCoin("del", sdk.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid buy amount",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdk.NewInt(0)),
			maxCoinToSell: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:   true,
		},
		{
			tag:           "invalid sell amount",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdk.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("btc", sdk.NewInt(0)),
			expectError:   true,
		},
		{
			tag:           "same coin",
			sender:        "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToBuy:     sdk.NewCoin("del", sdk.NewInt(1)),
			maxCoinToSell: sdk.NewCoin("del", sdk.NewInt(1)),
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
			coinToSell:   sdk.NewCoin("del", sdk.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:  false,
		},
		{
			tag:          "invalid sender",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinToSell:   sdk.NewCoin("del", sdk.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "invalid sell amount",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdk.NewInt(0)),
			minCoinToBuy: sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:  true,
		},
		{
			tag:          "invalid buy amount",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdk.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("btc", sdk.NewInt(0)),
			expectError:  true,
		},
		{
			tag:          "same coin",
			sender:       "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinToSell:   sdk.NewCoin("del", sdk.NewInt(1)),
			minCoinToBuy: sdk.NewCoin("del", sdk.NewInt(1)),
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
		tag              string
		sender           string
		coinSymbolToSell string
		minCoinToBuy     sdk.Coin
		expectError      bool
	}{
		{
			tag:              "valid sell",
			sender:           "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinSymbolToSell: "del",
			minCoinToBuy:     sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:      false,
		},
		{
			tag:              "invalid sender",
			sender:           "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coinSymbolToSell: "del",
			minCoinToBuy:     sdk.NewCoin("btc", sdk.NewInt(1)),
			expectError:      true,
		},
		{
			tag:              "invalid buy amount",
			sender:           "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinSymbolToSell: "del",
			minCoinToBuy:     sdk.NewCoin("btc", sdk.NewInt(0)),
			expectError:      true,
		},
		{
			tag:              "same coin",
			sender:           "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coinSymbolToSell: "del",
			minCoinToBuy:     sdk.NewCoin("del", sdk.NewInt(1)),
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		msg := MsgSellAllCoin{
			Sender:           tc.sender,
			CoinSymbolToSell: tc.coinSymbolToSell,
			MinCoinToBuy:     tc.minCoinToBuy,
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
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: false,
		},
		{
			tag:         "invalid sender",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn0",
			coin:        sdk.NewCoin("del", sdk.NewInt(1)),
			expectError: true,
		},
		{
			tag:         "invalid burn amount",
			sender:      "dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			coin:        sdk.NewCoin("del", sdk.NewInt(0)),
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

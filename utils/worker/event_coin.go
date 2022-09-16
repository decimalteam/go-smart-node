package worker

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TODO: get from blockchain
const baseCoinSymbol = "del"

type EventEditCoin struct {
	Volume  sdkmath.Int `json:"volume"`
	Reserve sdkmath.Int `json:"reserve"`
}

type EventUpdateCoin struct {
	LimitVolume sdkmath.Int `json:"limitVolume"`
	Avatar      string      `json:"avatar"` // identity
}

type EventCreateCoin struct {
	Symbol      string      `json:"symbol"`
	Title       string      `json:"title"`
	Volume      sdkmath.Int `json:"volume"`
	Reserve     sdkmath.Int `json:"reserve"`
	CRR         uint64      `json:"crr"`
	LimitVolume sdkmath.Int `json:"limitVolume"`
	Creator     string      `json:"creator"`
	Avatar      string      `json:"avatar"` // identity
	// can get from transactions
	TxHash  string `json:"txHash"`
	BlockID int64  `json:"blockId"`
	// ? priceUSD
	// ? burn
}

// decimal.coin.v1.EventCreateCoin
func processEventCreateCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string sender = 1;
	   string symbol = 2;
	   string title = 3;
	   uint64 crr = 4;
	   string initial_volume = 5;
	   string initial_reserve = 6;
	   string limit_volume = 7;
	   string identity = 8;
	   string commission_create_coin = 9;
	*/
	var ecc EventCreateCoin
	var err error
	var ok bool
	var sender string
	var commission sdk.Coin
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			sender = string(attr.Value)
			ecc.Creator = sender
		case "symbol":
			ecc.Symbol = string(attr.Value)
		case "title":
			ecc.Title = string(attr.Value)
		case "identity":
			ecc.Avatar = string(attr.Value)
		case "crr":
			ecc.CRR, err = strconv.ParseUint(string(attr.Value), 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse crr '%s': %s", string(attr.Value), err.Error())
			}
		case "initial_volume":
			ecc.Volume, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse initial_volume '%s'", string(attr.Value))
			}
		case "initial_reserve":
			ecc.Reserve, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse initial_reserve '%s'", string(attr.Value))
			}
		case "limit_volume":
			ecc.LimitVolume, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse limit_volume '%s'", string(attr.Value))
			}
		case "commission_create_coin":
			commission, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse commission_create_coin '%s': %s", string(attr.Value), err.Error())
			}
		}
	}
	ecc.TxHash = txHash
	ecc.BlockID = blockId
	ea.addBalanceChange(sender, baseCoinSymbol, ecc.Reserve.Neg())
	ea.addBalanceChange(sender, ecc.Symbol, ecc.Volume)
	ea.addBalanceChange(sender, commission.Denom, commission.Amount.Neg())

	ea.CoinsCreates = append(ea.CoinsCreates, ecc)
	return nil
}

// decimal.coin.v1.EventUpdateCoin
func processEventUpdateCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string sender = 1;
	   string symbol = 2;
	   string limit_volume = 3;
	   string identity = 4;
	*/
	var ok bool
	var euc EventUpdateCoin
	var symbol string
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "symbol":
			symbol = string(attr.Value)
		case "identity":
			euc.Avatar = string(attr.Value)
		case "limit_volume":
			euc.LimitVolume, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse limit_volume '%s'", string(attr.Value))
			}
		}
	}
	ea.CoinUpdates[symbol] = euc
	return nil
}

// decimal.coin.v1.EventEditCoin
func processEventEditCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string symbol = 1;
	   string volume = 2;
	   string reserve = 3;
	*/
	var ok bool
	var eec EventEditCoin
	var symbol string
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "symbol":
			symbol = string(attr.Value)
		case "volume":
			eec.Volume, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse volume '%s'", string(attr.Value))
			}
		case "reserve":
			eec.Reserve, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse reserve '%s'", string(attr.Value))
			}
		}
	}
	ea.CoinEdits[symbol] = eec
	return nil
}

// decimal.coin.v1.EventSendCoin
func processEventSendCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
		string sender = 1;
		string receiver = 2;
		string coin = 3;
	*/
	var err error
	var sender, receiver string
	var coin sdk.Coin
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			sender = string(attr.Value)
		case "receiver":
			receiver = string(attr.Value)
		case "coin":
			coin, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
	}

	ea.addBalanceChange(sender, coin.Denom, coin.Amount.Neg())
	ea.addBalanceChange(receiver, coin.Denom, coin.Amount)
	return nil
}

// decimal.coin.v1.EventBuySellCoin
func processEventBuySellCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string sender = 1;
	   string coin_to_buy = 2;
	   string coin_to_sell = 3;
	   string amount_in_base_coin = 4;
	*/
	var err error
	var sender string
	var coinToBuy, coinToSell sdk.Coin
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			sender = string(attr.Value)
		}
		if string(attr.Key) == "coin_to_buy" {
			coinToBuy, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
		if string(attr.Key) == "coin_to_sell" {
			coinToSell, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
	}

	ea.addBalanceChange(sender, coinToBuy.Denom, coinToBuy.Amount)
	ea.addBalanceChange(sender, coinToSell.Denom, coinToSell.Amount.Neg())
	return nil
}

// decimal.coin.v1.EventBurnCoin
func processEventBurnCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string sender = 1;
	   string coin = 2;
	*/
	var err error
	var sender string
	var coinToBurn sdk.Coin
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			sender = string(attr.Value)
		}
		if string(attr.Key) == "coin" {
			coinToBurn, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
	}

	ea.addBalanceChange(sender, coinToBurn.Denom, coinToBurn.Amount.Neg())
	return nil
}

// decimal.coin.v1.EventRedeemCheck
func processEventRedeemCheck(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string sender = 1;
	   string issuer = 2;
	   string coin = 3;
	   string nonce = 4;
	   string due_block = 5;
	   string commission_redeem_check = 6;
	*/
	var err error
	var sender, issuer string
	var coin, commission sdk.Coin
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			sender = string(attr.Value)
		}
		if string(attr.Key) == "issuer" {
			issuer = string(attr.Value)
		}
		if string(attr.Key) == "coin" {
			coin, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
		if string(attr.Key) == "commission_redeem_check" {
			commission, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse commission_redeem_check '%s': %s", string(attr.Value), err.Error())
			}
		}
	}

	ea.addBalanceChange(sender, coin.Denom, coin.Amount)
	ea.addBalanceChange(issuer, coin.Denom, coin.Amount.Neg())
	ea.addBalanceChange(issuer, commission.Denom, commission.Amount.Neg())
	return nil
}

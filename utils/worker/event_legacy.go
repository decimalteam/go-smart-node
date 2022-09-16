package worker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type LegacyReturnNFT struct {
	OldAddress string `json:"old_address"`
	NewAddress string `json:"new_address"`
	Denom      string `json:"denom"`
	ID         string `json:"id"`
}

type LegacyReturnWallet struct {
	OldAddress string `json:"old_address"`
	NewAddress string `json:"new_address"`
	Wallet     string `json:"wallet"`
}

// decimal.legacy.v1.EventLegacyReturnCoin
func processEventLegacyReturnCoin(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string old_address = 2;
	   string new_address = 3;
	   string coins = 4;
	*/
	var err error
	var oldAddress, newAddress string
	var coins sdk.Coins
	for _, attr := range event.Attributes {
		if string(attr.Key) == "old_address" {
			oldAddress = string(attr.Value)
		}
		if string(attr.Key) == "new_address" {
			newAddress = string(attr.Value)
		}
		if string(attr.Key) == "coins" {
			coins, err = sdk.ParseCoinsNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coins '%s': %s", string(attr.Value), err.Error())
			}
		}
	}
	for _, coin := range coins {
		ea.addBalanceChange(newAddress, coin.Denom, coin.Amount)
	}
	ea.LegacyReown[oldAddress] = newAddress
	return nil

}

// decimal.legacy.v1.EventLegacyReturnNFT
func processEventLegacyReturnNFT(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string old_address = 2;
	   string new_address = 3;
	   string denom = 4;
	   string token_id = 5;
	*/
	var ret LegacyReturnNFT
	for _, attr := range event.Attributes {
		if string(attr.Key) == "old_address" {
			ret.OldAddress = string(attr.Value)
		}
		if string(attr.Key) == "new_address" {
			ret.NewAddress = string(attr.Value)
		}
		if string(attr.Key) == "denom" {
			ret.Denom = string(attr.Value)
		}
		if string(attr.Key) == "token_id" {
			ret.ID = string(attr.Value)
		}
	}
	ea.LegacyReturnNFT = append(ea.LegacyReturnNFT, ret)
	return nil

}

// decimal.legacy.v1.EventLegacyReturnWallet
func processEventLegacyReturnWallet(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	   string old_address = 2;
	   string new_address = 3;
	   string wallet = 4;
	*/
	var ret LegacyReturnWallet
	for _, attr := range event.Attributes {
		if string(attr.Key) == "old_address" {
			ret.OldAddress = string(attr.Value)
		}
		if string(attr.Key) == "new_address" {
			ret.NewAddress = string(attr.Value)
		}
		if string(attr.Key) == "wallet" {
			ret.Wallet = string(attr.Value)
		}
	}
	ea.LegacyReturnWallet = append(ea.LegacyReturnWallet, ret)
	return nil

}

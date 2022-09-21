package worker

import (
	"encoding/json"
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

// decimal.legacy.v1.EventReturnLegacyCoins
func processEventReturnLegacyCoins(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	  string legacy_owner = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string owner = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  repeated cosmos.base.v1beta1.Coin coins = 3 [
	    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins",
	    (gogoproto.nullable) = false
	  ];
	*/
	var err error
	var oldAddress, newAddress string
	var coins sdk.Coins
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "legacy_owner":
			oldAddress = string(attr.Value)
		case "owner":
			newAddress = string(attr.Value)
		case "coins":
			err = json.Unmarshal(attr.Value, &coins)
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

// decimal.legacy.v1.EventReturnLegacySubToken
func processEventReturnLegacySubToken(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	  string legacy_owner = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string owner = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string denom = 3;
	  string id = 4 [ (gogoproto.customname) = "ID" ];
	  repeated uint32 sub_token_ids = 5 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/
	var ret LegacyReturnNFT
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "legacy_owner":
			ret.OldAddress = string(attr.Value)
		case "owner":
			ret.NewAddress = string(attr.Value)
		case "denom":
			ret.Denom = string(attr.Value)
		case "id":
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

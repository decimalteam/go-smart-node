package worker

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LegacyReturnNFT struct {
	LegacyOwner string `json:"legacy_owner"`
	Owner       string `json:"owner"`
	Denom       string `json:"denom"`
	Creator     string `json:"creator"`
	ID          string `json:"id"`
}

type LegacyReturnWallet struct {
	LegacyOwner string `json:"legacy_owner"`
	Owner       string `json:"owner"`
	Wallet      string `json:"wallet"`
}

// decimal.legacy.v1.EventReturnLegacyCoins
func processEventReturnLegacyCoins(ea *EventAccumulator, event abci.Event, txHash string) error {
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
			err = json.Unmarshal([]byte(attr.Value), &coins)
			if err != nil {
				return fmt.Errorf("can't parse coins '%s': %s", string(attr.Value), err.Error())
			}
		}
	}
	//for _, coin := range coins {
	//ea.addBalanceChange(newAddress, coin.Denom, coin.Amount)
	//}
	ea.LegacyReown[oldAddress] = newAddress
	return nil

}

// decimal.legacy.v1.EventReturnLegacySubToken
func processEventReturnLegacySubToken(ea *EventAccumulator, event abci.Event, txHash string) error {
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
			ret.LegacyOwner = string(attr.Value)
		case "owner":
			ret.Owner = string(attr.Value)
		case "denom":
			ret.Denom = string(attr.Value)
		case "id":
			ret.ID = string(attr.Value)
		}
	}
	ea.LegacyReturnNFT = append(ea.LegacyReturnNFT, ret)
	return nil

}

// decimal.legacy.v1.EventReturnMultisigWallet
func processEventReturnMultisigWallet(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string legacy_owner = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string owner = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string wallet = 3 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	*/
	var ret LegacyReturnWallet
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "legacy_owner":
			ret.LegacyOwner = string(attr.Value)
		case "owner":
			ret.Owner = string(attr.Value)
		case "wallet":
			ret.Wallet = string(attr.Value)
		}
	}
	ea.LegacyReturnWallet = append(ea.LegacyReturnWallet, ret)
	return nil

}

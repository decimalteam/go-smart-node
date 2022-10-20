package worker

import (
	"encoding/json"
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// decimal-models
type MultisigCreateWallet struct {
	Address   string          `json:"address"`
	Threshold uint32          `json:"threshold"`
	Creator   string          `json:"creator"`
	Owners    []MultisigOwner `json:"owners"`
}

type MultisigOwner struct {
	Address  string `json:"address"`
	Multisig string `json:"multisig"`
	Weight   uint32 `json:"weight"`
}

type MultisigCreateTx struct {
	Sender      string    `json:"sender"`
	Wallet      string    `json:"wallet"`
	Receiver    string    `json:"receiver"`
	Transaction string    `json:"transaction"`
	Coins       sdk.Coins `json:"coins"`
	TxHash      string    `json:"txHash"`
}

type MultisigSignTx struct {
	Sender        string `json:"sender"`
	Wallet        string `json:"wallet"`
	Transaction   string `json:"transaction"`
	SignerWeight  uint32 `json:"signer_weight"`
	Confirmations uint32 `json:"confirmations"`
	Confirmed     bool   `json:"confirmed"`
}

// decimal.multisig.v1.EventCreateWallet
func processEventCreateWallet(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  repeated string owners = 3;
	  repeated uint32 weights = 4;
	  uint32 threshold = 5;
	*/
	e := MultisigCreateWallet{}
	var owners []string
	var weights []uint32
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Creator = string(attr.Value)
		case "wallet":
			e.Address = string(attr.Value)
		case "threshold":
			thr, err := strconv.ParseUint(string(attr.Value), 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse threshold '%s': %s", string(attr.Value), err.Error())
			}
			e.Threshold = uint32(thr)
		case "owners":
			err := json.Unmarshal(attr.Value, &owners)
			if err != nil {
				return fmt.Errorf("can't unmarshal owners: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "weights":
			err := json.Unmarshal(attr.Value, &weights)
			if err != nil {
				return fmt.Errorf("can't unmarshal weights: %s", err.Error())
			}
		}
	}
	for i, owner := range owners {
		e.Owners = append(e.Owners, MultisigOwner{
			Address:  owner,
			Multisig: e.Address,
			Weight:   weights[i],
		})
	}

	ea.MultisigCreateWallets = append(ea.MultisigCreateWallets, e)
	return nil
}

// decimal.multisig.v1.EventCreateTransaction
func processEventCreateTransaction(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string receiver = 3 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string coins = 4;
		string transaction = 5;
	*/
	e := MultisigCreateTx{}
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "wallet":
			e.Wallet = string(attr.Value)
		case "receiver":
			e.Receiver = string(attr.Value)
		case "coins":
			err := json.Unmarshal(attr.Value, &e.Coins)
			if err != nil {
				return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "transaction":
			e.Transaction = string(attr.Value)
		}
	}

	e.TxHash = txHash
	ea.MultisigCreateTxs = append(ea.MultisigCreateTxs, e)
	return nil
}

// decimal.multisig.v1.EventSignTransaction
func processEventSignTransaction(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string transaction = 3;
	  uint32 signer_weight = 4;
	  uint32 confirmations = 5;
	  bool confirmed = 6;
	*/
	e := MultisigSignTx{}
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "wallet":
			e.Wallet = string(attr.Value)
		case "transaction":
			e.Transaction = string(attr.Value)
		case "signer_weight":
			err := json.Unmarshal(attr.Value, &e.SignerWeight)
			if err != nil {
				return fmt.Errorf("can't unmarshal signer_weight: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "confirmations":
			err := json.Unmarshal(attr.Value, &e.Confirmations)
			if err != nil {
				return fmt.Errorf("can't unmarshal confirmations: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "confirmed":
			err := json.Unmarshal(attr.Value, &e.Confirmed)
			if err != nil {
				return fmt.Errorf("can't unmarshal confirmed: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}
	ea.MultisigSignTxs = append(ea.MultisigSignTxs, e)
	return nil
}

// decimal.multisig.v1.EventConfirmTransaction
func processEventConfirmTransaction(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string wallet = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string receiver = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string transaction = 3;
	  repeated cosmos.base.v1beta1.Coin coins = 4
	      [ (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins", (gogoproto.nullable) = false ];
	*/
	//var wallet, receiver string
	//var coins = sdk.NewCoins()
	//for _, attr := range event.Attributes {
	//	switch string(attr.Key) {
	//	case "wallet":
	//		wallet = string(attr.Value)
	//	case "receiver":
	//		receiver = string(attr.Value)
	//	case "coins":
	//		err := json.Unmarshal(attr.Value, &coins)
	//		if err != nil {
	//			return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
	//		}
	//	}
	//for _, coin := range coins {
	//	ea.addBalanceChange(wallet, coin.Denom, coin.Amount.Neg())
	//	ea.addBalanceChange(receiver, coin.Denom, coin.Amount)
	//}
	return nil
}

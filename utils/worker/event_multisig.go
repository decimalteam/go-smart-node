package worker

import (
	"encoding/json"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
	Sender   string    `json:"sender"`
	Wallet   string    `json:"wallet"`
	Receiver string    `json:"receiver"`
	TxHash   string    `json:"transaction"`
	Coins    sdk.Coins `json:"coins"`
}

type MultisigSignTx struct {
	Sender        string `json:"sender"`
	Wallet        string `json:"wallet"`
	TxHash        string `json:"transaction"`
	SignerWeight  uint32 `json:"signer_weight"`
	Confirmations uint32 `json:"confirmations"`
	Confirmed     bool   `json:"confirmed"`
}

// decimal.multisig.v1.EventCreateWallet
func processEventCreateWallet(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  repeated string owners = 3;
	  repeated uint32 weights = 4;
	  uint32 threshold = 5;
	*/
	mcw := MultisigCreateWallet{}
	var owners []string
	var weights []uint32
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			mcw.Creator = string(attr.Value)
		case "wallet":
			mcw.Address = string(attr.Value)
		case "threshold":
			thr, err := strconv.ParseUint(string(attr.Value), 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse threshold '%s': %s", string(attr.Value), err.Error())
			}
			mcw.Threshold = uint32(thr)
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
		mcw.Owners = append(mcw.Owners, MultisigOwner{
			Address:  owner,
			Multisig: mcw.Address,
			Weight:   weights[i],
		})
	}

	ea.MultisigCreateWallets = append(ea.MultisigCreateWallets, mcw)
	return nil
}

// decimal.multisig.v1.EventCreateTransaction
func processEventCreateTransaction(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
		string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string receiver = 3 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string coins = 4;
		string transaction = 5;
	*/
	mct := MultisigCreateTx{}
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			mct.Sender = string(attr.Value)
		case "wallet":
			mct.Wallet = string(attr.Value)
		case "receiver":
			mct.Receiver = string(attr.Value)
		case "coins":
			err := json.Unmarshal(attr.Value, &mct.Coins)
			if err != nil {
				return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "transaction":
			mct.TxHash = string(attr.Value)
		}
	}
	ea.MultisigCreateTxs = append(ea.MultisigCreateTxs, mct)
	return nil
}

// decimal.multisig.v1.EventSignTransaction
func processEventSignTransaction(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string wallet = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string transaction = 3;
	  uint32 signer_weight = 4;
	  uint32 confirmations = 5;
	  bool confirmed = 6;
	*/
	mst := MultisigSignTx{}
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			mst.Sender = string(attr.Value)
		case "wallet":
			mst.Wallet = string(attr.Value)
		case "transaction":
			mst.TxHash = string(attr.Value)
		case "signer_weight":
			err := json.Unmarshal(attr.Value, &mst.SignerWeight)
			if err != nil {
				return fmt.Errorf("can't unmarshal signer_weight: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "confirmations":
			err := json.Unmarshal(attr.Value, &mst.Confirmations)
			if err != nil {
				return fmt.Errorf("can't unmarshal confirmations: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "confirmed":
			err := json.Unmarshal(attr.Value, &mst.Confirmed)
			if err != nil {
				return fmt.Errorf("can't unmarshal confirmed: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}
	ea.MultisigSignTxs = append(ea.MultisigSignTxs, mst)
	return nil
}

// decimal.multisig.v1.EventConfirmTransaction
func processEventConfirmTransaction(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
	  string wallet = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string receiver = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string transaction = 3;
	  repeated cosmos.base.v1beta1.Coin coins = 4
	      [ (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins", (gogoproto.nullable) = false ];
	*/
	var wallet, receiver string
	var coins = sdk.NewCoins()
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "wallet":
			wallet = string(attr.Value)
		case "receiver":
			receiver = string(attr.Value)
		case "coins":
			err := json.Unmarshal(attr.Value, &coins)
			if err != nil {
				return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}
	for _, coin := range coins {
		ea.addBalanceChange(wallet, coin.Denom, coin.Amount.Neg())
		ea.addBalanceChange(receiver, coin.Denom, coin.Amount)
	}
	return nil
}

package worker

import (
	"encoding/json"
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"
)

// decimal-models
type MultisigCreateWallet struct {
	Address   string          `json:"address"`
	Threshold uint32          `json:"threshold"`
	Creator   string          `json:"creator"`
	Owners    []MultisigOwner `json:""`
}

type MultisigOwner struct {
	Address  string `json:"address"`
	Multisig string `json:"multisig"`
	Weight   uint32 `json:"weight"`
}

// TODO: Transaction create+sign

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
	return nil
}

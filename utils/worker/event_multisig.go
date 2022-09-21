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
	Threshold uint64          `json:"threshold"`
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
		string sender = 1;
		string wallet = 2;
		repeated string owners = 3;
		repeated uint64 weights = 4;
		uint64 threshold = 5;
	*/
	mcw := MultisigCreateWallet{}
	var owners []string
	var weights []uint32
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			mcw.Creator = string(attr.Value)
		}
		if string(attr.Key) == "wallet" {
			mcw.Address = string(attr.Value)
		}
		if string(attr.Key) == "threshold" {
			thr, err := strconv.ParseUint(string(attr.Value), 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse threshold '%s': %s", string(attr.Value), err.Error())
			}
			mcw.Threshold = thr
		}
		if string(attr.Key) == "owners" {
			err := json.Unmarshal(attr.Value, &owners)
			if err != nil {
				return fmt.Errorf("can't unmarshal owners: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
		if string(attr.Key) == "weights" {
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

func processEventCreateTransaction(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	return nil
}

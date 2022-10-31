package worker

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type EventDelegate struct {
	Delegator string
	Validator string
	Stake     Stake
}

type EventUndelegateComplete struct {
	Delegator string
	Validator string
	Stake     Stake
}

type EventRedelegateComplete struct {
	Delegator    string
	ValidatorSrc string
	ValidatorDst string
	Stake        Stake
}

type Stake struct {
	Stake sdk.Coin
}

func processEventDelegate(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string delegator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string validator = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  Stake stake = 3 [ (gogoproto.nullable) = false ];
	  string amount_base = 4 [
	    (cosmos_proto.scalar) = "cosmos.Int",
	    (gogoproto.customtype) = "cosmossdk.io/math.Int",
	    (gogoproto.nullable) = false
	  ];
	*/

	var (
		e EventDelegate
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "delegator":
			e.Delegator = string(attr.Value)
		case "validator":
			e.Validator = string(attr.Value)
		case "stake":
			var stake Stake

			err := json.Unmarshal(attr.Value, &stake)
			if err != nil {
				panic(err)
			}
			e.Stake = stake
		}
	}

	if _, ok := pool[e.Delegator]; !ok {
		ea.addBalanceChange(e.Delegator, e.Stake.Stake.Denom, e.Stake.Stake.Amount.Neg())
	}

	return nil
}

func processEventUndelegateComplete(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string delegator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string validator = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		Stake stake = 3 [ (gogoproto.nullable) = false ];
	*/

	var (
		e EventUndelegateComplete
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "delegator":
			e.Delegator = string(attr.Value)
		case "validator":
			e.Validator = string(attr.Value)
		case "stake":
			var stake Stake

			err := json.Unmarshal(attr.Value, &stake)
			if err != nil {
				panic(err)
			}
			e.Stake = stake
		}
	}

	if _, ok := pool[e.Delegator]; !ok {
		ea.addBalanceChange(e.Delegator, e.Stake.Stake.Denom, e.Stake.Stake.Amount)
	}

	return nil
}

func processEventRedelegateComplete(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string delegator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		  string validator_src = 2 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		  string validator_dst = 3 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		  Stake stake = 4 [ (gogoproto.nullable) = false ];
	*/

	var (
		e EventRedelegateComplete
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "delegator":
			e.Delegator = string(attr.Value)
		case "validator_src":
			e.ValidatorSrc = string(attr.Value)
		case "validator_dst":
			e.ValidatorDst = string(attr.Value)
		case "stake":
			var stake Stake

			err := json.Unmarshal(attr.Value, &stake)
			if err != nil {
				panic(err)
			}
			e.Stake = stake
		}
	}

	//if _, ok := pool[e.Delegator]; !ok {
	//	ea.addBalanceChange(e.Delegator, e.Stake.Stake.Denom, e.Stake.Stake.Amount)
	//}

	return nil
}

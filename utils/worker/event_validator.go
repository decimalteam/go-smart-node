package worker

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
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

type EventUpdateCoinsStaked struct {
	denom  string
	amount sdkmath.Int
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

func processEventUpdateCoinsStaked(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string denom = 1;
	  string total_amount = 2 [
	    (cosmos_proto.scalar) = "cosmos.Int",
	    (gogoproto.customtype) = "cosmossdk.io/math.Int",
	    (gogoproto.nullable) = false
	  ];
	*/

	var (
		e  EventUpdateCoinsStaked
		ok bool
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "denom":
			e.denom = string(attr.Value)
		case "total_amount":
			e.amount, ok = sdk.NewIntFromString(string(attr.Value))
			if !ok {
				return fmt.Errorf("can't parse total_amount '%s'", string(attr.Value))
			}
		}
	}

	ea.addCoinsStaked(e)
	return nil
}

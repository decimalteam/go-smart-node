package worker

import (
	"encoding/json"
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EventUpdateCoinPrices struct {
	Oracle     string
	CoinPrices struct {
		Denom     string
		Quote     string
		Price     string
		UpdatedAt time.Time
	}
}

func processEventUpdatePrices(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string oracle = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		repeated CoinPrice prices = 2 [ (gogoproto.nullable) = false ];
	*/

	// TODO this event need handle?

	return nil
}

// decimal.fee.v1.EventPayCommission
func processEventPayCommission(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string sender = 1;
		string coin = 2;
	*/
	var err error
	var payer string
	var coins sdk.Coins
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			payer = string(attr.Value)
		}
		if string(attr.Key) == "coin" {
			err = json.Unmarshal(attr.Value, &coins)
			if err != nil {
				return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}

	for _, coin := range coins {
		ea.addBalanceChange(payer, coin.Denom, coin.Amount.Neg())
	}
	return nil

}

package worker

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

type EventPayCommission struct {
	Payer  string    `json:"payer"`
	Coins  sdk.Coins `json:"coins"`
	Burnt  sdk.Coins `json:"burnt"`
	TxHash string    `json:"txHash"`
}

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
	  string payer = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  repeated cosmos.base.v1beta1.Coin coins = 2
	  [ (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins", (gogoproto.nullable) = false ];
	*/
	var (
		err   error
		coins sdk.Coins
		burnt sdk.Coins
		e     EventPayCommission
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "payer":
			e.Payer = string(attr.Value)
		case "coins":
			err = json.Unmarshal(attr.Value, &coins)
			if err != nil {
				return fmt.Errorf("can't unmarshal coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		case "burnt":
			err = json.Unmarshal(attr.Value, &burnt)
			if err != nil {
				return fmt.Errorf("can't unmarshal burnt coins: %s, value: '%s'", err.Error(), string(attr.Value))
			}
		}
	}

	e.Coins = coins
	e.TxHash = txHash
	e.Burnt = burnt
	ea.PayCommission = append(ea.PayCommission, e)
	return nil

}

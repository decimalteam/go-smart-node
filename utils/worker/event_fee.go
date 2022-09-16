package worker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// decimal.fee.v1.EventPayCommission
func processEventPayCommission(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
		string sender = 1;
		string coin = 2;
	*/
	var err error
	var sender string
	var coin sdk.Coin
	for _, attr := range event.Attributes {
		if string(attr.Key) == "sender" {
			sender = string(attr.Value)
		}
		if string(attr.Key) == "coin" {
			coin, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse coin '%s': %s", string(attr.Value), err.Error())
			}
		}
	}

	ea.addBalanceChange(sender, coin.Denom, coin.Amount.Neg())
	return nil

}

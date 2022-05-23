package worker

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParseTask struct {
	height int64
	txNum  int
}

func (w *Worker) parseTxFromStd(tx sdk.Tx) ParsedTx {
	var parsedTx ParsedTx
	if tx == nil {
		return parsedTx
	}
	for _, rawMsg := range tx.GetMsgs() {
		var msg TxMsg
		msg.Type = sdk.MsgTypeURL(rawMsg)
		msg.Params = rawMsg
		for _, signer := range rawMsg.GetSigners() {
			msg.From = append(msg.From, signer.String())
		}
		parsedTx.Msgs = append(parsedTx.Msgs, msg)
	}
	parsedTx.Fee.Gas = tx.(sdk.FeeTx).GetGas()
	parsedTx.Fee.Amount = tx.(sdk.FeeTx).GetFee()
	parsedTx.Memo = tx.(sdk.TxWithMemo).GetMemo()
	return parsedTx
}

func (w *Worker) parseEvents(events []abci.Event) []Event {
	var newEvents []Event
	for _, ev := range events {
		newEvent := Event{
			Type:       ev.Type,
			Attributes: []Attribute{},
		}
		for _, attr := range ev.Attributes {
			newEvent.Attributes = append(newEvent.Attributes, Attribute{
				Key:   string(attr.Key),
				Value: string(attr.Value),
			})
		}
		newEvents = append(newEvents, newEvent)
	}
	return newEvents
}

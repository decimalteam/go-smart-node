package worker

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParseTask struct {
	height int64
	txNum  int
}

func (w *Worker) parseTxInfo(tx sdk.Tx) (txInfo TxInfo) {
	if tx == nil {
		return
	}
	for _, rawMsg := range tx.GetMsgs() {
		var msg TxMsg
		msg.Type = sdk.MsgTypeURL(rawMsg)
		msg.Params = string(w.cdc.Marshaler.MustMarshalJSON(rawMsg))
		for _, signer := range rawMsg.GetSigners() {
			msg.From = append(msg.From, signer.String())
		}
		txInfo.Msgs = append(txInfo.Msgs, msg)
	}
	txInfo.Fee.Gas = tx.(sdk.FeeTx).GetGas()
	txInfo.Fee.Amount = tx.(sdk.FeeTx).GetFee()
	txInfo.Memo = tx.(sdk.TxWithMemo).GetMemo()
	return
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

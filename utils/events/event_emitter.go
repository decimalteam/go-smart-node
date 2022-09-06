package events

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
)

// THIS IS WORKAROUND FORK OF Cosmos SDK TypedEventToEvent

// TypedEventToEvent takes typed event and converts to Event object
func TypedEventToEvent(tev proto.Message) (sdk.Event, error) {
	evtType := proto.MessageName(tev)
	evtJSON, err := codec.ProtoMarshalJSON(tev, nil)
	if err != nil {
		return sdk.Event{}, err
	}

	var attrMap map[string]json.RawMessage
	err = json.Unmarshal(evtJSON, &attrMap)
	if err != nil {
		return sdk.Event{}, err
	}

	attrs := make([]abci.EventAttribute, 0, len(attrMap))
	for k, v := range attrMap {
		// dumb workaround for strings
		value := v
		if v[0] == '"' && v[len(v)-1] == '"' {
			value = v[1 : len(v)-1]
		}
		//
		attrs = append(attrs, abci.EventAttribute{
			Key:   []byte(k),
			Value: value,
		})
	}

	return sdk.Event{
		Type:       evtType,
		Attributes: attrs,
	}, nil
}

func EmitTypedEvent(ctx sdk.Context, tev proto.Message) error {
	ev, err := TypedEventToEvent(tev)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(ev)
	return nil
}

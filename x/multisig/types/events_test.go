package types

import (
	"encoding/json"
	fmt "fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestEventMarshaling(t *testing.T) {
	tev := &EventCreateWallet{
		Sender:    "aaa",
		Wallet:    "bbb",
		Owners:    []string{"c", "d", "e"},
		Weights:   []uint64{1, 2, 3},
		Threshold: 40,
	}
	ev, err := sdk.TypedEventToEvent(tev)
	require.NoError(t, err)
	ev2, err := TypedEventToEvent(tev)
	require.NoError(t, err)

	for _, at := range ev.Attributes {
		fmt.Printf("%s = %s\n", string(at.Key), string(at.Value))
	}
	for _, at := range ev2.Attributes {
		fmt.Printf("%s = %s\n", string(at.Key), string(at.Value))
	}
	b, _ := json.MarshalIndent(ev, " ", " ")
	fmt.Printf("%s\n", string(b))
}

// TypedEventToEvent takes typed event and converts to Event object
func TypedEventToEvent(tev proto.Message) (sdk.Event, error) {
	evtType := proto.MessageName(tev)
	evtJSON, err := codec.ProtoMarshalJSON(tev, nil)
	if err != nil {
		return sdk.Event{}, err
	}

	// workaround for escaped quotes in string values
	var attrMapOriginal map[string]interface{}
	var attrMapRaw map[string]json.RawMessage
	err = json.Unmarshal(evtJSON, &attrMapOriginal)
	if err != nil {
		return sdk.Event{}, err
	}
	err = json.Unmarshal(evtJSON, &attrMapRaw)
	if err != nil {
		return sdk.Event{}, err
	}

	attrs := make([]abci.EventAttribute, 0, len(attrMapOriginal))
	for k, v := range attrMapOriginal {
		switch vv := v.(type) {
		case string:
			attrs = append(attrs, abci.EventAttribute{
				Key:   []byte(k),
				Value: []byte(vv),
			})
		default:
			attrs = append(attrs, abci.EventAttribute{
				Key:   []byte(k),
				Value: attrMapRaw[k],
			})
		}
	}

	return sdk.Event{
		Type:       evtType,
		Attributes: attrs,
	}, nil
}

package events

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	multisigTypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func TestCompare(t *testing.T) {
	var testCases = []struct {
		tev    proto.Message
		expect map[string]string
	}{
		{
			&multisigTypes.EventCreateWallet{
				Sender:    "",
				Wallet:    "b",
				Owners:    []string{"c", "d"},
				Weights:   []uint64{1, 2},
				Threshold: 3,
			},
			map[string]string{
				"sender":    "",
				"wallet":    "b",
				"owners":    "[\"c\",\"d\"]",
				"weights":   "[\"1\",\"2\"]",
				"threshold": "3",
			},
		},
	}

	for _, tc := range testCases {
		ev, err := TypedEventToEvent(tc.tev)
		require.NoError(t, err)
		compareResult(t, ev, tc.expect)
	}
}

func compareResult(t *testing.T, ev sdk.Event, expectValue map[string]string) {
	for _, att := range ev.Attributes {
		k := string(att.Key)
		v := string(att.Value)
		require.Equal(t, expectValue[k], v)
	}
}
package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLegacyReturnValidation(t *testing.T) {
	msg1 := NewMsgReturnLegacy([]byte{0}, []byte{0})
	require.Error(t, msg1.ValidateBasic())

	msg1 = NewMsgReturnLegacy([]byte{0}, []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		1, 2, 3})
	require.NoError(t, msg1.ValidateBasic())
}

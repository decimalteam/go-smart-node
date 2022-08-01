package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMnemonic(t *testing.T) {
	// words forth and back
	m, err := NewMnemonic("")
	require.NoError(t, err, "generate mnemonic")

	m2, err := NewMnemonicFromWords(m.Words(), "")
	require.NoError(t, err, "generate mnemonic2")

	require.NoError(t, err, "mnemonic from existing words")
	require.Equal(t, m.Entropy(), m2.Entropy(), "entropy must be same")
	require.Equal(t, m.Seed(), m2.Seed(), "seed must be same")
}

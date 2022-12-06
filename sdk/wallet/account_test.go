package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKnownAccount(t *testing.T) {
	var testData = []struct {
		mnemonic    string
		password    string
		address     string
		expectError bool
	}{
		// dscd --keyring-backend test keys add 111 --dry-run -i
		{
			"gasp history river forget aware wide dance velvet weather rain rail dry cliff assault coach jelly choose spirit shoulder isolate kidney outer trust message",
			"",
			"d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv",
			false,
		},
		{
			"section jeans evoke hockey result spell dish zero merge actress pink resource loan afford fitness install purity duck cannon ugly session stereo pattern spawn",
			"",
			"d018c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg4x0y6k",
			false,
		},
		{
			"citizen marine borrow just apology mistake trumpet border sauce drip smile current excuse sing shove puppy dial ticket margin fabric afraid identify rookie elite",
			"",
			"d0164ea54aqgsmp7dp6wzs0y8n6vjehudnkvlsqzx",
			false,
		},
		{
			"matter sketch program direct property attend humble any car develop useless mask like elevator garbage protect obvious boring vessel obscure wink raven fog flip",
			"123456",
			"d01d6n3s60lsp3cn9ddvtk5ctsfnmag0ceamlptk0",
			false,
		},
		{
			// same mnemonic, but other password
			"matter sketch program direct property attend humble any car develop useless mask like elevator garbage protect obvious boring vessel obscure wink raven fog flip",
			"12345",
			"d01cmku6x7jlpf4utdkpc63dfp508mhu28kxrhmq0",
			false,
		},
		{
			// invalid mnemonic
			"sketch program direct property attend humble any car develop useless mask like elevator garbage protect obvious boring vessel obscure wink raven fog flip",
			"",
			"d01cmku6x7jlpf4utdkpc63dfp508mhu28kxrhmq0",
			true,
		},
	}

	for _, td := range testData {
		acc, err := NewAccountFromMnemonicWords(td.mnemonic, td.password)
		if td.expectError {
			require.Error(t, err, "no error for '%s'", td.address)
		} else {
			require.NoError(t, err, "inexpected error for mnemonic for acc '%s'", td.address)
			require.Equal(t, td.address, acc.Address(), "address value for acc '%s'", td.address)
		}
	}
}

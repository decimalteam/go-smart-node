package redenom_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
)

// TestContractCenterFor pins the chain-id -> DecimalContractCenter resolution that the
// redenomination EVM scaling depends on. The addresses are the deployed-dsc_{main,test,dev}
// records from the decimal-smart-contracts repo.
//
// The resolver MUST key on the network family (EIP-155 prefix), the SAME predicate that
// gates the upgrade (helpers.IsMainnet/IsTestnet/IsDevnet) — NOT the exact chain-id. The
// live testnet runs decimal_202020-221213, but the old exact-match map was keyed to
// decimal_202020-1, so scaleEVM hard-failed and halted the upgrade on node-03.
func TestContractCenterFor(t *testing.T) {
	var (
		mainnetCC = common.HexToAddress("0xc108715A06f76CAA96fa2c943Ebf05159c29A87D")
		testnetCC = common.HexToAddress("0xbC96b61F137F28F0Da47Cc4Ef06e5f984B565A2E")
		devnetCC  = common.HexToAddress("0x481487AEafc60512233a08Da240FC7AF99c0f696")
	)

	cases := []struct {
		name    string
		chainID string
		wantCC  common.Address
		wantOK  bool
	}{
		// The incident: the live testnet revision must resolve.
		{"live testnet (incident)", "decimal_202020-221213", testnetCC, true},
		// Back-compat: the old map key must keep resolving too.
		{"stale testnet key", "decimal_202020-1", testnetCC, true},
		// Mainnet must resolve for ANY revision suffix, or the mainnet upgrade hard-fails
		// the same way testnet did (the upgrade is gated by the IsMainnet prefix).
		{"mainnet rev 1", "decimal_75-1", mainnetCC, true},
		{"mainnet rev bump", "decimal_75-2", mainnetCC, true},
		// Devnet must NOT be mistaken for testnet despite the shared "decimal_2020..." start.
		{"devnet", "decimal_20202020-1", devnetCC, true},
		// Unknown chains (forks, private nets, dry-run fixtures) must NOT resolve, so the
		// handler still hard-fails rather than scaling against the wrong registry.
		{"unknown cosmos chain", "cosmoshub-4", common.Address{}, false},
		{"unknown redenom fixture", "redenom-test_1-1", common.Address{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotCC, gotOK := redenom.ContractCenterFor(tc.chainID)
			require.Equal(t, tc.wantOK, gotOK, "ok mismatch for chain-id %q", tc.chainID)
			require.Equal(t, tc.wantCC, gotCC, "address mismatch for chain-id %q", tc.chainID)
		})
	}
}

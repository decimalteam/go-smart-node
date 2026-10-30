package keeper

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationv7"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestReliabilityValidationAfterLegacyNormalization(t *testing.T) {
	current, err := delegation.DelegationMetaData.GetAbi()
	require.NoError(t, err)
	legacy, err := delegationv7.DelegationMetaData.GetAbi()
	require.NoError(t, err)
	for name, raw := range map[string][]string{
		"StakeUpdated": logStakeUpdated, "StakeHolded": logStakeHolded,
		"StakeReset": logStakeReset, "WithdrawRequest": logWithdrawRequest,
		"TransferRequest": logTransferRequest,
	} {
		t.Run(name, func(t *testing.T) {
			input := mainnetLog(t, raw[0], raw[1])
			require.Error(t, validateCriticalEventEnvelope(current, input))
			normalized, err := normalizeDelegationLogs([]*ethtypes.Log{input}, mainnetDelegation, current, legacy)
			require.NoError(t, err)
			require.Len(t, normalized, 1)
			require.NoError(t, validateCriticalEventEnvelope(current, normalized[0]))
			broken := *normalized[0]
			broken.Data = broken.Data[:len(broken.Data)-1]
			require.Error(t, validateCriticalEventEnvelope(current, &broken))
		})
	}
}

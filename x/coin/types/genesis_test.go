package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		description  string
		genesisState *types.GenesisState
		valid        bool
	}{
		{
			description:  "default is valid",
			genesisState: types.DefaultGenesis(),
			valid:        true,
		},
		{
			description:  "valid genesis state",
			genesisState: &types.GenesisState{},
			valid:        true,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.genesisState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

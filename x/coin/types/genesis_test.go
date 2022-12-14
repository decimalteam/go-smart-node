package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

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
			genesisState: types.DefaultGenesisState(),
			valid:        true,
		},
		{
			description: "valid genesis state",
			genesisState: &types.GenesisState{
				Params: types.DefaultParams(),
				Coins: []types.Coin{
					{
						Denom:   "sco",
						Title:   "somecoin",
						Reserve: sdkTypes.NewIntFromUint64(100),
					},
				},
			},
			valid: true,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.genesisState.Validate()
			if tc.valid {
				require.NoError(t, err, tc.description)
			} else {
				require.Error(t, err, tc.description)
			}
		})
	}
}

package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
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
						Title:   "somecoin",
						Symbol:  "sco",
						Reserve: sdkTypes.NewIntFromUint64(100),
					},
				},
				LegacyBalances: []types.LegacyBalance{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sco",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
					},
				},
			},
			valid: true,
		},
		{
			description: "invalid coin in legacy",
			genesisState: &types.GenesisState{
				Params: types.DefaultParams(),
				Coins: []types.Coin{
					{
						Title:   "somecoin",
						Symbol:  "sco",
						Reserve: sdkTypes.NewIntFromUint64(100),
					},
				},
				LegacyBalances: []types.LegacyBalance{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid balance in legacy",
			genesisState: &types.GenesisState{
				Params: types.DefaultParams(),
				Coins: []types.Coin{
					{
						Title:   "somecoin",
						Symbol:  "sco",
						Reserve: sdkTypes.NewIntFromUint64(100),
					},
				},
				LegacyBalances: []types.LegacyBalance{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sco",
								Amount: sdkTypes.NewInt(-1),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid address in legacy",
			genesisState: &types.GenesisState{
				Params: types.DefaultParams(),
				Coins: []types.Coin{
					{
						Title:   "somecoin",
						Symbol:  "sco",
						Reserve: sdkTypes.NewIntFromUint64(100),
					},
				},
				LegacyBalances: []types.LegacyBalance{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depi",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sco",
								Amount: sdkTypes.NewInt(1),
							},
						},
					},
				},
			},
			valid: false,
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

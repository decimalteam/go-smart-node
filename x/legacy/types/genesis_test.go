package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	_config := sdkTypes.GetConfig()
	_config.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)

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
				Records: []types.Record{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						NFTs:    []string{"a1", "a2"},
						Wallets: []string{"d0108c4p0j7wqsawejfuuv43hj7nhyp36gttdx0g0", "d010fx59x9ytvf249axryvw0uh3eunwvgyf9ux8g7"},
					},
				},
			},
			valid: true,
		},
		{
			description: "invalid balance",
			genesisState: &types.GenesisState{
				Records: []types.Record{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewInt(-1),
							},
						},
						NFTs:    []string{"a1", "a2"},
						Wallets: []string{"d0108c4p0j7wqsawejfuuv43hj7nhyp36gttdx0g0", "d010fx59x9ytvf249axryvw0uh3eunwvgyf9ux8g7"},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid bech32 address",
			genesisState: &types.GenesisState{
				Records: []types.Record{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depi",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						NFTs:    []string{"a1", "a2"},
						Wallets: []string{"d0108c4p0j7wqsawejfuuv43hj7nhyp36gttdx0g0", "d010fx59x9ytvf249axryvw0uh3eunwvgyf9ux8g7"},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid wallet address",
			genesisState: &types.GenesisState{
				Records: []types.Record{
					{
						LegacyAddress: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						NFTs:    []string{"a1", "a2"},
						Wallets: []string{"d0108c4p0j7wqsawejfuuv43hj7nhyp36gt0296r0"},
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

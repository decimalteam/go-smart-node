package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
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
				LegacyRecords: []types.LegacyRecord{
					{
						Address: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						Nfts:    []types.NFTRecord{{Denom: "a", Id: "a1"}, {Denom: "b", Id: "b1"}},
						Wallets: []string{"dx108c4p0j7wqsawejfuuv43hj7nhyp36gt0296rs", "dx10fx59x9ytvf249axryvw0uh3eunwvgyfpm9jrp"},
					},
				},
			},
			valid: true,
		},
		{
			description: "invalid balance",
			genesisState: &types.GenesisState{
				LegacyRecords: []types.LegacyRecord{
					{
						Address: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewInt(-1),
							},
						},
						Nfts:    []types.NFTRecord{{Denom: "a", Id: "a1"}, {Denom: "b", Id: "b1"}},
						Wallets: []string{"dx108c4p0j7wqsawejfuuv43hj7nhyp36gt0296rs", "dx10fx59x9ytvf249axryvw0uh3eunwvgyfpm9jrp"},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid bech32 address",
			genesisState: &types.GenesisState{
				LegacyRecords: []types.LegacyRecord{
					{
						Address: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depi",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						Nfts:    []types.NFTRecord{{Denom: "a", Id: "a1"}, {Denom: "b", Id: "b1"}},
						Wallets: []string{"dx108c4p0j7wqsawejfuuv43hj7nhyp36gt0296rs", "dx10fx59x9ytvf249axryvw0uh3eunwvgyfpm9jrp"},
					},
				},
			},
			valid: false,
		},
		{
			description: "invalid wallet address",
			genesisState: &types.GenesisState{
				LegacyRecords: []types.LegacyRecord{
					{
						Address: "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
						Coins: sdkTypes.Coins{
							{
								Denom:  "sca",
								Amount: sdkTypes.NewIntFromUint64(100),
							},
						},
						Nfts:    []types.NFTRecord{{Denom: "a", Id: "a1"}, {Denom: "b", Id: "b1"}},
						Wallets: []string{"dx108c4p0j7wqsawejfuuv43hj7nhyp36gt0296r0"},
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

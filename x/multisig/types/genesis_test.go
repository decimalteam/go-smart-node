package types_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)

	var testCases = []struct {
		description  string
		genesisState *types.GenesisState
		valid        bool
	}{
		{
			"valid genesis",
			types.DefaultGenesisState(),
			true,
		},
		{
			"wallet with legacy address",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 2}).String(),
							"dx14elhyzmq95f98wrkvujtsr5cyudffp6q2hfkhs",
						},
						Weights:   []uint32{1, 1, 1},
						Threshold: 2,
					},
				},
			},
			true,
		},
		{
			"double wallet",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 2}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 2,
					},
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 2}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 2,
					},
				},
			},
			false,
		},
		{
			"invalid wallet address",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String() + "1",
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 2}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 2,
					},
				},
			},
			false,
		},
		{
			"double owner",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 2,
					},
				},
			},
			false,
		},
		{
			"empty owner",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							"",
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 2,
					},
				},
			},
			false,
		},
		{
			"threshold over sum of weights",
			&types.GenesisState{
				Wallets: []types.Wallet{
					{
						Address: sdk.AccAddress([]byte{1, 2, 3}).String(),
						Owners: []string{
							sdk.AccAddress([]byte{1, 2, 3, 1}).String(),
							sdk.AccAddress([]byte{1, 2, 3, 2}).String(),
						},
						Weights:   []uint32{1, 1},
						Threshold: 3,
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.genesisState.Validate()
		if tc.valid {
			require.NoError(t, err, tc.description)
		} else {
			require.Error(t, err, tc.description)
		}
	}
}

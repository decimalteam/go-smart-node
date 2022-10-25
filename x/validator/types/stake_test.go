package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestStakeEqual(t *testing.T) {
	var testCases = []struct {
		name     string
		stake1   Stake
		stake2   Stake
		expEqual bool
	}{
		{"equal coin", NewStakeCoin(sdk.NewInt64Coin("aaa", 1)), NewStakeCoin(sdk.NewInt64Coin("aaa", 1)), true},
		{"equal nft", NewStakeNFT("abc", []uint32{1, 2}, sdk.NewInt64Coin("aaa", 1)), NewStakeNFT("abc", []uint32{2, 1}, sdk.NewInt64Coin("aaa", 1)), true},
		{"diff coin 1", NewStakeCoin(sdk.NewInt64Coin("aaa", 2)), NewStakeCoin(sdk.NewInt64Coin("aaa", 1)), false},
		{"diff coin 2", NewStakeCoin(sdk.NewInt64Coin("aab", 1)), NewStakeCoin(sdk.NewInt64Coin("aaa", 1)), false},
		{"diff nft 1", NewStakeNFT("abc", []uint32{1, 2}, sdk.NewInt64Coin("aaa", 1)), NewStakeNFT("abcd", []uint32{2, 1}, sdk.NewInt64Coin("aaa", 1)), false},
		{"diff nft 2", NewStakeNFT("abc", []uint32{1, 2}, sdk.NewInt64Coin("aaa", 1)), NewStakeNFT("abc", []uint32{1}, sdk.NewInt64Coin("aaa", 1)), false},
		{"diff nft 3", NewStakeNFT("abc", []uint32{1, 2}, sdk.NewInt64Coin("aaa", 1)), NewStakeNFT("abc", []uint32{2, 3}, sdk.NewInt64Coin("aaa", 1)), false},
		{"diff type", Stake{Type: StakeType_Coin}, Stake{Type: StakeType_NFT}, false},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expEqual, tc.stake1.Equal(&tc.stake2), tc.name)
	}
}

func TestSetOperations(t *testing.T) {
	require.Equal(t, true, SetHasSubset([]uint32{1, 2, 3}, []uint32{}), "empty subset")
	require.Equal(t, true, SetHasSubset([]uint32{1, 2, 3}, []uint32{2, 3}), "valid subset")
	require.Equal(t, false, SetHasSubset([]uint32{1, 2, 3}, []uint32{3, 4}), "small invalid subset")
	require.Equal(t, false, SetHasSubset([]uint32{1, 2, 3}, []uint32{1, 2, 3, 4}), "small big subset")

	require.Equal(t, []uint32{2}, SetSubstract([]uint32{1, 2, 3}, []uint32{1, 3}), "valid substract")
	require.Equal(t, []uint32{}, SetSubstract([]uint32{1, 2, 3}, []uint32{1, 3, 2}), "valid substract")
}

func TestAddSubTokens(t *testing.T) {
	var subtokens []uint32
	var err error
	stake := NewStakeNFT("abc", []uint32{1, 2, 3}, sdk.NewInt64Coin("aaa", 1))
	_, err = stake.AddSubTokens([]uint32{3})
	require.Error(t, err, "repeated token")
	subtokens, err = stake.AddSubTokens([]uint32{5, 4, 6})
	require.NoError(t, err, "valid add")
	require.True(t, SetHasSubset([]uint32{1, 2, 3, 4, 5, 6}, subtokens) && SetHasSubset(subtokens, []uint32{1, 2, 3, 4, 5, 6}))

	stake = NewStakeNFT("abc", []uint32{1, 2, 3, 3}, sdk.NewInt64Coin("aaa", 1))
	_, err = stake.AddSubTokens([]uint32{4})
	require.Error(t, err, "repeated token 2")
}

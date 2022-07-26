package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNFTEditMetadata(t *testing.T) {
	nft := NewBaseNFT(
		firstID,
		firstAddress,
		firstTokenURI,
		firstReserve,
		firstAllowMint,
	)
	require.Equal(t, firstTokenURI, nft.TokenURI)

	nft.EditMetadata(secondTokenURI)
	require.Equal(t, firstTokenURI, nft.TokenURI)

	nft = nft.EditMetadata(secondTokenURI)
	require.Equal(t, secondTokenURI, nft.TokenURI)
}

func TestNFTSetOwners(t *testing.T) {
	nft := NewBaseNFT(
		firstID,
		firstAddress,
		firstTokenURI,
		firstReserve,
		firstAllowMint,
	)
	require.Nil(t, nft.Owners)

	firstTokenOwner := NewTokenOwner(firstAddress, []uint64{1, 2, 3})
	nft.GetOwners().SetOwner(firstTokenOwner)
	require.Len(t, nft.Owners, 0)

	nft.SetOwners(nft.GetOwners().SetOwner(firstTokenOwner))
	require.Len(t, nft.Owners, 0)

	nft = nft.SetOwners(nft.GetOwners().SetOwner(firstTokenOwner))
	require.Len(t, nft.Owners, 1)

	secondTokenOwner := NewTokenOwner(secondAddress, []uint64{4, 5, 6})
	nft.GetOwners().SetOwner(secondTokenOwner)
	require.Len(t, nft.Owners, 1)

	nft.SetOwners(nft.GetOwners().SetOwner(secondTokenOwner))
	require.Len(t, nft.Owners, 1)

	nft = nft.SetOwners(nft.GetOwners().SetOwner(secondTokenOwner))
	require.Len(t, nft.Owners, 2)
}

func TestNFTUpdateOwner(t *testing.T) {
	nft := NewBaseNFT(
		firstID,
		firstAddress,
		firstTokenURI,
		firstReserve,
		firstAllowMint,
	)
	require.Nil(t, nft.Owners)

	// Set first owner with sub tokens
	subTokenIDs := []uint64{1, 2, 3}
	firstTokenOwner := NewTokenOwner(firstAddress, subTokenIDs)
	nft = nft.SetOwners(nft.GetOwners().SetOwner(firstTokenOwner))
	require.Len(t, nft.Owners, 1)
	require.Len(t, nft.GetOwners().GetOwner(firstAddress).SubTokenIDs, len(subTokenIDs))

	// Update first owner sub tokens
	newSubTokenIDs := []uint64{1}
	newTokenOwner := NewTokenOwner(firstAddress, newSubTokenIDs)
	nft.GetOwners().SetOwner(newTokenOwner)
	require.Len(t, nft.Owners, 1)
	require.Len(t, nft.GetOwners().GetOwner(firstAddress).SubTokenIDs, len(newSubTokenIDs))
}

func TestNFTGenSubTokenIDs(t *testing.T) {
	nft := NewBaseNFT(
		firstID,
		firstAddress,
		firstTokenURI,
		firstReserve,
		firstAllowMint,
	)
	require.Nil(t, nft.Owners)

	subTokenIDs := nft.GenSubTokenIDs(5)
	require.Len(t, subTokenIDs, 5)
	require.Equal(t, SortedUintArray{1, 2, 3, 4, 5}, subTokenIDs)

	// Set owner with sub tokens
	firstTokenOwner := NewTokenOwner(firstAddress, subTokenIDs)
	nft = nft.SetOwners(nft.GetOwners().SetOwner(firstTokenOwner))
	require.Len(t, nft.GetOwners().GetOwner(firstAddress).SubTokenIDs, len(subTokenIDs))

	subTokenIDs = nft.GenSubTokenIDs(5)
	require.Len(t, subTokenIDs, 5)
	require.Equal(t, SortedUintArray{6, 7, 8, 9, 10}, subTokenIDs)

	// Set second owner with sub tokens
	secondTokenOwner := NewTokenOwner(secondAddress, subTokenIDs)
	nft = nft.SetOwners(nft.GetOwners().SetOwner(secondTokenOwner))
	require.Len(t, nft.GetOwners().GetOwner(secondAddress).SubTokenIDs, len(subTokenIDs))

	subTokenIDs = nft.GenSubTokenIDs(5)
	require.Len(t, subTokenIDs, 5)
	require.Equal(t, SortedUintArray{11, 12, 13, 14, 15}, subTokenIDs)
}

package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectionAddNFT(t *testing.T) {
	initialIDs := []string{"222", "111", "333"}

	collection := NewCollection(firstDenom, initialIDs)
	require.Len(t, collection.NFTs, len(initialIDs))

	collection.AddNFT("444")
	require.Len(t, collection.NFTs, len(initialIDs))

	newSubTokenIDs := SortedStringArray{"111", "222", "333", "444"}
	collection = collection.AddNFT("444")
	require.Len(t, collection.NFTs, len(newSubTokenIDs))
	require.Equal(t, newSubTokenIDs, collection.NFTs)
}

func TestCollectionSupply(t *testing.T) {
	initialIDs := []string{"222", "111", "333"}

	collection := NewCollection(firstDenom, initialIDs)
	require.Len(t, collection.NFTs, len(initialIDs))
	require.Equal(t, len(initialIDs), collection.Supply())
}

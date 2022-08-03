package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	firstID        string = "first_id"
	firstDenom     string = "first_test_denom"
	firstTokenURI  string = "first_token_uri"
	secondTokenURI string = "second_token_uri"
	firstAllowMint bool   = true
	firstAddress   string = "first_address"
	secondAddress  string = "second_address"
)

var (
	firstReserve sdk.Int = NewMinReserve2
)

func TestOwnerCollectionAdd(t *testing.T) {
	initialIDs := []string{"222", "111", "333"}

	ownerCollection := NewOwnerCollection(firstDenom, initialIDs)
	require.Equal(t, OwnerCollection{
		Denom: firstDenom,
		NFTs:  []string{"111", "222", "333"},
	}, ownerCollection)

	ownerCollection = ownerCollection.AddID("444")
	require.Equal(t, OwnerCollection{
		Denom: firstDenom,
		NFTs:  []string{"111", "222", "333", "444"},
	}, ownerCollection)
}

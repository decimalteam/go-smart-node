package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSubToken creates a new SubToken instance
func NewSubToken(subTokenID uint64, reserve sdk.Int) SubToken {
	return SubToken{
		ID:      subTokenID,
		Reserve: reserve,
	}
}

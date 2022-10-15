package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LegacyKeeper interface {
	IsLegacyAddress(ctx sdk.Context, address string) bool
	ActualizeLegacy(ctx sdk.Context, pubKeyBytes []byte) error
}

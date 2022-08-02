package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeKeeper interface {
	GetPrice(ctx sdk.Context) sdk.Dec
}

package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

type FeeKeeper interface {
	GetPrice(ctx sdk.Context) sdk.Dec
	GetModuleParams(ctx sdk.Context) Params
}

// interface from ethermint evm module
type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) feemarkettypes.Params
}

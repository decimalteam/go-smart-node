package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

type FeeKeeper interface {
	GetPrice(ctx sdk.Context, denom, quote string) (CoinPrice, error)
	GetModuleParams(ctx sdk.Context) Params
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}

// interface from ethermint evm module
type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) feemarkettypes.Params
	GetModuleParams(ctx sdk.Context) Params
	GetPrice(ctx sdk.Context) (sdk.Dec, error)
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}

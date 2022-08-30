package keeper

import (
	"math/big"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

// implementation of interface FeeMarketKeeper
// for evm module
var _ types.FeeMarketKeeper = Keeper{}

var defaultBase = sdk.NewInt(1000000000)

func (k Keeper) GetBaseFee(ctx sdk.Context) *big.Int {
	price, err := k.GetPrice(ctx)
	if err != nil {
		// fallback to default price
		return defaultBase.BigInt()
	}
	fee := sdk.OneDec().MulInt(defaultBase).Quo(price).RoundInt()
	return fee.BigInt()
}

func (k Keeper) GetParams(ctx sdk.Context) feemarkettypes.Params {
	// TODO: watch for new params
	return feemarkettypes.Params{
		// we always have base fee
		NoBaseFee: false,
		BaseFee:   sdk.NewIntFromBigInt(k.GetBaseFee(ctx)),
		// these parameters is using inside feemarket module
		BaseFeeChangeDenominator: 1,
		ElasticityMultiplier:     1,
		EnableHeight:             0,
		// see ethermint - x/feemarket/types/params.go
		MinGasPrice:      feemarkettypes.DefaultMinGasPrice,
		MinGasMultiplier: feemarkettypes.DefaultMinGasMultiplier,
	}
}

func (k Keeper) AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error) {
	// TODO: this function is used in NewGasWantedDecorator
	// Do we need implement?
	return 0, nil
}

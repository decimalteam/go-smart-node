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
	var maxGas int64
	cp := k.baseApp.GetConsensusParams(ctx)
	if !(cp == nil || cp.Block == nil) {
		maxGas = cp.Block.MaxGas
	}
	// *1000 is for right scaling maxGas to fee
	feeByConsensusParam := sdk.NewInt(maxGas).MulRaw(1000)

	price, err := k.GetPrice(ctx)
	if err != nil {
		// fallback to default price
		return defaultBase.BigInt()
	}
	feeByPrice := sdk.OneDec().MulInt(defaultBase).Quo(price).RoundInt()

	// maxGas may be -1 (nol limit) or 0 (something wrong with consensus param)
	if feeByConsensusParam.IsPositive() && feeByConsensusParam.LT(feeByPrice) {
		return feeByConsensusParam.BigInt()
	}
	return feeByPrice.BigInt()
}

func (k Keeper) GetParams(ctx sdk.Context) feemarkettypes.Params {
	// TODO: watch for new params
	return feemarkettypes.NewParams(
		false,                                  // noBaseFee bool,
		1,                                      // baseFeeChangeDenom,
		1,                                      // elasticityMultiplier uint32,
		k.GetBaseFee(ctx).Uint64(),             // baseFee uint64,
		0,                                      // enableHeight int64,
		feemarkettypes.DefaultMinGasPrice,      // minGasPrice sdk.Dec,
		feemarkettypes.DefaultMinGasMultiplier, // minGasPriceMultiplier sdk.Dec,
	)
}

func (k Keeper) AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error) {
	// TODO: this function is used in NewGasWantedDecorator
	// Do we need implement?
	return 0, nil
}

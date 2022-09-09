package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"math/big"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

// implementation of interface FeeMarketKeeper
// for evm module
var _ types.FeeMarketKeeper = &Keeper{}

func (k Keeper) GetBaseFee(ctx sdk.Context) *big.Int {
	return big.NewInt(0)
}

func (k Keeper) GetParams(ctx sdk.Context) feemarkettypes.Params {
	baseDenomPrice, err := k.GetPrice(ctx)
	if err != nil {
		panic(err)
	}

	evmGasPrice := helpers.DecToDecWithE18(k.GetModuleParams(ctx).EvmGasPrice)
	minGasPrice := evmGasPrice.Quo(baseDenomPrice)

	// TODO: watch for new params
	return feemarkettypes.NewParams(
		true,                                   //noBaseFee bool,
		1,                                      //baseFeeChangeDenom,
		1,                                      //elasticityMultiplier uint32,
		k.GetBaseFee(ctx).Uint64(),             //baseFee uint64,
		0,                                      //enableHeight int64,
		minGasPrice,                            //minGasPrice sdk.Dec,
		feemarkettypes.DefaultMinGasMultiplier, //minGasPriceMultiplier sdk.Dec,
	)
}

func (k Keeper) AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error) {
	// TODO: this function is used in NewGasWantedDecorator
	// Do we need implement?
	return 0, nil
}

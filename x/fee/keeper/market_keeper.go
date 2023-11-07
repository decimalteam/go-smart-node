package keeper

import (
	"math/big"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"
)

// implementation of interface FeeMarketKeeper
// for evm module
var _ types.FeeMarketKeeper = &Keeper{}

func (k Keeper) GetBaseFee(ctx sdk.Context) *big.Int {
	baseFee := k.GetMinGasPrice(ctx).TruncateInt()

	return baseFee.BigInt()
}

func (k Keeper) GetBaseFeeEnabled(ctx sdk.Context) bool {
	return true
}

func (k Keeper) GetParams(ctx sdk.Context) feemarkettypes.Params {
	minGasPrice := k.GetMinGasPrice(ctx)

	// TODO: watch for new params
	return feemarkettypes.NewParams(
		false,                                  //noBaseFee bool,
		1,                                      //baseFeeChangeDenom,
		1,                                      //elasticityMultiplier uint32,
		k.GetBaseFee(ctx).Uint64(),             //baseFee uint64,
		0,                                      //enableHeight int64,
		minGasPrice,                            //minGasPrice sdk.Dec,
		feemarkettypes.DefaultMinGasMultiplier, //minGasPriceMultiplier sdk.Dec,
	)
}

func (k Keeper) GetMinGasPrice(ctx sdk.Context) sdk.Dec {
	baseDenomPrice, err := k.GetPrice(ctx, config.BaseDenom, feeconfig.DefaultQuote)
	if err != nil {
		panic(err)
	}

	evmGasPrice := helpers.DecToDecWithE18(k.GetModuleParams(ctx).EvmGasPrice)

	res := evmGasPrice.Quo(baseDenomPrice.Price)

	return res
}

func (k Keeper) AddTransientGasWanted(_ sdk.Context, _ uint64) (uint64, error) {
	// TODO: this function is used in NewGasWantedDecorator
	// Do we need implement?
	return 0, nil
}

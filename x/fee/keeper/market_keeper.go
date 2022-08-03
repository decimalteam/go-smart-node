package keeper

import (
	"math/big"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

// implementation of interface FeeMarketKeeper
// for evm module
var _ types.FeeMarketKeeper = Keeper{}

func (k Keeper) GetBaseFee(ctx sdk.Context) *big.Int {
	// TODO: implement
	return big.NewInt(1000000000) // default base fee from feemarket module
}

func (k Keeper) GetParams(ctx sdk.Context) feemarkettypes.Params {
	return feemarkettypes.Params{
		// we always have base fee
		NoBaseFee: false,
		BaseFee:   sdk.NewIntFromBigInt(k.GetBaseFee(ctx)),
		// these parameters is using inside feemarket module
		BaseFeeChangeDenominator: 1,
		ElasticityMultiplier:     1,
		EnableHeight:             0,
	}
}

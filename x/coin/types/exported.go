package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ContextFeeKey defines special key type used to store fee info to the context.
type ContextFeeKey struct{}

// CoinKeeper defines the exported coin keeper.
type CoinKeeper interface {
	// Params
	GetParams(ctx sdk.Context) (params Params)
	SetParams(ctx sdk.Context, params Params)

	// Coins
	GetCoins(ctx sdk.Context) (coins []Coin)
	GetCoin(ctx sdk.Context, denom string) (coin Coin, err error)
	GetCoinByDrc20(ctx sdk.Context, drc20 string) (coin Coin, err error)
	SetCoin(ctx sdk.Context, coin Coin)
	UpdateCoinVR(ctx sdk.Context, denom string, volume sdkmath.Int, reserve sdkmath.Int) error

	// Checks
	IsCheckRedeemed(ctx sdk.Context, check *Check) bool
	GetChecks(ctx sdk.Context) (checks []Check)
	GetCheck(ctx sdk.Context, checkHash []byte) (check Check, err error)
	SetCheck(ctx sdk.Context, check *Check)

	// Module specific
	IsCoinBase(ctx sdk.Context, denom string) bool
	GetBaseDenom(ctx sdk.Context) string
	CheckFutureChanges(ctx sdk.Context, coinInfo Coin, amount sdkmath.Int) error
}

package types

//go:generate mockgen -destination=../testutil/expected_bank_keeper_mock.go -package=testutil "github.com/cosmos/cosmos-sdk/x/bank/keeper" Keeper
//go:generate mockgen -destination=../testutil/expected_coin_keeper_mock.go -package=testutil . CoinKeeper
//go:generate mockgen -destination=../testutil/expected_auth_keeper_mock.go -package=testutil . AccountKeeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

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
	GetPrice(ctx sdk.Context, denom string, quote string) (CoinPrice, error)
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}

type CoinKeeper interface {
	BurnPoolCoins(ctx sdk.Context, poolName string, coins sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

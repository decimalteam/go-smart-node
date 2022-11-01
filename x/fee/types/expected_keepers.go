package types

//go:generate mockgen -destination=../testutil/expected_bank_keeper_mock.go -package=testutil "github.com/cosmos/cosmos-sdk/x/bank/keeper" Keeper
//go:generate mockgen -destination=../testutil/expected_coin_keeper_mock.go -package=testutil . CoinKeeper
//go:generate mockgen -destination=../testutil/expected_auth_keeper_mock.go -package=testutil . AccountKeeper

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
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
	GetCoin(ctx sdk.Context, denom string) (coin cointypes.Coin, err error)
	GetCoins(ctx sdk.Context) (coins []cointypes.Coin)
	GetBaseDenom(ctx sdk.Context) string
	GetDecreasingFactor(ctx sdk.Context, coin sdk.Coin) (sdk.Dec, error)
	BurnPoolCoins(ctx sdk.Context, poolName string, coins sdk.Coins) error
	UpdateCoinVR(ctx sdk.Context, denom string, volume sdkmath.Int, reserve sdkmath.Int) error
	IsCoinExists(ctx sdk.Context, denom string) bool
}

type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

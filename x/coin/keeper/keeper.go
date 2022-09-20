package keeper

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.StoreKey
	ps       paramtypes.Subspace

	accountKeeper auth.AccountKeeperI
	bankKeeper    bank.Keeper

	// cached params value (for optimization)
	cacheParams types.Params
}

// NewKeeper creates new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	ps paramtypes.Subspace,
	ac auth.AccountKeeperI,
	bk bank.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	keeper := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		ps:            ps,
		accountKeeper: ac,
		bankKeeper:    bk,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the total set of the module parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.cacheParams = params
	k.ps.SetParamSet(ctx, &params)
}

// GetBaseDenom returns base coin denomination.
func (k *Keeper) GetBaseDenom(ctx sdk.Context) string {
	if len(k.cacheParams.BaseDenom) == 0 {
		k.cacheParams = k.GetParams(ctx)
	}
	return k.cacheParams.BaseDenom
}

// IsCoinBase returns true if specified denom is base icon.
func (k *Keeper) IsCoinBase(ctx sdk.Context, denom string) bool {
	return k.GetBaseDenom(ctx) == denom
}

func (k *Keeper) GetCommission(ctx sdk.Context, feeAmountBase sdkmath.Int) (feeAmount sdkmath.Int, denom string, err error) {
	baseCoinDenom := k.GetBaseDenom(ctx)

	fee, ok := ctx.Value(types.ContextFeeKey{}).(sdk.Coins)
	if !ok || len(fee) == 0 {
		feeAmount = feeAmountBase
		denom = baseCoinDenom
		return
	}

	denom = strings.ToLower(fee[0].Denom)
	if denom != baseCoinDenom {
		coin, err := k.GetCoin(ctx, denom)
		if err != nil {
			return sdkmath.Int{}, "", err
		}

		if coin.Reserve.LT(feeAmountBase) {
			return sdkmath.Int{}, "", errors.InsufficientCoinReserve
		}

		feeAmount = formulas.CalculateSaleAmount(coin.Volume, coin.Reserve, uint(coin.CRR), feeAmountBase)
	} else {
		feeAmount = feeAmountBase
	}

	return
}

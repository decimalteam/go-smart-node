package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkstore "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	storeKey sdkstore.StoreKey // Unexposed key to access store from sdk.Context
	cdc      codec.BinaryCodec // The amino codec for binary encoding/decoding.
	ps       paramtypes.Subspace

	bankKeeper keeper.Keeper
	coinKeeper types.CoinKeeper
	authKeeper types.AccountKeeper
	calcFunc   types.CalculateCommissionFunc

	baseDenom *string
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdkstore.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper keeper.Keeper,
	coinKeeper types.CoinKeeper,
	authKeeper types.AccountKeeper,
	baseDenom string,
	calcFunc types.CalculateCommissionFunc,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		ps:         ps,
		bankKeeper: bankKeeper,
		coinKeeper: coinKeeper,
		authKeeper: authKeeper,
		baseDenom:  &baseDenom,
		calcFunc:   calcFunc,
	}
}

////////////////////////////////////////////////////////////////
// Params
////////////////////////////////////////////////////////////////

// GetParams returns the total set of the module parameters.
func (k *Keeper) GetModuleParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters to the param space.
func (k *Keeper) SetModuleParams(ctx sdk.Context, params types.Params) {
	k.ps.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SavePrice(
	ctx sdk.Context,
	price types.CoinPrice,
) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceKey(price.Denom, price.Quote)
	value := k.cdc.MustMarshalLengthPrefixed(&price)
	store.Set(key, value)
	return nil
}

func (k *Keeper) GetPrice(
	ctx sdk.Context,
	denom string,
	quote string,
) (types.CoinPrice, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceKey(denom, quote)
	value := store.Get(key)
	if len(value) == 0 {
		return types.CoinPrice{}, errors.PriceNotFound
	}
	var price types.CoinPrice
	err := k.cdc.UnmarshalLengthPrefixed(value, &price)
	if err != nil {
		return types.CoinPrice{}, errors.Internal.Wrapf("err: %s", err.Error())
	}
	return price, nil
}

func (k *Keeper) GetPrices(
	ctx sdk.Context,
) ([]types.CoinPrice, error) {
	var result []types.CoinPrice

	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStorePrefixIterator(store, types.PriceKeyPrefix)
	for ; it.Valid(); it.Next() {
		var price types.CoinPrice
		k.cdc.MustUnmarshalLengthPrefixed(it.Value(), &price)
		result = append(result, price)
	}

	err := it.Close()
	if err != nil {
		return []types.CoinPrice{}, err
	}

	return result, nil
}

package keeper

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	storeKey store.StoreKey    // Unexposed key to access store from sdk.Context
	cdc      codec.BinaryCodec // The amino codec for binary encoding/decoding.
	ps       paramtypes.Subspace

	bankKeeper keeper.Keeper

	baseDenom *string
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper keeper.Keeper,
	baseDenom string,
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
		baseDenom:  &baseDenom,
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
	value, err := price.Marshal()
	if err != nil {
		return err
	}
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
	err := price.Unmarshal(value)
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
		price.Unmarshal(it.Value())
		result = append(result, price)
	}

	err := it.Close()
	if err != nil {
		return []types.CoinPrice{}, err
	}

	return result, nil
}

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"fmt"
	"strings"
	"sync"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// Keeper implements the module data storaging.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.StoreKey
	ps       paramtypes.Subspace

	accountKeeper auth.AccountKeeperI
	bankKeeper    bank.Keeper

	baseDenom string

	coinCache      map[string]bool
	coinCacheMutex *sync.Mutex
}

// NewKeeper creates new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper auth.AccountKeeperI,
	bankKeeper bank.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	keeper := &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		ps:             ps,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		coinCache:      make(map[string]bool),
		coinCacheMutex: &sync.Mutex{},
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

////////////////////////////////////////////////////////////////
// Coin
////////////////////////////////////////////////////////////////

// GetCoin returns the coin if exists in KVStore.
func (k *Keeper) GetCoin(ctx sdk.Context, symbol string) (coin types.Coin, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCoin, []byte(strings.ToLower(symbol))...)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &coin)
	return
}

// GetCoins returns all coins existing in KVStore.
func (k *Keeper) GetCoins(ctx sdk.Context) (coins []types.Coin) {
	it := k.GetCoinsIterator(ctx)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var coin types.Coin
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &coin)
		if err != nil {
			panic(err)
		}
		coins = append(coins, coin)
	}

	return coins
}

// GetCoinsIterator returns iterator over all coins existing in KVStore.
func (k *Keeper) GetCoinsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyPrefixCoin)
}

// SetCoin writes coin to KVStore.
func (k *Keeper) SetCoin(ctx sdk.Context, coin types.Coin) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalLengthPrefixed(&coin)
	key := append(types.KeyPrefixCoin, []byte(strings.ToLower(coin.Symbol))...)
	store.Set(key, value)
}

// Edit updates current coin reserve and volume and writes coin to KVStore.
func (k *Keeper) EditCoin(ctx sdk.Context, coin types.Coin, reserve sdk.Int, volume sdk.Int) error {
	if !k.IsCoinBase(ctx, coin.Symbol) {
		k.SetCachedCoin(coin.Symbol)
	}

	// Update coin reserve and volume
	coin.Reserve = reserve
	coin.Volume = volume
	k.SetCoin(ctx, coin)

	// Emit event
	err := ctx.EventManager().EmitTypedEvent(&types.EventEditCoin{
		Symbol:  coin.Symbol,
		Volume:  coin.Volume.String(),
		Reserve: coin.Reserve.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}
	return nil
}

////////////////////////////////////////////////////////////////
// Check
////////////////////////////////////////////////////////////////

func (k *Keeper) IsCheckRedeemed(ctx sdk.Context, check *types.Check) bool {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCheck, checkHash[:]...)
	value := store.Get(key)
	if len(value) == 0 {
		return false
	}
	var c types.Check
	return k.cdc.UnmarshalLengthPrefixed(value, &c) == nil
}

func (k *Keeper) GetCheck(ctx sdk.Context, checkHash []byte) (check types.Check, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCheck, checkHash...)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CheckDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &check)
	return
}

// GetChecks returns all checks existing in KVStore.
func (k *Keeper) GetChecks(ctx sdk.Context) (checks []types.Check) {
	it := k.GetChecksIterator(ctx)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var check types.Check
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &check)
		if err != nil {
			panic(err)
		}
		checks = append(checks, check)
	}

	return checks
}

// GetChecksIterator returns iterator over all checks existing in KVStore.
func (k *Keeper) GetChecksIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyPrefixCheck)
}

// SetCheck writes check to KVStore.
func (k *Keeper) SetCheck(ctx sdk.Context, check *types.Check) {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCheck, checkHash[:]...)
	value := k.cdc.MustMarshalLengthPrefixed(check)
	store.Set(key, value)
}

////////////////////////////////////////////////////////////////
// Params
////////////////////////////////////////////////////////////////

// GetParams returns the total set of the module parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ps.SetParamSet(ctx, &params)
}

////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////

func (k *Keeper) GetBaseDenom(ctx sdk.Context) (symbol string) {
	k.ps.Get(ctx, types.ParamStoreKeyBaseSymbol, &symbol)
	return symbol
}

func (k *Keeper) IsCoinBase(ctx sdk.Context, symbol string) bool {
	return k.GetBaseDenom(ctx) == symbol
}

func (k *Keeper) GetCommission(ctx sdk.Context, feeAmountBase sdk.Int) (sdk.Int, string, error) {
	baseCoinDenom := k.GetBaseDenom(ctx)

	var feeDenom string
	fee, ok := ctx.Value("fee").(sdk.Coins)
	if !ok || len(fee) == 0 {
		feeDenom = baseCoinDenom
		return feeAmountBase, feeDenom, nil
	}

	feeDenom = strings.ToLower(fee[0].Denom)
	feeAmount := feeAmountBase
	if feeDenom != baseCoinDenom {
		coin, err := k.GetCoin(ctx, feeDenom)
		if err != nil {
			return sdk.Int{}, "", err
		}

		if coin.Reserve.LT(feeAmountBase) {
			return sdk.Int{}, "", errors.InsufficientCoinReserve
		}

		feeAmount = formulas.CalculateSaleAmount(coin.Volume, coin.Reserve, uint(coin.CRR), feeAmountBase)
	}

	return feeAmount, feeDenom, nil
}

////////////////////////////////////////////////////////////////
// Coin cache
////////////////////////////////////////////////////////////////

func (k *Keeper) GetCoinCache(symbol string) bool {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	_, ok := k.coinCache[symbol]
	return ok
}

func (k *Keeper) SetCachedCoin(coin string) {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	k.coinCache[coin] = true
}

func (k *Keeper) ClearCoinCache() {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	for key := range k.coinCache {
		delete(k.coinCache, key)
	}
}

package keeper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
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
	storeKey sdk.StoreKey
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
	storeKey sdk.StoreKey,
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
		err = fmt.Errorf("coin %s is not found in the key-value store", strings.ToLower(symbol))
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
func (k *Keeper) EditCoin(ctx sdk.Context, coin types.Coin, reserve sdk.Int, volume sdk.Int) {
	if !k.IsCoinBase(ctx, coin.Symbol) {
		k.SetCachedCoin(coin.Symbol)
	}

	// Update coin reserve and volume
	coin.Reserve = reserve
	coin.Volume = volume
	k.SetCoin(ctx, coin)

	// Emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateCoin,
		sdk.NewAttribute(types.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types.AttributeVolume, coin.Volume.String()),
		sdk.NewAttribute(types.AttributeReserve, coin.Reserve.String()),
	))
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
		err = fmt.Errorf("check with hash %X is not found in the key-value store", checkHash)
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
			return sdk.Int{}, "", fmt.Errorf(
				"coin reserve balance is not sufficient for transaction. Has: %s, required %s",
				coin.Reserve.String(), feeAmountBase.String())
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

////////////////////////////////////////////////////////////////
// Legacy balances
////////////////////////////////////////////////////////////////

// GetLegacyBalance returns balance for legacy address if exists in KVStore.
func (k *Keeper) GetLegacyBalance(ctx sdk.Context, legacyAddress string) (balance types.LegacyBalance, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixLegacy, []byte(legacyAddress)...)
	value := store.Get(key)
	if len(value) == 0 {
		err = fmt.Errorf("legacy address %s is not found in the key-value store", strings.ToLower(legacyAddress))
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &balance)
	return
}

// SetLegacyBalance store legacy balance for legacy address. Must call only in genesis
func (k *Keeper) SetLegacyBalance(ctx sdk.Context, balance types.LegacyBalance) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixLegacy, []byte(balance.LegacyAddress)...)
	value := k.cdc.MustMarshalLengthPrefixed(&balance)
	store.Set(key, value)
}

// DeleteLegacyBalance delete balance for old address. Must call in return transaction
func (k *Keeper) DeleteLegacyBalance(ctx sdk.Context, legacyAddress string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixLegacy, []byte(legacyAddress)...)
	store.Delete(key)
}

// GetLegacyBalancesIterator returns iterator over all legacy balances existing in KVStore.
func (k *Keeper) GetLegacyBalancesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyPrefixLegacy)
}

// GetLegacyBalance returns balance for old address if exists in KVStore.
func (k *Keeper) GetLegacyBalances(ctx sdk.Context) (balances []types.LegacyBalance) {
	it := k.GetLegacyBalancesIterator(ctx)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var balance types.LegacyBalance
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &balance)
		if err != nil {
			panic(err)
		}
		balances = append(balances, balance)
	}

	return balances
}

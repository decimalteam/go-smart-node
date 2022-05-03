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
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramStore paramtypes.Subspace

	AccountKeeper auth.AccountKeeperI
	BankKeeper    bank.Keeper

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
		paramStore:     ps,
		AccountKeeper:  accountKeeper,
		BankKeeper:     bankKeeper,
		coinCache:      make(map[string]bool),
		coinCacheMutex: &sync.Mutex{},
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCoin returns the coin if exists in KVStore.
func (k *Keeper) GetCoin(ctx sdk.Context, symbol string) (coin types.Coin, err error) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCoin, []byte(strings.ToLower(symbol))...)
	value := store.Get(key)
	if value == nil {
		err = fmt.Errorf("coin %s is not found in the key-value store", strings.ToLower(symbol))
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &coin)
	if err != nil {
		return
	}
	return
}

// GetAllCoins returns all coins existing in KVStore.
func (k *Keeper) GetAllCoins(ctx sdk.Context) (coins []types.Coin) {
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

func (k *Keeper) GetBaseCoin(ctx sdk.Context) string {
	return k.GetParams(ctx).BaseSymbol
}

func (k *Keeper) IsCoinBase(ctx sdk.Context, symbol string) bool {
	return k.GetParams(ctx).BaseSymbol == symbol
}

func (k *Keeper) UpdateCoin(ctx sdk.Context, coin types.Coin, reserve sdk.Int, volume sdk.Int) {
	if !k.IsCoinBase(ctx, coin.Symbol) {
		k.SetCachedCoin(coin.Symbol)
	}
	coin.Reserve = reserve
	coin.Volume = volume
	k.SetCoin(ctx, coin)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateCoin,
		sdk.NewAttribute(types.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types.AttributeVolume, coin.Volume.String()),
		sdk.NewAttribute(types.AttributeReserve, coin.Reserve.String()),
	))
}

func (k *Keeper) IsCheckRedeemed(ctx sdk.Context, check *types.Check) bool {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCheck, checkHash[:]...)
	value := store.Get(key)
	return len(value) > 0 && value[0] > 0
}

func (k *Keeper) SetCheckRedeemed(ctx sdk.Context, check *types.Check) {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCheck, checkHash[:]...)
	store.Set(key, []byte{1})
}

func (k *Keeper) GetCommission(ctx sdk.Context, commissionInBaseCoin sdk.Int) (sdk.Int, string, error) {
	var feeCoin string
	fee, ok := ctx.Value("fee").(sdk.Coins)
	if !ok || fee == nil {
		feeCoin = k.GetBaseCoin(ctx)
		return commissionInBaseCoin, feeCoin, nil
	}

	commission := sdk.ZeroInt()

	coin := fee[0]

	feeCoin = coin.Denom
	if feeCoin != k.GetBaseCoin(ctx) {
		coinInfo, err := k.GetCoin(ctx, feeCoin)
		if err != nil {
			return sdk.Int{}, "", err
		}

		if coinInfo.Reserve.LT(commissionInBaseCoin) {
			return sdk.Int{}, "", fmt.Errorf("coin reserve balance is not sufficient for transaction. Has: %s, required %s",
				coinInfo.Reserve.String(),
				commissionInBaseCoin.String())
		}

		commission = formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), commissionInBaseCoin)
	}

	return commission, feeCoin, nil
}

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

// // Updating balances
// func (k Keeper) UpdateBalance(ctx sdk.Context, symbol string, amount sdk.Int, address sdk.AccAddress) error {
// 	// Get account coin information
// 	coin := k.BankKeeper.GetBalance(ctx, address, symbol)
// 	// Updating coin information
// 	if amount.IsNegative() {
// 		coin = coin.Sub(sdk.NewCoin(symbol, amount.Neg()))
// 	} else {
// 		coin = coin.Add(sdk.NewCoin(symbol, amount))
// 	}
// 	// Get account instance
// 	acc := k.AccountKeeper.GetAccount(ctx, address)
// 	if acc == nil {
// 		acc = k.AccountKeeper.NewAccountWithAddress(ctx, address)
// 	}
// 	// Update coin information
// 	err := acc.SetCoins(coins)
// 	if err != nil {
// 		return err
// 	}
// 	// Update account information
// 	k.AccountKeeper.SetAccount(ctx, acc)
// 	return nil
// }

package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// GetCoins returns all coins existing in KVStore.
func (k *Keeper) GetCoins(ctx sdk.Context) (coins []types.Coin) {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.GetCoinsKey())
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var coin types.Coin
		// request coin
		err := k.cdc.UnmarshalLengthPrefixed(it.Value(), &coin)
		if err != nil {
			panic(err)
		}
		// NOTE: special needed step to avoid migration to add coin min emission
		if coin.MinVolume.IsNil() {
			coin.MinVolume = sdkmath.ZeroInt()
		}
		// request volume and reserve separately
		volume, reserve, err := k.getCoinVR(store, coin.Denom)
		if err != nil {
			panic(err)
		}
		coin.Volume = volume
		coin.Reserve = reserve
		coins = append(coins, coin)
	}

	return
}

// GetCoin returns the coin if exists in KVStore.
func (k *Keeper) GetCoin(ctx sdk.Context, denom string) (coin types.Coin, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCoinKey(denom)
	// request coin
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &coin)
	if err != nil {
		return
	}
	// NOTE: special needed step to avoid migration to add coin min emission
	if coin.MinVolume.IsNil() {
		coin.MinVolume = sdkmath.ZeroInt()
	}
	// request volume and reserve separately
	volume, reserve, err := k.getCoinVR(store, coin.Denom)
	if err != nil {
		return
	}
	coin.Volume = volume
	coin.Reserve = reserve
	return
}

// GetCoin returns the coin if exists in KVStore.
func (k *Keeper) GetCoinByDrc20(ctx sdk.Context, drc20Address string) (coin types.Coin, err error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCoinDrcKey(drc20Address)
	// request coin
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &coin)
	if err != nil {
		return
	}
	// NOTE: special needed step to avoid migration to add coin min emission
	if coin.MinVolume.IsNil() {
		coin.MinVolume = sdkmath.ZeroInt()
	}
	// request volume and reserve separately
	volume, reserve, err := k.getCoinVR(store, coin.Denom)
	if err != nil {
		return
	}
	coin.Volume = volume
	coin.Reserve = reserve
	return
}

// GetCoin returns the coin if exists in KVStore.
func (k *Keeper) IsCoinExists(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetCoinKey(denom))
}

// SetCoin writes coin to KVStore.
func (k *Keeper) SetCoin(ctx sdk.Context, coin types.Coin) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCoinKey(coin.Denom)
	c := coin
	// never write coin volume and reserve to main coin record
	c.Volume = sdkmath.ZeroInt()
	c.Reserve = sdkmath.ZeroInt()
	value := k.cdc.MustMarshalLengthPrefixed(&c)
	// write coin
	store.Set(key, value)
	// write volume and reserve separately
	k.setCoinVR(store, coin.Denom, coin.Volume, coin.Reserve)
}

// UpdateCoinVR updates current coin reserve and volume and writes coin to KVStore.
func (k *Keeper) UpdateCoinVR(ctx sdk.Context, denom string, volume sdkmath.Int, reserve sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	// write volume and reserve separately
	k.setCoinVR(store, denom, volume, reserve)
	err := events.EmitTypedEvent(ctx, &types.EventUpdateCoinVR{
		Denom:   denom,
		Volume:  volume.String(),
		Reserve: reserve.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}
	return nil
}

// getCoinVR returns volume and reserve of the coin if exists in KVStore.
func (k *Keeper) getCoinVR(store sdk.KVStore, denom string) (volume sdkmath.Int, reserve sdkmath.Int, err error) {
	key := types.GetCoinVRKey(denom)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	var coinVR types.CoinVR
	err = k.cdc.UnmarshalLengthPrefixed(value, &coinVR)
	if err != nil {
		return
	}
	volume = coinVR.Volume
	reserve = coinVR.Reserve
	return
}

// setCoinVR writes coin volume and reserve to KVStore.
func (k *Keeper) setCoinVR(store sdk.KVStore, denom string, volume sdkmath.Int, reserve sdkmath.Int) {
	key := types.GetCoinVRKey(denom)
	value := k.cdc.MustMarshalLengthPrefixed(&types.CoinVR{
		Volume:  volume,
		Reserve: reserve,
	})
	store.Set(key, value)
}

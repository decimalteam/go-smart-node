package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"

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

// GetCoinByDRC returns the coin if exists in KVStore.
func (k *Keeper) GetCoinByDRC(ctx sdk.Context, addressDRC string) (coin types.Coin, err error) {

	addressDRC = strings.ToLower(addressDRC)
	store := ctx.KVStore(k.storeKey)
	key := types.GetCoinDRCKey(addressDRC)
	// request coin
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	var coinDRC types.CoinDRC
	err = k.cdc.UnmarshalLengthPrefixed(value, &coinDRC)
	if err != nil {
		return
	}

	coin, err = k.GetCoin(ctx, coinDRC.Denom)
	if err != nil {
		err = errors.CoinDoesNotExist
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
	coin.DRC20Contract = coinDRC.DRC20Contract
	return
}

// IsCoinExists returns the coin if exists in KVStore.
func (k *Keeper) IsCoinExists(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetCoinKey(denom))
}

// IsCoinExistsByDRC returns the coin if exists in KVStore.
func (k *Keeper) IsCoinExistsByDRC(ctx sdk.Context, addressDRC string) bool {
	addressDRC = strings.ToLower(addressDRC)
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetCoinDRCKey(addressDRC))
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

	// write DRC address separately
	k.setCoinDRC(store, coin.Denom, coin.DRC20Contract)
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

// UpdateCoinDRC updates current coin reserve and volume and writes coin to KVStore.
func (k *Keeper) UpdateCoinDRC(ctx sdk.Context, denom string, addressDRC string) error {
	store := ctx.KVStore(k.storeKey)
	// write volume and reserve separately
	k.setCoinDRC(store, denom, addressDRC)
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

// getCoinDRC returns volume and reserve of the coin if exists in KVStore.
func (k *Keeper) getCoinDRC(store sdk.KVStore, addressDRC string) (coinDRC types.CoinDRC, err error) {
	addressDRC = strings.ToLower(addressDRC)
	key := types.GetCoinDRCKey(addressDRC)
	value := store.Get(key)
	if len(value) == 0 {
		err = errors.CoinDoesNotExist
		return
	}
	err = k.cdc.UnmarshalLengthPrefixed(value, &coinDRC)
	if err != nil {
		return
	}
	return
}

// setCoinDRC writes coin volume and reserve to KVStore.
func (k *Keeper) setCoinDRC(store sdk.KVStore, denom string, addressDRC string) {
	addressDRC = strings.ToLower(addressDRC)
	key := types.GetCoinDRCKey(addressDRC)
	value := k.cdc.MustMarshalLengthPrefixed(&types.CoinDRC{
		Denom:         denom,
		DRC20Contract: addressDRC,
	})
	store.Set(key, value)
}

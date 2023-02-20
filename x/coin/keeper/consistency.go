package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

// Check than buy/sell/fee deduct operations will not violate coin constants:
// minimal volume, limit volume, minimal reserve
// It includes check balance of auth.FeeCollectorName
// Positive amount = buy = increase volume and reserve
// Negative amount = sell/deduct = decrease volume and reserve
func (k *Keeper) CheckFutureChanges(ctx sdk.Context, coinInfo types.Coin, amount sdkmath.Int) error {
	// no need to check base coin
	if coinInfo.Denom == k.GetBaseDenom(ctx) {
		return nil
	}

	// simple check new volume
	newVolume := coinInfo.Volume.Add(amount)
	if newVolume.LT(config.MinCoinSupply) {
		return errors.TxBreaksMinVolumeLimit
	}
	if newVolume.GT(coinInfo.LimitVolume) {
		return errors.TxBreaksVolumeLimit
	}
	// for sell/deduct need include auth.FeeCollectorName balance for coin
	// because this balance will be burned
	if amount.IsNegative() {
		coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), coinInfo.Denom)
		coinToBurn := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(feetypes.BurningPool), coinInfo.Denom)
		futureAmountToBurn := coinInCollector.Amount.Add(coinToBurn.Amount).Add(amount.Neg())
		// check for minimal volume
		newVolume = coinInfo.Volume.Sub(futureAmountToBurn)
		if newVolume.LT(config.MinCoinSupply) {
			return errors.TxBreaksMinVolumeLimit
		}
		// check for minimal reserve
		futureReserveToDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), futureAmountToBurn)
		if coinInfo.Reserve.Sub(futureReserveToDecrease).LT(config.MinCoinReserve) {
			return errors.TxBreaksMinReserveRule
		}
		// check for minimal emission
		if !coinInfo.MinVolume.IsZero() {
			if newVolume.LT(coinInfo.MinVolume) {
				return errors.TxBreaksMinEmission
			}
		}
	}
	return nil
}

// same as above, but check only volume
// need for burn operation, because this doest not change reserve
func (k *Keeper) CheckFutureVolumeChanges(ctx sdk.Context, coinInfo types.Coin, amount sdkmath.Int) error {
	// no need to check base coin
	if coinInfo.Denom == k.GetBaseDenom(ctx) {
		return nil
	}

	// simple check new volume
	newVolume := coinInfo.Volume.Add(amount)
	if newVolume.LT(config.MinCoinSupply) {
		return errors.TxBreaksMinVolumeLimit
	}
	if newVolume.GT(coinInfo.LimitVolume) {
		return errors.TxBreaksVolumeLimit
	}
	// for sell/deduct need include auth.FeeCollectorName balance for coin
	// because this balance will be burned
	if amount.IsNegative() {
		coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), coinInfo.Denom)
		coinToBurn := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(feetypes.BurningPool), coinInfo.Denom)
		futureAmountToBurn := coinInCollector.Amount.Add(coinToBurn.Amount).Add(amount.Neg())
		// check for minimal volume
		newVolume = coinInfo.Volume.Sub(futureAmountToBurn)
		if newVolume.LT(config.MinCoinSupply) {
			return errors.TxBreaksMinVolumeLimit
		}
	}
	return nil
}

// Special burn mainly for slashing in validator module
// It decrease both volume and reserve
// pool must be exists and must have burning right
func (k *Keeper) BurnPoolCoins(ctx sdk.Context, poolName string, coins sdk.Coins) error {
	for _, coin := range coins {
		coinInfo, err := k.GetCoin(ctx, coin.Denom)
		if err != nil {
			return err
		}
		if coin.Denom == k.GetBaseDenom(ctx) {
			err = k.UpdateCoinVR(ctx, coin.Denom, coinInfo.Volume.Sub(coin.Amount), coinInfo.Reserve)
			continue
		}
		// because BurnPoolCoins is used in EndBlocker, error in CheckFutureChanges will cause panic
		// this is the reason to disable check and allow break limits for coin in some case
		// but check for transactions still work
		/*
			err = k.CheckFutureChanges(ctx, coinInfo, coin.Amount.Neg())
			if err != nil && !goerrors.Is(err, errors.TxBreaksMinReserveRule) {
				return err
			}
		*/
		futureReserveToDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve,
			uint(coinInfo.CRR), coin.Amount)
		coinInfo.Volume = coinInfo.Volume.Sub(coin.Amount)
		coinInfo.Reserve = coinInfo.Reserve.Sub(futureReserveToDecrease)
		err = k.UpdateCoinVR(ctx, coin.Denom, coinInfo.Volume, coinInfo.Reserve)
		if err != nil {
			return err
		}
	}
	return k.bankKeeper.BurnCoins(ctx, poolName, coins)
}

func (k *Keeper) GetDecreasingFactor(ctx sdk.Context, coin sdk.Coin) (sdk.Dec, error) {
	coinInfo, err := k.GetCoin(ctx, coin.Denom)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), coin.Denom)
	return CalculateDecreasingFactor(coinInfo, coinInCollector.Amount, coin.Amount), nil
}

// Helper function for slashing in validator module
// CalculateDecreasingFactor checks future parameters for coin burn
func CalculateDecreasingFactor(coinInfo types.Coin, amountInCollector sdkmath.Int, amountToBurn sdkmath.Int) sdk.Dec {
	newAmount := amountToBurn
	futureAmountToBurn := amountInCollector.Add(newAmount)
	// check for minimal volume
	if coinInfo.Volume.Sub(futureAmountToBurn).LT(config.MinCoinSupply) {
		newAmount = coinInfo.Volume.Sub(config.MinCoinSupply).Sub(amountInCollector)
		if newAmount.IsNegative() {
			return sdk.ZeroDec()
		}
		futureAmountToBurn = amountInCollector.Add(newAmount)
	}
	reserveToDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve,
		uint(coinInfo.CRR), futureAmountToBurn)
	if coinInfo.Reserve.Sub(reserveToDecrease).LT(config.MinCoinReserve) {
		reserveToDecrease = coinInfo.Reserve.Sub(config.MinCoinReserve)
		newAmount = formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve,
			uint(coinInfo.CRR), reserveToDecrease)
		newAmount = newAmount.Sub(amountInCollector)
	}
	if newAmount.IsNegative() {
		return sdk.ZeroDec()
	}
	return sdk.NewDecFromInt(newAmount).Quo(sdk.NewDecFromInt(amountToBurn))
}

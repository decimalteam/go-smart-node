package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Check than buy/sell/fee deduct operations will not violate coin constants:
// minimal volume, limit volume, minimal reserve
// It includes check balance of auth.FeeCollectorName
// Positive amount = buy = increase volume and reserve
// Negative amount = sell/deduct = decrease volume and reserve
func (k *Keeper) CheckFutureChanges(ctx sdk.Context, coinInfo types.Coin, amount sdk.Int) error {
	// no need to chech base coin
	if coinInfo.Symbol == k.GetBaseDenom(ctx) {
		return nil
	}

	// simple check new volume
	newVolume := coinInfo.Volume.Add(amount)
	if newVolume.LT(types.MinCoinSupply) {
		return errors.TxBreaksMinVolumeLimit
	}
	if newVolume.GT(coinInfo.LimitVolume) {
		return errors.TxBreaksVolumeLimit
	}
	// for sell/deduct need include auth.FeeCollectorName balance for coin
	// because this balance will be burned
	if amount.IsNegative() {
		coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), coinInfo.Symbol)
		futureAmountToBurn := coinInCollector.Amount.Add(amount.Neg())
		// check for minimal volume
		newVolume = coinInfo.Volume.Sub(futureAmountToBurn)
		if newVolume.LT(types.MinCoinSupply) {
			return errors.TxBreaksMinVolumeLimit
		}
		// check for minimal reserve
		futureReserveToDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve,
			uint(coinInfo.CRR), futureAmountToBurn)
		if coinInfo.Reserve.Sub(futureReserveToDecrease).LT(types.MinCoinReserve) {
			return errors.TxBreaksVolumeLimit
		}
	}
	return nil
}

// same as above, but check only volume
// need for burn operation, because this doest not change reserve
func (k *Keeper) CheckFutureVolumeChanges(ctx sdk.Context, coinInfo types.Coin, amount sdk.Int) error {
	// no need to chech base coin
	if coinInfo.Symbol == k.GetBaseDenom(ctx) {
		return nil
	}

	// simple check new volume
	newVolume := coinInfo.Volume.Add(amount)
	if newVolume.LT(types.MinCoinSupply) {
		return errors.TxBreaksMinVolumeLimit
	}
	if newVolume.GT(coinInfo.LimitVolume) {
		return errors.TxBreaksVolumeLimit
	}
	// for sell/deduct need include auth.FeeCollectorName balance for coin
	// because this balance will be burned
	if amount.IsNegative() {
		coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), coinInfo.Symbol)
		futureAmountToBurn := coinInCollector.Amount.Add(amount.Neg())
		// check for minimal volume
		newVolume = coinInfo.Volume.Sub(futureAmountToBurn)
		if newVolume.LT(types.MinCoinSupply) {
			return errors.TxBreaksMinVolumeLimit
		}
	}
	return nil
}

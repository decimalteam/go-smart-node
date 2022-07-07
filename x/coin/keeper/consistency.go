package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Check than buy/sell/fee deduct operations will not violate coin constants:
// minimal volume, limit volume, minimal reserve
// It includes check balance of auth.FeeCollectorName
// Positive amount = buy = increase volume and reserve
// Negative amount = sell/deduct = decrease volume and reserve
func (k *Keeper) CheckFutureChanges(ctx sdk.Context, symbol string, amount sdk.Int) error {
	// no need to chech base coin
	if symbol == k.GetBaseDenom(ctx) {
		return nil
	}

	coinInfo, err := k.GetCoin(ctx, symbol)
	if err != nil {
		return types.ErrCoinDoesNotExist(symbol)
	}
	// simple check new volume
	newVolume := coinInfo.Volume.Add(amount)
	if newVolume.LT(types.MinCoinSupply) {
		return types.ErrTxBreaksMinVolumeLimit(newVolume.String(), types.MinCoinSupply.String())
	}
	if newVolume.GT(coinInfo.LimitVolume) {
		return types.ErrTxBreaksVolumeLimit(newVolume.String(), coinInfo.LimitVolume.String())
	}
	// for sell/deduct need include auth.FeeCollectorName balance for coin
	// because this balance will be burned
	if amount.IsNegative() {
		coinInCollector := k.bankKeeper.GetBalance(ctx, sdkAuthTypes.NewModuleAddress(sdkAuthTypes.FeeCollectorName), symbol)
		futureAmountToBurn := coinInCollector.Amount.Add(amount.Neg())
		// check for minimal volume
		newVolume = coinInfo.Volume.Sub(futureAmountToBurn)
		if newVolume.LT(types.MinCoinSupply) {
			return types.ErrTxBreaksMinVolumeLimit(newVolume.String(), types.MinCoinSupply.String())
		}
		// check for minimal reserve
		futureReserveToDecrease := formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve,
			uint(coinInfo.CRR), futureAmountToBurn)
		if coinInfo.Reserve.Sub(futureReserveToDecrease).LT(types.MinCoinReserve) {
			return types.ErrTxBreaksMinReserveRule(types.MinCoinReserve.String(),
				coinInfo.Reserve.Sub(futureReserveToDecrease).String())
		}
	}
	return nil
}

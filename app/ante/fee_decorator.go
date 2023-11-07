package ante

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	evmTypes "github.com/decimalteam/ethermint/x/evm/types"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinconfig "bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	feeerrors "bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

type FeeDecorator struct {
	coinKeeper    cointypes.CoinKeeper
	bankKeeper    evmTypes.BankKeeper
	accountKeeper evmTypes.AccountKeeper
	feeKeeper     feetypes.FeeKeeper
	cdc           codec.BinaryCodec
}

// NewFeeDecorator creates new FeeDecorator to deduct fee
func NewFeeDecorator(ck cointypes.CoinKeeper, bk evmTypes.BankKeeper, ak evmTypes.AccountKeeper, fk feetypes.FeeKeeper, cdc codec.BinaryCodec) FeeDecorator {
	return FeeDecorator{
		coinKeeper:    ck,
		bankKeeper:    bk,
		accountKeeper: ak,
		feeKeeper:     fk,
		cdc:           cdc,
	}
}

// AnteHandle implements sdk.AnteHandler function.
func (fd FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// no fee on blockchain start
	if ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate)
	}

	// initial check for transaction
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, NotFeeTxType
	}

	if addr := fd.accountKeeper.GetModuleAddress(sdkAuthTypes.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", sdkAuthTypes.FeeCollectorName))
	}

	params := fd.feeKeeper.GetModuleParams(ctx)

	delPrice, err := fd.feeKeeper.GetPrice(ctx, config.BaseDenom, feeconfig.DefaultQuote)
	if err != nil {
		return ctx, err
	}

	commissionInBaseCoin, err := CalculateFee(fd.cdc, tx.GetMsgs(), int64(len(ctx.TxBytes())), delPrice.Price, params)
	if err != nil {
		return ctx, err
	}

	feeFromTx := feeTx.GetFee()
	feePayerAcc := fd.accountKeeper.GetAccount(ctx, feeTx.FeePayer())
	baseDenom := fd.coinKeeper.GetBaseDenom(ctx)

	ctx = ctx.WithValue(cointypes.ContextFeeKey{}, feeFromTx)

	if feePayerAcc == nil {
		return ctx, FeePayerAddressDoesNotExist
	}

	// fee from transaction is zero (empty), we deduct calculated fee from payer
	if feeFromTx.IsZero() {
		// NOTE: commissionInBaseCoin may be zero in case of RedeemCheck
		// do not remove this condition
		if commissionInBaseCoin.IsZero() {
			return next(ctx, tx, simulate)
		}
		// deduct the fees
		err = DeductFees(ctx, fd.bankKeeper, fd.coinKeeper, feeTx.FeePayer(), sdk.NewCoin(baseDenom, commissionInBaseCoin), params.CommissionBurnFactor)
		if err != nil {
			return ctx, err
		}

		// need for calculate TotalGasUsed for tendermint block results
		// TODO: special gas for validator.Delegate (*10)
		gasUsed := helpers.WeiToFinney(commissionInBaseCoin).Uint64()
		ctx = SetGasMeter(simulate, ctx, gasUsed)
		ctx.GasMeter().ConsumeGas(gasUsed, GasCommissionDesc)

		return next(ctx, tx, simulate)
	}

	// fee from transaction not empty
	feeInBaseCoin := feeFromTx[0].Amount

	// calculate fee in base coin to check enough amount of fee
	if feeFromTx[0].Denom != baseDenom {
		coinInfo, err := fd.coinKeeper.GetCoin(ctx, feeFromTx[0].Denom)
		if err != nil {
			return ctx, err
		}

		if coinInfo.Reserve.LT(commissionInBaseCoin) {
			return ctx, CoinReserveInsufficient
		}

		feeInBaseCoin = formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), feeFromTx[0].Amount)

		if coinInfo.Reserve.Sub(feeInBaseCoin).LT(coinconfig.MinCoinReserve) {
			return ctx, CoinReserveBecomeInsufficient
		}
	}

	if feeInBaseCoin.LT(commissionInBaseCoin) {
		return ctx, FeeLessThanCommission
	}

	// deduct the fees
	err = DeductFees(ctx, fd.bankKeeper, fd.coinKeeper, feeTx.FeePayer(), feeFromTx[0], params.CommissionBurnFactor)
	if err != nil {
		return ctx, err
	}

	// need for calculate TotalGasUsed for tendermint block results
	// TODO: special gas for validator.Delegate (*10)
	gasUsed := helpers.WeiToFinney(feeInBaseCoin).Uint64()
	ctx = SetGasMeter(simulate, ctx, gasUsed)
	ctx.GasMeter().ConsumeGas(gasUsed, GasCommissionDesc)

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
func DeductFees(ctx sdk.Context, bankKeeper evmTypes.BankKeeper, coinKeeper cointypes.CoinKeeper,
	feePayerAddress sdk.AccAddress, fee sdk.Coin, burningFactor sdk.Dec) error {

	if !fee.IsValid() {
		return InvalidFeeAmount
	}

	// verify the account has enough funds to pay fee
	balance := bankKeeper.GetBalance(ctx, feePayerAddress, fee.Denom)
	if balance.Amount.LT(fee.Amount) {
		return InsufficientFundsToPayFee
	}

	// verify for future coin burning: we must keep minimal volume and reserve
	if !coinKeeper.IsCoinBase(ctx, fee.Denom) {
		// .Neg() becausee coins will be burn in fee collector
		feeCoin, err := coinKeeper.GetCoin(ctx, fee.Denom)
		if err != nil {
			return err
		}
		err = coinKeeper.CheckFutureChanges(ctx, feeCoin, fee.Amount.Neg())
		if err != nil {
			return err
		}
	}

	// split to burning and collected part
	amountToBurn := sdk.NewDecFromInt(fee.Amount).Mul(burningFactor).RoundInt()
	amountToCollect := fee.Amount.Sub(amountToBurn)

	// NOTE: this functionality is dublicated in x/coin/keeper/msg_server.go:RedeemCheck
	// send to burn
	if amountToBurn.IsPositive() {
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, feePayerAddress, feetypes.BurningPool,
			sdk.NewCoins(sdk.NewCoin(fee.Denom, amountToBurn)))
		if err != nil {
			return FailedToSendCoins
		}
	}

	// send to collect
	if amountToCollect.IsPositive() {
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, feePayerAddress, sdkAuthTypes.FeeCollectorName,
			sdk.NewCoins(sdk.NewCoin(fee.Denom, amountToCollect)))
		if err != nil {
			return FailedToSendCoins
		}
	}

	// Emit fee deduction event
	// need for correct balance calculation for external services
	err := events.EmitTypedEvent(ctx, &feetypes.EventPayCommission{
		Payer: feePayerAddress.String(),
		Coins: sdk.NewCoins(fee),
		Burnt: sdk.NewCoins(sdk.NewCoin(fee.Denom, amountToBurn)),
	})
	if err != nil {
		return feeerrors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}

package ante

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmTypes "github.com/tharsis/ethermint/x/evm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeDecorator struct {
	coinKeeper    coinTypes.CoinKeeper
	bankKeeper    evmTypes.BankKeeper
	accountKeeper evmTypes.AccountKeeper
}

// NewFeeDecorator creates new FeeDecorator to deduct fee
func NewFeeDecorator(ck coinTypes.CoinKeeper, bk evmTypes.BankKeeper, ak evmTypes.AccountKeeper) FeeDecorator {
	return FeeDecorator{
		coinKeeper:    ck,
		bankKeeper:    bk,
		accountKeeper: ak,
	}
}

// AnteHandle implements sdk.AnteHandler function.
func (fd FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// no fee on blockchain start
	if ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate)
	}

	// initial check for transaction
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, ErrNotFeeTxType()
	}

	if addr := fd.accountKeeper.GetModuleAddress(sdkAuthTypes.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", sdkAuthTypes.FeeCollectorName))
	}

	commissionInBaseCoin, err := CalculateFee(tx.GetMsgs(), int64(len(ctx.TxBytes())), sdk.OneDec())
	if err != nil {
		return ctx, err
	}

	feeFromTx := feeTx.GetFee()
	feePayerAcc := fd.accountKeeper.GetAccount(ctx, feeTx.FeePayer())
	baseDenom := fd.coinKeeper.GetBaseDenom(ctx)

	if feePayerAcc == nil {
		return ctx, ErrFeePayerAddressDoesNotExist(feeTx.FeePayer().String())
	}

	// fee from transaction is zero (empty), we deduct calculated fee from payer
	if feeFromTx.IsZero() {
		// deduct the fees
		if commissionInBaseCoin.IsZero() {
			return next(ctx, tx, simulate)
		}
		err = DeductFees(ctx, fd.bankKeeper, fd.coinKeeper, feeTx.FeePayer(), sdk.NewCoin(baseDenom, commissionInBaseCoin))
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
			return ctx, ErrCoinReserveInsufficient(coinInfo.Reserve.String(), commissionInBaseCoin.String())
		}

		feeInBaseCoin = formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve,
			uint(coinInfo.CRR), feeFromTx[0].Amount)

		if coinInfo.Reserve.Sub(feeInBaseCoin).LT(coinTypes.MinCoinReserve) {
			return ctx, ErrCoinReserveBecomeInsufficient(coinInfo.Reserve.String(), feeInBaseCoin.String(),
				coinTypes.MinCoinReserve.String())
		}
	}

	if feeInBaseCoin.LT(commissionInBaseCoin) {
		return ctx, ErrFeeLessThanCommission(feeInBaseCoin.String(), commissionInBaseCoin.String())
	}

	// deduct the fees
	err = DeductFees(ctx, fd.bankKeeper, fd.coinKeeper, feeTx.FeePayer(), feeFromTx[0])
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
func DeductFees(ctx sdk.Context, bankKeeper evmTypes.BankKeeper, coinKeeper coinTypes.CoinKeeper,
	feePayerAddress sdk.AccAddress, fee sdk.Coin) error {

	if !fee.IsValid() {
		return ErrInvalidFeeAmount(fee.String())
	}

	// verify the account has enough funds to pay fee
	balance := bankKeeper.GetBalance(ctx, feePayerAddress, fee.Denom)
	if balance.Amount.LT(fee.Amount) {
		return ErrInsufficientFundsToPayFee(balance.String(), fee.String())
	}

	// verify for future coin burning: we must keep minimal volume and reserve
	if !coinKeeper.IsCoinBase(ctx, fee.Denom) {
		// .Neg() becausee coins will be burn in fee collector
		feeCoin, err := coinKeeper.GetCoin(ctx, fee.Denom)
		if err != nil {
			return ErrCoinDoesNotExist(fee.Denom)
		}
		err = coinKeeper.CheckFutureChanges(ctx, feeCoin, fee.Amount.Neg())
		if err != nil {
			return err
		}
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, feePayerAddress, sdkAuthTypes.FeeCollectorName, sdk.NewCoins(fee))
	if err != nil {
		return ErrFailedToSendCoins(err.Error())
	}

	return nil
}

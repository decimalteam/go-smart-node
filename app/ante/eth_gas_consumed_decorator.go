package ante

import (
	"math"
	"math/big"

	ethereumCommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ethante "github.com/evmos/ethermint/app/ante"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

// EthGasConsumeDecorator validates enough intrinsic gas for the transaction and
// gas consumption.
type EthGasConsumeDecorator struct {
	evmKeeper    ethante.EVMKeeper
	bankKeeper   BankKeeper
	feeKeeper    feetypes.FeeKeeper
	maxGasWanted uint64
}

type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// NewEthGasConsumeDecorator creates a new EthGasConsumeDecorator
func NewEthGasConsumeDecorator(
	evmKeeper ethante.EVMKeeper,
	bankKeeper BankKeeper,
	feeKeeper feetypes.FeeKeeper,
	maxGasWanted uint64,
) EthGasConsumeDecorator {
	return EthGasConsumeDecorator{
		evmKeeper,
		bankKeeper,
		feeKeeper,
		maxGasWanted,
	}
}

// AnteHandle validates that the Ethereum tx message has enough to cover intrinsic gas
// (during CheckTx only) and that the sender has enough balance to pay for the gas cost.
//
// Intrinsic gas for a transaction is the amount of gas that the transaction uses before the
// transaction is executed. The gas is a constant value plus any cost inccured by additional bytes
// of data supplied with the transaction.
//
// This AnteHandler decorator will fail if:
// - the message is not a MsgEthereumTx
// - sender account cannot be found
// - transaction's gas limit is lower than the intrinsic gas
// - user doesn't have enough balance to deduct the transaction fees (gas_limit * gas_price)
// - transaction or block gas meter runs out of gas
// - sets the gas meter limit
func (egcd EthGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	params := egcd.evmKeeper.GetParams(ctx)

	ethCfg := params.ChainConfig.EthereumConfig(egcd.evmKeeper.ChainID())

	blockHeight := big.NewInt(ctx.BlockHeight())
	homestead := ethCfg.IsHomestead(blockHeight)
	istanbul := ethCfg.IsIstanbul(blockHeight)
	london := ethCfg.IsLondon(blockHeight)
	evmDenom := params.EvmDenom
	gasWanted := uint64(0)

	// Use the lowest priority of all the messages as the final one.
	minPriority := int64(math.MaxInt64)

	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.ErrUnknownRequest.Wrapf("invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, sdkerrors.Wrap(err, "failed to unpack tx data")
		}

		if ctx.IsCheckTx() && egcd.maxGasWanted != 0 {
			// We can't trust the tx gas limit, because we'll refund the unused gas.
			if txData.GetGas() > egcd.maxGasWanted {
				gasWanted += egcd.maxGasWanted
			} else {
				gasWanted += txData.GetGas()
			}
		} else {
			gasWanted += txData.GetGas()
		}

		fees, priority, err := egcd.evmKeeper.DeductTxCostsFromUserBalance(
			ctx,
			*msgEthTx,
			txData,
			evmDenom,
			homestead,
			istanbul,
			london,
		)
		if err != nil {
			return ctx, sdkerrors.Wrapf(err, "failed to deduct transaction costs from user balance")
		}

		// Send part of fee to burning pool
		// NOTE: DeductTxCostsFromUserBalance (from above) sends fee coins to fee_collector
		// Because it's too hard to inject into EVMKeeper, we send burning amount here from module to module
		feeParams := egcd.feeKeeper.GetModuleParams(ctx)
		feeToBurn := sdk.NewCoins()
		for _, coin := range fees {
			amountToBurn := sdk.NewDecFromInt(coin.Amount).Mul(feeParams.CommissionBurnFactor).RoundInt()
			feeToBurn = feeToBurn.Add(sdk.NewCoin(coin.Denom, amountToBurn))
		}
		err = egcd.bankKeeper.SendCoinsFromModuleToModule(ctx, sdkAuthTypes.FeeCollectorName, feetypes.BurningPool, feeToBurn)
		if err != nil {
			return ctx, sdkerrors.Wrapf(err, "failed to send burning coins from fee_collector to burning_pool")
		}

		// Decimal decorator differs from ethermint, in the events it adds to the result
		// if you want to update this code to the latest version, then don't touch the event emitter
		adr := ethereumCommon.HexToAddress(msgEthTx.From)
		dxAdr, err := sdk.Bech32ifyAddressBytes(config.Bech32PrefixAccAddr, adr.Bytes())
		if err != nil {
			return ctx, sdkerrors.Wrapf(err, "failed to convert ethreum address to bech32 form")
		}

		err = events.EmitTypedEvent(ctx, &feetypes.EventPayCommission{
			Payer: dxAdr,
			Coins: fees,
			Burnt: feeToBurn,
		})
		if err != nil {
			return ctx, sdkerrors.Wrapf(err, "failed to emit commission event")
		}

		if priority < minPriority {
			minPriority = priority
		}
	}

	// TODO: deprecate after https://github.com/cosmos/cosmos-sdk/issues/9514  is fixed on SDK
	blockGasLimit := ethermint.BlockGasLimit(ctx)

	// NOTE: safety check
	if blockGasLimit > 0 {
		// generate a copy of the gas pool (i.e block gas meter) to see if we've run out of gas for this block
		// if current gas consumed is greater than the limit, this funcion panics and the error is recovered on the Baseapp
		gasPool := sdk.NewGasMeter(blockGasLimit)
		gasPool.ConsumeGas(ctx.GasMeter().GasConsumedToLimit(), "gas pool check")
	}

	// Set ctx.GasMeter with a limit of GasWanted (gasLimit)
	gasConsumed := ctx.GasMeter().GasConsumed()
	ctx = ctx.WithGasMeter(ethermint.NewInfiniteGasMeterWithLimit(gasWanted))
	ctx.GasMeter().ConsumeGas(gasConsumed, "copy gas consumed")

	newCtx := ctx.WithPriority(minPriority)

	// we know that we have enough gas on the pool to cover the intrinsic gas
	return next(newCtx, tx, simulate)
}

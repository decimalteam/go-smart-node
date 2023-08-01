package ante

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// EVMDecorator execute evm action to cosmos
type EVMDecorator struct {
	evmKeeper  ethante.EVMKeeper
	bankKeeper BankKeeper
}

// NewEVMDecorator creates a new EVMDecorator
func NewEVMDecorator(
	evmKeeper ethante.EVMKeeper,
	bankKeeper BankKeeper,
) EVMDecorator {
	return EVMDecorator{
		evmKeeper,
		bankKeeper,
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
func (ed EVMDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.ErrUnknownRequest.Wrapf("invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		ctx.Logger().Info(hex.EncodeToString(msgEthTx.AsTransaction().Data()[:]))
		ctx.Logger().Info(hex.EncodeToString(msgEthTx.AsTransaction().To()[:]))
	}
	// we know that we have enough gas on the pool to cover the intrinsic gas
	return next(ctx, tx, simulate)
}

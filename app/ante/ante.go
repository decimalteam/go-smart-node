package ante

import (
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	updatetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethante "github.com/decimalteam/ethermint/app/ante"
)

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer ethante.Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					anteHandler = newEthAnteHandler(options)
				//case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
				//	// handle as normal Decimal SDK tx, except signature is checked for EIP712 representation
				//	anteHandler = newCosmosAnteHandlerEip712(options)
				//case "/ethermint.types.v1.ExtensionOptionDynamicFeeTx":
				//	// cosmos-sdk tx with dynamic fee extension
				//	anteHandler = newCosmosAnteHandler(options)
				default:
					return ctx, sdkerrors.ErrUnknownExtensionOptions
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		switch tx.(type) {
		case sdk.Tx:
			for _, msg := range tx.GetMsgs() {
				if _, ok := msg.(*updatetypes.MsgSoftwareUpgrade); ok {
					anteHandler = newCosmosAnteHandler(options)
					return anteHandler(ctx, tx, sim)
				}
				if _, ok := msg.(*updatetypes.MsgCancelUpgrade); ok {
					anteHandler = newCosmosAnteHandler(options)
					return anteHandler(ctx, tx, sim)
				}
				if _, ok := msg.(*validatortypes.MsgSetOnline); ok {
					anteHandler = newCosmosAnteHandler(options)
					return anteHandler(ctx, tx, sim)
				}
			}
		default:
			return ctx, sdkerrors.ErrUnknownRequest
		}

		return ctx, sdkerrors.ErrUnknownRequest
	}
}

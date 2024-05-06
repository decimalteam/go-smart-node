package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"

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

		return ctx, sdkerrors.ErrUnknownRequest

		// handle as totally normal Decimal SDK tx
		//switch tx.(type) {
		//case sdk.Tx:
		//	anteHandler = newCosmosAnteHandler(options)
		//default:
		//	return ctx, sdkerrors.ErrUnknownRequest
		//}
		//
		//return anteHandler(ctx, tx, sim)
	}
}

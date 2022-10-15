package ante

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ibcante "github.com/cosmos/ibc-go/v5/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	ethante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacytypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
)

// HandlerOptions defines the list of module keepers required to run the Decimal AnteHandler decorators.
type HandlerOptions struct {
	AccountKeeper          evmtypes.AccountKeeper
	BankKeeper             evmtypes.BankKeeper
	IBCKeeper              *ibckeeper.Keeper
	FeeMarketKeeper        evmtypes.FeeMarketKeeper
	EvmKeeper              ethante.EVMKeeper
	FeegrantKeeper         authante.FeegrantKeeper
	CoinKeeper             cointypes.CoinKeeper
	LegacyKeeper           legacytypes.LegacyKeeper
	FeeKeeper              feetypes.FeeKeeper
	SignModeHandler        authsigning.SignModeHandler
	SigGasConsumer         func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
	Cdc                    codec.BinaryCodec
	MaxTxGasWanted         uint64
	ExtensionOptionChecker authante.ExtensionOptionChecker
	TxFeeChecker           authante.TxFeeChecker
}

// Validate checks if the keepers are defined
func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.ErrLogic.Wrapf("sign mode handler is required for ante builder")
	}
	if options.FeeMarketKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("evm keeper is required for AnteHandler")
	}
	if options.FeegrantKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("feegrant keeper is required for AnteHandler")
	}
	if options.CoinKeeper == nil {
		return sdkerrors.ErrLogic.Wrapf("coin keeper is required for AnteHandler")
	}
	if options.LegacyKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "legacy keeper is required for AnteHandler")
	}
	return nil
}

// newCosmosAnteHandler creates the default ante handler for Ethereum transactions
func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.NewEthSetUpContextDecorator(options.EvmKeeper), // outermost AnteDecorator. SetUpContext must be called first
		NewCountMsgDecorator(),
		ethante.NewEthMempoolFeeDecorator(options.EvmKeeper),                   // Check eth effective gas price against minimal-gas-prices
		NewEthMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper), // Check eth effective gas price against the global MinGasPrice
		ethante.NewEthValidateBasicDecorator(options.EvmKeeper),
		ethante.NewEthSigVerificationDecorator(options.EvmKeeper),
		ethante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		ethante.NewCanTransferDecorator(options.EvmKeeper),
		NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		ethante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
		ethante.NewEthEmitEventDecorator(options.EvmKeeper), // emit eth tx hash and index at the very last ante handler.
	)
}

// newCosmosAnteHandler creates the default ante handler for Cosmos transactions
// keep in sync with newCosmosAnteHandlerEip712, except signature verification
func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewSetUpContextDecorator(),
		NewCountMsgDecorator(),
		authante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewPreCreateAccountDecorator(options.AccountKeeper), // should be before SetPubKeyDecorator
		NewFeeDecorator(options.CoinKeeper, options.BankKeeper, options.AccountKeeper, options.FeeKeeper),
		NewValidatorCommissionDecorator(options.Cdc),
		// SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewSetPubKeyDecorator(options.AccountKeeper),
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		authante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewPostCreateAccountDecorator(options.AccountKeeper), // should be after SigVerificationDecorator
		NewLegacyActualizerDecorator(options.LegacyKeeper),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}

// newCosmosAnteHandlerEip712 creates the ante handler for transactions signed with EIP-712
// keep in sync with newCosmosAnteHandler, except signature verification
func newCosmosAnteHandlerEip712(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewSetUpContextDecorator(),
		NewCountMsgDecorator(),
		// NOTE: extensions option decorator removed
		// authante.NewRejectExtensionOptionsDecorator(),
		ethante.NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewPreCreateAccountDecorator(options.AccountKeeper), // should be before SetPubKeyDecorator
		NewFeeDecorator(options.CoinKeeper, options.BankKeeper, options.AccountKeeper, options.FeeKeeper),
		NewValidatorCommissionDecorator(options.Cdc),
		// SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewSetPubKeyDecorator(options.AccountKeeper),
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		// Note: signature verification uses EIP instead of the cosmos signature validator
		ethante.NewEip712SigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewPostCreateAccountDecorator(options.AccountKeeper), // should be after SigVerificationDecorator
		NewLegacyActualizerDecorator(options.LegacyKeeper),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}

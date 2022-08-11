package ante

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	evmante "github.com/tharsis/ethermint/app/ante"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	legacytypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
)

// HandlerOptions defines the list of module keepers required to run the Decimal AnteHandler decorators.
type HandlerOptions struct {
	AccountKeeper evmtypes.AccountKeeper
	BankKeeper    evmtypes.BankKeeper
	//IBCKeeper       *ibckeeper.Keeper
	FeeMarketKeeper evmtypes.FeeMarketKeeper
	EvmKeeper       evmante.EVMKeeper
	FeegrantKeeper  authante.FeegrantKeeper
	CoinKeeper      cointypes.CoinKeeper
	LegacyKeeper    legacytypes.LegacyKeeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
	Cdc             codec.BinaryCodec
	MaxTxGasWanted  uint64
}

// Validate checks if the keepers are defined
func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.FeeMarketKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for AnteHandler")
	}
	if options.FeegrantKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "feegrant keeper is required for AnteHandler")
	}
	if options.CoinKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "coin keeper is required for AnteHandler")
	}
	if options.LegacyKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "legacy keeper is required for AnteHandler")
	}
	return nil
}

// newCosmosAnteHandler creates the default ante handler for Ethereum transactions
func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		evmante.NewEthSetUpContextDecorator(options.EvmKeeper), // outermost AnteDecorator. SetUpContext must be called first
		evmante.NewEthMempoolFeeDecorator(options.EvmKeeper),   // Check eth effective gas price against minimal-gas-prices
		evmante.NewEthValidateBasicDecorator(options.EvmKeeper),
		evmante.NewEthSigVerificationDecorator(options.EvmKeeper),
		evmante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		evmante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		//evmante.NewCanTransferDecorator(options.EvmKeeper),
		evmante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
	)
}

// newCosmosAnteHandler creates the default ante handler for Cosmos transactions
// keep in sync with newCosmosAnteHandlerEip712, except signature verification
func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		evmante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewSetUpContextDecorator(),
		authante.NewRejectExtensionOptionsDecorator(),
		authante.NewMempoolFeeDecorator(),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewFeeDecorator(options.CoinKeeper, options.BankKeeper, options.AccountKeeper),
		NewValidatorCommissionDecorator(options.Cdc),
		// SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewSetPubKeyDecorator(options.AccountKeeper),
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		authante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewLegacyActualizerDecorator(options.LegacyKeeper),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		//ibcante.NewAnteDecorator(options.IBCKeeper),
	)
}

// newCosmosAnteHandlerEip712 creates the ante handler for transactions signed with EIP712
// keep in sync with newCosmosAnteHandler, except signature verification
func newCosmosAnteHandlerEip712(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		evmante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewSetUpContextDecorator(),
		authante.NewMempoolFeeDecorator(),
		authante.NewValidateBasicDecorator(),
		authante.NewTxTimeoutHeightDecorator(),
		authante.NewValidateMemoDecorator(options.AccountKeeper),
		authante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewFeeDecorator(options.CoinKeeper, options.BankKeeper, options.AccountKeeper),
		NewValidatorCommissionDecorator(options.Cdc),
		// SetPubKeyDecorator must be called before all signature verification decorators
		authante.NewSetPubKeyDecorator(options.AccountKeeper),
		authante.NewValidateSigCountDecorator(options.AccountKeeper),
		authante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		// Note: signature verification uses EIP instead of the cosmos signature validator
		evmante.NewEip712SigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewLegacyActualizerDecorator(options.LegacyKeeper),
		authante.NewIncrementSequenceDecorator(options.AccountKeeper),
		//ibcante.NewAnteDecorator(options.IBCKeeper),
	)
}

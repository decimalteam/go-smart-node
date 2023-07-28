package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	PSKeyTxByteFee = []byte("TxByteFee")
	// coin transactions fees
	PSKeyCoinCreate      = []byte("CoinCreate")
	PSKeyCoinUpdate      = []byte("CoinUpdate")
	PSKeyCoinSend        = []byte("CoinSend")
	PSKeyCoinSendAdd     = []byte("CoinSendAdd")
	PSKeyCoinBuy         = []byte("CoinBuy")
	PSKeyCoinSell        = []byte("CoinSell")
	PSKeyCoinRedeemCheck = []byte("CoinRedeemCheck")
	PSKeyCoinBurn        = []byte("CoinCoinBurn")
	// special commission depends on coin symbol length
	PSKeyCoinCreateTicker3 = []byte("CoinCreateTicker3")
	PSKeyCoinCreateTicker4 = []byte("CoinCreateTicker4")
	PSKeyCoinCreateTicker5 = []byte("CoinCreateTicker5")
	PSKeyCoinCreateTicker6 = []byte("CoinCreateTicker6")
	PSKeyCoinCreateTicker7 = []byte("CoinCreateTicker7")
	// multisignature wallets
	PSKeyMultisigCreateWallet      = []byte("MultisigCreateWallet")
	PSKeyMultisigCreateTransaction = []byte("MultisigCreateTransaction")
	PSKeyMultisigSignTransaction   = []byte("MultisigSignTransaction")
	// nft
	PSKeyNftMintToken     = []byte("NftMintToken")
	PSKeyNftUpdateToken   = []byte("NftUpdateToken")
	PSKeyNftUpdateReserve = []byte("NftUpdateReserve")
	PSKeyNftSendToken     = []byte("NftSendToken")
	PSKeyNftBurnToken     = []byte("NftBurnToken")
	// swap
	PSKeySwapActivateChain   = []byte("SwapActivateChain")
	PSKeySwapDeactivateChain = []byte("SwapDeactivateChain")
	PSKeySwapInitialize      = []byte("SwapInitialize")
	PSKeySwapRedeem          = []byte("SwapRedeem")
	// validator operations
	PSKeyValidatorCreateValidator = []byte("ValidatorCreateValidator")
	PSKeyValidatorEditValidator   = []byte("ValidatorEditValidator")
	PSKeyValidatorDelegate        = []byte("ValidatorDelegate")
	PSKeyValidatorDelegateNFT     = []byte("ValidatorDelegateNFT")
	PSKeyValidatorRedelegate      = []byte("ValidatorRedelegate")
	PSKeyValidatorRedelegateNFT   = []byte("ValidatorRedelegateNFT")
	PSKeyValidatorUndelegate      = []byte("ValidatorUndelegate")
	PSKeyValidatorUndelegateNFT   = []byte("ValidatorUndelegateNFT")
	PSKeyValidatorSetOnline       = []byte("ValidatorSetOnline")
	PSKeyValidatorSetOffline      = []byte("ValidatorSetOffline")
	// commission burn factor
	PSKeyCommissionBurnFactor = []byte("CommissionBurnFactor")
	// oracle key
	PSKeyOracle = []byte("Oracle")
	// evm tx keys
	PSKeyEvmGasPrice = []byte("EvmGasPrice")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{

		TxByteFee: sdk.MustNewDecFromStr("0.00008"),
		// coin transactions fees
		CoinCreate:      sdk.MustNewDecFromStr("0.004"),
		CoinUpdate:      sdk.MustNewDecFromStr("0.004"),
		CoinSend:        sdk.MustNewDecFromStr("0.0004"),
		CoinSendAdd:     sdk.MustNewDecFromStr("0.0002"),
		CoinBuy:         sdk.MustNewDecFromStr("0.004"),
		CoinSell:        sdk.MustNewDecFromStr("0.004"),
		CoinRedeemCheck: sdk.MustNewDecFromStr("0.0012"),
		CoinBurn:        sdk.MustNewDecFromStr("0.0004"),
		// special commission depends on coin symbol length
		CoinCreateTicker3: sdk.NewDec(100_000),
		CoinCreateTicker4: sdk.NewDec(10_000),
		CoinCreateTicker5: sdk.NewDec(1_000),
		CoinCreateTicker6: sdk.NewDec(100),
		CoinCreateTicker7: sdk.NewDec(10),
		// multisignature wallets
		MultisigCreateWallet:      sdk.MustNewDecFromStr("0.004"),
		MultisigCreateTransaction: sdk.MustNewDecFromStr("0.004"),
		MultisigSignTransaction:   sdk.MustNewDecFromStr("0.004"),
		// nft
		NftMintToken:     sdk.MustNewDecFromStr("0.004"),
		NftUpdateToken:   sdk.MustNewDecFromStr("0.004"),
		NftUpdateReserve: sdk.MustNewDecFromStr("0.004"),
		NftSendToken:     sdk.MustNewDecFromStr("0.0004"),
		NftBurnToken:     sdk.MustNewDecFromStr("0.0004"),
		// swap
		SwapActivateChain:   sdk.MustNewDecFromStr("0.04"),
		SwapDeactivateChain: sdk.MustNewDecFromStr("0.004"),
		SwapInitialize:      sdk.MustNewDecFromStr("0.004"),
		SwapRedeem:          sdk.MustNewDecFromStr("0.0012"),
		// validator operations
		ValidatorCreateValidator: sdk.MustNewDecFromStr("0.04"),
		ValidatorEditValidator:   sdk.MustNewDecFromStr("0.04"),
		ValidatorDelegate:        sdk.MustNewDecFromStr("0.08"),
		ValidatorDelegateNFT:     sdk.MustNewDecFromStr("0.08"),
		ValidatorRedelegate:      sdk.MustNewDecFromStr("0.08"),
		ValidatorRedelegateNFT:   sdk.MustNewDecFromStr("0.08"),
		ValidatorUndelegate:      sdk.MustNewDecFromStr("0.08"),
		ValidatorUndelegateNFT:   sdk.MustNewDecFromStr("0.08"),
		ValidatorSetOnline:       sdk.MustNewDecFromStr("0.04"),
		ValidatorSetOffline:      sdk.MustNewDecFromStr("0.04"),
		//
		CommissionBurnFactor: sdk.MustNewDecFromStr("0.5"),
		// oracle
		// NOTE: default address is []byte{0}
		Oracle: "d01gczphl4h9aqrzy237jfm97elu66dam2wtn9kg8",
		// evm min gas price in usd*10^-18
		EvmGasPrice: sdk.MustNewDecFromStr("0.000190476"),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(PSKeyTxByteFee, &p.TxByteFee, validateDec),
		// coin transactions fees
		paramtypes.NewParamSetPair(PSKeyCoinCreate, &p.CoinCreate, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinUpdate, &p.CoinUpdate, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinSend, &p.CoinSend, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinSendAdd, &p.CoinSendAdd, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinBuy, &p.CoinBuy, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinSell, &p.CoinSell, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinRedeemCheck, &p.CoinRedeemCheck, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinBurn, &p.CoinBurn, validateDec),
		// special commission depends on coin symbol length
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker3, &p.CoinCreateTicker3, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker4, &p.CoinCreateTicker4, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker5, &p.CoinCreateTicker5, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker6, &p.CoinCreateTicker6, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker7, &p.CoinCreateTicker7, validateDec),
		// multisignature wallets
		paramtypes.NewParamSetPair(PSKeyMultisigCreateWallet, &p.MultisigCreateWallet, validateDec),
		paramtypes.NewParamSetPair(PSKeyMultisigCreateTransaction, &p.MultisigCreateTransaction, validateDec),
		paramtypes.NewParamSetPair(PSKeyMultisigSignTransaction, &p.MultisigSignTransaction, validateDec),
		// nft
		paramtypes.NewParamSetPair(PSKeyNftMintToken, &p.NftMintToken, validateDec),
		paramtypes.NewParamSetPair(PSKeyNftUpdateToken, &p.NftUpdateToken, validateDec),
		paramtypes.NewParamSetPair(PSKeyNftUpdateReserve, &p.NftUpdateReserve, validateDec),
		paramtypes.NewParamSetPair(PSKeyNftSendToken, &p.NftSendToken, validateDec),
		paramtypes.NewParamSetPair(PSKeyNftBurnToken, &p.NftBurnToken, validateDec),
		// swap
		paramtypes.NewParamSetPair(PSKeySwapActivateChain, &p.SwapActivateChain, validateDec),
		paramtypes.NewParamSetPair(PSKeySwapDeactivateChain, &p.SwapDeactivateChain, validateDec),
		paramtypes.NewParamSetPair(PSKeySwapInitialize, &p.SwapInitialize, validateDec),
		paramtypes.NewParamSetPair(PSKeySwapRedeem, &p.SwapRedeem, validateDec),
		// validator operations
		paramtypes.NewParamSetPair(PSKeyValidatorCreateValidator, &p.ValidatorCreateValidator, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorEditValidator, &p.ValidatorEditValidator, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorDelegate, &p.ValidatorDelegate, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorDelegateNFT, &p.ValidatorDelegateNFT, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorRedelegate, &p.ValidatorRedelegate, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorRedelegateNFT, &p.ValidatorRedelegateNFT, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorUndelegate, &p.ValidatorUndelegate, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorUndelegateNFT, &p.ValidatorUndelegateNFT, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorSetOnline, &p.ValidatorSetOnline, validateDec),
		paramtypes.NewParamSetPair(PSKeyValidatorSetOffline, &p.ValidatorSetOffline, validateDec),
		// burn factor
		paramtypes.NewParamSetPair(PSKeyCommissionBurnFactor, &p.CommissionBurnFactor, validateLimitDec),
		// oracle
		paramtypes.NewParamSetPair(PSKeyOracle, &p.Oracle, validateAddress),
		// evm
		paramtypes.NewParamSetPair(PSKeyEvmGasPrice, &p.EvmGasPrice, validateDec),
	}
}

// Validate validates the set of params.
func (p *Params) Validate() error {
	if _, err := sdk.AccAddressFromBech32(p.Oracle); err != nil {
		return err
	}
	// all parameters are uint64, i.e. >= 0
	// and currently there is no limits
	return nil
}

func validateAddress(i interface{}) error {
	addr, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := sdk.AccAddressFromBech32(addr); err != nil {
		return fmt.Errorf("invalid address '%s': %s", addr, err.Error())
	}
	return nil
}

func validateDec(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("negative fee: %s", v.String())
	}
	return nil
}

func validateLimitDec(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("negative value: %s", v.String())
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("value over 1.0: %s", v.String())
	}
	return nil
}

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
	// special commission depends on coin symbol length
	PSKeyCoinCreateTicker_3 = []byte("CoinCreateTicker_3")
	PSKeyCoinCreateTicker_4 = []byte("CoinCreateTicker_4")
	PSKeyCoinCreateTicker_5 = []byte("CoinCreateTicker_5")
	PSKeyCoinCreateTicker_6 = []byte("CoinCreateTicker_6")
	PSKeyCoinCreateTicker_7 = []byte("CoinCreateTicker_7")
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
	//
	PSKeyOracle = []byte("Oracle")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		TxByteFee: sdk.NewDec(2),
		// coin transactions fees
		CoinCreate:      sdk.NewDec(100),
		CoinUpdate:      sdk.ZeroDec(),
		CoinSend:        sdk.NewDec(10),
		CoinSendAdd:     sdk.NewDec(5),
		CoinBuy:         sdk.NewDec(100),
		CoinSell:        sdk.NewDec(100),
		CoinRedeemCheck: sdk.ZeroDec(),
		// special commission depends on coin symbol length
		CoinCreateTicker_3: sdk.NewDec(1_000_000),
		CoinCreateTicker_4: sdk.NewDec(100_000),
		CoinCreateTicker_5: sdk.NewDec(10_000),
		CoinCreateTicker_6: sdk.NewDec(1_000),
		CoinCreateTicker_7: sdk.NewDec(100),
		// multisignature wallets
		MultisigCreateWallet:      sdk.NewDec(100),
		MultisigCreateTransaction: sdk.NewDec(100),
		MultisigSignTransaction:   sdk.NewDec(100),
		// nft
		NftMintToken:     sdk.ZeroDec(),
		NftUpdateToken:   sdk.ZeroDec(),
		NftUpdateReserve: sdk.ZeroDec(),
		NftSendToken:     sdk.ZeroDec(),
		NftBurnToken:     sdk.ZeroDec(),
		// swap
		SwapActivateChain:   sdk.ZeroDec(),
		SwapDeactivateChain: sdk.ZeroDec(),
		SwapInitialize:      sdk.ZeroDec(),
		SwapRedeem:          sdk.ZeroDec(),
		// validator operations
		ValidatorCreateValidator: sdk.NewDec(10_000),
		ValidatorEditValidator:   sdk.NewDec(10_000),
		ValidatorDelegate:        sdk.NewDec(200),
		ValidatorDelegateNFT:     sdk.NewDec(200),
		ValidatorRedelegate:      sdk.NewDec(200),
		ValidatorRedelegateNFT:   sdk.NewDec(200),
		ValidatorUndelegate:      sdk.NewDec(200),
		ValidatorUndelegateNFT:   sdk.NewDec(200),
		ValidatorSetOnline:       sdk.NewDec(100),
		ValidatorSetOffline:      sdk.NewDec(100),
		// oracle
		// NOTE: default address is []byte{0}
		Oracle: "dx1qqjrdrw8",
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
		// special commission depends on coin symbol length
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker_3, &p.CoinCreateTicker_3, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker_4, &p.CoinCreateTicker_4, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker_5, &p.CoinCreateTicker_5, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker_6, &p.CoinCreateTicker_6, validateDec),
		paramtypes.NewParamSetPair(PSKeyCoinCreateTicker_7, &p.CoinCreateTicker_7, validateDec),
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
		// oracle
		paramtypes.NewParamSetPair(PSKeyOracle, &p.Oracle, validateAddress),
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

func validateUint64(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
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
	if !v.IsPositive() {
		return fmt.Errorf("negative fee: %s", v.String())
	}
	return nil
}

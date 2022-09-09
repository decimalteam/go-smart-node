package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	PSKeyByteFee = []byte("ByteFee")
	// coin transactions fees
	PSKeyCoinSend              = []byte("CoinSend")
	PSKeyCoinSendMultiAddition = []byte("CoinSendMultiAddition")
	PSKeyCoinBuy               = []byte("CoinBuy")
	PSKeyCoinSell              = []byte("CoinSell")
	// common transaction commission
	PSKeyCoinCreate = []byte("CoinCreate")
	// special commission depends on coin symbol length
	PSKeyCoinCreateLength3     = []byte("CoinCreateLength3")
	PSKeyCoinCreateLength4     = []byte("CoinCreateLength4")
	PSKeyCoinCreateLength5     = []byte("CoinCreateLength5")
	PSKeyCoinCreateLength6     = []byte("CoinCreateLength6")
	PSKeyCoinCreateLengthOther = []byte("CoinCreateLengthOther")
	// multisignature wallets
	PSKeyMultisigCreateWallet      = []byte("MultisigCreateWallet")
	PSKeyMultisigCreateTransaction = []byte("MultisigCreateTransaction")
	PSKeyMultisigSignTransaction   = []byte("MultisigSignTransaction")
	// validator operations
	PSKeyValidatorDeclareCandidate = []byte("ValidatorDeclareCandidate")
	PSKeyValidatorEditCandidate    = []byte("ValidatorEditCandidate")
	PSKeyValidatorDelegate         = []byte("ValidatorDelegate")
	PSKeyValidatorUnbond           = []byte("ValidatorUnbond")
	PSKeyValidatorSetOnline        = []byte("ValidatorSetOnline")
	PSKeyValidatorSetOffline       = []byte("ValidatorSetOffline")
	// oracle key
	PSKeyOracleAddress = []byte("OracleAddress")
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
		// byte fee in usd*10^-18
		ByteFee: sdk.MustNewDecFromStr("0.001"),
		// coin transactions fees
		// byte fee in usd*10^-3
		CoinSend:              sdk.MustNewDecFromStr("0.082"),
		CoinSendMultiAddition: sdk.MustNewDecFromStr("0.04"),
		CoinBuy:               sdk.MustNewDecFromStr("0.8"),
		CoinSell:              sdk.MustNewDecFromStr("0.8"),
		// common transaction commission
		CoinCreate: sdk.MustNewDecFromStr("0.008"), // x8
		// special commission depends on coin symbol length
		CoinCreateLength3:     sdk.MustNewDecFromStr("100000"),
		CoinCreateLength4:     sdk.MustNewDecFromStr("10000"),
		CoinCreateLength5:     sdk.MustNewDecFromStr("1000"),
		CoinCreateLength6:     sdk.MustNewDecFromStr("100"),
		CoinCreateLengthOther: sdk.MustNewDecFromStr("10"),
		// multisignature wallets
		MultisigCreateWallet:      sdk.MustNewDecFromStr("0.1"),
		MultisigCreateTransaction: sdk.MustNewDecFromStr("0.1"),
		MultisigSignTransaction:   sdk.MustNewDecFromStr("0.1"),
		// validator operations
		ValidatorDeclareCandidate: sdk.MustNewDecFromStr("10"),
		ValidatorEditCandidate:    sdk.MustNewDecFromStr("10"),
		ValidatorDelegate:         sdk.MustNewDecFromStr("0.2"),
		ValidatorUnbond:           sdk.MustNewDecFromStr("0.2"),
		ValidatorSetOnline:        sdk.MustNewDecFromStr("0.1"),
		ValidatorSetOffline:       sdk.MustNewDecFromStr("0.1"),
		// oracle
		// NOTE: default address is []byte{0}
		OracleAddress: "dx1qqjrdrw8",
		// evm min gas price in usd*10^-18
		EvmGasPrice: sdk.MustNewDecFromStr("0.000019047619047619"),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(PSKeyByteFee, &p.ByteFee, validateUint64),
		// coin transactions fees
		paramtypes.NewParamSetPair(PSKeyCoinSend, &p.CoinSend, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinSendMultiAddition, &p.CoinSendMultiAddition, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinBuy, &p.CoinBuy, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinSell, &p.CoinSell, validateUint64),
		// common transaction commission
		paramtypes.NewParamSetPair(PSKeyCoinCreate, &p.CoinCreate, validateUint64),
		// special commission depends on coin symbol length
		paramtypes.NewParamSetPair(PSKeyCoinCreateLength3, &p.CoinCreateLength3, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinCreateLength4, &p.CoinCreateLength4, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinCreateLength5, &p.CoinCreateLength5, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinCreateLength6, &p.CoinCreateLength6, validateUint64),
		paramtypes.NewParamSetPair(PSKeyCoinCreateLengthOther, &p.CoinCreateLengthOther, validateUint64),
		// multisignature wallets
		paramtypes.NewParamSetPair(PSKeyMultisigCreateWallet, &p.MultisigCreateWallet, validateUint64),
		paramtypes.NewParamSetPair(PSKeyMultisigCreateTransaction, &p.MultisigCreateTransaction, validateUint64),
		paramtypes.NewParamSetPair(PSKeyMultisigSignTransaction, &p.MultisigSignTransaction, validateUint64),
		// validator operations
		paramtypes.NewParamSetPair(PSKeyValidatorDeclareCandidate, &p.ValidatorDeclareCandidate, validateUint64),
		paramtypes.NewParamSetPair(PSKeyValidatorEditCandidate, &p.ValidatorEditCandidate, validateUint64),
		paramtypes.NewParamSetPair(PSKeyValidatorDelegate, &p.ValidatorDelegate, validateUint64),
		paramtypes.NewParamSetPair(PSKeyValidatorUnbond, &p.ValidatorUnbond, validateUint64),
		paramtypes.NewParamSetPair(PSKeyValidatorSetOnline, &p.ValidatorSetOnline, validateUint64),
		paramtypes.NewParamSetPair(PSKeyValidatorSetOffline, &p.ValidatorSetOffline, validateUint64),
		// oracle
		paramtypes.NewParamSetPair(PSKeyOracleAddress, &p.OracleAddress, validateAddress),
		// evm
		paramtypes.NewParamSetPair(PSKeyEvmGasPrice, &p.EvmGasPrice, validateUint64),
	}
}

// Validate validates the set of params.
func (p *Params) Validate() error {
	if _, err := sdk.AccAddressFromBech32(p.OracleAddress); err != nil {
		return err
	}
	// all parameters are uint64, i.e. >= 0
	// and currently there is no limits
	return nil
}

// String implements the Stringer interface.
func (p *Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
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

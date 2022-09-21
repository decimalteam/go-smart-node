package types

import (
	"encoding/hex"
	fmt "fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyLockedTimeOut   = []byte("LockedTimeOut")
	KeyLockedTimeIn    = []byte("LockedTimeIn")
	KeyServiceAddress  = []byte("ServiceAddress")
	KeyCheckingAddress = []byte("CheckingAddress")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		LockedTimeOut:   DefaultLockedTimeOut,
		LockedTimeIn:    DefaultLockedTimeIn,
		ServiceAddress:  DefaultSwapServiceAddress,
		CheckingAddress: DefaultCheckingAddress,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLockedTimeOut, &p.LockedTimeOut, validateLockedTime),
		paramtypes.NewParamSetPair(KeyLockedTimeIn, &p.LockedTimeIn, validateLockedTime),
		paramtypes.NewParamSetPair(KeyServiceAddress, &p.ServiceAddress, validateSdkAddress),
		paramtypes.NewParamSetPair(KeyCheckingAddress, &p.CheckingAddress, validateHexAddress),
	}
}

func validateLockedTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("locked time must be positive: %d", v)
	}

	return nil
}

func validateHexAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	_, err := hex.DecodeString(v)
	if err != nil {
		return err
	}
	return nil
}

func validateSdkAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := sdk.AccAddressFromBech32(v); err != nil {
		return err
	}
	return nil
}

func (p Params) Validate() error {
	return nil
}

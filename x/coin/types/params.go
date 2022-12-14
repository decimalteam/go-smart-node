package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	cmdconfig "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
)

// Coin params default values.
var (
	DefaultBaseDenom  string      = cmdconfig.BaseDenom
	DefaultBaseTitle  string      = "Decimal coin"
	DefaultBaseVolume sdkmath.Int = helpers.EtherToWei(sdkmath.NewInt(340_000_000))
)

// Parameter store keys.
var (
	KeyBaseDenom  = []byte("BaseDenom")
	KeyBaseTitle  = []byte("BaseTitle")
	KeyBaseVolume = []byte("BaseVolume")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		BaseDenom:  DefaultBaseDenom,
		BaseTitle:  DefaultBaseTitle,
		BaseVolume: DefaultBaseVolume,
	}
}

// NewParams creates a new Params instance.
func NewParams(baseDenom string, baseTitle string, baseVolume sdkmath.Int) Params {
	return Params{
		BaseDenom:  baseDenom,
		BaseTitle:  baseTitle,
		BaseVolume: baseVolume,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateBaseDenom),
		paramtypes.NewParamSetPair(KeyBaseTitle, &p.BaseTitle, validateBaseTitle),
		paramtypes.NewParamSetPair(KeyBaseVolume, &p.BaseVolume, validateBaseVolume),
	}
}

// Validate validates the set of params.
func (p *Params) Validate() (err error) {
	if err = validateBaseDenom(p.BaseDenom); err != nil {
		return
	}
	if err = validateBaseTitle(p.BaseTitle); err != nil {
		return
	}
	if err = validateBaseVolume(p.BaseVolume); err != nil {
		return
	}
	return
}

func validateBaseDenom(i interface{}) error {
	denom, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !config.CoinDenomValidator.MatchString(denom) {
		return errors.InvalidCoinDenom
	}
	// TODO
	return nil
}

func validateBaseTitle(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}

func validateBaseVolume(i interface{}) error {
	_, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}

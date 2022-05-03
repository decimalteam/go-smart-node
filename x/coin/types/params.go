package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	ParamStoreKeyBaseTitle         = []byte("BaseTitle")
	ParamStoreKeyBaseSymbol        = []byte("BaseSymbol")
	ParamStoreKeyBaseInitialVolume = []byte("BaseInitialVolume")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		BaseTitle:         "Decimal coin",
		BaseSymbol:        "del",
		BaseInitialVolume: sdk.NewInt(0),
	}
}

// NewParams creates a new Params instance.
func NewParams(
	baseTitle string,
	baseSymbol string,
	baseInitialVolume sdk.Int,
) Params {
	return Params{
		BaseTitle:         baseTitle,
		BaseSymbol:        baseSymbol,
		BaseInitialVolume: baseInitialVolume,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyBaseTitle, &p.BaseTitle, validateBaseTitle),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseSymbol, &p.BaseSymbol, validateBaseSymbol),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseInitialVolume, &p.BaseInitialVolume, validateBaseInitialVolume),
	}
}

// Validate validates the set of params.
func (p *Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p *Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateBaseTitle(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}

func validateBaseSymbol(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}

func validateBaseInitialVolume(i interface{}) error {
	_, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}

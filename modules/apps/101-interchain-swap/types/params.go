package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultSwapEnabled = true
	// DefaultMaxFeeRate is 0.003
	DefaultMaxFeeRate = 300
)

var (
	KeySwapEnabled    = []byte("SwapEnabled")
	KeySwapMaxFeeRate = []byte("MaxFeeRate")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(enable bool, feeRate uint32) Params {
	return Params{
		SwapEnabled: enable,
		MaxFeeRate:  feeRate,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultSwapEnabled, DefaultMaxFeeRate)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySwapEnabled, p.SwapEnabled, validateEnabled),
		paramtypes.NewParamSetPair(KeySwapMaxFeeRate, p.MaxFeeRate, validateMaxFeeRate),
	}
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxFeeRate(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

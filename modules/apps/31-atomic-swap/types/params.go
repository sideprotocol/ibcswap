package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultSwapEnabled = true
	// DefaultMaxFeeRate is 0.0010
	DefaultMaxFeeRate = 10
)

var (
	KeySwapEnabled    = []byte("SwapEnabled")
	KeySwapMaxFeeRate = []byte("MaxFeeRate")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the ibc transfer module
func NewParams(enable bool, feeRate uint32) Params {
	return Params{
		SwapEnabled: enable,
		MaxFeeRate:  feeRate,
	}
}

// DefaultParams is the default parameter configuration for the ibc-transfer module
func DefaultParams() Params {
	return NewParams(DefaultSwapEnabled, DefaultMaxFeeRate)
}

// Validate all ibc-swap module parameters
func (p Params) Validate() error {

	if err := validateMaxFeeRate(p.MaxFeeRate); err != nil {
		return err
	}
	return validateEnabled(p.SwapEnabled)
}

// ParamSetPairs implements params.ParamSet
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

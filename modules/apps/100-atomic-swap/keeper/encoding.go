package keeper

import (
	"github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"
)

// UnmarshalLimitOrder attempts to decode and return an LimitOrder object from
// raw encoded bytes.
func (k Keeper) UnmarshalLimitOrder(bz []byte) (types.LimitOrder, error) {
	var order types.LimitOrder
	if err := k.cdc.Unmarshal(bz, &order); err != nil {
		return types.LimitOrder{}, err
	}

	return order, nil
}

// MustUnmarshalLimitOrder attempts to decode and return an LimitOrder object from
// raw encoded bytes. It panics on error.
func (k Keeper) MustUnmarshalLimitOrder(bz []byte) types.LimitOrder {
	var order types.LimitOrder
	k.cdc.MustUnmarshal(bz, &order)
	return order
}

// MarshalLimitOrder attempts to encode an LimitOrder object and returns the
// raw encoded bytes.
func (k Keeper) MarshalLimitOrder(order types.LimitOrder) ([]byte, error) {
	return k.cdc.Marshal(&order)
}

// MustMarshalLimitOrder attempts to encode an LimitOrder object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalLimitOrder(order types.LimitOrder) []byte {
	return k.cdc.MustMarshal(&order)
}

/// OTC Order

// UnmarshalOTCOrder attempts to decode and return an OTCOrder object from
// raw encoded bytes.
func (k Keeper) UnmarshalOTCOrder(bz []byte) (types.OTCOrder, error) {
	var order types.OTCOrder
	if err := k.cdc.Unmarshal(bz, &order); err != nil {
		return types.OTCOrder{}, err
	}

	return order, nil
}

// MustUnmarshalOTCOrder attempts to decode and return an OTCOrder object from
// raw encoded bytes. It panics on error.
func (k Keeper) MustUnmarshalOTCOrder(bz []byte) types.OTCOrder {
	var order types.OTCOrder
	k.cdc.MustUnmarshal(bz, &order)
	return order
}

// MarshalOTCOrder attempts to encode an OTCOrder object and returns the
// raw encoded bytes.
func (k Keeper) MarshalOTCOrder(order types.OTCOrder) ([]byte, error) {
	return k.cdc.Marshal(&order)
}

// MustMarshalOTCOrder attempts to encode an OTCOrder object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalOTCOrder(order types.OTCOrder) []byte {
	return k.cdc.MustMarshal(&order)
}

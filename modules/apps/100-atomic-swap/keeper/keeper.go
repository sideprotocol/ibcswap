package keeper

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the IBC Swap keeper
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	ics4Wrapper   types.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ics4Wrapper types.ICS4Wrapper, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	authKeeper types.AccountKeeper, bankKeeper types.BankKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	// ensure ibc transfer module account is set
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the IBC swap module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace,
		ics4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}

// IsBound checks if the transfer module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the swap module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the swap module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the swap module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// GetLimitOrder returns the LimitOrder for the swap module.
func (k Keeper) GetLimitOrder(ctx sdk.Context, orderId string) (types.LimitOrder, bool) {
	key, err := hex.DecodeString(orderId)
	if err != nil {
		return types.LimitOrder{}, false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LimitOrderBookKey)
	bz := store.Get(key)
	if bz == nil {
		return types.LimitOrder{}, false
	}

	order := k.MustUnmarshalLimitOrder(bz)
	return order, true
}

// HasLimitOrder checks if a the key with the given id exists on the store.
func (k Keeper) HasLimitOrder(ctx sdk.Context, orderId string) bool {
	key, err := hex.DecodeString(orderId)
	if err != nil {
		return false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LimitOrderBookKey)
	return store.Has(key)
}

// SetLimitOrder sets a new LimitOrder to the store.
func (k Keeper) SetLimitOrder(ctx sdk.Context, order types.LimitOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LimitOrderBookKey)
	bz := k.MustMarshalLimitOrder(order)
	store.Set([]byte(order.Id), bz)
}

// IterateLimitOrders iterates over the limit orders in the store
// and performs a callback function.
func (k Keeper) IterateLimitOrders(ctx sdk.Context, cb func(order types.LimitOrder) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LimitOrderBookKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		order := k.MustUnmarshalLimitOrder(iterator.Value())
		if cb(order) {
			break
		}
	}
}

// GetAllLimitOrders returns the information for all the limit orders.
func (k Keeper) GetAllLimitOrders(ctx sdk.Context) []types.LimitOrder {
	var orders []types.LimitOrder
	k.IterateLimitOrders(ctx, func(order types.LimitOrder) bool {
		orders = append(orders, order)
		return false
	})

	return orders
}

/// OTC orders

// GetOTCOrder returns the OTCOrder for the swap module.
func (k Keeper) GetOTCOrder(ctx sdk.Context, orderId string) (types.OTCOrder, bool) {
	key, err := hex.DecodeString(orderId)
	if err != nil {
		return types.OTCOrder{}, false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	bz := store.Get(key)
	if bz == nil {
		return types.OTCOrder{}, false
	}

	order := k.MustUnmarshalOTCOrder(bz)
	return order, true
}

// HasOTCOrder checks if a the key with the given id exists on the store.
func (k Keeper) HasOTCOrder(ctx sdk.Context, orderId string) bool {
	key, err := hex.DecodeString(orderId)
	if err != nil {
		return false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	return store.Has(key)
}

// SetOTCOrder sets a new OTCOrder to the store.
func (k Keeper) SetOTCOrder(ctx sdk.Context, order types.OTCOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	bz := k.MustMarshalOTCOrder(order)
	store.Set([]byte(order.Id), bz)
}

// IterateOTCOrders iterates over the limit orders in the store
// and performs a callback function.
func (k Keeper) IterateOTCOrders(ctx sdk.Context, cb func(order types.OTCOrder) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OTCOrderBookKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		order := k.MustUnmarshalOTCOrder(iterator.Value())
		if cb(order) {
			break
		}
	}
}

// GetAllOTCOrders returns the information for all the limit orders.
func (k Keeper) GetAllOTCOrders(ctx sdk.Context) []types.OTCOrder {
	var orders []types.OTCOrder
	k.IterateOTCOrders(ctx, func(order types.OTCOrder) bool {
		orders = append(orders, order)
		return false
	})

	return orders
}

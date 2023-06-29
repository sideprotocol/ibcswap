package keeper

import (
	"encoding/hex"

	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the IBC Swap keeper
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	ics4Wrapper   porttypes.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ics4Wrapper porttypes.ICS4Wrapper, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
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

/// Atomic orders

// GetAtomicOrder returns the OTCOrder for the swap module.
func (k Keeper) GetAtomicOrder(ctx sdk.Context, orderId string) (types.Order, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	bz := store.Get([]byte(orderId))
	if bz == nil {
		return types.Order{}, false
	}

	order := k.MustUnmarshalOrder(bz)
	return order, true
}

// HasAtomicOrder checks if a the key with the given id exists on the store.
func (k Keeper) HasAtomicOrder(ctx sdk.Context, orderId string) bool {
	key, err := hex.DecodeString(orderId)
	if err != nil {
		return false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	return store.Has(key)
}

// SetAtomicOrder sets a new OTCOrder to the store.
func (k Keeper) SetAtomicOrder(ctx sdk.Context, order types.Order) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	bz := k.MustMarshalOrder(order)
	store.Set([]byte(order.Id), bz)
}

// // IterateAtomicOrders iterates over the limit orders in the store
// // and performs a callback function.
// func (k Keeper) IterateAtomicOrders(ctx sdk.Context, req *types.QueryOrdersRequest, cb func(order types.Order) bool) {
// 	store := ctx.KVStore(k.storeKey)
// 	iterator := sdk.KVStorePrefixIterator(store, types.OTCOrderBookKey)

// 	defer iterator.Close()
// 	for ; iterator.Valid(); iterator.Next() {
// 		order := k.MustUnmarshalOrder(iterator.Value())
// 		if cb(order) {
// 			break
// 		}
// 	}
// }

// // GetAllAtomicOrders returns the information for all the limit orders.
// func (k Keeper) GetAllAtomicOrders(ctx sdk.Context, req *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {

// 	orderStore := ctx.KVStore(k.storeKey)
// 	iterator := sdk.KVStorePrefixIterator(orderStore, types.OTCOrderBookKey)
// 	var orders []*types.Order
// 	pageRes, err := query.Paginate(orderStore, req.Pagination, func(key, value []byte) error {
// 		order := k.MustUnmarshalOrder(iterator.Value())
// 		orders = append(orders, &order)
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
// 	}

// 	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, err
// }

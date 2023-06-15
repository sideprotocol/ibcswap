package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/tendermint/tendermint/libs/log"

	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		paramstore    paramtypes.Subspace
		ics4Wrapper   porttypes.ICS4Wrapper
		channelKeeper types.ChannelKeeper
		portKeeper    types.PortKeeper
		scopedKeeper  capabilitykeeper.ScopedKeeper
		bankKeeper    types.BankKeeper
		authKeeper    types.AccountKeeper
		msgRouter     types.MessageRouter
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ics4Wrapper porttypes.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	bankKeeper types.BankKeeper,
	authKeeper types.AccountKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	msgRouter types.MessageRouter,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		ics4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
		bankKeeper:    bankKeeper,
		authKeeper:    authKeeper,
		msgRouter:     msgRouter,
	}
}

// ----------------------------------------------------------------------------
// IBC Keeper Logic
// ----------------------------------------------------------------------------

// ChanCloseInit defines a wrapper function for the channel Keeper's function.
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// IsBound checks if the IBC app module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the IBC app module to claim a capability that core IBC
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetBalance(ctx sdk.Context, sender sdk.AccAddress) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, sender, sdk.DefaultBondDenom)
}

func (k Keeper) SingleDepositTest(ctx sdk.Context, sender sdk.AccAddress) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, sender, sdk.DefaultBondDenom)
}

// You may need to adjust the function signature, return types, and parameter types based on your module's implementation
func (k Keeper) EscrowAddress(ctx context.Context, req *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error) {
	escrowAddress := types.GetEscrowAddress(req.PortId, req.ChannelId)
	return &types.QueryEscrowAddressResponse{
		EscrowAddress: escrowAddress.String(),
	}, nil
}

// You may need to adjust the function signature, return types, and parameter types based on your module's implementation
func (k Keeper) MultiDepositOrders(ctx context.Context, req *types.QueryMultiDepositOrdersRequest) (*types.QueryMultiDepositOrdersResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)
    orders := k.GetAllMultiDepositOrder(sdkCtx, req.PoolId)
    
    ordersPtr := make([]*types.MultiAssetDepositOrder, len(orders))
    for i := range orders {
        ordersPtr[i] = &orders[i]
    }
    
    return &types.QueryMultiDepositOrdersResponse{
        Orders: ordersPtr,
    },nil
}

func (k Keeper) validateCoins(ctx sdk.Context, pool *types.InterchainLiquidityPool, sender string, tokensIn []*sdk.Coin) ([]sdk.Coin, error) {
	// Deposit token to Escrow account
	coins := []sdk.Coin{}
	for _, coin := range tokensIn {
		accAddress := sdk.MustAccAddressFromBech32(sender)
		balance := k.bankKeeper.GetBalance(ctx, accAddress, coin.Denom)
		if balance.Amount.Equal(sdk.NewInt(0)) {
			return nil, types.ErrInvalidAmount
		}
		coins = append(coins, *coin)
		if pool.Status == types.PoolStatus_INITIALIZED {
			poolAsset, err := pool.FindAssetByDenom(coin.Denom)
			if err == nil {
				if !poolAsset.Balance.Amount.Equal(coin.Amount) {
					return nil, types.ErrInvalidInitialDeposit
				}
			}
		}
	}
	return coins, nil
}

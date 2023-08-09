package types

// IBC events
const (
	EventTypeTimeout                 = "timeout"
	EventTypePacket                  = "interchain_ibc_packet"
	EventTypeMakePool                = "make_pool"
	EventTypeCancelPool              = "cancel_pool"
	EventTypeTakePool                = "take_pool"
	EventTypeMakeMultiDepositOrder   = "make_multi_asset_deposit"
	EventTypeCancelMultiDepositOrder = "cancel_multi_deposit_order"
	EventTypeTakeMultiDepositOrder   = "take_multi_asset_deposit"
	EventTypeSingleDepositOrder      = "single_deposit"
	EventTypeLiquidityWithdraw       = "liquidity_withdraw"
	EventTypeSwap                    = "swap_assets"
	EventTypeIBCStep                 = ""

	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"

	AttributeKeyPoolId              = "pool_id"
	AttributeKeyMultiDepositOrderId = "order_id"
	AttributeKeyTokenIn             = "token_in"
	AttributeKeyTokenOut            = "token_out"
	AttributeKeyLpToken             = "liquidity_pool_token"
	AttributeKeyLpSupply            = "Liquidity_pool_token_supply"
	AttributeIBCStep                = "ibc_step"
)

const (
	ON_RECEIVE  = "on_receive"
	ACKNOWLEDGE = "acknowledge"
)

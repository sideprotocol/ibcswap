package types

// IBC events
const (
	EventTypeTimeout               = "timeout"
	EventTypePacket                = "interchain_ibc_packet"
	EventTypeMakePool              = "make_pool"
	EventTypeTakePool              = "take_pool"
	EventTypeMakeMultiDepositOrder = "make_multi_asset_deposit"
	EventTypeTakeMultiDepositOrder = "take_multi_asset_deposit"
	EventTypeSingleDepositOrder    = "single_deposit"
	EventTypeLiquidityWithdraw     = "liquidity_withdraw"
	EventTypeSwap                  = "swap_assets"

	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"

	AttributeKeyPoolId              = "poolId"
	AttributeKeyMultiDepositOrderId = "orderId"
	AttributeKeyTokenIn             = "tokenIn"
	AttributeKeyTokenOut            = "tokenOut"
	AttributeKeyLpToken             = "liquidity_pool_token"
	AttributeKeyLpSupply            = "Liquidity_pool_token_supply"
)

package types

// IBC events
const (
	EventTypeTimeout = "timeout"
	EventTypePacket  = "interchain_ibc_packet"

	EventTypeSingleDepositOrder = "single_deposit"
	EventTypeLiquidityWithdraw  = "liquidity_withdraw"
	EventTypeSwap               = "swap_assets"
	EventTypeIBCStep            = ""

	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"

	AttributeKeyAction              = "action"
	AttributeKeyPoolId              = "pool_id"
	AttributeKeyMultiDepositOrderId = "order_id"
	AttributeKeyMultiDeposits       = "deposit"
	AttributeKeyTokenIn             = "token_in"
	AttributeKeyTokenOut            = "token_out"
	AttributeKeyLpToken             = "liquidity_pool_token"
	AttributeKeyLpSupply            = "liquidity_pool_token_supply"
	AttributeKeyPoolCreator         = "pool_creator"
	AttributeKeyOrderCreator        = "order_creator"
	AttributeKeyName                = "name"
	AttributeKeyPoolStatus          = "pool_status"
	AttributeKeyMsgSender           = "msg_sender"
)

const (
	EventValueActionMakePool   = "make_pool"
	EventValueActionTakePool   = "take_pool"
	EventValueActionCancelPool = "cancel_pool"

	EventValueActionMakeOrder            = "make_order"
	EventValueActionTakeOrder            = "take_order"
	EventValueActionCancelOrder          = "cancel_order"
	EventValueActionSingleDeposit        = "single_deposit"
	EventValueActionMakeMultiDeposit     = "make_multi_deposit_order"
	EventValueActionTakeMultiDeposit     = "take_multi_deposit_order"
	EventValueActionCancelMultiDeposit   = "cancel_multi_deposit_order"
	EventValueActionWithdrawMultiDeposit = "withdraw_multi_deposit_order"
	EventValueActionSwap                 = "swap"
	EventOwner                           = "interchain_swap"
)

const (
	EventValueSuffixReceived     = "received"
	EventValueSuffixAcknowledged = "acknowledged"
)

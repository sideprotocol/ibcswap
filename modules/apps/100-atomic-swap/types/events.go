package types

// IBC transfer events
const (
	EventTypeTimeout      = "timeout"
	EventTypePacket       = "fungible_token_packet"
	EventTypeSwap         = "ibc_swap"
	EventTypeChannelClose = "channel_closed"
	EventTypeMakeSwap     = "make_swap"
	EventTypeTakeSwap     = "take_swap"
	EventTypeCancelSwap   = "cancel_swap"

	AttributeType          = "type"
	AttributeData          = "data"
	AttributeMemo          = "memo"
	AttributeKeyAmount     = "amount"
	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
	AttributeOrderId       = "order_id"
	AttributeAction        = "action"
	AttributeName          = "name"
)

const (
	EventValueActionMakeOrder   = "make_order"
	EventValueActionTakeOrder   = "take_order"
	EventValueActionCancelOrder = "cancel_order"
	EventOwner                  = "atomic_swap"
)

const (
	EventValueSuffixReceived     = "received"
	EventValueSuffixAcknowledged = "acknowledged"
)

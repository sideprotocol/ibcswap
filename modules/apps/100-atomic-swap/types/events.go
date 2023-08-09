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
	AttributeIBCStep       = "ibc_step"
)

const (
	ON_RECEIVE  = "on_receive"
	ACKNOWLEDGE = "acknowledge"
)

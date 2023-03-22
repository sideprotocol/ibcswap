package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// msg types
const (
	TypeMsgMakeSwap   = "make_swap"
	TypeMsgTakeSwap   = "take_swap"
	TypeMsgCancelSwap = "cancel_swap"
)

// NewMsgMakeSwap creates a new MsgMakeSwapRequest instance
func NewMsgMakeSwap(
	sourcePort, sourceChannel string,
	sellToken, buyToken sdk.Coin,
	senderAddress, senderReceivingAddress, desiredTaker string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
	createdTimestamp int64,
) *MakeSwapMsg {
	return &MakeSwapMsg{
		SourcePort:            sourcePort,
		SourceChannel:         sourceChannel,
		SellToken:             sellToken,
		BuyToken:              buyToken,
		MakerAddress:          senderAddress,
		MakerReceivingAddress: senderReceivingAddress,
		DesiredTaker:          desiredTaker,
		TimeoutHeight:         timeoutHeight,
		TimeoutTimestamp:      timeoutTimestamp,
		CreateTimestamp:       createdTimestamp,
	}
}

// Route implements sdk.Msg
func (*MakeSwapMsg) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MakeSwapMsg) Type() string {
	return TypeMsgMakeSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg *MakeSwapMsg) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if !msg.SellToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.SellToken.String())
	}
	if !msg.SellToken.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.SellToken.String())
	}
	if !msg.BuyToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.BuyToken.String())
	}
	if !msg.BuyToken.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.BuyToken.String())
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if strings.TrimSpace(msg.MakerReceivingAddress) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	// return ValidateIBCDenom(msg.SendingToken.Denom)
	return nil
}

// GetSignBytes implements sdk.Msg.
func (msg *MakeSwapMsg) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MakeSwapMsg) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgTakeSwap creates a new MsgTakeSwapRequest instance
func NewMsgTakeSwap(
	orderId string,
	sellToken sdk.Coin,
	senderAddress, senderReceivingAddress string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	createdTimestamp int64,
) *TakeSwapMsg {
	return &TakeSwapMsg{
		OrderId:               orderId,
		SellToken:             sellToken,
		TakerAddress:          senderAddress,
		TakerReceivingAddress: senderReceivingAddress,
		TimeoutHeight:         timeoutHeight,
		TimeoutTimestamp:      timeoutTimestamp,
		CreateTimestamp:       createdTimestamp,
	}
}

// Route implements sdk.Msg
func (*TakeSwapMsg) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*TakeSwapMsg) Type() string {
	return TypeMsgTakeSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg *TakeSwapMsg) ValidateBasic() error {
	//if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
	//	return sdkerrors.Wrap(err, "invalid source port ID")
	//}
	//if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
	//	return sdkerrors.Wrap(err, "invalid source channel ID")
	//}
	if len(msg.OrderId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, msg.OrderId)
	}
	if !msg.SellToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.SellToken.String())
	}
	if !msg.SellToken.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.SellToken.String())
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.TakerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if strings.TrimSpace(msg.TakerReceivingAddress) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (msg *TakeSwapMsg) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *TakeSwapMsg) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.TakerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgCancelSwap creates a new MsgCancelSwapRequest instance
func NewMsgCancelSwap(
	senderAddress, orderId string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *CancelSwapMsg {
	return &CancelSwapMsg{
		MakerAddress:     senderAddress,
		OrderId:          orderId,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*CancelSwapMsg) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*CancelSwapMsg) Type() string {
	return TypeMsgCancelSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg *CancelSwapMsg) ValidateBasic() error {
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if strings.TrimSpace(msg.OrderId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "OrderId is required")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (msg *CancelSwapMsg) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *CancelSwapMsg) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
)

// msg types
const (
	TypeMsgAtomicSwap = "swap"
)

// NewMsgSwap creates a new MsgSwap instance
//nolint:interfacer
func NewMsgSwap(
	sourcePort, sourceChannel string,
	sendingToken, receivingToken sdk.Coin,
	senderAddress, senderReceivingAddress, expectedCounterpartyAddress string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgSwap {
	return &MsgSwap{
		SourcePort:                  sourcePort,
		SourceChannel:               sourceChannel,
		SendingToken:                sendingToken,
		ReceivingToken:              receivingToken,
		SenderAddress:               senderAddress,
		SenderReceivingAddress:      senderReceivingAddress,
		ExpectedCounterpartyAddress: expectedCounterpartyAddress,
		TimeoutHeight:               timeoutHeight,
		TimeoutTimestamp:            timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (MsgSwap) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgSwap) Type() string {
	return TypeMsgAtomicSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg MsgSwap) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if !msg.SendingToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.SendingToken.String())
	}
	if !msg.SendingToken.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.SendingToken.String())
	}
	if !msg.ReceivingToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.ReceivingToken.String())
	}
	if !msg.ReceivingToken.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.ReceivingToken.String())
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.SenderAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if strings.TrimSpace(msg.SenderReceivingAddress) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	// return ValidateIBCDenom(msg.SendingToken.Denom)
	return nil
}

// GetSignBytes implements sdk.Msg.
func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.SenderAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

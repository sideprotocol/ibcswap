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
	TypeMsgCreatPool     = "create_pool"
	TypeMsgSingleDeposit = "single_deposit"
	TypeMsgWithdraw      = "withdraw"
	TypeMsgLeftSwap      = "left_swap"
	TypeMsgRightSwap     = "right_swap"
)

// NewMsgCreatePoolRequest creates a new MsgCreatePoolRequest instance
func NewMsgCreatePoolRequest(
	sourcePort, sourceChannel string,
	sender string,
	denoms []string,
	decimals []uint32,
	weight string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgCreatePoolRequest {
	return &MsgCreatePoolRequest{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Sender:           sender,
		Denoms:           denoms,
		Decimals:         decimals,
		Weight:           weight,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*MsgCreatePoolRequest) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MsgCreatePoolRequest) Type() string {
	return TypeMsgCreatPool
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
func (m *MsgCreatePoolRequest) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if len(m.Denoms) == 0 || len(m.Decimals) == 0 {
		return ErrInvalidPairLength
	}

	// length of weight, denominations, decimals should be equal
	if len(m.Denoms) != len(m.Decimals) {
		return ErrInvalidPairLength
	}

	if _, err := ParseWeight(m.Weight, len(m.Denoms)); err != nil {
		return err
	}

	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}

	return nil
}

// GetSignBytes implements sdk.Msg.
func (m *MsgCreatePoolRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgCreatePoolRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgSingleDepositRequest creates a new MsgSingleDepositRequest instance
func NewMsgSingleDepositRequest(
	sourcePort, sourceChannel string,
	tokens []*sdk.Coin,
	sender string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgSingleDepositRequest {
	return &MsgSingleDepositRequest{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Tokens:           tokens,
		Sender:           sender,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*MsgSingleDepositRequest) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MsgSingleDepositRequest) Type() string {
	return TypeMsgSingleDeposit
}

// ValidateBasic performs a basic check of the MsgDepositRequest fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
func (m *MsgSingleDepositRequest) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}

	length := len(m.Tokens)
	if length == 0 {
		return ErrInvalidPairLength
	}

	for i := 0; i < length; i++ {
		if !m.Tokens[i].IsValid() {
			return ErrInvalidToken
		}
	}

	if len(m.PoolId) == 0 {
		return ErrInvalidPoolId
	}

	return nil
}

// GetSignBytes implements sdk.Msg.
func (m *MsgSingleDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgSingleDepositRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgWithdrawRequest creates a new MsgWithdrawRequest instance
func NewMsgWithdrawRequest(
	sourcePort, sourceChannel string,
	sender string,
	poolToken *sdk.Coin,
	denomOut string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Sender:           sender,
		PoolToken:        poolToken,
		DenomOut:         denomOut,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*MsgWithdrawRequest) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MsgWithdrawRequest) Type() string {
	return TypeMsgWithdraw
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
func (m *MsgWithdrawRequest) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if !m.PoolToken.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool token is invalid")
	}
	if strings.TrimSpace(m.DenomOut) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "DenomOut is required")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (m *MsgWithdrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgLeftSwapRequest creates a new MsgLeftSwapRequest instance
func NewMsgLeftSwapRequest(
	sourcePort, sourceChannel string,
	sender string,
	tokenIn *sdk.Coin,
	denomOut string,
	slippage uint32,
	recipient string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgLeftSwapRequest {
	return &MsgLeftSwapRequest{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Sender:           sender,
		TokenIn:          tokenIn,
		DenomOut:         denomOut,
		Slippage:         slippage,
		Recipient:        recipient,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*MsgLeftSwapRequest) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MsgLeftSwapRequest) Type() string {
	return TypeMsgLeftSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
func (m *MsgLeftSwapRequest) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if !m.TokenIn.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool token is invalid")
	}
	if m.Slippage <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "maximum slippage is invalid")
	}
	if strings.TrimSpace(m.DenomOut) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "DenomOut is required")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (m *MsgLeftSwapRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgLeftSwapRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgRightSwapRequest creates a new MsgRightSwapRequest instance
func NewMsgRightSwapRequest(
	sourcePort, sourceChannel string,
	sender string,
	tokenIn *sdk.Coin,
	tokenOut *sdk.Coin,
	slippage uint32,
	recipient string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgRightSwapRequest {
	return &MsgRightSwapRequest{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Sender:           sender,
		TokenIn:          tokenIn,
		TokenOut:         tokenOut,
		Slippage:         slippage,
		Recipient:        recipient,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}

// Route implements sdk.Msg
func (*MsgRightSwapRequest) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (*MsgRightSwapRequest) Type() string {
	return TypeMsgRightSwap
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
func (m *MsgRightSwapRequest) ValidateBasic() error {
	if err := host.PortIdentifierValidator(m.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if !m.TokenIn.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "input token is invalid")
	}
	if !m.TokenOut.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "output token is invalid")
	}
	if m.Slippage <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "maximum slippage is invalid")
	}
	return nil
}

// GetSignBytes implements sdk.Msg.
func (m *MsgRightSwapRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRightSwapRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

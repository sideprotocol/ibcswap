package types

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSwap_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapRequest
		err  error
	}{
		{
			name: "invalid sender address",
			msg: MsgSwapRequest{
				Sender: "invalid_address",
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "invalid recipient address",
			msg: MsgSwapRequest{
				Sender:    "invalid_address",
				Recipient: "invalid_address",
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "dropped token denoms",
			msg: MsgSwapRequest{
				Sender:    sample.AccAddress(),
				Recipient: sample.AccAddress(),
			},
			err: errorsmod.Wrapf(ErrEmptyDenom, "missed token denoms (%s)", ""),
		},
		{
			name: "invalid token-in amounts",
			msg: MsgSwapRequest{
				Sender:    sample.AccAddress(),
				Recipient: sample.AccAddress(),
				TokenIn:   &types.Coin{Denom: "atom", Amount: math.NewInt(0)},
			},
			err: errorsmod.Wrapf(ErrEmptyDenom, "missed token denoms (%s)", ""),
		},
		{
			name: "invalid token-out amounts",
			msg: MsgSwapRequest{
				Sender:    sample.AccAddress(),
				Recipient: sample.AccAddress(),
				TokenIn:   &types.Coin{Denom: "atom", Amount: math.NewInt(10)},
				TokenOut:  &types.Coin{Denom: "atom", Amount: math.NewInt(0)},
			},
			err: errorsmod.Wrapf(ErrInvalidAmount, "invalid token amounts (%s)", ""),
		},

		{
			name: "valid message",
			msg: MsgSwapRequest{
				Sender:    sample.AccAddress(),
				Recipient: sample.AccAddress(),
				TokenIn:   &types.Coin{Denom: "atom", Amount: math.NewInt(100)},
				TokenOut:  &types.Coin{Denom: "marscoin", Amount: math.NewInt(100)},
				Slippage:  100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

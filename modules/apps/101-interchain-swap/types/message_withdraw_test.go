package types

import (
	"testing"

	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdraw_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgWithdrawRequest
		err  error
	}{
		{
			name: "invalid sender address",
			msg: MsgWithdrawRequest{
				Sender: "invalid_address",
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "invalid denomout",
			msg: MsgWithdrawRequest{
				Sender: sample.AccAddress(),
			},
			err: errorsmod.Wrapf(ErrEmptyDenom, "none exist denom (%s)", ""),
		},
		{
			name: "invalid pool-coin amount",
			msg: MsgWithdrawRequest{
				Sender:   sample.AccAddress(),
				DenomOut: "atom",
				PoolCoin: &types.Coin{
					Denom:  "atm",
					Amount: types.NewInt(0),
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", ""),
		},
		{
			name: "valid message",
			msg: MsgWithdrawRequest{
				Sender: sample.AccAddress(),
				PoolCoin: &types.Coin{
					Denom:  "atm",
					Amount: types.NewInt(100),
				},
				DenomOut: "marscoin",
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

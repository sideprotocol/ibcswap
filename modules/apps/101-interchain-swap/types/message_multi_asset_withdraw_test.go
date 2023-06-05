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
		msg  MsgMultiAssetWithdrawRequest
		err  error
	}{
		{
			name: "invalid sender address",

			msg: MsgMultiAssetWithdrawRequest{
				PoolId: "test",
				Sender: "invalid address",
				Withdraws: []*WithdrawAsset{
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "atm",
							Amount: types.NewInt(0),
						},
					},
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "btm",
							Amount: types.NewInt(0),
						},
					},
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "invalid pool-coin amount",
			msg: MsgMultiAssetWithdrawRequest{
				Sender: sample.AccAddress(),
				Withdraws: []*WithdrawAsset{
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "atm",
							Amount: types.NewInt(0),
						},
					},
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "btm",
							Amount: types.NewInt(0),
						},
					},
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", ""),
		},
		{
			name: "valid message",
			msg: MsgMultiAssetWithdrawRequest{
				Sender: sample.AccAddress(),
				Withdraws: []*WithdrawAsset{
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "atm",
							Amount: types.NewInt(1000),
						},
					},
					{
						Receiver: sample.AccAddress(),
						Balance: &types.Coin{
							Denom:  "btm",
							Amount: types.NewInt(1000),
						},
					},
				},
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

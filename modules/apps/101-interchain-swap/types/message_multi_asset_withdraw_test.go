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
				LocalWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender: "invalid_address",
				},
				RemoteWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender: "invalid_address",
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "invalid denomout",
			msg: MsgMultiAssetWithdrawRequest{
				LocalWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender: sample.AccAddress(),
				},
				RemoteWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender: "invalid_address",
				},
			},
			err: errorsmod.Wrapf(ErrEmptyDenom, "none exist denom (%s)", ""),
		},
		{
			name: "invalid pool-coin amount",
			msg: MsgMultiAssetWithdrawRequest{
				LocalWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender:   sample.AccAddress(),
					DenomOut: types.DefaultBondDenom,
					PoolCoin: &types.Coin{
						Denom:  "atm",
						Amount: types.NewInt(0),
					},
				},
				RemoteWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender:   "invalid_address",
					DenomOut: types.DefaultBondDenom,
					PoolCoin: &types.Coin{
						Denom:  "btm",
						Amount: types.NewInt(0),
					},
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", ""),
		},
		{
			name: "valid message",
			msg: MsgMultiAssetWithdrawRequest{
				LocalWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender:   sample.AccAddress(),
					DenomOut: types.DefaultBondDenom,
					PoolCoin: &types.Coin{
						Denom:  "atm",
						Amount: types.NewInt(0),
					},
				},
				RemoteWithdraw: &MsgSingleAssetWithdrawRequest{
					Sender:   "invalid_address",
					DenomOut: types.DefaultBondDenom,
					PoolCoin: &types.Coin{
						Denom:  "btm",
						Amount: types.NewInt(0),
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

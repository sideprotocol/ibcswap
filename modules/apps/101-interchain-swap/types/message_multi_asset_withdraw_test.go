package types

import (
	"testing"

	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/testing/testutil/sample"
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
				Receiver: "invalid address",
				PoolToken: &types.Coin{
					Denom:  "atm",
					Amount: types.NewInt(0),
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		},
		{
			name: "invalid pool-coin amount",
			msg: MsgMultiAssetWithdrawRequest{
				Receiver: sample.AccAddress(),
				PoolToken: &types.Coin{
					Denom:  "atm",
					Amount: types.NewInt(0),
				},
			},
			err: errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", ""),
		},
		{
			name: "valid message",
			msg: MsgMultiAssetWithdrawRequest{
				Receiver: sample.AccAddress(),
				PoolToken: &types.Coin{
					Denom:  "atm",
					Amount: types.NewInt(0),
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

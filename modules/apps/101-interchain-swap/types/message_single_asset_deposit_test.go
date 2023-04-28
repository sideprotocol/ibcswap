package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgDeposit_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSingleAssetDepositRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSingleAssetDepositRequest{
				Sender: "invalid_address",
			},
			err: errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", ""),
		}, {
			name: "valid address",
			msg: MsgSingleAssetDepositRequest{
				PoolId: "test",
				Sender: sample.AccAddress(),
				Token: &types.Coin{
					Denom:  "atom",
					Amount: types.NewInt(100),
				},
			},
		},

		{
			name: "invalid denom length",
			msg: MsgSingleAssetDepositRequest{
				Sender: sample.AccAddress(),
			},
			err: errorsmod.Wrapf(ErrInvalidTokenLength, "invalid token length (%d)", 1),
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

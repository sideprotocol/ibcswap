package types

import (
	"testing"

	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePoolRequest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePoolRequest{
				Sender:        "invalid_address",
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100:99",
				Denoms:        []string{"atom", "marscoin"},
				Decimals:      []uint32{10, 10},
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePoolRequest{
				Sender:        sample.AccAddress(),
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100:99",
				Denoms:        []string{"atom", "marscoin"},
				Decimals:      []uint32{10, 10},
			},
		},
		{
			name: "invalid denom length",
			msg: MsgCreatePoolRequest{
				Sender:        sample.AccAddress(),
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100:99",
				Denoms:        []string{"atom"},
				Decimals:      []uint32{10, 10},
			},
			err: ErrInvalidDenomPair,
		},
		{
			name: "invalid decimal pair",
			msg: MsgCreatePoolRequest{
				Sender:        sample.AccAddress(),
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100:99",
				Denoms:        []string{"atom", "marscoin"},
				Decimals:      []uint32{10},
			},
			err: ErrInvalidDecimalPair,
		},
		{
			name: "invalid weight type",
			msg: MsgCreatePoolRequest{
				Sender:        sample.AccAddress(),
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100,323",
				Denoms:        []string{"atom", "marscoin"},
				Decimals:      []uint32{10, 10},
			},
			err: ErrInvalidWeightPair,
		},
		{
			name: "invalid weight length",
			msg: MsgCreatePoolRequest{
				Sender:        sample.AccAddress(),
				SourcePort:    "interchainswap",
				SourceChannel: "interchainswap-1",
				Weight:        "100:323:200",
				Denoms:        []string{"atom", "marscoin"},
				Decimals:      []uint32{10, 10},
			},
			err: ErrInvalidWeightPair,
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

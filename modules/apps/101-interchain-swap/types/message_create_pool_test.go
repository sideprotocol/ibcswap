package types

import (
	"testing"

	"github.com/sideprotocol/ibcswap/v6/testing/testutil/sample"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
				Creator:             "invalid_address",
				CounterPartyCreator: "invalid_address",
				SourcePort:          "interchainswap",
				SourceChannel:       "interchainswap-1",
				Liquidity: []*PoolAsset{
					{
						Balance: &sdk.Coin{
							Denom:  "aside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
				},
			},
			err: ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgCreatePoolRequest{
				Creator:             sample.AccAddress(),
				CounterPartyCreator: sample.AccAddress(),
				SourcePort:          "interchainswap",
				SourceChannel:       "interchainswap-1",
				Liquidity: []*PoolAsset{
					{
						Balance: &sdk.Coin{
							Denom:  "aside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
				},
			},
		},
		{
			name: "invalid denom length",
			msg: MsgCreatePoolRequest{
				Creator:             sample.AccAddress(),
				CounterPartyCreator: sample.AccAddress(),
				SourcePort:          "interchainswap",
				SourceChannel:       "interchainswap-1",
				Liquidity: []*PoolAsset{
					{
						Balance: &sdk.Coin{
							Denom:  "aside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  20,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
				},
			},
			err: ErrInvalidDenomPair,
		},
		{
			name: "invalid decimal pair",
			msg: MsgCreatePoolRequest{
				Creator:             sample.AccAddress(),
				CounterPartyCreator: sample.AccAddress(),
				SourcePort:          "interchainswap",
				SourceChannel:       "interchainswap-1",
				Liquidity: []*PoolAsset{
					{
						Balance: &sdk.Coin{
							Denom:  "aside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 20,
					},
				},
			},
			err: ErrInvalidDecimalPair,
		},

		{
			name: "invalid weight length",
			msg: MsgCreatePoolRequest{
				Creator:             sample.AccAddress(),
				CounterPartyCreator: sample.AccAddress(),
				SourcePort:          "interchainswap",
				SourceChannel:       "interchainswap-1",
				Liquidity: []*PoolAsset{
					{
						Balance: &sdk.Coin{
							Denom:  "aside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
					{
						Balance: &sdk.Coin{
							Denom:  "bside",
							Amount: sdk.NewInt(1000),
						},
						Weight:  50,
						Decimal: 6,
					},
				},
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

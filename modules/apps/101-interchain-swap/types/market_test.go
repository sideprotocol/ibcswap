package types

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestLeftSwap(t *testing.T) {

	// create mock pool
	demons := []string{"a", "b"}
	poolId := GetPoolId(demons)
	assets := []*PoolAsset{
		{
			Side: PoolSide_NATIVE,
			Balance: &types.Coin{
				Amount: types.NewInt(100000000),
				Denom:  demons[0],
			},
			Weight:  1,
			Decimal: 6,
		},
		{
			Side: PoolSide_REMOTE,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  demons[1],
			},
			Weight:  1,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(100000000),
			Denom:  poolId,
		},

		Status:                PoolStatus_POOL_STATUS_READY,
		EncounterPartyPort:    "test",
		EncounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		pool,
		DefaultMaxFeeRate,
	)

	// mock swap message
	msg := MsgSwapRequest{
		SwapType:  SwapMsgType_LEFT,
		Sender:    "",
		TokenIn:   &types.Coin{Denom: demons[1], Amount: types.NewInt(1000)},
		TokenOut:  &types.Coin{Denom: demons[0], Amount: types.NewInt(1000)},
		Slippage:  10,
		Recipient: "",
	}
	outToken, err := amm.LeftSwap(*msg.TokenIn, msg.TokenOut.Denom)
	fmt.Println(outToken.Amount.Uint64())
	require.NoError(t, err)
}

func TestUpdatePoolAsset(t *testing.T) {

	// create mock pool
	demons := []string{"a", "b"}
	poolId := GetPoolId(demons)
	assets := []*PoolAsset{
		{
			Side: PoolSide_NATIVE,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  demons[0],
			},
			Weight:  1,
			Decimal: 6,
		},
		{
			Side: PoolSide_REMOTE,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  demons[1],
			},
			Weight:  1,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(0),
			Denom:  poolId,
		},

		Status:                PoolStatus_POOL_STATUS_READY,
		EncounterPartyPort:    "test",
		EncounterPartyChannel: "test",
	}

	// mock swap message
	msg := MsgDepositRequest{
		PoolId: poolId,
		Sender: "",
		Tokens: []*types.Coin{
			{
				Amount: types.NewInt(100000000),
				Denom:  demons[0],
			},
			{
				Amount: types.NewInt(100000000),
				Denom:  demons[1],
			},
		},
	}
	for _, token := range msg.Tokens {
		pool.AddAsset(*token)
	}
	fmt.Println("Pool:", pool)
	// require.NoError(t, pool.Assets[])
	require.NotEqual(t, pool.Assets[0].Balance.Amount.Uint64(), uint(0))

}

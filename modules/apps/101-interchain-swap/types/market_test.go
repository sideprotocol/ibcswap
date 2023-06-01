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
		&pool,
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
	msg := MsgSingleAssetDepositRequest{
		PoolId: poolId,
		Sender: "",
		Token: &types.Coin{
			Amount: types.NewInt(100000000),
			Denom:  demons[0],
		},
	}
	pool.AddAsset(*msg.Token)
	fmt.Println("Pool:", pool)
	// require.NoError(t, pool.Assets[])
	require.NotEqual(t, pool.Assets[0].Balance.Amount.Uint64(), uint(0))

}

func TestSingleDeposit(t *testing.T) {
	const initialX = 2_000_000_000_000 // USDT
	const initialY = 1000_000_000      // ETH
	// create mock pool
	denoms := []string{"a", "b"}
	poolId := GetPoolId(denoms)
	assets := []*PoolAsset{
		{
			Side: PoolSide_NATIVE,
			Balance: &types.Coin{
				Amount: types.NewInt(initialX),
				Denom:  denoms[0],
			},
			Weight:  50,
			Decimal: 6,
		},
		{
			Side: PoolSide_REMOTE,
			Balance: &types.Coin{
				Amount: types.NewInt(initialY),
				Denom:  denoms[1],
			},
			Weight:  50,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(initialX + initialY),
			Denom:  poolId,
		},

		Status:                PoolStatus_POOL_STATUS_READY,
		EncounterPartyPort:    "test",
		EncounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		&pool,
		DefaultMaxFeeRate,
	)

	pool.PoolPrice = float32(amm.LpPrice())

	newDeposit := &types.Coin{
		Amount: types.NewInt(initialY),
		Denom:  denoms[0],
	}

	_, err := amm.DepositSingleAsset(*newDeposit)
	require.Error(t, err)
}

func TestSingleWithdraw(t *testing.T) {
	const initialX = 2_000_000 // USDT
	const initialY = 1000      // ETH
	// create mock pool
	denoms := []string{"a", "b"}
	poolId := GetPoolId(denoms)
	assets := []*PoolAsset{
		{
			Side: PoolSide_NATIVE,
			Balance: &types.Coin{
				Amount: types.NewInt(20000000),
				Denom:  denoms[0],
			},
			Weight:  50,
			Decimal: 6,
		},
		{
			Side: PoolSide_REMOTE,
			Balance: &types.Coin{
				Amount: types.NewInt(1000),
				Denom:  denoms[1],
			},
			Weight:  50,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(initialX + initialY),
			Denom:  poolId,
		},

		Status:                PoolStatus_POOL_STATUS_READY,
		EncounterPartyPort:    "test",
		EncounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		&pool,
		DefaultMaxFeeRate,
	)

	redeem := types.NewCoin(poolId, types.NewInt(initialX+initialY))
	outToken, err := amm.SingleWithdraw(redeem, denoms[0])
	fmt.Println(outToken.Amount.Uint64())
	require.NoError(t, err)
}

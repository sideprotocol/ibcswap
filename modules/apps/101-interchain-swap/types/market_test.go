package types

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestLeftSwap(t *testing.T) {

	// create mock pool
	denoms := []string{"a", "b"}
	poolId := GetPoolId("test", "test1", denoms)
	assets := []*PoolAsset{
		{
			Side: PoolAssetSide_SOURCE,
			Balance: &types.Coin{
				Amount: types.NewInt(100000000),
				Denom:  denoms[0],
			},
			Weight:  1,
			Decimal: 6,
		},
		{
			Side: PoolAssetSide_DESTINATION,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  denoms[1],
			},
			Weight:  1,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		Id:     poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(100000000),
			Denom:  poolId,
		},
		SwapFee:             300,
		Status:              PoolStatus_ACTIVE,
		CounterPartyPort:    "test",
		CounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		&pool,
	)

	// mock swap message
	msg := MsgSwapRequest{
		SwapType:  SwapMsgType_LEFT,
		Sender:    "",
		TokenIn:   &types.Coin{Denom: denoms[1], Amount: types.NewInt(1000)},
		TokenOut:  &types.Coin{Denom: denoms[0], Amount: types.NewInt(1000)},
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
	poolId := GetPoolId("test", "test", demons)
	assets := []*PoolAsset{
		{
			Side: PoolAssetSide_SOURCE,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  demons[0],
			},
			Weight:  1,
			Decimal: 6,
		},
		{
			Side: PoolAssetSide_DESTINATION,
			Balance: &types.Coin{
				Amount: types.NewInt(0),
				Denom:  demons[1],
			},
			Weight:  1,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		Id:     poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(0),
			Denom:  poolId,
		},
		SwapFee:             300,
		Status:              PoolStatus_ACTIVE,
		CounterPartyPort:    "test",
		CounterPartyChannel: "test",
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
	poolId := GetPoolId("test", "test", denoms)
	assets := []*PoolAsset{
		{
			Side: PoolAssetSide_SOURCE,
			Balance: &types.Coin{
				Amount: types.NewInt(initialX),
				Denom:  denoms[0],
			},
			Weight:  50,
			Decimal: 6,
		},
		{
			Side: PoolAssetSide_DESTINATION,
			Balance: &types.Coin{
				Amount: types.NewInt(initialY),
				Denom:  denoms[1],
			},
			Weight:  50,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		Id:     poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(initialX + initialY),
			Denom:  poolId,
		},
		SwapFee:             300,
		Status:              PoolStatus_ACTIVE,
		CounterPartyPort:    "test",
		CounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		&pool,
	)

	//pool.PoolPrice = amm.LpPrice()

	newDeposit := &types.Coin{
		Amount: types.NewInt(initialY),
		Denom:  denoms[0],
	}

	_, err := amm.DepositSingleAsset(*newDeposit)
	require.Error(t, err)
}

func TestSingleWithdraw(t *testing.T) {
	const initialX = 1_000_000 // USDT
	const initialY = 1_000_000 // ETH
	// create mock pool
	denoms := []string{"a", "b"}
	poolId := GetPoolId("test", "test", denoms)
	assets := []*PoolAsset{
		{
			Side: PoolAssetSide_SOURCE,
			Balance: &types.Coin{
				Amount: types.NewInt(initialX),
				Denom:  denoms[0],
			},
			Weight:  50,
			Decimal: 6,
		},
		{
			Side: PoolAssetSide_DESTINATION,
			Balance: &types.Coin{
				Amount: types.NewInt(initialY),
				Denom:  denoms[1],
			},
			Weight:  50,
			Decimal: 6,
		},
	}

	pool := InterchainLiquidityPool{
		Id:     poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(initialX + initialY),
			Denom:  poolId,
		},
		SwapFee:             300,
		Status:              PoolStatus_INITIALIZED,
		CounterPartyPort:    "test",
		CounterPartyChannel: "test",
	}

	// create mock liquidity pool.
	amm := NewInterchainMarketMaker(
		&pool,
	)

	redeem := types.NewCoin(poolId, types.NewInt(initialX))
	outToken, err := amm.SingleWithdraw(redeem, denoms[0])
	fmt.Println(outToken.Amount.Uint64())
	require.NoError(t, err)
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateLiquidityBasic(liquidity []*PoolAsset) error {
	if len(liquidity) != 2 {
		return ErrInvalidDenomPair
	}

	weightSum := 0
	for _, asset := range liquidity {
		if asset.Balance.Amount.Equal(sdk.NewInt(0)) {
			return ErrInvalidAmount
		}
		if asset.Decimal > 18 {
			return ErrInvalidDecimalPair
		}
		if asset.Weight >= 100 {
			return ErrInvalidWeight
		}
		weightSum += int(asset.Weight)
	}

	if weightSum != 100 {
		return ErrInvalidWeightPair
	}
	return nil
}

func GetCoinsFromDepositAssets(assets []*DepositAsset) []*sdk.Coin {
	var coins []*sdk.Coin
	for _, asset := range assets {
		coins = append(coins, asset.Balance)
	}
	return coins
}

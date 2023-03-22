package types

import (
	"strconv"
	"strings"

	mathtool "math"

	"github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

// create new liquidity pool
func NewInterchainLiquidityPool(
	ctx types.Context,
	store BankKeeper,
	denoms []string,
	decimals []uint32,
	weight string,
	portId string,
	channelId string,
) *InterchainLiquidityPool {

	//generate poolId
	poolId := GetPoolId(denoms)

	weights := strings.Split(weight, ":")
	weightSize := len(weights)
	denomSize := len(denoms)
	decimalSize := len(decimals)
	assets := []*PoolAsset{}

	if denomSize == weightSize && decimalSize == weightSize {
		for index, denom := range denoms {
			side := PoolSide_NATIVE
			if !store.HasSupply(ctx, denom) {
				side = PoolSide_REMOTE
			}
			weight, _ := strconv.ParseUint(weights[index], 10, 32)
			asset := PoolAsset{
				Side: side,
				Balance: &types.Coin{
					Amount: types.NewInt(0),
					Denom:  denom,
				},
				Weight:  uint32(weight),
				Decimal: decimals[index],
			}
			assets = append(assets, &asset)
		}
	}

	return &InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: types.NewInt(0),
			Denom:  poolId,
		},
		Status:                PoolStatus_POOL_STATUS_READY,
		EncounterPartyPort:    portId,
		EncounterPartyChannel: channelId,
	}
}

// find pool asset by denom
func (ilp *InterchainLiquidityPool) FindAssetByDenom(denom string) (*PoolAsset, error) {
	for _, asset := range ilp.Assets {
		if asset.Balance.Denom == denom {
			return asset, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// update denom
func (ilp *InterchainLiquidityPool) UpdateAssetPoolSide(denom string, side PoolSide) (*PoolAsset, error) {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == denom {
			ilp.Assets[index].Side = side
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// update denom
func (ilp *InterchainLiquidityPool) AddAsset(token types.Coin) error {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == token.Denom {
			updatedCoin := ilp.Assets[index].Balance.Add(token)
			ilp.Assets[index].Balance = &updatedCoin
		}
	}
	return ErrNotFoundDenomInPool
}

// update denom
func (ilp *InterchainLiquidityPool) SubAsset(token types.Coin) error {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == token.Denom {
			updatedCoin := ilp.Assets[index].Balance.Sub(token)
			ilp.Assets[index].Balance = &updatedCoin
		}
	}
	return ErrNotFoundDenomInPool
}

// update pool suppy
func (ilp *InterchainLiquidityPool) AddPoolSupply(token types.Coin) error {
	if token.Denom != ilp.PoolId {
		return ErrInvalidDenom
	}
	updatedCoin := ilp.Supply.Add(token)
	ilp.Supply = &updatedCoin
	return nil
}

// update pool suppy
func (ilp *InterchainLiquidityPool) SubPoolSupply(token types.Coin) error {
	if token.Denom != ilp.PoolId {
		return ErrInvalidDenom
	}
	updatedCoin := ilp.Supply.Sub(token)
	ilp.Supply = &updatedCoin
	return nil
}

//create new market maker

func NewInterchainMarketMaker(
	pool InterchainLiquidityPool,
	feeRate uint32,
) *InterchainMarketMaker {
	return &InterchainMarketMaker{
		Pool:    &pool,
		FeeRate: uint64(feeRate),
	}
}

// MarketPrice Bi / Wi / (Bo / Wo)
func (imm *InterchainMarketMaker) MarketPrice(denomIn, denomOut string) (*float64, error) {
	tokenIn, err := imm.Pool.FindAssetByDenom(denomIn)
	if err != nil {
		return nil, err
	}

	tokenOut, err := imm.Pool.FindAssetByDenom(denomOut)
	if err != nil {
		return nil, err
	}

	balanceIn := tokenIn.Balance.Amount
	balanceOut := tokenOut.Balance.Amount
	weightIn := tokenIn.Weight
	weightOut := tokenOut.Weight
	balance := float64(balanceIn.Uint64()) / float64(weightIn) / float64(float64(balanceOut.Uint64())/float64(weightOut))
	return &balance, nil
}

// P_issued = P_supply * ((1 + At/Bt) ** Wt -1)
func (imm *InterchainMarketMaker) DepositSingleAsset(token types.Coin) (*types.Coin, error) {
	asset, err := imm.Pool.FindAssetByDenom(token.Denom)
	if err != nil {
		return nil, err
	}

	amount := float64(token.Amount.Uint64())
	supply := float64(imm.Pool.Supply.Amount.Uint64())
	weight := float64(asset.Weight) / 100
	var issueAmount float64
	if supply == 0 {
		issueAmount = float64(token.Amount.Uint64())
	} else {
		issueAmount = supply * mathtool.Pow(1+amount/float64(asset.Balance.Amount.Uint64()), float64(weight)-1)
	}

	return &types.Coin{
		Amount: types.NewInt(int64(issueAmount)),
		Denom:  imm.Pool.Supply.Denom,
	}, nil
}

// input the supply token, output the expected token.
// At = Bt * (1 - (1 - P_redeemed / P_supply) ** 1/Wt)
func (imm *InterchainMarketMaker) Withdraw(redeem types.Coin, denomOut string) (*types.Coin, error) {
	asset, err := imm.Pool.FindAssetByDenom(denomOut)

	if err != nil {
		return nil, err
	}
	err = asset.Balance.Validate()
	if err != nil {
		return nil, err
	}

	if imm.Pool.Status != PoolStatus_POOL_STATUS_READY {
		return nil, ErrNotReadyForSwap
	}

	if redeem.Amount.GT(asset.Balance.Amount) {
		return nil, errorsmod.Wrapf(err, "bigger redeem amount than asset balance(%d).", asset.Balance.Amount)
	}

	if redeem.Denom != imm.Pool.Supply.Denom {
		return nil, ErrInvalidDenomPair
	}

	balance := float64(asset.Balance.Amount.Uint64())
	supply := float64(imm.Pool.Supply.Amount.Uint64())
	weight := float64(asset.Weight) / 100
	amountOut := balance * mathtool.Pow((1-(1-float64(redeem.Amount.Uint64())/supply)), 1/weight)

	return &types.Coin{
		Amount: types.NewInt(int64(amountOut)),
		Denom:  denomOut,
	}, nil
}

// LeftSwap implements OutGivenIn
// Input how many coins you want to sell, output an amount you will receive
// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
func (imm *InterchainMarketMaker) LeftSwap(amountIn types.Coin, denomOut string) (*types.Coin, error) {

	assetIn, err := imm.Pool.FindAssetByDenom(amountIn.Denom)
	if err != nil {
		return nil, err
	}

	assetOut, err := imm.Pool.FindAssetByDenom(denomOut)
	if err != nil {
		return nil, err
	}

	// redeem.weight is percentage
	balanceOut := float64(assetOut.Balance.Amount.Uint64())
	balanceIn := float64(assetIn.Balance.Amount.Uint64())
	weightIn := float64(assetIn.Weight) / 100
	weightOut := float64(assetOut.Weight) / 100
	amount := imm.MinusFees(amountIn.Amount)

	amountOut := balanceOut * (1 - balanceIn/mathtool.Pow((balanceIn+amount), weightIn/weightOut))
	return &types.Coin{
		Amount: types.NewInt(int64(amountOut)),
		Denom:  imm.Pool.Supply.Denom,
	}, nil
}

// RightSwap implements InGivenOut
// Input how many coins you want to buy, output an amount you need to pay
// Ai = Bi * ((Bo/(Bo - Ao)) ** Wo/Wi -1)
func (imm *InterchainMarketMaker) RightSwap(amountIn types.Coin, amountOut types.Coin) (*types.Coin, error) {
	assetIn, err := imm.Pool.FindAssetByDenom(amountIn.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "right swap failed! %s", imm.PoolId)
	}

	assetOut, err := imm.Pool.FindAssetByDenom(amountOut.Denom)
	if err != nil {
		return nil, nil //errorsmod.Wrapf(err, "right swap failed because of %s")
	}

	// redeem.weight is percentage
	balanceOut := float64(assetOut.Balance.Amount.Uint64())
	balanceIn := float64(assetIn.Balance.Amount.Uint64())
	weightIn := float64(assetIn.Weight) / 100
	weightOut := float64(assetOut.Weight) / 100

	amount := types.NewInt(int64(balanceIn * (balanceOut/mathtool.Pow(balanceOut-float64(amountOut.Amount.Uint64()), weightOut/weightIn) - 1)))

	if amountIn.Amount.LT(amount) {
		return nil, nil //errorsmod.Wrapf(ErrInvalidAmount, "right swap failed because of %s")
	}
	return &types.Coin{
		Amount: amount,
		Denom:  imm.Pool.Supply.Denom,
	}, nil
}

func (imm *InterchainMarketMaker) MinusFees(amount types.Int) float64 {
	return float64(amount.Uint64()) * (1 - float64(imm.FeeRate)/10000)
}

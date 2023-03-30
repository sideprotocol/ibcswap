package types

import (
	math "math"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

// create new liquidity pool
func NewInterchainLiquidityPool(
	ctx types.Context,
	creator string,
	store BankKeeper,
	tokens []*types.Coin,
	decimals []uint32,
	weight string,
	portId string,
	channelId string,
) *InterchainLiquidityPool {

	//generate poolId
	poolId := GetPoolIdWithTokens(tokens)

	weights := strings.Split(weight, ":")
	weightSize := len(weights)
	denomSize := len(tokens)
	decimalSize := len(decimals)
	assets := []*PoolAsset{}

	if denomSize == weightSize && decimalSize == weightSize {
		for index, token := range tokens {
			side := PoolSide_NATIVE
			if !store.HasSupply(ctx, token.Denom) {
				side = PoolSide_REMOTE
			}
			weight, _ := strconv.ParseUint(weights[index], 10, 32)
			asset := PoolAsset{
				Side:    side,
				Balance: token,
				Weight:  uint32(weight),
				Decimal: decimals[index],
			}
			assets = append(assets, &asset)
		}
	} else {
		return nil
	}

	return &InterchainLiquidityPool{
		PoolId:  poolId,
		Creator: creator,
		Assets:  assets,
		Supply: &types.Coin{
			Amount: types.NewInt(0),
			Denom:  poolId,
		},
		Status:                PoolStatus_POOL_STATUS_INITIAL,
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
	pool *InterchainLiquidityPool,
	feeRate uint32,
) *InterchainMarketMaker {
	return &InterchainMarketMaker{
		Pool:    pool,
		FeeRate: feeRate,
	}
}

// MarketPrice Bi / Wi / (Bo / Wo)
func (imm *InterchainMarketMaker) MarketPrice(denomIn, denomOut string) (*types.Dec, error) {
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

	// Convert all values to Dec type
	balanceInDec := types.NewDecFromBigInt(balanceIn.BigInt())
	balanceOutDec := types.NewDecFromBigInt(balanceOut.BigInt())
	weightInDec := types.NewDecFromInt(types.NewInt(int64(weightIn)))
	weightOutDec := types.NewDecFromInt(types.NewInt(int64(weightOut)))

	// Perform calculations using Dec type
	ratioIn := balanceInDec.Quo(weightInDec)
	ratioOut := balanceOutDec.Quo(weightOutDec)
	marketPrice := ratioIn.Quo(ratioOut)

	return &marketPrice, nil
}

// P_issued = P_supply * ((1 + At/Bt) ** Wt -1)
func (imm *InterchainMarketMaker) DepositSingleAsset(token types.Coin) (*types.Coin, error) {

	asset, err := imm.Pool.FindAssetByDenom(token.Denom)
	if err != nil {
		return nil, err
	}

	amountDec := types.NewDecFromBigInt(token.Amount.BigInt())
	supplyDec := types.NewDecFromBigInt(imm.Pool.Supply.Amount.BigInt())
	weightDec := types.NewDecFromInt(types.NewInt(int64(asset.Weight))).Quo(types.NewDec(100))

	var issueAmount types.Int

	if imm.Pool.Status == PoolStatus_POOL_STATUS_INITIAL {
		totalInitialLpAmount := types.NewInt(0)
		for _, asset := range imm.Pool.Assets {
			totalInitialLpAmount = totalInitialLpAmount.Add(asset.Balance.Amount)
		}
		issueAmount = totalInitialLpAmount.Mul(types.NewInt(int64(asset.Weight))).Quo(types.NewInt(100))
	} else {
		//const issueAmount = supply * (math.pow(1+amount/asset.balance, weight) - 1)

		balanceDec := types.NewDecFromBigInt(asset.Balance.Amount.BigInt())
		ratio := amountDec.Quo(balanceDec)

		// Calculate 1 + ratio / power
		factorBase := types.NewDec(1).Add(ratio)
		factorBaseFloat64, _ := factorBase.Float64()
		powerFloat64, _ := weightDec.Float64()

		// Convert Dec to float64, perform math.Pow, and convert back to Dec
		factorFloat64 := math.Pow(factorBaseFloat64, powerFloat64)
		factorDec := types.NewDecFromInt(types.NewInt(int64(factorFloat64 * 1e18))).Quo(types.NewDecFromInt(types.NewInt(1e18)))

		issueAmountDec := supplyDec.Mul(factorDec.Sub(types.NewDec(1)))
		issueAmount = types.NewIntFromBigInt(issueAmountDec.TruncateInt().BigInt())
	}
	outputToken := &types.Coin{
		Amount: issueAmount,
		Denom:  imm.Pool.Supply.Denom,
	}
	return outputToken, nil
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

	if redeem.Amount.GT(imm.Pool.Supply.Amount) {
		return nil, errorsmod.Wrapf(err, "bigger redeem amount than asset balance(%d).", asset.Balance.Amount)
	}

	if redeem.Denom != imm.Pool.Supply.Denom {
		return nil, ErrInvalidDenomPair
	}

	balance := types.NewDecFromInt(asset.Balance.Amount)
	supply := types.NewDecFromInt(imm.Pool.Supply.Amount)
	weight := types.NewDec(int64(asset.Weight)).Quo(types.NewDec(100))

	// At = Bt * (1 - (1 - P_redeemed / P_supply) ** 1/Wt)
	redeemDec := types.NewDecFromInt(redeem.Amount)
	oneMinusRatio := types.NewDec(1).Sub(redeemDec.Quo(supply))
	power := oneMinusRatio.Power(types.NewDecFromInt(types.NewInt(1)).Quo(weight).TruncateInt().Uint64())
	oneMinusPower := types.NewDec(1).Sub(power)
	amountOut := balance.Mul(oneMinusPower)

	return &types.Coin{
		Amount: amountOut.RoundInt(),
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

	balanceOut := types.NewDecFromInt(assetOut.Balance.Amount)
	balanceIn := types.NewDecFromInt(assetIn.Balance.Amount)
	weightIn := types.NewDec(int64(assetIn.Weight)).Quo(types.NewDec(100))
	weightOut := types.NewDec(int64(assetOut.Weight)).Quo(types.NewDec(100))
	amount := imm.MinusFees(amountIn.Amount)

	// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
	balanceInPlusAmount := balanceIn.Add(amount)
	ratio := balanceIn.Quo(balanceInPlusAmount)
	oneMinusRatio := types.NewDec(1).Sub(ratio)
	power := weightIn.Quo(weightOut)
	factor := oneMinusRatio.Power(power.TruncateInt().Uint64())
	amountOut := balanceOut.Mul(factor)

	return &types.Coin{
		Amount: amountOut.RoundInt(),
		Denom:  denomOut,
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
		return nil, errorsmod.Wrapf(err, "right swap failed because of %s")
	}

	balanceOut := types.NewDecFromInt(assetOut.Balance.Amount)
	balanceIn := types.NewDecFromInt(assetIn.Balance.Amount)
	weightIn := types.NewDec(int64(assetIn.Weight)).Quo(types.NewDec(100))
	weightOut := types.NewDec(int64(assetOut.Weight)).Quo(types.NewDec(100))

	one := types.NewDec(1)
	amountOutDec := types.NewDecFromInt(amountOut.Amount)

	numerator := balanceOut.Sub(amountOutDec)
	power := weightOut.Quo(weightIn)
	denominator, _ := one.Sub(numerator.Quo(balanceOut)).ApproxSqrt()
	factor := one.Sub(denominator.Power(power.TruncateInt().Uint64()))

	amountRequired := balanceIn.Mul(factor).TruncateInt()

	if amountIn.Amount.LT(amountRequired) {
		return nil, errorsmod.Wrapf(ErrInvalidAmount, "right swap failed because of %s")
	}
	return &types.Coin{
		Amount: amountRequired,
		Denom:  amountIn.Denom,
	}, nil
}

func (imm *InterchainMarketMaker) MinusFees(amount types.Int) types.Dec {
	amountDec := types.NewDecFromInt(amount)
	feeRate := types.NewDec(int64(imm.FeeRate)).Quo(types.NewDec(10000))
	fees := amountDec.Mul(feeRate)
	amountMinusFees := amountDec.Sub(fees)
	return amountMinusFees
}

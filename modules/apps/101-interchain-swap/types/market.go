package types

import (
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/types"
)

// create new liquidity pool
func NewInterchainLiquidityPool(
	ctx types.Context,
	sourceCreator string,
	destinationCreator string,
	store BankKeeper,
	poolId string,
	assets []*PoolAsset,
	swapFee uint32,
	portId string,
	channelId string,

) *InterchainLiquidityPool {

	initialLiquidity := types.NewInt(0)
	liquidity := []*PoolAsset{}
	
	for _, asset := range assets {
		
		initialLiquidity = initialLiquidity.Add(asset.Balance.Amount)
		if store.HasSupply(ctx, asset.Balance.Denom) {
			asset.Side = PoolAssetSide_SOURCE
		} else {
			asset.Side = PoolAssetSide_DESTINATION
		}
		liquidity = append(liquidity, asset)
	}

	pool := InterchainLiquidityPool{
		Id:                 poolId,
		SourceCreator:      sourceCreator,
		DestinationCreator: destinationCreator,
		Assets:             liquidity,
		Supply: &types.Coin{
			Denom:  poolId,
			Amount: initialLiquidity,
		},
		Status:              PoolStatus_INITIALIZED,
		PoolPrice:           0,
		CounterPartyPort:    portId,
		CounterPartyChannel: channelId,
		SwapFee:             swapFee,
	}
	return &pool
}

// find pool asset by denom
func (ilp *InterchainLiquidityPool) SetSupply(amount types.Int) {
	ilp.Supply = &types.Coin{Denom: ilp.Id, Amount: amount}
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

// find pool asset by denom
func (ilp *InterchainLiquidityPool) FindDenomBySide(side PoolAssetSide) (*string, error) {
	for _, asset := range ilp.Assets {
		if asset.Side == side {
			return &asset.Balance.Denom, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// find pool asset by denom
func (ilp *InterchainLiquidityPool) FindAssetBySide(side PoolAssetSide) (*types.Coin, error) {
	for _, asset := range ilp.Assets {
		if asset.Side == side {
			return asset.Balance, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// find pool asset by denom
func (ilp *InterchainLiquidityPool) FindPoolAssetBySide(side PoolAssetSide) (*PoolAsset, error) {
	for _, asset := range ilp.Assets {
		if asset.Side == side {
			return asset, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// update denom
func (ilp *InterchainLiquidityPool) UpdateAssetPoolSide(denom string, side PoolAssetSide) (*PoolAsset, error) {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == denom {
			ilp.Assets[index].Side = side
			return ilp.Assets[index], nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// add assets
func (ilp *InterchainLiquidityPool) AddAsset(token types.Coin) error {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == token.Denom {
			updatedCoin := ilp.Assets[index].Balance.Add(token)
			ilp.Assets[index].Balance = &updatedCoin
			return nil
		}
	}
	return ErrNotFoundDenomInPool
}

// Update denom
func (ilp *InterchainLiquidityPool) SubtractAsset(token types.Coin) (*types.Coin, error) {
	for index, asset := range ilp.Assets {
		if asset.Balance.Denom == token.Denom {
			updatedCoin := ilp.Assets[index].Balance.Sub(token)
			ilp.Assets[index].Balance = &updatedCoin
			return ilp.Assets[index].Balance, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

// Increase pool suppy
func (ilp *InterchainLiquidityPool) AddPoolSupply(token types.Coin) error {
	if token.Denom != ilp.Id {
		return ErrInvalidDenom
	}
	updatedCoin := ilp.Supply.Add(token)
	ilp.Supply = &updatedCoin
	return nil
}

// Decrease pool suppy
func (ilp *InterchainLiquidityPool) SubtractPoolSupply(token types.Coin) error {
	if token.Denom != ilp.Id {
		return ErrInvalidDenom
	}
	updatedCoin := ilp.Supply.Sub(token)
	ilp.Supply = &updatedCoin
	return nil
}

// Decrease pool suppy
func (ilp *InterchainLiquidityPool) AllAssetsWithdrawn() bool {
	allAssetsWithdrawn := true
	for _, asset := range ilp.Assets {
		if !asset.Balance.Amount.IsZero() {
			allAssetsWithdrawn = false
		}
	}
	return allAssetsWithdrawn
}

// Calculate total amount of assets in pool
func (ilp *InterchainLiquidityPool) SumOfPoolAssets() types.Int {
	totalAssets := types.NewInt(0)
	for _, asset := range ilp.Assets {
		totalAssets = totalAssets.Add(asset.Balance.Amount)
	}
	return totalAssets
}

// Create new market maker
func NewInterchainMarketMaker(
	pool *InterchainLiquidityPool,
) *InterchainMarketMaker {
	return &InterchainMarketMaker{
		Pool: pool,
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

	decToken := (types.NewDecCoinFromCoin(token))
	decAsset := types.NewDecCoinFromCoin(*asset.Balance)

	var issueAmount types.Int
	if imm.Pool.Status != PoolStatus_ACTIVE {
		return nil, err
	} else {
		weight := types.NewDec(int64(asset.Weight)).Quo(types.NewDec(100)) // divide by 100
		ratio := decToken.Amount.Quo(decAsset.Amount).Add(types.NewDec(1))
		exponent := 1 - math.Pow(ratio.MustFloat64(), weight.MustFloat64())*Multiplier //Ln(ratio).Mul(weight)
		factor, err := types.NewDecFromStr(fmt.Sprintf("%f", exponent/1e8))
		if err != nil {
			return nil, err
		}
		issueAmount = imm.Pool.Supply.Amount.Mul(factor.RoundInt()).Quo(types.NewInt(1e8))
	}

	outputToken := &types.Coin{
		Amount: issueAmount,
		Denom:  imm.Pool.Supply.Denom,
	}
	return outputToken, nil
}

// P_issued = P_supply * Wt * Dt/Bt
func (imm *InterchainMarketMaker) DepositMultiAsset(coins types.Coins) ([]*types.Coin, error) {
	outTokens := []*types.Coin{}
	for _, coin := range coins {
		asset, err := imm.Pool.FindAssetByDenom(coin.Denom)
		if err != nil {
			return nil, err
		}
		var issueAmount types.Dec
		if imm.Pool.Status == PoolStatus_INITIALIZED {
			totalAssetAmount := types.NewDec(0)
			for _, asset := range imm.Pool.Assets {
				decAssetAmount := types.NewDecFromBigIntWithPrec(asset.Balance.Amount.BigInt(), int64(asset.Decimal)) // Convert the amount considering decimal places
				totalAssetAmount = totalAssetAmount.Add(decAssetAmount)
			}
			issueAmount = totalAssetAmount.Mul(types.NewDec(int64(asset.Weight))).Quo(types.NewDec(100))
		} else {

			decToken := types.NewDecCoinFromCoin(coin)
			decAsset := types.NewDecCoinFromCoin(*asset.Balance)
			decSupply := types.NewDecCoinFromCoin(*imm.Pool.Supply)

			ratio := decToken.Amount.Quo(decAsset.Amount).Mul(types.NewDec(Multiplier))
			issueAmount = (decSupply.Amount.Mul(types.NewDec(int64(asset.Weight))).Mul(ratio).Quo(types.NewDec(100)).Quo(types.NewDec(Multiplier)))
		}

		outputToken := &types.Coin{
			Amount: issueAmount.RoundInt(),
			Denom:  imm.Pool.Supply.Denom,
		}
		outTokens = append(outTokens, outputToken)
	}

	return outTokens, nil
}

// input the supply token, output the expected token.
// At = Bt * (1 - (1 - P_redeemed / P_supply) ** 1/Wt)
// func (imm *InterchainMarketMaker) SingleWithdraw(redeem types.Coin, denomOut string) (*types.Coin, error) {
// 	asset, err := imm.Pool.FindAssetByDenom(denomOut)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = asset.Balance.Validate()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if redeem.Amount.GT(imm.Pool.Supply.Amount) {
// 		return nil, fmt.Errorf("bigger than balance")
// 	}

// 	if redeem.Denom != imm.Pool.Supply.Denom {
// 		return nil, fmt.Errorf("invalid denom pair")
// 	}

// 	w := types.NewDec(int64(asset.Weight)).Quo(types.NewDec(100)) // divide by 100
// 	decSupply := types.NewDecCoinFromCoin(*imm.Pool.Supply)
// 	decRedeem := types.NewDecCoinFromCoin(redeem)
// 	decAsset := types.NewDecCoinFromCoin(*asset.Balance)
// 	ratio := decSupply.Amount.Sub(decRedeem.Amount).Mul(types.NewDec(Multiplier)).Quo(decSupply.Amount)

// 	exponent := types.NewDec(1).Quo(w)
// 	factor := types.NewDec(1).Sub(Exp(Ln(ratio).Mul(exponent))).Mul(types.NewDec(Multiplier))

// 	amountOut := decAsset.Amount.Mul(factor).Quo(types.NewDec(Multiplier))
// 	return &types.Coin{
// 		Amount: amountOut.RoundInt(),
// 		Denom:  denomOut,
// 	}, nil
// }

// input the supply token, output the expected token.
// At = Bt * (P_redeemed / P_supply)/Wt
func (imm *InterchainMarketMaker) MultiAssetWithdraw(redeem types.Coin) ([]*types.Coin, error) {
	outs := []*types.Coin{}
	if redeem.Amount.GT(imm.Pool.Supply.Amount) {
		return nil, ErrOverflowAmount
	}
	for _, asset := range imm.Pool.Assets {
		out := asset.Balance.Amount.Mul(redeem.Amount).Quo(imm.Pool.Supply.Amount)
		outs = append(outs, &types.Coin{
			Denom:  asset.Balance.Denom,
			Amount: out,
		})
	}
	return outs, nil
}

// LeftSwap implements OutGivenIn
// Input how many coins you want to sell, output an amount you will receive
// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
func (imm *InterchainMarketMaker) LeftSwap(amountIn types.Coin, denomOut string) (*types.Coin, error) {
	assetIn, err := imm.Pool.FindAssetByDenom(amountIn.Denom)
	if err != nil {
		return nil, fmt.Errorf("left swap failed: could not find asset in by denom")
	}

	assetOut, err := imm.Pool.FindAssetByDenom(denomOut)
	if err != nil {
		return nil, fmt.Errorf("left swap failed: could not find asset out by denom")
	}

	balanceOut := types.NewDecFromBigInt(assetOut.Balance.Amount.BigInt())
	balanceIn := types.NewDecFromBigInt(assetIn.Balance.Amount.BigInt())
	weightIn := types.NewDec(int64(assetIn.Weight)).Quo(types.NewDec(100))
	weightOut := types.NewDec(int64(assetOut.Weight)).Quo(types.NewDec(100))
	amount := imm.MinusFees(amountIn.Amount)

	// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
	balanceInPlusAmount := balanceIn.Add(amount)
	ratio := balanceIn.Quo(balanceInPlusAmount)
	oneMinusRatio := types.NewDec(1).Sub(ratio)

	power := weightIn.Quo(weightOut)
	factor := math.Pow(oneMinusRatio.MustFloat64(), power.MustFloat64()) * Multiplier
	finalFactor := factor / 1e8

	amountOut := balanceOut.Mul(types.MustNewDecFromStr(fmt.Sprintf("%f", finalFactor))).Quo(types.NewDec(1e10))
	return &types.Coin{
		Amount: amountOut.RoundInt(),
		Denom:  denomOut,
	}, nil
}

// / RightSwap implements InGivenOut
// Input how many coins you want to buy, output an amount you need to pay
// Ai = Bi * ((Bo/(Bo - Ao)) ** Wo/Wi -1)
func (imm *InterchainMarketMaker) RightSwap(amountIn types.Coin, amountOut types.Coin) (*types.Coin, error) {
	assetIn, err := imm.Pool.FindAssetByDenom(amountIn.Denom)
	if err != nil {
		return nil, fmt.Errorf("right swap failed: could not find asset in by denom")
	}

	assetOut, err := imm.Pool.FindAssetByDenom(amountOut.Denom)

	decAmountOut := types.NewDecCoinFromCoin(amountOut)
	decAssetIn := types.NewDecCoinFromCoin(*assetIn.Balance)
	decAssetOut := types.NewDecCoinFromCoin(*assetOut.Balance)

	if err != nil {
		return nil, fmt.Errorf("right swap failed: could not find asset out by denom")
	}

	// Ai = Bi * ((Bo/(Bo - Ao)) ** Wo/Wi -1)
	balanceIn := decAssetIn.Amount
	weightIn := types.NewDec(int64(assetIn.Weight)).Quo(types.NewDec(100))
	weightOut := types.NewDec(int64(assetOut.Weight)).Quo(types.NewDec(100))

	numerator := decAssetOut.Amount
	power := weightOut.Quo(weightIn)
	denominator := decAssetOut.Amount.Sub(decAmountOut.Amount)
	base := numerator.Quo(denominator)
	factor := Exp(Ln(base).Mul(power)).Sub(types.NewDec(1)).Mul(types.NewDec(Multiplier))
	amountRequired := balanceIn.Mul(factor).Quo(types.NewDec(Multiplier)).RoundInt()

	if amountIn.Amount.LT(amountRequired) {
		return nil, fmt.Errorf("right swap failed: insufficient amount")
	}
	return &types.Coin{
		Amount: amountRequired,
		Denom:  amountIn.Denom,
	}, nil
}

func (imm *InterchainMarketMaker) MinusFees(amount types.Int) types.Dec {
	amountDec := types.NewDecFromInt(amount)
	feeRate := types.NewDec(int64(imm.Pool.SwapFee)).Quo(types.NewDec(10000))
	fees := amountDec.Mul(feeRate)
	amountMinusFees := amountDec.Sub(fees)
	return amountMinusFees
}

// Worth Function V=M)
func (imm *InterchainMarketMaker) Invariant() types.Dec {
	v := types.NewDec(1)
	totalBalance := types.NewDec(0)
	for _, asset := range imm.Pool.Assets {
		decAssetAmount := types.NewDecFromBigIntWithPrec(asset.Balance.Amount.BigInt(), int64(asset.Decimal)) // Convert the amount considering decimal places
		totalBalance = totalBalance.Add(decAssetAmount)
	}

	for _, asset := range imm.Pool.Assets {
		w := types.NewDec(int64(asset.Weight)).Quo(types.NewDec(100))                                  // divide by 100
		balance := types.NewDecFromBigIntWithPrec(asset.Balance.Amount.BigInt(), int64(asset.Decimal)) //types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt())

		// Raise balance to the power of w using logarithm and exponential functions
		exponent := Ln(balance).Mul(w)
		v = v.Mul(Exp(exponent))
	}
	return v
}

// Natural logarithm with the Newton-Raphson method
func Ln(dec types.Dec) types.Dec {
	// You can adjust the precision and initial guess based on your requirements
	const maxIterations = 500
	guess, _ := types.NewDecFromStr("1.0")

	two := types.NewDec(2)
	for i := 0; i < maxIterations; i++ {
		guessSquared := guess.Mul(guess)
		guess = guess.Add((dec.Quo(guess).Sub(guessSquared)).Quo(two))
	}
	return guess
}

// Exponential function
func Exp(dec types.Dec) types.Dec {
	// Again, adjust the precision and the number of terms in the series based on your needs
	const maxIterations = 500
	result, _ := types.NewDecFromStr("1.0")
	term := result

	for i := 1; i < maxIterations; i++ {
		term = term.Mul(dec).QuoInt64(int64(i))
		result = result.Add(term)
	}
	return result
}

// func (imm *InterchainMarketMaker) InvariantWithInput(tokenIn types.Coin) types.Dec {
// 	v := types.NewDec(1)
// 	totalBalance := types.NewDec(0)
// 	for _, asset := range imm.Pool.Assets {
// 		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
// 		totalBalance = totalBalance.Add(types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt()))
// 		if asset.Balance.Denom == tokenIn.Denom {
// 			totalBalance = totalBalance.Add(types.NewDecFromBigInt(tokenIn.Amount.BigInt()))
// 		}
// 	}
// 	for _, asset := range imm.Pool.Assets {
// 		w := types.NewDec(int64(asset.Weight)).Quo(types.NewDec(100)) // divide by 100
// 		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
// 		var balance types.Dec
// 		if tokenIn.Denom != asset.Balance.Denom {
// 			balance = types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt())
// 		} else {
// 			balance = types.NewDecFromBigInt(asset.Balance.Amount.Add(tokenIn.Amount).Quo(decimal).BigInt())
// 		}

// 		// Raise balance to the power of w using logarithm and exponential functions
// 		exponent := Ln(balance).Mul(w)
// 		v = v.Mul(Exp(exponent))
// 	}
// 	return v
// }

// func (imm *InterchainMarketMaker) LpPrice() uint64 {
// 	decSupply := types.NewDecCoinFromCoin(*imm.Pool.Supply)
// 	lpPrice := imm.Invariant().Quo(decSupply.Amount)
// 	return lpPrice.BigInt().Uint64()
// }

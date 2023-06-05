package types

import (
	"fmt"
	math "math"

	"github.com/cosmos/cosmos-sdk/types"
)

// create new liquidity pool
func NewInterchainLiquidityPool(
	ctx types.Context,
	creator string,
	store BankKeeper,
	poolId string,
	assets []*PoolAsset,
	swapFee uint32,
	portId string,
	channelId string,

) *InterchainLiquidityPool {

	initialLiquidity := types.NewInt(0)
	for _, asset := range assets {
		initialLiquidity = initialLiquidity.Add(asset.Balance.Amount)
	}

	pool := InterchainLiquidityPool{
		Id:      poolId,
		Creator: creator,
		Assets:  assets,
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

	var issueAmount types.Int
	if imm.Pool.Status != PoolStatus_ACTIVE {
		return nil, err
	} else {

		weight := float64(asset.Weight) / 100
		ratio := 1 + float64(token.Amount.Int64())/float64(asset.Balance.Amount.Int64())
		factor := (math.Pow(ratio, float64(weight)) - 1) * Multiplier
		issueAmount = imm.Pool.Supply.Amount.Mul(types.NewInt(int64(factor))).Quo(types.NewInt(Multiplier))
	}

	outputToken := &types.Coin{
		Amount: issueAmount,
		Denom:  imm.Pool.Supply.Denom,
	}
	return outputToken, nil
}

// P_issued = P_supply * ((1 + At/Bt) ** Wt -1)
func (imm *InterchainMarketMaker) DepositMultiAsset(tokens []*types.Coin) ([]*types.Coin, error) {
	outTokens := []*types.Coin{}
	for _, token := range tokens {
		asset, err := imm.Pool.FindAssetByDenom(token.Denom)
		if err != nil {
			return nil, err
		}
		var issueAmount types.Int
		if imm.Pool.Status == PoolStatus_INITIALIZED {
			issueAmount = asset.Balance.Amount
		} else {
			ratio := float64(token.Amount.Int64()) / float64(asset.Balance.Amount.Int64()) * Multiplier
			issueAmount = imm.Pool.Supply.Amount.Mul(types.NewInt(int64(ratio))).Quo(types.NewInt(Multiplier))
		}

		outputToken := &types.Coin{
			Amount: issueAmount,
			Denom:  imm.Pool.Supply.Denom,
		}
		outTokens = append(outTokens, outputToken)
	}

	return outTokens, nil
}

// input the supply token, output the expected token.
// At = Bt * (1 - (1 - P_redeemed / P_supply) ** 1/Wt)
func (imm *InterchainMarketMaker) SingleWithdraw(redeem types.Coin, denomOut string) (*types.Coin, error) {
	asset, err := imm.Pool.FindAssetByDenom(denomOut)
	if err != nil {
		return nil, err
	}
	// err = asset.Balance.Validate()
	// if err != nil {
	// 	return nil, err
	// }

	if redeem.Amount.GT(imm.Pool.Supply.Amount) {
		return nil, fmt.Errorf("bigger than balance")
	}

	if redeem.Denom != imm.Pool.Supply.Denom {
		return nil, fmt.Errorf("invalid denom pair")
	}

	w := float64(asset.Weight) / 100
	ratio := imm.Pool.Supply.Amount.Sub(redeem.Amount).Mul(types.NewInt(Multiplier)).Quo(imm.Pool.Supply.Amount)
	ratioFloat := float64(ratio.Int64()) / Multiplier

	exponent := 1 / w
	factor := (1 - math.Pow(ratioFloat, exponent)) * Multiplier
	amountOut := asset.Balance.Amount.Mul(types.NewInt(int64(factor))).Quo(types.NewInt(Multiplier))
	return &types.Coin{
		Amount: amountOut,
		Denom:  denomOut,
	}, nil
}

// input the supply token, output the expected token.
// At = Bt * (P_redeemed / P_supply)
func (imm *InterchainMarketMaker) MultiAssetWithdraw(redeem types.Coin, denomOut string) (*types.Coin, error) {

	asset, err := imm.Pool.FindAssetByDenom(denomOut)
	if err != nil {
		return nil, err
	}
	out := asset.Balance.Amount.Mul(redeem.Amount).Quo(imm.Pool.Supply.Amount)
	return &types.Coin{
		Denom:  denomOut,
		Amount: out,
	}, nil
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
	amountOut := balanceOut.Mul(types.NewDec(int64(factor))).Quo(types.NewDec(Multiplier))
	return &types.Coin{
		Amount: amountOut.TruncateInt(),
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
	if err != nil {
		return nil, fmt.Errorf("right swap failed: could not find asset out by denom")
	}

	// Ai = Bi * ((Bo/(Bo - Ao)) ** Wo/Wi -1)
	balanceIn := types.NewDecFromBigInt(assetIn.Balance.Amount.BigInt())
	weightIn := types.NewDec(int64(assetIn.Weight)).Quo(types.NewDec(100))
	weightOut := types.NewDec(int64(assetOut.Weight)).Quo(types.NewDec(100))

	numerator := types.NewDecFromBigInt(assetOut.Balance.Amount.BigInt())
	power := weightOut.Quo(weightIn)
	denominator := types.NewDecFromBigInt(assetOut.Balance.Amount.Sub(amountOut.Amount).BigInt())
	base := numerator.Quo(denominator)
	factor := math.Pow(base.MustFloat64(), power.MustFloat64()) * Multiplier
	amountRequired := balanceIn.Mul(types.NewDec(int64(factor))).Quo(types.NewDec(Multiplier)).RoundInt()

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
func (imm *InterchainMarketMaker) Invariant() float64 {
	v := 1.0
	totalBalance := types.NewDec(0)
	for _, asset := range imm.Pool.Assets {
		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
		totalBalance = totalBalance.Add(types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt()))
	}
	for _, asset := range imm.Pool.Assets {
		w := float64(asset.Weight) / 100.0
		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
		balance := types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt())
		v *= math.Pow(balance.MustFloat64(), w)
	}
	return v
}

func (imm *InterchainMarketMaker) InvariantWithInput(tokenIn types.Coin) float64 {
	v := 1.0
	totalBalance := types.NewDec(0)
	for _, asset := range imm.Pool.Assets {
		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
		totalBalance = totalBalance.Add(types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt()))
		if asset.Balance.Denom == tokenIn.Denom {
			totalBalance.Add(types.NewDecFromBigInt(tokenIn.Amount.BigInt()))
		}
	}
	for _, asset := range imm.Pool.Assets {
		w := float64(asset.Weight) / 100.0
		decimal := types.NewInt(int64(math.Pow10(int(asset.Decimal))))
		var balance types.Dec
		if tokenIn.Denom != asset.Balance.Denom {
			balance = types.NewDecFromBigInt(asset.Balance.Amount.Quo(decimal).BigInt())
		} else {
			balance = types.NewDecFromBigInt(asset.Balance.Amount.Add(tokenIn.Amount).Quo(decimal).BigInt())
		}
		v *= math.Pow(balance.MustFloat64(), w)
	}
	return v
}

func (imm *InterchainMarketMaker) LpPrice() float64 {
	lpPrice := imm.Invariant() / float64(imm.Pool.Supply.Amount.Int64())
	return lpPrice
}

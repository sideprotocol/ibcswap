package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BalancerAMM struct {
	Pool *BalancerLiquidityPool
	// basis point
	FeeRate int64
}

func NewBalanceAMM(pool *BalancerLiquidityPool, feeRate int64) BalancerAMM {
	return BalancerAMM{
		Pool:    pool,
		FeeRate: feeRate,
	}
}

// MarketPrice Bi / Wi / (Bo / Wo)
func (p *BalancerAMM) MarketPrice(denomIn, denomOut string) float64 {
	//Ti, err := p.Pool.findAssetByDenom(denomIn)
	//if err != nil {
	//	return 0
	//}
	//To, err := p.Pool.findAssetByDenom(denomOut)
	//if err != nil {
	//	return 0
	//}
	//Bi := Ti.Balance.Amount
	//Bo := To.Balance.Amount
	//Wi := Ti.Weight
	//Wo := To.Weight
	//
	//return Bi.ToDec().MustFloat64() / float64(Wi) / (Bo.ToDec().MustFloat64() / float64(Wo))
	return 0
}

func (p *BalancerAMM) Deposit(amount []*sdk.Coin) (sdk.Coin, error) {
	switch len(amount) {
	case 1:
		return p.depositSingleAsset(amount[0])
	case 2:
		return p.depositMultiAssets(amount)
	default:
		return sdk.Coin{}, ErrInvalidPairLength
	}
}

func (p *BalancerAMM) Withdraw(amount *sdk.Coin, denomOut string) (sdk.Coin, error) {

	//if amount.Denom != denomOut {
	//	return sdk.Coin{}, ErrInvalidToken
	//}
	//
	//redeem, err := p.Pool.findAssetByDenom(denomOut)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//
	//// redeem.weight is percentage
	//outNum := float64(redeem.Balance.Amount.Int64()) *
	//	(1 - math.Pow(
	//		1-amount.Amount.Quo(p.Pool.PoolToken.Amount).ToDec().MustFloat64(),
	//		100/float64(redeem.Weight)))
	return sdk.NewCoin(denomOut, sdk.NewInt(int64(0))), nil
}

// LeftSwap implements OutGivenIn
// Input how many coins you want to sell, output an amount you will receive
// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
func (p *BalancerAMM) LeftSwap(Ai *sdk.Coin, denomOut string) (sdk.Coin, error) {

	//Bi, err := p.Pool.findAssetByDenom(Ai.Denom)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//Bo, err := p.Pool.findAssetByDenom(denomOut)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//
	//// redeem.weight is percentage
	//Ao := float64(Bo.Balance.Amount.Int64()) *
	//	math.Pow(
	//		1-Bi.Balance.Amount.Quo(Bi.Balance.Amount.Add(Ai.Amount)).ToDec().MustFloat64(),
	//		float64(Bi.Weight)/float64(Bo.Weight))
	return sdk.NewCoin(denomOut, sdk.NewInt(int64(0))), nil
}

// RightSwap implements InGivenOut
// Input how many coins you want to buy, output an amount you need to pay
// Ai = Bi * ((Bo/(Bo - Ao)) ** Wo/Wi -1)
func (p *BalancerAMM) RightSwap(Ai, Ao *sdk.Coin) (sdk.Coin, error) {

	//Bi, err := p.Pool.findAssetByDenom(Ai.Denom)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//Bo, err := p.Pool.findAssetByDenom(Ao.Denom)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//
	//amount := float64(Bi.Balance.Amount.Int64()) *
	//	(math.Pow(
	//		Bo.Balance.Amount.Quo(Bo.Balance.Amount.Sub(Ao.Amount)).ToDec().MustFloat64(),
	//		float64(Bo.Weight)/float64(Bi.Weight)) - 1)
	//
	//if Ai.Amount.LT(sdk.NewInt(int64(amount))) {
	//	return sdk.Coin{}, ErrAmountInsufficient
	//}

	return sdk.NewCoin(Ai.Denom, sdk.NewInt(int64(0))), nil
}

// amount - amount * feeRate / 10000
func (p *BalancerAMM) minusFees(amount sdk.Int) sdk.Int {
	return amount.Sub(amount.Mul(sdk.NewInt(p.FeeRate / 10000)))
}

// P_issued = P_supply * ((1 + At/Bt) ** Wt -1)
func (p *BalancerAMM) depositSingleAsset(token *sdk.Coin) (sdk.Coin, error) {
	//Bt, err := p.Pool.findAssetByDenom(token.Denom)
	//if err != nil {
	//	return sdk.Coin{}, err
	//}
	//At := p.minusFees(token.Amount)
	//Bt.Balance.AddAmount(At) // update pool states
	//
	//issue := CalIssueAmount(p.Pool.PoolToken.Amount, At, Bt.Balance.Amount, Bt.Weight)
	//Bt.Balance.Amount = Bt.Balance.Amount.Add(token.Amount) // update Balance to tokens
	return sdk.NewCoin(p.Pool.PoolToken.Denom, math.NewInt(0)), nil
}

// Dt = ((P_supply + P_issued) / P_supply - 1) * Bt
func (p *BalancerAMM) depositMultiAssets(tokens []*sdk.Coin) (sdk.Coin, error) {
	panic("unimplemented")
}

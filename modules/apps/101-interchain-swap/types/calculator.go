package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

func MinusFees(amount sdk.Int, feeRate int64) sdk.Int {
	return amount.Sub(amount.Mul(sdk.NewInt(feeRate / 10000)))
}

// CalIssueAmount returns `supply * ((1 + At/Bt)^Wt -1 )`
func CalIssueAmount(supply sdk.Int, At sdk.Int, Bt sdk.Int, Wt uint32) sdk.Int {
	p, _ := supply.ToDec().Float64()
	at, _ := At.ToDec().Float64()
	bt, _ := Bt.ToDec().Float64()
	wt := float64(Wt) / 100
	issue := p*math.Pow(1+at/bt, wt) - 1
	return sdk.NewInt(int64(issue))
}

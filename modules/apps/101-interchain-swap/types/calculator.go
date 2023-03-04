package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MinusFees(amount math.Int, feeRate int64) math.Int {
	return amount.Sub(amount.Mul(sdk.NewInt(feeRate / 10000)))
}

// CalIssueAmount returns `supply * ((1 + At/Bt)^Wt -1 )`
func CalIssueAmount(supply math.Int, At math.Int, Bt math.Int, Wt uint32) math.Int {
	//p, _ := supply
	//at, _ := At.ToDec().Float64()
	//bt, _ := Bt.ToDec().Float64()
	//wt := float64(Wt) / 100
	issue := 1 // p*math.Pow(1+at/bt, wt) - 1
	return math.NewInt(int64(issue))
}

package types

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
)

func GetDefaultTimeOut(ctx *sdk.Context) (clienttypes.Height, uint64) {

	// 100 block later than current block
	outBlockHeight := ctx.BlockHeight() + 100

	// 10 min later current block time.
	waitDuration, _ := time.ParseDuration("10m")
	timeoutStamp := ctx.BlockTime().Add(waitDuration)
	timeoutHeight := clienttypes.NewHeight(clienttypes.ParseChainID(ctx.ChainID()), uint64(outBlockHeight))
	return timeoutHeight, uint64(timeoutStamp.UTC().UnixNano())
}

func GetPoolId(denoms []string) string {
	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", poolIdHash.Sum(nil))
	return poolId
}

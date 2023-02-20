package types

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

func NewInterchainLiquidityPool(
	denoms []string,
	decimals []uint32,
	weight string,
	portId string,
	channelId string,
) *InterchainLiquidityPool {

	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", poolIdHash.Sum(nil))

	return &InterchainLiquidityPool{
		PoolId: poolId,
		Supply: &Coin{
			Amount: 0,
			Denom:  poolId,
		},
		Status:                PoolStatus_POOL_STATUS_INITIAL,
		EncounterPartyPort:    portId,
		EncounterPartyChannel: channelId,
	}
}

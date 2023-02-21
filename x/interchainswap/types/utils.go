package types

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"time"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func GetDefaultTimeOut() (clienttypes.Height, uint64) {
	timeoutHeight := clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: 10,
	}
	timeoutStamp := time.Now().UTC().Unix()
	return timeoutHeight, uint64(timeoutStamp)
}

func GetPoolId(denoms []string) string {
	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", poolIdHash.Sum(nil))
	return poolId
}

package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

func GetDefaultTimeOut(ctx *sdk.Context) (clienttypes.Height, uint64) {

	// 100 block later than current block
	outBlockHeight := ctx.BlockHeight() + 200

	// 10 min later current block time.
	waitDuration, _ := time.ParseDuration("10m")
	timeoutStamp := ctx.BlockTime().Add(waitDuration)
	timeoutHeight := clienttypes.NewHeight(clienttypes.ParseChainID(ctx.ChainID()), uint64(outBlockHeight))
	return timeoutHeight, uint64(timeoutStamp.UTC().UnixNano())
}

func GetPoolIdWithTokens(tokens []*sdk.Coin) string {

	denoms := []string{}
	for _, token := range tokens {
		denoms = append(denoms, token.Denom)
	}
	return GetPoolId(denoms)
}
func GetPoolId(denoms []string) string {
	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIdHash.Sum(nil)))
	//poolId := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIdHash.Sum(nil)))
	return poolId
}

func GetEscrowAddress(portID, channelID string) sdk.AccAddress {
	// a slash is used to create domain separation between port and channel identifiers to
	// prevent address collisions between escrow addresses created for different channels
	contents := fmt.Sprintf("%s/%s", portID, channelID)

	// ADR 028 AddressHash construction
	preImage := []byte(Version)
	preImage = append(preImage, 0)
	preImage = append(preImage, contents...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}

func GetEscrowAddressWithModuleName(name string) sdk.AccAddress {
	// a slash is used to create domain separation between port and channel identifiers to
	// prevent address collisions between escrow addresses created for different channels

	// ADR 028 AddressHash construction
	preImage := []byte(Version)
	preImage = append(preImage, 0)
	preImage = append(preImage, name...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}

func GetEscrowModuleName(portID, channelID string) string {
	return fmt.Sprintf("%s:-%s-%s", ModuleName, portID, channelID)

}

func CreateEscrowAccount(portID, channelID string) {
	name := fmt.Sprintf("%s/%s", portID, channelID)
	acc := authtypes.NewEmptyModuleAccount(name)
	pubAddr := GetEscrowAddress(portID, channelID)
	acc.SetAddress(pubAddr)
}

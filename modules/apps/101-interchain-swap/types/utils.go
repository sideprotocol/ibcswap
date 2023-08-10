package types

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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
	timeoutHeight := clienttypes.NewHeight(0, uint64(outBlockHeight))
	return timeoutHeight, uint64(timeoutStamp.UTC().UnixNano())
}

func GetPoolId(sourceChainId, destinationChainId string, denoms []string) string {
	connectionId := GetConnectID(sourceChainId, destinationChainId)
	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	//salt := GenerateRandomString(chainID, 10)
	denoms = append(denoms, connectionId)
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIdHash.Sum(nil)))
	return poolId
}

func GetOrderId(maker string, sequence uint64) string {
	orderIdHash := sha256.New()
	orderIdHash.Write([]byte(strings.Join([]string{maker, fmt.Sprintf("%s", sequence)}, "-")))
	orderId := "multi_deposit_order" + fmt.Sprintf("%v", hex.EncodeToString(orderIdHash.Sum(nil)))
	return orderId
}

func GetConnectID(chainIds ...string) string {
	//generate poolId
	sort.Strings(chainIds)
	return strings.Join(chainIds, "/")
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

func GenerateRandomString(chainID string, n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return chainID + base64.URLEncoding.EncodeToString(b)
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

func BytesToUint(b []byte) (uint, error) {
	buf := bytes.NewReader(b)
	var num uint
	err := binary.Read(buf, binary.LittleEndian, &num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func UintToBytes(num uint) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// slippage value have to be in 0~10000
func CheckSlippage(expect, result sdk.Dec, desireSlippage int64) error {
	if desireSlippage > 10000 {
		return ErrInvalidSlippage
	}
	// Define the slippage tolerance (1% in this example)
	slippageTolerance := sdk.NewDec(desireSlippage)

	// Calculate the absolute difference between the ratios
	ratioDiff := expect.Sub(result).Abs()

	// Calculate slippage percentage (slippage = ratioDiff/expect * 10000)
	slippage := ratioDiff.Mul(sdk.NewDec(10000)).Quo(expect)

	// Check if the slippage is within the tolerance
	if slippage.GT(slippageTolerance) {
		return ErrInvalidPairRatio
	}
	return nil
}

func GetEventValueWithSuffix(value, suffix string) string {
	return value + "_" + suffix
}

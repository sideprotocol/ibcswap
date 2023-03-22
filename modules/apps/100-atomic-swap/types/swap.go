package types

import (
	"crypto/sha256"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
	"strconv"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func NewAtomicOrder(maker *SwapMaker, channelId string) AtomicSwapOrder {
	buf, _ := proto.Marshal(maker)
	id := Hash(buf).String()
	return AtomicSwapOrder{
		Id:                id,
		Maker:             maker,
		Status:            Status_INITIAL,
		ChannelId:         channelId,
		Takers:            nil,
		CancelTimestamp:   0,
		CompleteTimestamp: 0,
	}
}

func NewMakerFromMsg(msg *MakeSwapMsg) *SwapMaker {
	return &SwapMaker{
		SourcePort:            msg.SourcePort,
		SourceChannel:         msg.SourceChannel,
		SellToken:             msg.SellToken,
		BuyToken:              msg.BuyToken,
		MakerAddress:          msg.MakerAddress,
		MakerReceivingAddress: msg.MakerReceivingAddress,
		DesiredTaker:          msg.DesiredTaker,
		CreateTimestamp:       msg.CreateTimestamp,
	}
}

func NewTakerFromMsg(msg *TakeSwapMsg) *SwapTaker {
	return &SwapTaker{
		OrderId:               msg.OrderId,
		SellToken:             msg.SellToken,
		TakerAddress:          msg.TakerAddress,
		TakerReceivingAddress: msg.TakerReceivingAddress,
		CreateTimestamp:       msg.CreateTimestamp,
	}
}

func Hash(content []byte) tmbytes.HexBytes {
	hash := sha256.Sum256(content)
	return hash[:]
}

func CreateOrder(msg *MakeSwapMsg, packet channeltypes.Packet) AtomicSwapOrder {
	//path := orderPath(packet)
	//return AtomicSwapOrder{
	//	Id:     GenerateOrderId(packet),
	//	Status: Status_INITIAL,
	//	Path:   path,
	//	Maker:  msg,
	//}
	return AtomicSwapOrder{}
}

func orderPath(packet channeltypes.Packet) string {
	return "channel/" + packet.SourceChannel +
		"/port/" + packet.SourcePort +
		"/channel/" + packet.DestinationChannel +
		"/port/" + packet.DestinationPort +
		"/sequence/" + strconv.FormatUint(packet.Sequence, 10)
}

// GenerateOrderId id is a global unique string, since packet contains sourceChannel, SourcePort, distChannel, distPort, sequence and msg data
func GenerateOrderId(packet channeltypes.Packet) string {
	bytes, _ := proto.Marshal(&packet)
	hash := sha256.Sum256(bytes)
	return string(hash[:])
}

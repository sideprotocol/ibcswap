package types

import (
	"crypto/sha256"
	"strconv"

	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto"
)

func CreateOrder(msg *MakeSwapMsg) Order {
	path := orderPath(msg.Packet)
	return Order{
		Id:     GenerateOrderId(msg.Packet),
		Status: Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}
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

package types

import (
	"crypto/sha256"
	"github.com/gogo/protobuf/proto"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func NewAtomicOrder(msg *MsgMakeSwapRequest, channelId string) AtomicSwapOrder {
	maker := NewMakerFromMsg(msg)
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

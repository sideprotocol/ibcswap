package types

import (
	"crypto/sha256"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

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

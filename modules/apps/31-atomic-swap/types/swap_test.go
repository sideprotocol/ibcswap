package types

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestOrderPath(t *testing.T) {

	packet := channeltypes.Packet{
		SourcePort:         "validPort",
		SourceChannel:      "validChannel",
		DestinationPort:    "validDestination",
		DestinationChannel: "DestinationChannel",
		Sequence:           7,
	}

	msg := orderPath(packet)
	expected := "channel/validChannel/port/validPort/channel/DestinationChannel/port/validDestination/sequence/7"
	require.Equal(t, expected, msg)
}

func TestGenerateOrderID(t *testing.T) {
	packet := channeltypes.Packet{
		SourcePort:         "validPort",
		SourceChannel:      "validChannel",
		DestinationPort:    "validDestination",
		DestinationChannel: "DestinationChannel",
		Sequence:           7,
	}
	id := GenerateOrderId(packet)
	expected := "\xaa\x8d>WR\xdd\xcd\xc5r%\x16\xb5]\x14\xf6\xef\xf4]\xe5;\xeby\x9c\x8b\xb7\xc3\xcd\x12\xea\xe8\x7f\xe4"
	require.Equal(t, expected, id)

}

func TestCreateOrder(t *testing.T) {

	packet := channeltypes.Packet{
		SourcePort:         "validPort",
		SourceChannel:      "validChannel",
		DestinationPort:    "validDestination",
		DestinationChannel: "DestinationChannel",
		Sequence:           7,
	}
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr.String(), addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())
	msg.Packet = packet

	expectedOrder := Order{
		Id:     "\xaa\x8d>WR\xdd\xcd\xc5r%\x16\xb5]\x14\xf6\xef\xf4]\xe5;\xeby\x9c\x8b\xb7\xc3\xcd\x12\xea\xe8\x7f\xe4",
		Path:   "channel/validChannel/port/validPort/channel/DestinationChannel/port/validDestination/sequence/7",
		Status: Status_INITIAL,
		Maker:  msg,
	}

	order := CreateOrder(msg)
	require.Equal(t, expectedOrder, order)

}

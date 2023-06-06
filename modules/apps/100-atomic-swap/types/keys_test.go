package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// Test that there is domain separation between the port id and the channel id otherwise an
// escrow address may overlap with another channel end
func TestGetEscrowAddress(t *testing.T) {
	var (
		port1    = "ibcswap"
		channel1 = "channel"
		port2    = "swapcha"
		channel2 = "nnel"
	)

	escrow1 := types.GetEscrowAddress(port1, channel1)
	escrow2 := types.GetEscrowAddress(port2, channel2)
	require.NotEqual(t, escrow1, escrow2)
}

package types

import (
	"testing"
)

const (
	denom              = "transfer/gaiachannel/atom"
	amount             = "100"
	denom1             = "uosmo"
	amount1            = "500"
	counterparty       = ""
	largeAmount        = "18446744073709551616"                                                           // one greater than largest uint64 (^uint64(0))
	invalidLargeAmount = "115792089237316195423570985008687907853269984665640564039457584007913129639936" // 2^256
)

// TestFungibleTokenPacketDataValidateBasic tests ValidateBasic for FungibleTokenPacketData
func TestFungibleTokenPacketDataValidateBasic(t *testing.T) {
	// testCases := []struct {
	// 	name       string
	// 	packetData SwapPacketData
	// 	expPass    bool
	// }{
	// 	{"valid packet", NewAtomicSwapPacketData(denom, amount, denom1, amount1, addr1, addr2, counterparty), true},
	// 	{"valid packet with large amount", NewAtomicSwapPacketData(denom, largeAmount, denom1, amount1, addr1, addr2, counterparty), true},
	// 	{"invalid denom", NewAtomicSwapPacketData("", amount, denom1, amount1, addr1, addr2, counterparty), false},
	// 	{"invalid empty amount", NewAtomicSwapPacketData(denom, "", denom1, amount1, addr1, addr2, counterparty), false},
	// 	{"invalid zero amount", NewAtomicSwapPacketData(denom, "0", denom1, amount1, addr1, addr2, counterparty), false},
	// 	{"invalid negative amount", NewAtomicSwapPacketData(denom, "-1", denom1, amount1, addr1, addr2, counterparty), false},
	// 	{"invalid large amount", NewAtomicSwapPacketData(denom, invalidLargeAmount, denom1, amount1, addr1, addr2, counterparty), false},
	// 	{"missing sender address", NewAtomicSwapPacketData(denom, amount, denom1, amount1, emptyAddr, addr2, counterparty), false},
	// 	{"missing recipient address", NewAtomicSwapPacketData(denom, amount, denom1, amount1, addr1, emptyAddr, counterparty), false},
	// }

	// for i, tc := range testCases {
	// 	err := tc.packetData.ValidateBasic()
	// 	if tc.expPass {
	// 		require.NoError(t, err, "valid test case %d failed: %v", i, err)
	// 	} else {
	// 		require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
	// 	}
	// }
}

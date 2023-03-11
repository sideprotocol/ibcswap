package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

// define constants used for testing
const (
	validPort        = "testportid"
	invalidPort      = "(invalidport1)"
	invalidShortPort = "p"
	// 195 characters
	invalidLongPort = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eros neque, ultricies vel ligula ac, convallis porttitor elit. Maecenas tincidunt turpis elit, vel faucibus nisl pellentesque sodales"

	validChannel        = "testchannel"
	invalidChannel      = "(invalidchannel1)"
	invalidShortChannel = "invalid"
	invalidLongChannel  = "invalidlongchannelinvalidlongchannelinvalidlongchannelinvalidlongchannel"
)

var (
	addr1     = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addr2     = sdk.AccAddress("testaddr2").String()
	emptyAddr string

	coin             = sdk.NewCoin("atom", sdk.NewInt(100))
	coin2            = sdk.NewCoin("osmo", sdk.NewInt(500))
	ibcCoin          = sdk.NewCoin("ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", sdk.NewInt(100))
	invalidIBCCoin   = sdk.NewCoin("ibc/7F1D3FCF4AE79E1554", sdk.NewInt(100))
	invalidDenomCoin = sdk.Coin{Denom: "0atom", Amount: sdk.NewInt(100)}
	zeroCoin         = sdk.Coin{Denom: "atoms", Amount: sdk.NewInt(0)}
	orderId          = "validOder"
	timeoutHeight    = clienttypes.NewHeight(0, 10)
)

// TestMsgSwapRoute tests Route for MsgSwap
func TestMsgSwapRoute(t *testing.T) {
	testCases := []struct {
		name     string
		route    string
		expRoute string
	}{
		{
			"new msg make swap route",
			NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()).Route(),
			RouterKey,
		},
		{
			"new msg take swap route",
			NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()).Route(),
			RouterKey,
		},
		{
			"new msg cancel swap route",
			NewMsgCancelSwap(validPort, validChannel, addr1, "", timeoutHeight, 0).Route(),
			RouterKey,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expRoute, tc.route, fmt.Sprintf("Test passed to %s", tc.name))
	}
}

// TestMsgSwapType tests Type for MsgSwap
func TestMsgSwapType(t *testing.T) {
	testCases := []struct {
		name       string
		msgType    string
		expMsgType string
	}{
		{
			"new msg make swap message type",
			NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()).Type(),
			"make_swap",
		},
		{
			"new msg take swap message type",
			NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()).Type(),
			"take_swap",
		},
		{
			"new msg cancel swap message type",
			NewMsgCancelSwap(validPort, validChannel, addr1, "", timeoutHeight, 0).Type(),
			"cancel_swap",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expMsgType, tc.msgType, fmt.Sprintf("Test passed to %s", tc.name))
	}
}

func TestMsgMakeSwapGetSignBytes(t *testing.T) {
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())
	expected := fmt.Sprintf(`{"type":"cosmos-sdk/MsgMakeSwap","value":{"receiver":"%s","sender":"%s","source_channel":"testchannel","source_port":"testportid","timeout_height":{"revision_height":"10"},"token":{"amount":"100","denom":"atom"}}}`, addr2, addr1)
	require.NotPanics(t, func() {
		res := msg.GetSignBytes()
		require.Equal(t, expected, string(res))
	})
}

func TestMsgTakeSwapGetSignBytes(t *testing.T) {
	timestamp := time.Now().UTC().Unix()
	msg := NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, timestamp)
	expected := fmt.Sprintf(`{"type":"cosmos-sdk/MsgTakeSwap","value":{"create_timestamp":"%v","sell_token":{"amount":"100","denom":"atom"},"source_channel":"testchannel","source_port":"testportid","taker_address":"%s","taker_receiving_address":"%s","timeout_height":{"revision_height":"10"}}}`, timestamp, addr2, addr1)
	require.NotPanics(t, func() {
		res := msg.GetSignBytes()
		require.Equal(t, expected, string(res))
	})
}

func TestMsgCancelSwapGetSignBytes(t *testing.T) {
	msg := NewMsgCancelSwap(validPort, validChannel, addr1, orderId, timeoutHeight, 0)
	expected := fmt.Sprintf(`{"type":"cosmos-sdk/MsgCancelSwap","value":{"maker_address":"%s","order_id":"%s","source_channel":"testchannel","source_port":"testportid","timeout_height":{"revision_height":"10"}}}`, addr1, orderId)
	require.NotPanics(t, func() {
		res := msg.GetSignBytes()
		require.Equal(t, expected, string(res))
	})
}

// TestMsgSwapValidation tests ValidateBasic for MsgTransfer
func TestMsgMakeSwapValidation(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *MsgMakeSwapRequest
		expPass bool
	}{
		{"valid msg with base denom", NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"valid msg with trace hash", NewMsgMakeSwap(validPort, validChannel, ibcCoin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"too short port id", NewMsgMakeSwap(invalidShortPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too long port id", NewMsgMakeSwap(invalidLongPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"port id contains non-alpha", NewMsgMakeSwap(invalidPort, validChannel, coin, coin2, addr1, "", addr2, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too short channel id", NewMsgMakeSwap(validPort, invalidShortChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too long channel id", NewMsgMakeSwap(validPort, invalidLongChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"channel id contains non-alpha", NewMsgMakeSwap(validPort, invalidChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"invalid denom", NewMsgMakeSwap(validPort, validChannel, invalidDenomCoin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"zero coin", NewMsgMakeSwap(validPort, validChannel, zeroCoin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"missing sender address", NewMsgMakeSwap(validPort, validChannel, coin, coin2, emptyAddr, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"missing recipient address", NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, "", "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"empty coin", NewMsgMakeSwap(validPort, validChannel, sdk.Coin{}, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %s", i, tc.name)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

func TestMsgTakeSwapValidation(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *MsgTakeSwapRequest
		expPass bool
	}{
		{"valid msg with base denom", NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"valid msg with trace hash", NewMsgTakeSwap(validPort, validChannel, ibcCoin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"too short port id", NewMsgTakeSwap(invalidShortPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too long port id", NewMsgTakeSwap(invalidLongPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"port id contains non-alpha", NewMsgTakeSwap(invalidPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too short channel id", NewMsgTakeSwap(validPort, invalidShortChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"too long channel id", NewMsgTakeSwap(validPort, invalidLongChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"channel id contains non-alpha", NewMsgTakeSwap(validPort, invalidChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"invalid denom", NewMsgTakeSwap(validPort, validChannel, invalidDenomCoin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"zero coin", NewMsgTakeSwap(validPort, validChannel, zeroCoin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"missing sender address", NewMsgTakeSwap(validPort, validChannel, coin, emptyAddr, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"missing recipient address", NewMsgTakeSwap(validPort, validChannel, coin, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
		{"empty coin", NewMsgTakeSwap(validPort, validChannel, sdk.Coin{}, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()), false},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %s", i, tc.name)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

// TestMsgCancelSwapValidation tests ValidateBasic for MsgTransfer
func TestMsgCancelSwapValidation(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *MsgCancelSwapRequest
		expPass bool
	}{
		{"valid msg with base denom", NewMsgCancelSwap(validPort, validChannel, addr1, orderId, timeoutHeight, 0), true},
		{"valid msg with trace hash", NewMsgCancelSwap(validPort, validChannel, addr1, orderId, timeoutHeight, 0), true},
		{"too short port id", NewMsgCancelSwap(invalidShortPort, validChannel, addr1, orderId, timeoutHeight, 0), false},
		{"too long port id", NewMsgCancelSwap(invalidLongPort, validChannel, addr1, orderId, timeoutHeight, 0), false},
		{"port id contains non-alpha", NewMsgCancelSwap(invalidPort, validChannel, addr1, orderId, timeoutHeight, 0), false},
		{"too short channel id", NewMsgCancelSwap(validPort, invalidShortChannel, addr1, orderId, timeoutHeight, 0), false},
		{"too long channel id", NewMsgCancelSwap(validPort, invalidLongChannel, addr1, orderId, timeoutHeight, 0), false},
		{"channel id contains non-alpha", NewMsgCancelSwap(validPort, invalidChannel, addr1, orderId, timeoutHeight, 0), false},
		{"missing sender address", NewMsgCancelSwap(validPort, validChannel, emptyAddr, orderId, timeoutHeight, 0), false},
		{"OrderId is required", NewMsgCancelSwap(validPort, validChannel, addr1, "", timeoutHeight, 0), false},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %s", i, tc.name)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

// TestMsgSwapGetSigners tests GetSigners for MsgTransfer
func TestMsgSwapGetSigners(t *testing.T) {
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	testCases := []struct {
		name      string
		signer    []sdk.AccAddress
		expSigner []sdk.AccAddress
	}{
		{
			"new msg make swap get signers",
			NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr.String(), addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()).GetSigners(),
			[]sdk.AccAddress{addr},
		},
		{
			"new msg take swap get signers",
			NewMsgTakeSwap(validPort, validChannel, coin, addr.String(), addr1, timeoutHeight, 0, time.Now().UTC().Unix()).GetSigners(),
			[]sdk.AccAddress{addr},
		},
		{
			"new msg cancel swap get signers",
			NewMsgCancelSwap(validPort, validChannel, addr.String(), "", timeoutHeight, 0).GetSigners(),
			[]sdk.AccAddress{addr},
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expSigner, tc.signer, fmt.Sprintf("Test passed to %s", tc.name))
	}
}

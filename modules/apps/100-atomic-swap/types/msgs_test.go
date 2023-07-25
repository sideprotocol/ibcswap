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

	timeoutHeight = clienttypes.NewHeight(0, 10)
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
		//{
		//	"new msg take swap route",
		//	NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()).Route(),
		//	RouterKey,
		//},
		{
			"new msg cancel swap route",
			NewMsgCancelSwap(addr1, "", timeoutHeight, 0).Route(),
			RouterKey,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expRoute, tc.route, fmt.Sprintf("Test passed to %s", tc.name))
	}
}

// TestMsgSwapType tests Type for MsgSwap
func TestMsgSwapType(t *testing.T) {
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())

	require.Equal(t, "make_swap", msg.Type())
}

func TestMsgSwapGetSignBytes(t *testing.T) {
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())
	expected := fmt.Sprintf(`{"type":"cosmos-sdk/MsgMakeSwap","value":{"receiver":"%s","sender":"%s","source_channel":"testchannel","source_port":"testportid","timeout_height":{"revision_height":"10"},"token":{"amount":"100","denom":"atom"}}}`, addr2, addr1)
	require.NotPanics(t, func() {
		res := msg.GetSignBytes()
		require.Equal(t, expected, string(res))
	})
}

// TestMsgSwapValidation tests ValidateBasic for MsgTransfer
func TestMsgSwapValidation(t *testing.T) {
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
		//{
		//	"new msg take swap message type",
		//	NewMsgTakeSwap(validPort, validChannel, coin, addr2, addr1, timeoutHeight, 0, time.Now().UTC().Unix()).Type(),
		//	"take_swap",
		//},
		{
			"new msg cancel swap message type",
			NewMsgCancelSwap(addr1, "", timeoutHeight, 0).Type(),
			"cancel_swap",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expMsgType, tc.msgType, fmt.Sprintf("Test passed to %s", tc.name))
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
		//{
		//	"new msg take swap get signers",
		//	NewMsgTakeSwap(validPort, validChannel, coin, addr.String(), addr1, timeoutHeight, 0, time.Now().UTC().Unix()).GetSigners(),
		//	[]sdk.AccAddress{addr},
		//},
		{
			"new msg cancel swap get signers",
			NewMsgCancelSwap(addr.String(), "", timeoutHeight, 0).GetSigners(),
			[]sdk.AccAddress{addr},
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expSigner, tc.signer, fmt.Sprintf("Test passed to %s", tc.name))
	}
}

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
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgSwapType tests Type for MsgSwap
func TestMsgSwapType(t *testing.T) {
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())

	require.Equal(t, "swap", msg.Type())
}

func TestMsgSwapGetSignBytes(t *testing.T) {
	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())
	expected := fmt.Sprintf(`{"type":"cosmos-sdk/MsgTransfer","value":{"receiver":"%s","sender":"%s","source_channel":"testchannel","source_port":"testportid","timeout_height":{"revision_height":"10"},"token":{"amount":"100","denom":"atom"}}}`, addr2, addr1)
	require.NotPanics(t, func() {
		res := msg.GetSignBytes()
		require.Equal(t, expected, string(res))
	})
}

// TestMsgSwapValidation tests ValidateBasic for MsgTransfer
func TestMsgSwapValidation(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *MsgMakeSwapRequest
		expPass bool
	}{
		{"valid msg with base denom", NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"valid msg with trace hash", NewMsgMakeSwap(validPort, validChannel, ibcCoin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), true},
		{"invalid ibc denom", NewMsgMakeSwap(validPort, validChannel, invalidIBCCoin, coin2, addr1, addr2, "", timeoutHeight, 0, time.Now().UTC().Unix()), false},
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

// TestMsgSwapGetSigners tests GetSigners for MsgTransfer
func TestMsgSwapGetSigners(t *testing.T) {
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	msg := NewMsgMakeSwap(validPort, validChannel, coin, coin2, addr.String(), addr2, "", timeoutHeight, 0, time.Now().UTC().Unix())
	res := msg.GetSigners()

	require.Equal(t, []sdk.AccAddress{addr}, res)
}

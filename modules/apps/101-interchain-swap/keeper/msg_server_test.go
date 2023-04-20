package keeper_test

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/go-bip39"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestMsgInterchainSwap() {

	// var msg *types.MsgMakeSwapRequest
	// testCases := []struct {
	// 	name     string
	// 	malleate func()
	// 	expPass  bool
	// }{
	// 	{
	// 		"success",
	// 		func() {},
	// 		true,
	// 	},
	// 	// {
	// 	// 	"invalid sender",
	// 	// 	func() {
	// 	// 		msg.MakerAddress = "address"
	// 	// 	},
	// 	// 	false,
	// 	// },
	// 	//{
	// 	//	"sender is a blocked address",
	// 	//	func() {
	// 	//		msg.SenderAddress = suite.chainA.GetSimApp().AccountKeeper.GetModuleAddress(types.ModuleName).String()
	// 	//	},
	// 	//	false,
	// 	//},
	// 	// {
	// 	// 	"channel does not exist",
	// 	// 	func() {
	// 	// 		msg.SourceChannel = "channel-100"
	// 	// 	},
	// 	// 	false,
	// 	// },
	// }

	// for _, tc := range testCases {
	// 	suite.SetupTest()

	// 	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	// 	suite.coordinator.Setup(path)

	// 	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
	// 	msg = types.NewMsgMakeSwap(
	// 		path.EndpointA.ChannelConfig.PortID,
	// 		path.EndpointA.ChannelID,
	// 		coin, coin,
	// 		suite.chainA.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(),
	// 		suite.chainB.SenderAccount.GetAddress().String(),
	// 		suite.chainB.GetTimeoutHeight(), 0, // only use timeout height
	// 		time.Now().UTC().Unix(),
	// 	)

	// 	tc.malleate()

	// 	res, err := suite.chainA.GetSimApp().AtomicSwapKeeper.MakeSwap(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

	// 	if tc.expPass {
	// 		suite.Require().NoError(err)
	// 		suite.Require().NotNil(res)
	// 	} else {
	// 		suite.Require().Error(err)
	// 		suite.Require().Nil(res)
	// 	}
	// }
}

func (suite *KeeperTestSuite) TestMsgDoubleDeposit() {
	denomPair := []string{"aside", "bside"}
	pooId, err := suite.SetupPoolWithDenomPair(denomPair)
	suite.Require().NoError(err)
	fmt.Println(pooId)

	ctx := suite.chainB.GetContext()

	// Add stake tokens to the sender accounts
	bankAKeeper := suite.chainA.GetSimApp().BankKeeper
	bankBKeeper := suite.chainB.GetSimApp().BankKeeper
	stakeATokens := sdk.NewCoins(sdk.NewCoin(denomPair[0], sdk.NewInt(2000)))
	stakeBTokens := sdk.NewCoins(sdk.NewCoin(denomPair[1], sdk.NewInt(2000)))
	bankAKeeper.MintCoins(suite.chainA.GetContext(), types.ModuleName, stakeATokens)
	bankBKeeper.MintCoins(suite.chainB.GetContext(), types.ModuleName, stakeBTokens)
	bankAKeeper.SendCoinsFromModuleToAccount(suite.chainA.GetContext(), types.ModuleName, suite.chainA.SenderAccount.GetAddress(), stakeATokens)
	bankBKeeper.SendCoinsFromModuleToAccount(suite.chainB.GetContext(), types.ModuleName, suite.chainB.SenderAccount.GetAddress(), stakeBTokens)

	nonce := suite.chainB.SenderAccount.GetSequence()

	remoteDepositTx := types.RemoteDeposit{
		Sequence: nonce,
		Sender:   suite.chainB.SenderAccount.GetAddress().String(),
		Token:    &sdk.Coin{Denom: denomPair[1], Amount: sdk.NewInt(1000)},
	}

	rawDepositTx := types.ModuleCdc.MustMarshal(&remoteDepositTx)

	if err != nil {
		fmt.Println("Marshal Error:", err)
	}
	signedDepositTx, err := suite.chainB.SenderPrivKey.Sign(rawDepositTx)
	suite.Require().NoError(err)
	isValid := verifySignedMessage(rawDepositTx, signedDepositTx, suite.chainB.SenderAccount.GetPubKey())
	suite.Require().Equal(isValid, true)

	msg := types.NewMsgDoubleDeposit(
		*pooId,
		[]string{suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String()},
		[]*sdk.Coin{{Denom: denomPair[0], Amount: sdk.NewInt(1000)}, {Denom: denomPair[1], Amount: sdk.NewInt(1000)}},
		signedDepositTx,
	)

	res, err := suite.chainB.GetSimApp().InterchainSwapKeeper.OnDoubleDepositReceived(
		ctx, msg,
	)
	suite.Require().NoError(err)
	_ = res
}

func (suite *KeeperTestSuite) TestVerifyRawTx() {

	// create pool first of all.
	denomPair := []string{"aside", "bside"}
	bankTxMsg := banktypes.NewMsgSend(
		suite.chainA.SenderAccount.GetAddress(),
		suite.chainA.SenderAccount.GetAddress(),
		sdk.NewCoins(sdk.Coin{Denom: denomPair[1], Amount: sdk.NewInt(1000)}),
	)

	rawData, err := json.Marshal([]interface{}{bankTxMsg, 1})
	if err != nil {
		fmt.Println("Marshal Error:", err)
	}

	suite.Require().NoError(err)
	_, _, priv, err := createNewWallet()
	suite.Require().NoError(err)
	signedTx, err := priv.Sign(rawData)
	suite.Require().NoError(err)
	pubKey := priv.PubKey()
	suite.Require().Equal(verifySignedMessage(rawData, signedTx, pubKey), true)
}

func createNewWallet() (sdk.AccAddress, string, crypto.PrivKey, error) {
	// Generate a new random mnemonic
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to generate new entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to generate new mnemonic: %w", err)
	}

	// Create a BIP39 seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Create a BIP32 master key from the seed
	masterKey, _ := hd.ComputeMastersFromSeed(seed)

	// Create a secp256k1 private key from the master key

	privKey := secp256k1.GenPrivKeyFromSecret(masterKey[:])

	// Get the address associated with the private key
	addr := sdk.AccAddress(privKey.PubKey().Address().Bytes())

	return addr, mnemonic, privKey, nil
}

func verifySignedMessage(rawTx []byte, signedMessage []byte, publicKey crypto.PubKey) bool {
	// Replace this with the actual rawTx bytes you want to verify
	return publicKey.VerifySignature(rawTx, signedMessage)
}

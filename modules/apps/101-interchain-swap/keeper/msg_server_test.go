package keeper_test

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
	// 	{
	// 		"invalid sender",
	// 		func() {
	// 			msg.MakerAddress = "address"
	// 		},
	// 		false,
	// 	},
	// 	//{
	// 	//	"sender is a blocked address",
	// 	//	func() {
	// 	//		msg.SenderAddress = suite.chainA.GetSimApp().AccountKeeper.GetModuleAddress(types.ModuleName).String()
	// 	//	},
	// 	//	false,
	// 	//},
	// 	{
	// 		"channel does not exist",
	// 		func() {
	// 			msg.SourceChannel = "channel-100"
	// 		},
	// 		false,
	// 	},
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

	// 	res, err := suite.chainA.GetSimApp().IBCSwapKeeper.MakeSwap(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

	// 	if tc.expPass {
	// 		suite.Require().NoError(err)
	// 		suite.Require().NotNil(res)
	// 	} else {
	// 		suite.Require().Error(err)
	// 		suite.Require().Nil(res)
	// 	}
	// }
}

package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	//ibctesting "github.com/sideprotocol/ibcswap/v4/testing"
)

func (suite *KeeperTestSuite) TestMsgCreatePool() {
	var msg *types.MsgCreatePoolRequest

	testCases := []struct {
		name    string
		expPass bool
	}{
		{
			"success",
			true,
		},
		// {
		// 	"channel does not exist",
		// 	false,
		// },
	}

	for _, tc := range testCases {
		suite.SetupTest()

		path := NewInterchainSwapPath(suite.chainA, suite.chainB)
		suite.coordinator.Setup(path)

		msg = types.NewMsgCreatePool(
			path.EndpointA.ChannelConfig.PortID,
			path.EndpointA.ChannelID,
			suite.chainA.SenderAccount.GetAddress().String(),
			"1:2",
			[]string{sdk.DefaultBondDenom, "venuscoin"},
			[]uint32{10, 100},
		)
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().IBCInterchainSwapKeeper)
		res, err := msgSrv.CreatePool(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}

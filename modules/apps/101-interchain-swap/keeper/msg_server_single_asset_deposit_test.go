package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
)

func (suite *KeeperTestSuite) SetupPool() (*string, error) {
	suite.SetupTest()
	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	msg := types.NewMsgCreatePool(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		suite.chainA.SenderAccount.GetAddress().String(),
		"50:50",
		[]*sdk.Coin{{
			Denom:  sdk.DefaultBondDenom,
			Amount: sdk.NewInt(1000),
		}, {
			Denom:  "bside",
			Amount: sdk.NewInt(1000),
		}},
		[]uint32{10, 100},
	)

	ctx := suite.chainA.GetContext()
	suite.chainA.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctx, msg)
	poolId := types.GetPoolIdWithTokens(msg.Tokens)
	return &poolId, nil
}

func (suite *KeeperTestSuite) SetupPoolWithDenomPair(denomPair []string) (*string, error) {
	if len(denomPair) != 2 {
		return nil, fmt.Errorf("invalid denom length")
	}
	suite.SetupTest()
	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	msg := types.NewMsgCreatePool(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		suite.chainA.SenderAccount.GetAddress().String(),
		"50:50",
		[]*sdk.Coin{{
			Denom:  denomPair[0],
			Amount: sdk.NewInt(1000),
		}, {
			Denom:  denomPair[1],
			Amount: sdk.NewInt(1000),
		}},
		[]uint32{10, 100},
	)

	ctxA := suite.chainA.GetContext()
	suite.chainA.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctxA, msg)
	ctxB := suite.chainB.GetContext()
	suite.chainB.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctxB, msg)
	poolId := types.GetPoolIdWithTokens(msg.Tokens)
	return &poolId, nil
}

func (suite *KeeperTestSuite) TestMsgDeposit() {
	var msg *types.MsgSingleAssetDepositRequest
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"invalid address",
			func() {
				msg.Sender = "invalid address"
			},
			false,
		},
		{
			"invalid amount",
			func() {
				msg.Sender = sample.AccAddress()
			},
			false,
		},
	}

	for _, tc := range testCases {
		// create pool first of all.
		pooId, err := suite.SetupPool()
		suite.Require().NoError(err)
		fmt.Println(pooId)

		//
		msg = types.NewMsgSingleAssetDeposit(
			*pooId,
			suite.chainA.SenderAccount.GetAddress().String(),
			&sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)},
		)

		tc.malleate()
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().InterchainSwapKeeper)

		res, err := msgSrv.SingleAssetDeposit(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}

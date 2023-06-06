package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/sideprotocol/ibcswap/v6/testing/testutil/sample"
)

func (suite *KeeperTestSuite) SetupPool() (*string, error) {
	suite.SetupTest()
	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	msg := types.NewMsgCreatePool(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		[]byte("0"),
		types.PoolAsset{
			Side: types.PoolAssetSide_SOURCE,
			Balance: &sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(1000),
			},
			Weight:  50,
			Decimal: 6,
		},

		types.PoolAsset{
			Side: types.PoolAssetSide_SOURCE,
			Balance: &sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(1000),
			},
			Weight:  50,
			Decimal: 6,
		},
		300,
	)

	ctx := suite.chainA.GetContext()
	suite.chainA.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctx, msg)
	poolId := types.GetPoolId(msg.GetLiquidityDenoms())
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
		suite.chainB.SenderAccount.GetAddress().String(),
		[]byte("0"),
		types.PoolAsset{
			Side: types.PoolAssetSide_SOURCE,
			Balance: &sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(1000),
			},
			Weight:  50,
			Decimal: 6,
		},

		types.PoolAsset{
			Side: types.PoolAssetSide_SOURCE,
			Balance: &sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(1000),
			},
			Weight:  50,
			Decimal: 6,
		},
		300,
	)

	ctxA := suite.chainA.GetContext()
	suite.chainA.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctxA, msg)
	ctxB := suite.chainB.GetContext()
	suite.chainB.GetSimApp().InterchainSwapKeeper.OnCreatePoolAcknowledged(ctxB, msg)
	poolId := types.GetPoolId(msg.GetLiquidityDenoms())
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

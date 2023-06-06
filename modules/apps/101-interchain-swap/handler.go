package interchainswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func NewMarketFeeUpdateProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.MarketFeeUpdateProposal:
			return keeper.HandleMarketFeeUpdateProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}

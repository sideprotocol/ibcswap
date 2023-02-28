package testsuite

import (
	"context"

	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	types "github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/ibc"
)

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryClientStatus(ctx context.Context, chain ibc.Chain, clientID string) (string, error) {
	queryClient := s.GetChainGRCPClients(chain).ClientQueryClient
	res, err := queryClient.ClientStatus(ctx, &clienttypes.QueryClientStatusRequest{
		ClientId: clientID,
	})
	if err != nil {
		return "", err
	}
	return res.Status, nil
}

// QueryPacketCommitment queries the packet commitment on the given chain for the provided channel and sequence.
func (s *E2ETestSuite) QueryPacketCommitment(ctx context.Context, chain ibc.Chain, portID, channelID string, sequence uint64) ([]byte, error) {
	queryClient := s.GetChainGRCPClients(chain).ChannelQueryClient
	res, err := queryClient.PacketCommitment(ctx, &channeltypes.QueryPacketCommitmentRequest{
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  sequence,
	})
	if err != nil {
		return nil, err
	}
	return res.Commitment, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryInterchainswapPool(ctx context.Context, chain ibc.Chain, poolID string) (*types.QueryGetInterchainLiquidityPoolResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).InterchainQueryClient
	res, err := queryClient.InterchainLiquidityPool(
		ctx,
		&types.QueryGetInterchainLiquidityPoolRequest{
			PoolId: poolID,
		},
	)

	if err != nil {
		return res, err
	}
	return res, nil
}

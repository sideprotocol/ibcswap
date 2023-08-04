package testsuite

import (
	"context"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesbeta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	intertxtypes "github.com/cosmos/interchain-accounts/x/inter-tx/types"
	atomicswaptypes "github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v6/ibc"

	controllertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	feetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// QueryClientState queries the client state on the given chain for the provided clientID.
func (s *E2ETestSuite) QueryClientState(ctx context.Context, chain ibc.Chain, clientID string) (ibcexported.ClientState, error) {
	queryClient := s.GetChainGRCPClients(chain).ClientQueryClient
	res, err := queryClient.ClientState(ctx, &clienttypes.QueryClientStateRequest{
		ClientId: clientID,
	})
	if err != nil {
		return nil, err
	}

	clientState, err := clienttypes.UnpackClientState(res.ClientState)
	if err != nil {
		return nil, err
	}

	return clientState, nil
}

// QueryChannel queries the channel on a given chain for the provided portID and channelID
func (s *E2ETestSuite) QueryChannel(ctx context.Context, chain ibc.Chain, portID, channelID string) (channeltypes.Channel, error) {
	queryClient := s.GetChainGRCPClients(chain).ChannelQueryClient
	res, err := queryClient.Channel(ctx, &channeltypes.QueryChannelRequest{
		PortId:    portID,
		ChannelId: channelID,
	})
	if err != nil {
		return channeltypes.Channel{}, err
	}

	return *res.Channel, nil
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

// QueryPacketAck queries the packet ack on the given chain for the provided channel and sequence.
func (s *E2ETestSuite) QueryPacketAcknowledged(ctx context.Context, chain ibc.Chain, portID, channelID string, sequence uint64) ([]byte, error) {
	queryClient := s.GetChainGRCPClients(chain).ChannelQueryClient
	res, err := queryClient.PacketAcknowledgement(ctx, &channeltypes.QueryPacketAcknowledgementRequest{
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  sequence,
	})
	if err != nil {
		return nil, err
	}
	return res.Acknowledgement, nil
}

// QueryInterchainAccount queries the interchain account for the given owner and connectionID.
func (s *E2ETestSuite) QueryInterchainAccount(ctx context.Context, chain ibc.Chain, owner, connectionID string) (string, error) {
	queryClient := s.GetChainGRCPClients(chain).ICAQueryClient
	res, err := queryClient.InterchainAccount(ctx, &controllertypes.QueryInterchainAccountRequest{
		Owner:        owner,
		ConnectionId: connectionID,
	})
	if err != nil {
		return "", err
	}
	return res.Address, nil
}

// QueryInterchainAccountLegacy queries the interchain account for the given owner and connectionID using the intertx module.
func (s *E2ETestSuite) QueryInterchainAccountLegacy(ctx context.Context, chain ibc.Chain, owner, connectionID string) (string, error) {
	queryClient := s.GetChainGRCPClients(chain).InterTxQueryClient
	res, err := queryClient.InterchainAccount(ctx, &intertxtypes.QueryInterchainAccountRequest{
		Owner:        owner,
		ConnectionId: connectionID,
	})
	if err != nil {
		return "", err
	}

	return res.InterchainAccountAddress, nil
}

// QueryIncentivizedPacketsForChannel queries the incentivized packets on the specified channel.
func (s *E2ETestSuite) QueryIncentivizedPacketsForChannel(
	ctx context.Context,
	chain *cosmos.CosmosChain,
	portId,
	channelId string,
) ([]*feetypes.IdentifiedPacketFees, error) {
	queryClient := s.GetChainGRCPClients(chain).FeeQueryClient
	res, err := queryClient.IncentivizedPacketsForChannel(ctx, &feetypes.QueryIncentivizedPacketsForChannelRequest{
		PortId:    portId,
		ChannelId: channelId,
	})
	if err != nil {
		return nil, err
	}
	return res.IncentivizedPackets, err
}

// QueryCounterPartyPayee queries the counterparty payee of the given chain and relayer address on the specified channel.
func (s *E2ETestSuite) QueryCounterPartyPayee(ctx context.Context, chain ibc.Chain, relayerAddress, channelID string) (string, error) {
	queryClient := s.GetChainGRCPClients(chain).FeeQueryClient
	res, err := queryClient.CounterpartyPayee(ctx, &feetypes.QueryCounterpartyPayeeRequest{
		ChannelId: channelID,
		Relayer:   relayerAddress,
	})
	if err != nil {
		return "", err
	}
	return res.CounterpartyPayee, nil
}

// QueryProposal queries the governance proposal on the given chain with the given proposal ID.
func (s *E2ETestSuite) QueryProposal(ctx context.Context, chain ibc.Chain, proposalID uint64) (govtypesbeta1.Proposal, error) {
	queryClient := s.GetChainGRCPClients(chain).GovQueryClient
	res, err := queryClient.Proposal(ctx, &govtypesbeta1.QueryProposalRequest{
		ProposalId: proposalID,
	})
	if err != nil {
		return govtypesbeta1.Proposal{}, err
	}

	return res.Proposal, nil
}

func (s *E2ETestSuite) QueryProposalV1(ctx context.Context, chain ibc.Chain, proposalID uint64) (govtypesv1.Proposal, error) {
	queryClient := s.GetChainGRCPClients(chain).GovQueryClientV1
	res, err := queryClient.Proposal(ctx, &govtypesv1.QueryProposalRequest{
		ProposalId: proposalID,
	})
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	return *res.Proposal, nil
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

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryInterchainswapPools(ctx context.Context, chain ibc.Chain) (*types.QueryAllInterchainLiquidityPoolResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).InterchainQueryClient
	res, err := queryClient.InterchainLiquidityPoolAll(ctx, &types.QueryAllInterchainLiquidityPoolRequest{})

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryLatestMultiDepositOrder(ctx context.Context, chain ibc.Chain, poolId, srcMaker string) (*types.QueryGetInterchainMultiDepositOrderResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).InterchainQueryClient
	res, err := queryClient.InterchainLatestMultiDepositOrderByCreator(ctx, &types.QueryLatestInterchainMultiDepositOrderBySourceMakerRequest{
		PoolId:      poolId,
		SourceMaker: srcMaker,
	})

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryInterchainMultiDepositOrders(ctx context.Context, chain ibc.Chain, poolId string) (*types.QueryAllInterchainMultiDepositOrdersResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).InterchainQueryClient
	res, err := queryClient.InterchainMultiDepositOrdersAll(ctx, &types.QueryAllInterchainMultiDepositOrdersRequest{
		PoolId: poolId,
	})

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryAtomicswapOrders(ctx context.Context, chain ibc.Chain) (*atomicswaptypes.QueryOrdersResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).AtomicQueryClient
	res, err := queryClient.GetAllOrders(
		ctx,
		&atomicswaptypes.QueryOrdersRequest{},
	)

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryAtomicswapByOrders(ctx context.Context, chain ibc.Chain, orderType atomicswaptypes.OrderType) (*atomicswaptypes.QueryOrdersResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).AtomicQueryClient
	res, err := queryClient.GetAllOrdersByType(
		ctx,
		&atomicswaptypes.QueryOrdersByRequest{
			OrderType: orderType,
		},
	)

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QuerySubmittedAtomicswap(ctx context.Context, chain ibc.Chain, makeAddress string) (*atomicswaptypes.QueryOrdersResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).AtomicQueryClient
	res, err := queryClient.GetSubmittedOrders(
		ctx,
		&atomicswaptypes.QuerySubmittedOrdersRequest{
			MakerAddress: makeAddress,
		},
	)

	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryBalance(ctx context.Context, chain ibc.Chain, addr string, denom string) (*banktypes.QueryBalanceResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).BankQueryClient

	res, err := queryClient.Balance(
		ctx,
		&banktypes.QueryBalanceRequest{
			Address: addr,
			Denom:   denom,
		},
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) GetBalance(ctx context.Context, chain ibc.Chain, addr string, denom string) (*sdk.Int, error) {
	res, err := s.QueryBalance(ctx, chain, addr, denom)
	if err != nil {
		return nil, err
	}
	return &res.Balance.Amount, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryAccount(ctx context.Context, chain ibc.Chain, addr string) (*authtypes.QueryAccountResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).AccountQueryClient

	res, err := queryClient.Account(
		ctx,
		&authtypes.QueryAccountRequest{
			Address: addr,
		},
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

// QueryClientStatus queries the status of the client by clientID
func (s *E2ETestSuite) QueryModuleAccount(ctx context.Context, chain ibc.Chain, name string) (*authtypes.QueryModuleAccountByNameResponse, error) {
	queryClient := s.GetChainGRCPClients(chain).AccountQueryClient

	res, err := queryClient.ModuleAccountByName(
		ctx,
		&authtypes.QueryModuleAccountByNameRequest{
			Name: name,
		},
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

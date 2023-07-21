package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channelutils "github.com/cosmos/ibc-go/v6/modules/core/04-channel/client/utils"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func GetTimeOuts(clientCtx client.Context, srcPort, scrChannel, timeoutHeightStr string, timeoutTimestamp uint64, absoluteTimeouts bool) (*clienttypes.Height, *uint64, error) {
	// if the timeouts are not absolute, retrieve latest block height and block timestamp
	// for the consensus state connected to the destination port/channel
	timeoutHeight, err := clienttypes.ParseHeight(timeoutHeightStr)
	if err != nil {
		return nil, nil, err
	}
	if !absoluteTimeouts {
		consensusState, height, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, scrChannel)
		if err != nil {
			return nil, nil, err
		}

		if !timeoutHeight.IsZero() {
			absoluteHeight := height
			absoluteHeight.RevisionNumber += timeoutHeight.RevisionNumber
			absoluteHeight.RevisionHeight += timeoutHeight.RevisionHeight
			timeoutHeight = absoluteHeight
		}

		if timeoutTimestamp != 0 {
			// use local clock time as reference time if it is later than the
			// consensus state timestamp of the counter party chain, otherwise
			// still use consensus state timestamp as reference
			now := time.Now().UnixNano()
			consensusStateTimestamp := consensusState.GetTimestamp()
			if now > 0 {
				now := uint64(now)
				if now > consensusStateTimestamp {
					timeoutTimestamp = now + timeoutTimestamp
				} else {
					timeoutTimestamp = consensusStateTimestamp + timeoutTimestamp
				}
			} else {
				return nil, nil, errors.New("local clock time is not greater than Jan 1st, 1970 12:00 AM")
			}
		}
	}
	return &timeoutHeight, &timeoutTimestamp, nil
}

// QueryPool fetches the pool information from the chain using the client context
func QueryPool(clientCtx client.Context, poolId string) (*types.InterchainLiquidityPool, error) {
	fmt.Println("PoolID:", poolId)
	queryClient := types.NewQueryClient(clientCtx)
	params := &types.QueryGetInterchainLiquidityPoolRequest{
		PoolId: poolId,
	}
	res, err := queryClient.InterchainLiquidityPool(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return &res.InterchainLiquidityPool, nil
}

package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

func CmdShowInterchainMultiDepositOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-interchain-multi-deposit-order [pool-id] [order-id]",
		Short: "lists all InterchainMultiDepositOrder according to the poolId",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPoolId := args[0]
			argOrderId := args[1]
			orderId, err := strconv.Atoi(argOrderId)
			if err != nil {
				return err
			}

			params := &types.QueryGetInterchainMultiDepositOrderRequest{
				PoolId:  argPoolId,
				OrderId: uint64(orderId),
			}

			res, err := queryClient.InterchainMultiDepositOrder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetInterchainMultiDepositOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-interchain-multi-deposit-order [pool-id]",
		Short: "lists all InterchainMultiDepositOrder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPoolId := args[0]

			params := &types.QueryAllInterchainMultiDepositOrdersRequest{
				PoolId: argPoolId,
			}

			res, err := queryClient.InterchainMultiDepositOrdersAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

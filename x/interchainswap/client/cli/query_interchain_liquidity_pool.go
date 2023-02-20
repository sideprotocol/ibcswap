package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
	"github.com/spf13/cobra"
)

func CmdListInterchainLiquidityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-interchain-liquidity-pool",
		Short: "list all InterchainLiquidityPool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInterchainLiquidityPoolRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InterchainLiquidityPoolAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowInterchainLiquidityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-interchain-liquidity-pool [pool-id]",
		Short: "shows a InterchainLiquidityPool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPoolId := args[0]

			params := &types.QueryGetInterchainLiquidityPoolRequest{
				PoolId: argPoolId,
			}

			res, err := queryClient.InterchainLiquidityPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

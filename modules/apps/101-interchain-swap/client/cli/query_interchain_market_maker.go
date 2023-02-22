package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

func CmdListInterchainMarketMaker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-interchain-market-maker",
		Short: "list all InterchainMarketMaker",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInterchainMarketMakerRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InterchainMarketMakerAll(context.Background(), params)
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

func CmdShowInterchainMarketMaker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-interchain-market-maker [pool-id]",
		Short: "shows a InterchainMarketMaker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPoolId := args[0]

			params := &types.QueryGetInterchainMarketMakerRequest{
				PoolId: argPoolId,
			}

			res, err := queryClient.InterchainMarketMaker(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"
)

// GetCmdParams returns the command handler for atomic swap parameter querying.
func GetCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current ibc-swap parameters",
		Long:    "Query the current ibc-swap parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query ibc-swap params", version.AppName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryEscrowAddress returns the command handler for ibc-swap parameter querying.
func GetCmdQueryEscrowAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "escrow-address",
		Short:   "Get the escrow address for a channel",
		Long:    "Get the escrow address for a channel",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query ibc-swap escrow-address [port] [channel-id]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			port := args[0]
			channel := args[1]
			addr := types.GetEscrowAddress(port, channel)
			return clientCtx.PrintString(fmt.Sprintf("%s\n", addr.String()))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

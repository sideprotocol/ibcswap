package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSingleAssetWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "single withdraw [demon] [pool coins]",
		Short: "Broadcast message Withdraw",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argOutDenoms := args[1]
			argCoin := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			coins, err := GetTokens(argCoin)
			if err != nil {
				return nil
			}
			if len(coins) != 2 {
				return fmt.Errorf("invalid token length! : %d", len(coins))
			}

			msg := types.NewMsgSingleAssetWithdraw(
				clientCtx.GetFromAddress().String(),
				argOutDenoms,
				coins[0],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

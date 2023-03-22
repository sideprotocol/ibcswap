package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [swap_type] [sender] [slippage] [recipient]",
		Short: "Broadcast message Swap",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			swapTypeArg := args[0]

			swapType := types.SwapMsgType_LEFT
			switch swapTypeArg {
			case "right":
				swapType = types.SwapMsgType_RIGHT
			case "left":
				swapType = types.SwapMsgType_LEFT
			default:
				return fmt.Errorf("invalid swap type:: %s, please try 'left' or 'right' only", swapTypeArg)
			}

			argSender := args[1]
			argSlippage, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argRecipient := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(argSender)
			argTokenIn := args[3]
			argTokenOut := args[4]

			tokenIn, err := GetTokens(argTokenIn)
			if err != nil {
				return err
			}

			tokenOut, err := GetTokens(argTokenOut)
			if err != nil {
				return err
			}

			msg := types.NewMsgSwap(
				swapType,
				clientCtx.GetFromAddress().String(),
				argSlippage,
				argRecipient,
				tokenIn[0],
				tokenOut[0],
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

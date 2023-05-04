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
		Use:   "swap [swap_type] [sender] [slippage] [recipient] [tokenIn] [tokenOut]",
		Short: "Broadcast message Swap",
		Args:  cobra.ExactArgs(6),
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
			argSlippage, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argRecipient := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(argSender)
			argTokenIn := args[4]
			argTokenOut := args[5]

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

			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")

			poolId := types.GetPoolId([]string{tokenIn[0].Denom, tokenOut[0].Denom})
			pool, err := QueryPool(clientCtx, poolId)
			if err != nil {
				return err
			}

			if err1 == nil && err2 == nil {
				timeoutHeight, timeoutTimestamp, err := GetTimeOuts(clientCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, packetTimeoutHeight, uint64(packetTimeoutTimestamp), false)
				fmt.Println("Timeout Height:", timeoutHeight)
				fmt.Println("Timeout Timestamp:", timeoutTimestamp)
				fmt.Println("Timeouts Err:", err)
				if err == nil {
					msg.TimeoutHeight = timeoutHeight
					msg.TimeoutTimeStamp = *timeoutTimestamp
				}
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("packet-timeout-height", "", "Packet timeout height")
	cmd.Flags().Uint("packet-timeout-timestamp", 0, "Packet timeout timestamp (in nanoseconds)")

	return cmd
}

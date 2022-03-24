package cmd

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func ContractEmergencyWithdrawFund() *cobra.Command {
	cmd := &cobra.Command{
		Use: "contract-emergency-withdraw-fund",
		Long: `Emergency withdraw funds if found any security risk.
Usage:
contract-emergency-withdraw-fund --chain [Chain] --contract-hash [Liquidity pool address] --tokens [List of tokens] --new-owner [New owner address] --index [Index of this message]

Example:
./sisu contract-emergency-withdraw-fund --chain ganache1 --contract-hash 0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C
--tokens 0xf0D676183dD5ae6b370adDdbE770235F23546f9d --new-owner 0x2d532C099CA476780c7703610D807948ae47856A 
--index=0 --from=node0 --keyring-backend test --chain-id=eth-sisu-local -y
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, _ := cmd.Flags().GetString(flags.Chain)
			liquidityAddress, _ := cmd.Flags().GetString(flags.ContractHash)
			tokenStr, _ := cmd.Flags().GetString(flags.Tokens)
			newOwner, _ := cmd.Flags().GetString(flags.NewOwner)
			index, _ := cmd.Flags().GetInt32(flags.Index)

			if len(chain) == 0 {
				return fmt.Errorf("invalid chain")
			}

			if len(liquidityAddress) == 0 {
				return fmt.Errorf("invalid name")
			}

			if len(tokenStr) == 0 {
				return fmt.Errorf("invalid tokens")
			}

			if len(newOwner) == 0 {
				return fmt.Errorf("invalid newOwner")

			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens := strings.Split(tokenStr, ",")
			msg := types.NewEmergencyWithdrawFundMsg(clientCtx.GetFromAddress().String(), chain, liquidityAddress, tokens, newOwner, index)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	sdkflags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(sdkflags.FlagChainID, "", "name of the sisu chain")
	cmd.Flags().String(flags.Chain, "", "target chain of the command")
	cmd.Flags().String(flags.NewOwner, "", "new owner address")
	cmd.Flags().Int32(flags.Index, 0, "index of the command. This index is used to differentiate calling this contract multiple times")
	cmd.Flags().String(flags.Tokens, "", "array of needed emergency withdraw tokens")
	cmd.Flags().String(flags.ContractHash, "", "liquidity pool address")

	return cmd
}

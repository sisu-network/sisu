package dev

import (
	"context"
	"time"

	"github.com/sisu-network/sisu/contracts/eth/vault"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sisu-network/lib/log"
	"github.com/spf13/cobra"

	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
)

type remoteCallCommand struct{}

func RemoteCall() *cobra.Command {
	cmd := &cobra.Command{
		Use: "remote-call",
		Long: `trigger remoteCall.
Usage:
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString(flags.Mnemonic)
			dst, _ := cmd.Flags().GetString(flags.Dst)
			genesisFolder, _ := cmd.Flags().GetString(flags.GenesisFolder)
			src, _ := cmd.Flags().GetString(flags.Src)
			recipient, _ := cmd.Flags().GetString(flags.Recipient)
			sisuRpc, _ := cmd.Flags().GetString(flags.SisuRpc)
			message, _ := cmd.Flags().GetString(flags.Message)
			gasLimit, _ := cmd.Flags().GetUint64(flags.GasLimit)

			clients := getEthClients([]string{src}, genesisFolder)
			if len(clients) == 0 {
				panic("There is no healthy client")
			}

			client := clients[0]
			command := &remoteCallCommand{}
			command.callRemoteCall(client, mnemonic, sisuRpc, src, dst, recipient, message, gasLimit)

			return nil
		},
	}

	cmd.Flags().String(flags.Mnemonic, "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum", "Mnemonic used to deploy the contract.")
	cmd.Flags().String(flags.Src, "ganache1", "Caller's chain where the caller sends message")
	cmd.Flags().String(flags.Dst, "ganache2", "App chain where the message is called")
	cmd.Flags().String(flags.Recipient, "0x123", "App address where the message is called")
	cmd.Flags().String(flags.Message, "hello world", "The message content")
	cmd.Flags().Uint64(flags.GasLimit, 50_000, "The gas limit of destination call")
	cmd.Flags().String(flags.SisuRpc, "0.0.0.0:9090", "URL to connect to Sisu. Please do NOT include http:// prefix")
	cmd.Flags().String(flags.GenesisFolder, "./misc/dev", "Genesis folder that contains configuration files.")

	return cmd
}

// swapFromEth creates an ETH transaction and sends to gateway contract.
func (c *remoteCallCommand) callRemoteCall(client *ethclient.Client, mnemonic string,
	sisuRpc, callerChain, appChain, app, message string, gasLimit uint64,
) {
	v := common.HexToAddress(getEthVaultAddress(context.Background(), callerChain, sisuRpc))
	contract, err := vault.NewVault(v, client)
	if err != nil {
		panic(err)
	}

	opts, err := getAuthTransactor(client, mnemonic)
	if err != nil {
		panic(err)
	}

	appChainId := libchain.GetChainIntFromId(appChain)
	appAddr := common.HexToAddress(app)

	var tx *ethtypes.Transaction
	tx, err = contract.RemoteCall(opts, appChainId, appAddr, []byte(message), gasLimit)
	if err != nil {
		panic(err)
	}

	waitTx, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	log.Info("txHash = ", waitTx.TxHash.Hex())

	time.Sleep(time.Second * 3)
}

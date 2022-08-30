package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/sisu-network/lib/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	sisu "github.com/sisu-network/sisu/x/sisu"
)

func decodeTxParams(abi abi.ABI, callData []byte) (map[string]interface{}, error) {
	if len(callData) < 4 {
		return nil, fmt.Errorf("decodeTxParams: call data size is smaller than 4")
	}

	txParams := map[string]interface{}{}
	m, err := abi.MethodById(callData[:4])
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := m.Inputs.UnpackIntoMap(txParams, callData[4:]); err != nil {
		log.Error(err)
		return nil, err
	}

	return txParams, nil
}

func parseTx() {
	client, err := ethclient.Dial("https://rpc.testnet.fantom.network")
	if err != nil {
		panic(err)
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(10009972))
	if err != nil {
		panic(err)
	}

	for _, ethTx := range block.Transactions() {
		if ethTx.Hash().String() != "0x9a30d8db1880149deab188721c1d151bb1403b34ab36544e2402d5661272e7ed" {
			continue
		}

		erc20gatewayContract := sisu.SupportedContracts[sisu.ContractVault]
		gwAbi := erc20gatewayContract.Abi

		callData := ethTx.Data()
		txParams, err := decodeTxParams(gwAbi, callData)
		if err != nil {
			panic(fmt.Errorf("Failed to decode tx params, err = %v", err))
		}

		_, err = ethTx.AsMessage(ethtypes.NewLondonSigner(ethTx.ChainId()), nil)
		if err != nil {
			panic(fmt.Errorf("Failed to convert to messages, err = %v", err))
		}

		_, ok := txParams["token"].(ethcommon.Address)
		if !ok {
			err := fmt.Errorf("cannot convert token to type ethcommon.Address: %v", txParams)
			panic(err)
		}

		_, ok = txParams["dstChain"].(string)
		if !ok {
			err := fmt.Errorf("cannot convert dstChain to type string: %v", txParams)
			panic(err)
		}
	}
}

func main() {
	parseTx()
}

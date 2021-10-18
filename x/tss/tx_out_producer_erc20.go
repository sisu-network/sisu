package tss

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *DefaultTxOutputProducer) createErc20ContractResponse(ethTx *ethTypes.Transaction, fromChain string) (*types.TxResponse, error) {
	erc20Contract := SupportedContracts[ContractErc20]

	contract := p.db.GetContractFromAddress(fromChain, ethTx.To().String())
	if contract == nil {
		utils.LogError("Cannot find contract at address", ethTx.To())
		return nil, fmt.Errorf("cannot find contract for tx sent to %s", ethTx.To())
	}

	fmt.Println("contract name = ", contract.Name)
	payload := ethTx.Data()
	fmt.Println("Hex payload = ", hex.EncodeToString(payload))

	abiMethods := erc20Contract.Abi.Methods
	var methodName string
	for name, method := range abiMethods {
		hash := crypto.Keccak256([]byte(method.Sig))
		if bytes.Compare(hash[:4], payload[:4]) == 0 {
			utils.LogVerbose("Found method: ", name, hex.EncodeToString(hash[:4]))
			methodName = name
			break
		}
	}

	if methodName == "" {
		return nil, fmt.Errorf("cannot find funcName")
	}

	params, err := abiMethods[methodName].Inputs.Unpack(payload[4:])
	if err != nil {
		return nil, err
	}

	fmt.Println("AAAAAA 000000")

	switch methodName {
	case MethodTransferOutFromContract:
		if len(params) != 4 {
			return nil, fmt.Errorf("transferOutFromContract expects 4 params")
		}

		fmt.Println("AAAAAA 1111111")

		fmt.Println("params =  ", params)

		// Creates a transferIn function in the other chain.
		toChain := params[1].(string)

		// TODO: Creates tx out for other chains.
		if utils.IsETHBasedChain(toChain) {
			toChainContract := p.db.GetContractFromHash(toChain, erc20Contract.AbiHash)
			if toChainContract == nil {
				return nil, fmt.Errorf("cannot find erc20 contract for toChain %s", toChain)
			}

			fmt.Println("toChain = ", toChain)

			assetId := toChain + "__" + (params[0].(ethcommon.Address)).Hex()
			fmt.Println("assetId = ", assetId)

			recipient := params[2].(string)
			fmt.Printf("recipient type = %T\n", recipient)

			amount := params[3]

			fmt.Println("recipient = ", recipient)
			fmt.Println("amount = ", amount)

			input, err := erc20Contract.Abi.Pack(MethodTransferIn, assetId, ethcommon.HexToAddress(recipient), amount)
			if err != nil {
				return nil, err
			}

			fmt.Println("AAAAAA 22222222")

			nonce := p.worldState.UseAndIncreaseNonce(toChain)
			if nonce < 0 {
				return nil, fmt.Errorf("cannont find nonce for chain %s", toChain)
			}

			fmt.Println("nonce in tx producer = ", nonce)

			rawTx := ethTypes.NewTransaction(
				uint64(nonce),
				ethcommon.HexToAddress(toChainContract.Address),
				big.NewInt(0),
				p.getGasLimit(toChain),
				p.getGasPrice(toChain),
				input,
			)
			bz, err := rawTx.MarshalBinary()
			if err != nil {
				return nil, err
			}

			return &types.TxResponse{
				OutChain: toChain,
				EthTx:    rawTx,
				RawBytes: bz,
			}, nil
		}

	}

	return nil, fmt.Errorf("unhandle case")
}

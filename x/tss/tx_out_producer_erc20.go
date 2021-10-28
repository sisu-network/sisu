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

	utils.LogInfo("Creating response erc20 transaction, contract name =", contract.Name)

	payload := ethTx.Data()

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

	utils.LogInfo("Found method name = ", methodName)

	params, err := abiMethods[methodName].Inputs.Unpack(payload[4:])
	if err != nil {
		return nil, err
	}

	utils.LogInfo("params = ", params)

	switch methodName {
	case MethodTransferOutFromContract:
		if len(params) != 4 {
			return nil, fmt.Errorf("transferOutFromContract expects 4 params")
		}

		// Creates a transferIn function in the other chain.
		toChain := params[1].(string)
		utils.LogInfo("toChain = ", toChain)

		// TODO: Creates tx out for other chains.
		if utils.IsETHBasedChain(toChain) {
			toChainContract := p.db.GetContractFromHash(toChain, erc20Contract.AbiHash)
			if toChainContract == nil {
				return nil, fmt.Errorf("cannot find erc20 contract for toChain %s", toChain)
			}

			assetId := fromChain + "__" + (params[0].(ethcommon.Address)).Hex()
			recipient := params[2].(string)
			amount := params[3]

			utils.LogInfo("assetId = ", assetId)
			utils.LogInfo("recipient = ", recipient)
			utils.LogInfo("amount = ", amount)

			input, err := erc20Contract.Abi.Pack(MethodTransferIn, assetId, ethcommon.HexToAddress(recipient), amount)
			if err != nil {
				return nil, err
			}

			nonce := p.worldState.UseAndIncreaseNonce(toChain)
			if nonce < 0 {
				return nil, fmt.Errorf("cannont find nonce for chain %s", toChain)
			}

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
	case MethodTransferOut:

	}

	return nil, fmt.Errorf("unhandle case")
}

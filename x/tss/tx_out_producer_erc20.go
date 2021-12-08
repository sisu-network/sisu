package tss

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *DefaultTxOutputProducer) createErc20ContractResponse(ctx sdk.Context, ethTx *ethTypes.Transaction, fromChain string) (*types.TxResponse, error) {
	erc20Contract := SupportedContracts[ContractErc20]

	contract := p.db.GetContractFromAddress(fromChain, ethTx.To().String())
	if contract == nil {
		log.Error("Cannot find contract at address", ethTx.To())
		return nil, fmt.Errorf("cannot find contract for tx sent to %s", ethTx.To())
	}

	log.Info("Creating response erc20 transaction, contract name =", contract.Name)

	payload := ethTx.Data()

	abiMethods := erc20Contract.Abi.Methods
	var methodName string
	for name, method := range abiMethods {
		hash := crypto.Keccak256([]byte(method.Sig))
		if bytes.Compare(hash[:4], payload[:4]) == 0 {
			log.Verbose("Found method: ", name, hex.EncodeToString(hash[:4]))
			methodName = name
			break
		}
	}

	if methodName == "" {
		log.Error("cannot find funcName")
		return nil, fmt.Errorf("cannot find funcName")
	}

	log.Info("Found method name = ", methodName)

	params, err := abiMethods[methodName].Inputs.Unpack(payload[4:])
	if err != nil {
		log.Error("cannot unpack data", err)
		return nil, err
	}

	log.Info("params = ", params)

	switch methodName {
	case MethodTransferOutFromContract:
		if len(params) != 4 {
			log.Error("transferOutFromContract expects 4 params")
			return nil, fmt.Errorf("transferOutFromContract expects 4 params")
		}

		// Creates a transferIn function in the other chain.
		toChain := params[1].(string)
		log.Info("toChain = ", toChain)

		// TODO: Creates tx out for other chains.
		if libchain.IsETHBasedChain(toChain) {
			toChainContract := p.kvStore.GetContractFromHash(ctx, toChain, erc20Contract.AbiHash)
			if toChainContract == nil {
				log.Error("cannot find erc20 contract for toChain %s", toChain)
				return nil, fmt.Errorf("cannot find erc20 contract for toChain %s", toChain)
			}

			assetId := fromChain + "__" + (params[0].(ethcommon.Address)).Hex()
			recipient := params[2].(string)
			amount := params[3]

			log.Info("assetId = ", assetId)
			log.Info("recipient = ", recipient)
			log.Info("amount = ", amount)

			input, err := erc20Contract.Abi.Pack(MethodTransferIn, assetId, ethcommon.HexToAddress(recipient), amount)
			if err != nil {
				log.Error("cannot pack abi", err)
				return nil, err
			}

			nonce := p.worldState.UseAndIncreaseNonce(toChain)
			if nonce < 0 {
				log.Error("cannot find nonce for chain %s", toChain)
				return nil, fmt.Errorf("cannot find nonce for chain %s", toChain)
			}

			rawTx := ethTypes.NewTransaction(
				uint64(nonce),
				ethcommon.HexToAddress(toChainContract.Address),
				big.NewInt(0),
				p.getGasLimit(toChain),
				p.getGasPrice(toChain),
				input,
			)

			log.Verbose("ERC20 producer: rawTx hash = ", rawTx.Hash())

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

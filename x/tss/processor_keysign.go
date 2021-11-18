package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	eTypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/chain"
	libchain "github.com/sisu-network/lib/chain"
)

// This function is called after dheart sends Sisu keysign result.
func (p *Processor) OnKeysignResult(result *htypes.KeysignResult) {
	// Post the keysign result to cosmos chain.
	msg := types.NewKeysignResult(p.appKeys.GetSignerAddress().String(), result.OutChain, result.OutHash, result.Success, result.Signature)
	go p.txSubmit.SubmitMessage(msg)

	// Sends it to deyes for deployment.
	if result.Success {
		// Find the tx in txout table

		txEntity := p.db.GetTxOutWithHash(result.OutChain, result.OutHash, false)
		if txEntity == nil {
			log.Error("Cannot find tx out with hash", result.OutHash)
		}

		tx := &etypes.Transaction{}
		if err := tx.UnmarshalBinary(txEntity.BytesWithoutSig); err != nil {
			log.Error("cannot unmarshal tx, err =", err)
			return
		}

		// Create full tx with signature.
		chainId := libchain.GetChainIntFromId(result.OutChain)
		signedTx, err := tx.WithSignature(etypes.NewEIP2930Signer(chainId), result.Signature)
		if err != nil {
			log.Error("cannot set signatuer for tx, err =", err)
			return
		}

		bz, err := signedTx.MarshalBinary()
		if err != nil {
			log.Error("cannot marshal tx")
			return
		}

		// Add the signature to txOuts
		p.db.UpdateTxOutSig(
			result.OutChain,
			result.OutHash,
			utils.KeccakHash32(string(bz)),
			result.Signature,
		)

		// If this is a contract deployment transaction, update the contract table with the hash of the
		// deployment tx bytes.
		isContractDeployment := chain.IsETHBasedChain(result.OutChain) && p.db.IsContractDeployTx(result.OutChain, result.OutHash)
		err = p.deploySignedTx(bz, result, isContractDeployment)
		if err != nil {
			log.Error("deployment error: ", err)
			return
		}
	} else {
		// TODO: handle failure case here.
	}
}

func (p *Processor) CheckKeysignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) DeliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}

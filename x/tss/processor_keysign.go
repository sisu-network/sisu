package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/chain"
	libchain "github.com/sisu-network/lib/chain"
)

// This function is called after dheart sends Sisu keysign result.
func (p *Processor) OnKeysignResult(result *htypes.KeysignResult) {
	// Post the keysign result to cosmos chain.
	request := result.Request
	msg := types.NewKeysignResult(p.appKeys.GetSignerAddress().String(), request.OutChain, request.OutHash, result.Success, result.Signature)
	go p.txSubmit.SubmitMessage(msg)

	// Sends it to deyes for deployment.
	if result.Success {
		// Find the tx in txout table
		txOut := p.privateDb.GetTxOut(request.OutChain, request.OutHash)
		if txOut == nil {
			log.Error("Cannot find tx out with hash", request.OutHash)
		}

		tx := &etypes.Transaction{}
		if err := tx.UnmarshalBinary(txOut.OutBytes); err != nil {
			log.Error("cannot unmarshal tx, err =", err)
			return
		}

		// Create full tx with signature.
		chainId := libchain.GetChainIntFromId(request.OutChain)
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

		// TODO: Broadcast the keysign result that includes this TxOutSig.
		// Save this to TxOutSig
		p.privateDb.SaveTxOutSig(&types.TxOutSig{
			Chain:       result.Request.OutChain,
			HashWithSig: signedTx.Hash().String(),
			HashNoSig:   result.Request.OutHash,
		})

		log.Info("signedTx hash = ", signedTx.Hash().String())

		// If this is a contract deployment transaction, update the contract table with the hash of the
		// deployment tx bytes.
		isContractDeployment := chain.IsETHBasedChain(request.OutChain) && txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT
		err = p.deploySignedTx(bz, result, isContractDeployment)
		if err != nil {
			log.Error("deployment error: ", err)
			return
		}
	} else {
		// TODO: handle failure case here.
	}
}

func (p *Processor) checkKeysignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) deliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}

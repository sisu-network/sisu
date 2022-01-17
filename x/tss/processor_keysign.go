package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/chain"
	libchain "github.com/sisu-network/lib/chain"
)

// This function is called after dheart sends Sisu keysign result.
func (p *Processor) OnKeysignResult(result *htypes.KeysignResult) {
	if result.Outcome == htypes.OutcometNotSelected {
		// Do nothing here if this node is not selected.
		return
	}

	// Post the keysign result to cosmos chain.
	request := result.Request

	for i, keysignMsg := range request.KeysignMessages {
		msg := types.NewKeysignResult(
			p.appKeys.GetSignerAddress().String(),
			keysignMsg.OutChain,
			keysignMsg.OutHash,
			result.Outcome == htypes.OutcomeSuccess,
			result.Signatures[i],
		)
		go p.txSubmit.SubmitMessage(msg)

		// Sends it to deyes for deployment.
		if result.Outcome == htypes.OutcomeSuccess {
			// Find the tx in txout table
			txOut := p.publicDb.GetTxOut(keysignMsg.OutChain, keysignMsg.OutHash)
			if txOut == nil {
				log.Error("Cannot find tx out with hash", keysignMsg.OutHash)
			}

			tx := &etypes.Transaction{}
			if err := tx.UnmarshalBinary(txOut.OutBytes); err != nil {
				log.Error("cannot unmarshal tx, err =", err)
				return
			}

			// Create full tx with signature.
			chainId := libchain.GetChainIntFromId(keysignMsg.OutChain)
			signedTx, err := tx.WithSignature(etypes.NewEIP2930Signer(chainId), result.Signatures[i])
			if err != nil {
				log.Error("cannot set signatuer for tx, err =", err)
				return
			}

			bz, err := signedTx.MarshalBinary()
			if err != nil {
				log.Error("cannot marshal tx")
				return
			}

			// // TODO: Broadcast the keysign result that includes this TxOutSig.
			// // Save this to TxOutSig
			p.privateDb.SaveTxOutSig(&types.TxOutSig{
				Chain:       keysignMsg.OutChain,
				HashWithSig: signedTx.Hash().String(),
				HashNoSig:   keysignMsg.OutHash,
			})

			log.Info("signedTx hash = ", signedTx.Hash().String())

			// If this is a contract deployment transaction, update the contract table with the hash of the
			// deployment tx bytes.
			isContractDeployment := chain.IsETHBasedChain(keysignMsg.OutChain) && txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT
			err = p.deploySignedTx(bz, keysignMsg.OutChain, result.Request.KeysignMessages[i].OutHash, isContractDeployment)
			if err != nil {
				log.Error("deployment error: ", err)
				return
			}

			// TODO: Check if we have any pending confirm tx that is waiting for this tx.
		} else {
			// TODO: handle failure case here.
		}
	}
}

func (p *Processor) deliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}

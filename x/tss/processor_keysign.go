package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	eTypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
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
			utils.LogError("Cannot find tx out with hash", result.OutHash)
		}

		tx := &etypes.Transaction{}
		if err := tx.UnmarshalBinary(txEntity.BytesWithoutSig); err != nil {
			utils.LogError("cannot unmarshal tx, err =", err)
			return
		}

		// Create full tx with signature.
		chainId := libchain.GetChainIntFromId(result.OutChain)
		signedTx, err := tx.WithSignature(etypes.NewEIP2930Signer(chainId), result.Signature)
		if err != nil {
			utils.LogError("cannot set signatuer for tx, err =", err)
			return
		}

		bz, err := signedTx.MarshalBinary()
		if err != nil {
			utils.LogError("cannot marshal tx")
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
		deployedResult, err := p.deploySignedTx(bz, result, isContractDeployment)
		if err != nil {
			utils.LogError("deployment error: ", err)
			return
		}

		if deployedResult == nil {
			utils.LogError("deployment result is nil")
			return
		}

		if deployedResult.Success {
			p.onTxDeployed(result.OutChain, result.OutHash, deployedResult, isContractDeployment)
		}
	} else {
		// TODO: handle failure case here.
	}
}

func (p *Processor) deploySignedTx(bz []byte, keysignResult *htypes.KeysignResult, isContractDeployment bool) (*eTypes.DispatchedTxResult, error) {
	utils.LogDebug("Sending final tx to the deyes for deployment for chain", keysignResult.OutChain)
	deyeClient := p.deyesClients[keysignResult.OutChain]

	pubkey := p.db.GetPubKey(keysignResult.OutChain)
	if pubkey == nil {
		return nil, fmt.Errorf("Cannot get pubkey for chain %s", keysignResult.OutChain)
	}

	if deyeClient != nil {
		return deyeClient.Dispatch(&eTypes.DispatchedTxRequest{
			IsEthContractDeployment: isContractDeployment,
			Chain:                   keysignResult.OutChain,
			Tx:                      bz,
			PubKey:                  pubkey,
		})
	} else {
		err := fmt.Errorf("Cannot find deyes client for chain %s", keysignResult.OutChain)
		return nil, err
	}
}

func (p *Processor) CheckKeysignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) DeliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}

func (p *Processor) onTxDeployed(chain, outHash string, deployResult *eTypes.DispatchedTxResult, isContractDeployment bool) {
	// If this is a ETH contract deployment, add the deployed address to the watch list.

	if isContractDeployment {
		// Add this to the watcher address.
		utils.LogInfo("Adding the deployment address to the watch addresss", deployResult.DeployedAddr)
		deyeClient := p.deyesClients[chain]
		deyeClient.AddWatchAddresses(chain, []string{deployResult.DeployedAddr})

		// Update database with the deployed address
		p.db.UpdateContractAddress(chain, outHash, deployResult.DeployedAddr)
	}
}

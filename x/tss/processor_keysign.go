package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	eTypes "github.com/sisu-network/deyes/types"
	dhTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"

	coreTypes "github.com/sisu-network/dcore/core/types"
)

// This function is called after dheart sends Sisu keysign result.
func (p *Processor) OnKeysignResult(result *dhTypes.KeysignResult) {
	// Post the keysign result to cosmos chain.
	msg := types.NewKeysignResult(p.appKeys.GetSignerAddress().String(), result.OutChain, result.OutHash, result.Success, result.Signature)
	go p.txSubmit.SubmitMessage(msg)

	// Sends it to deyes for deployment.
	if result.Success {
		tx := &coreTypes.Transaction{}
		if err := tx.UnmarshalBinary(result.OutBytes); err != nil {
			utils.LogError("cannot unmarshal tx, err =", err)
			return
		}

		chainId := utils.GetChainIntFromId(result.OutChain)
		signedTx, err := tx.WithSignature(coreTypes.NewEIP2930Signer(chainId), result.Signature)
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
		p.db.UpdateTxOut(result.OutChain, result.OutHash, bz, result.Signature)

		// If this is a contract deployment transaction, update the contract table with the hash of the
		// deployment tx bytes.
		isContractDeployment := utils.IsETHBasedChain(result.OutChain) && p.db.IsContractDeployTx(result.OutChain, result.OutHash)
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
	}
}

func (p *Processor) deploySignedTx(bz []byte, keysignResult *dhTypes.KeysignResult, isContractDeployment bool) (*eTypes.DispatchedTxResult, error) {
	utils.LogDebug("Sending final tx to the deyes for deployment")
	deyeClient := p.deyesClients[keysignResult.OutChain]

	if deyeClient != nil {
		return deyeClient.Dispatch(&eTypes.DispatchedTxRequest{
			IsEthContractDeployment: isContractDeployment,
			Chain:                   keysignResult.OutChain,
			Tx:                      bz,
			PubKey:                  p.storage.GetPubKey(keysignResult.OutChain),
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
	}
}

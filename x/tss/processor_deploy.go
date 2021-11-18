package tss

import (
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
)

// deploySignedTx creates a deployment request and sends it to deyes.
func (p *Processor) deploySignedTx(bz []byte, keysignResult *htypes.KeysignResult, isContractDeployment bool) error {
	log.Debug("Sending final tx to the deyes for deployment for chain", keysignResult.OutChain)
	deyeClient := p.deyesClients[keysignResult.OutChain]

	pubkey := p.db.GetPubKey(keysignResult.OutChain)
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", keysignResult.OutChain)
	}

	if deyeClient != nil {
		request := &eTypes.DispatchedTxRequest{
			Chain:                   keysignResult.OutChain,
			TxHash:                  keysignResult.OutHash,
			Tx:                      bz,
			PubKey:                  pubkey,
			IsEthContractDeployment: isContractDeployment,
		}

		go deyeClient.Dispatch(request)
	} else {
		err := fmt.Errorf("Cannot find deyes client for chain %s", keysignResult.OutChain)
		return err
	}

	return nil
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (p *Processor) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	chain := result.Chain
	outHash := result.TxHash
	isContractDeployment := result.IsEthContractDeployment

	log.Info("There is a deployment result")

	if result.Success {
		// If this is a ETH contract deployment, add the deployed address to the watch list.
		if isContractDeployment {
			// Add this to the watcher address.
			log.Info("Adding the deployment address to the watch addresss", result.DeployedAddr)
			deyeClient := p.deyesClients[chain]
			deyeClient.AddWatchAddresses(chain, []string{result.DeployedAddr})

			// Update database with the deployed address
			p.db.UpdateContractAddress(chain, outHash, result.DeployedAddr)
		}
	} else {
		// TODO: Handle deployment failure
	}
}

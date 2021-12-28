package tss

import (
	"fmt"

	etypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"

	libchain "github.com/sisu-network/lib/chain"
)

// deploySignedTx creates a deployment request and sends it to deyes.
func (p *Processor) deploySignedTx(bz []byte, keysignResult *htypes.KeysignResult, isContractDeployment bool) error {
	request := keysignResult.Request
	log.Debug("Sending final tx to the deyes for deployment for chain", request.OutChain)
	deyeClient := p.deyesClients[request.OutChain]

	pubkey := p.privateDb.GetKeygenPubkey(libchain.GetKeyTypeForChain(request.OutChain))
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", request.OutChain)
	}

	if deyeClient != nil {
		request := &etypes.DispatchedTxRequest{
			Chain:                   request.OutChain,
			TxHash:                  request.OutHash,
			Tx:                      bz,
			PubKey:                  pubkey,
			IsEthContractDeployment: isContractDeployment,
		}

		go deyeClient.Dispatch(request)
	} else {
		err := fmt.Errorf("Cannot find deyes client for chain %s", request.OutChain)
		return err
	}

	return nil
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (p *Processor) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	chain := result.Chain
	// outHash := result.TxHash
	isContractDeployment := result.IsEthContractDeployment

	log.Info("There is a deployment result")

	if result.Success {
		// If this is a ETH contract deployment, add the deployed address to the watch list.
		if isContractDeployment {
			// Add this to the watcher address.
			log.Info("Adding the deployment address to the watch addresss", result.DeployedAddr)
			deyeClient := p.deyesClients[chain]
			deyeClient.AddWatchAddresses(chain, []string{result.DeployedAddr})

			// TODO: Update database with the deployed address
			// p.db.UpdateContractAddress(chain, outHash, result.DeployedAddr)
		}
	} else {
		// TODO: Handle deployment failure
	}
}

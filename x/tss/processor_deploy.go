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
	log.Info("The transaction has been sent to blockchain (but not included in a block yet). chain = ",
		result.Chain, ", address = ", result.DeployedAddr)
}

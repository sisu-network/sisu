package sisu

import (
	"fmt"

	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"

	libchain "github.com/sisu-network/lib/chain"
)

// deploySignedTx creates a deployment request and sends it to deyes.
func (p *Processor) deploySignedTx(bz []byte, outChain string, outHash string, isContractDeployment bool) error {
	log.Debug("Sending final tx to the deyes for deployment for chain ", outChain)

	pubkey := p.publicDb.GetKeygenPubkey(libchain.GetKeyTypeForChain(outChain))
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", outChain)
	}

	request := &etypes.DispatchedTxRequest{
		Chain:                   outChain,
		TxHash:                  outHash,
		Tx:                      bz,
		PubKey:                  pubkey,
		IsEthContractDeployment: isContractDeployment,
	}

	go p.deyesClient.Dispatch(request)

	return nil
}

// OnTxDeploymentResult is a callback after there is a deployment result from deyes.
func (p *Processor) OnTxDeploymentResult(result *etypes.DispatchedTxResult) {
	log.Info("The transaction has been sent to blockchain (but not included in a block yet). chain = ",
		result.Chain, ", address = ", result.DeployedAddr)
}

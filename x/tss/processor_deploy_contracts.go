package tss

// This function
func (p *Processor) deployEthContracts(chain string) {
	// Prepare ETH deployment transactions
	rawTx := p.logic.PrepareEthContractDeployment(chain, 0)

	// Broadcast this rawTx to the network for voting.
	if rawTx != nil {

	}
}

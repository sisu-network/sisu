package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *Processor) deliverContracts(ctx sdk.Context, signerMsg *types.ContractsWithSigner) ([]byte, error) {
	if process, hash := p.shouldProcessMsg(ctx, signerMsg); process {
		p.doContracts(ctx, signerMsg)
		p.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

func (p *Processor) doContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	// TODO: Don't do duplicated delivery
	log.Info("Deliver pending contracts")

	for _, contract := range wrappedMsg.Data.Contracts {
		if p.publicDb.IsContractExisted(contract) {
			log.Infof("Contract %s has been processed", contract.Name)
			return nil, nil
		}
	}

	log.Info("Saving contracts, contracts length = ", len(wrappedMsg.Data.Contracts))

	// Save into KVStore & private db
	p.publicDb.SaveContracts(wrappedMsg.Data.Contracts, true)

	return nil, nil
}

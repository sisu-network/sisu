package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

// createPendingContracts creates and broadcast pending contracts. All nodes need to agree what
// contracts to deploy on what chains.
func (p *Processor) createPendingContracts(ctx sdk.Context, msg *types.Keygen) {
}

func (p *Processor) checkContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) error {
	for _, contract := range wrappedMsg.Data.Contracts {
		if !p.privateDb.IsContractExisted(contract) {
			return ErrCannotFindMessage
		}
	}

	// TODO: Check with KVStore

	return nil
}

func (p *Processor) deliverContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	// TODO: Don't do duplicated delivery
	log.Info("Deliver pending contracts")

	for _, contract := range wrappedMsg.Data.Contracts {
		if p.keeper.IsContractExisted(ctx, contract) {
			log.Infof("Contract %s has been processed", contract.Name)
			return nil, nil
		}
	}

	log.Info("Saving contracts, contracts length = ", len(wrappedMsg.Data.Contracts))

	// Save into KVStore & private db
	p.keeper.SaveContracts(ctx, wrappedMsg.Data.Contracts, true)
	p.privateDb.SaveContracts(wrappedMsg.Data.Contracts, true)

	return nil, nil
}

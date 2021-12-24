package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

// createPendingContracts creates and broadcast pending contracts. All nodes need to agree what
// contracts to deploy on what chains.
func (p *Processor) createPendingContracts(ctx sdk.Context, msg *types.KeygenResult) {
	contracts := make([]*types.Contract, 0)
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Symbol
		if libchain.GetKeyTypeForChain(chain) == msg.Keygen.KeyType {
			log.Info("Saving contracts for chain ", chain)

			for name, c := range SupportedContracts {
				contract := &types.Contract{
					Chain: chain,
					Hash:  c.AbiHash,
					Name:  name,
				}

				contracts = append(contracts, contract)
			}
		}
	}

	go func() {
		signer := p.appKeys.GetSignerAddress()
		p.txSubmit.SubmitMessage(types.NewContractsWithSigner(
			signer.String(),
			contracts,
		))
	}()
}

func (p *Processor) checkContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) error {
	// TODO: validate contracts to deploy here.
	return nil
}

func (p *Processor) deliverContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	// TODO: Don't do duplicated delivery
	log.Info("Deliver pending contracts")

	// Save this contract to keeper and our db.
	p.db.InsertContracts(wrappedMsg.Data.Contracts)

	// Save into KVStore
	p.keeper.SaveContracts(ctx, wrappedMsg.Data.Contracts)

	return nil, nil
}

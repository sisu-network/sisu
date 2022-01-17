package tss

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

// createPendingContracts creates and broadcast pending contracts. All nodes need to agree what
// contracts to deploy on what chains.
func (p *Processor) createPendingContracts(ctx sdk.Context, msg *types.Keygen) {
	log.Info("Create and broadcast contracts...")

	// We want the final contracts array to be deterministic. We need to sort the list of chains
	// and list of contracts alphabetically.
	// Sort all chains alphabetically.
	chains := make([]string, len(p.config.SupportedChains))
	i := 0
	for chain := range p.config.SupportedChains {
		chains[i] = chain
		i += 1
	}
	sort.Strings(chains)

	// Sort all contracts name alphabetically
	names := make([]string, len(SupportedContracts))
	i = 0
	for contract := range SupportedContracts {
		names[i] = contract
		i += 1
	}
	sort.Strings(names)

	// Create contracts
	contracts := make([]*types.Contract, 0)
	for _, chain := range chains {
		if libchain.GetKeyTypeForChain(chain) == msg.KeyType {
			log.Info("Saving contracts for chain ", chain)

			for _, name := range names {
				c := SupportedContracts[name]
				contract := &types.Contract{
					Chain:     chain,
					Hash:      c.AbiHash,
					Name:      name,
					ByteCodes: []byte(c.Bin),
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

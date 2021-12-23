package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	dhtypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"

	libchain "github.com/sisu-network/lib/chain"
)

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result dhtypes.KeygenResult) {
	// 1. Post result to the cosmos chain
	signer := p.appKeys.GetSignerAddress()

	resultEnum := types.KeygenResult_FAILURE
	if result.Success {
		resultEnum = types.KeygenResult_SUCCESS
	}

	wrappedMsg := types.NewKeygenResultWithSigner(signer.String(), result.KeyType, resultEnum, result.PubKeyBytes, result.Address)
	p.txSubmit.SubmitMessage(wrappedMsg)
	msg := wrappedMsg.Data

	// Update the address and pubkey of the keygen database.
	p.db.UpdateKeygenAddress(result.KeyType, result.Address, result.PubKeyBytes)

	// 2. Add the address to the watch list.
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Symbol
		deyesClient := p.deyesClients[chain]

		if libchain.GetKeyTypeForChain(chain) != msg.KeyType {
			continue
		}

		if deyesClient == nil {
			log.Critical("Cannot find deyes client for chain", chain)
		} else {
			log.Verbose("adding watcher address ", result.Address, " for chain ", chain)
			deyesClient.AddWatchAddresses(chain, []string{result.Address})
		}
	}
}

func (p *Processor) checkKeygenResult(ctx sdk.Context, wrappedMsg *types.KeygenResultWithSigner) error {
	msg := wrappedMsg.Data

	if msg.Result == types.KeygenResult_SUCCESS {
		if p.keeper.IsKeygenExisted(ctx, msg) {
			return ErrMessageHasBeenProcessed
		}

		return nil
	} else {
		// TODO: Process failure case. For failure case, we allow multiple message as each node can have
		// different blames.
	}

	return nil
}

func (p *Processor) deliverKeygenResult(ctx sdk.Context, wrappedMsg *types.KeygenResultWithSigner) ([]byte, error) {
	msg := wrappedMsg.Data

	if msg.Result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		// Save to KVStore
		p.keeper.SaveKeygen(ctx, msg)

		// We need to add this new watched address even though we are still catching up with blockchain.
		p.addWatchAddressAfterKeygen(ctx, msg)
	} else {
		// TODO: handle failure case
	}

	return nil, nil
}

func (p *Processor) addWatchAddressAfterKeygen(ctx sdk.Context, msg *types.KeygenResult) {
	// Check and see if we need to deploy some contracts. If we do, push them into the contract
	// queue for deployment later (after we receive some funding like ether to execute contract
	// deployment).
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Symbol
		if libchain.GetKeyTypeForChain(chain) == msg.KeyType {
			log.Info("Saving contracts for chain ", chain)
			p.txOutputProducer.SaveContractsToDeploy(chain)
		}
	}
}

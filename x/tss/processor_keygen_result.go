package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"

	libchain "github.com/sisu-network/lib/chain"
)

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

func (p *Processor) checkKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) error {
	if signerMsg.Data.Result == types.KeygenResult_SUCCESS {
		// Check if we have this data in our private db.
		if !p.privateDb.IsKeygenResultSuccess(signerMsg, p.appKeys.GetSignerAddress().String()) {
			log.Verbosef("Value does not match, data = %s %d %s", signerMsg.Keygen.KeyType, int(signerMsg.Keygen.Index), signerMsg.Data.From)
			return ErrValueDoesNotMatch
		}

		// TODO: Check if we have processed this message before.

		return nil
	} else {
		// TODO: Process failure case. For failure case, we allow multiple message as each node can have
		// different blames.
	}

	return nil
}

func (p *Processor) deliverKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) ([]byte, error) {
	msg := signerMsg.Data

	log.Info("Delivering keygen result, result = ", msg.Result)

	if msg.Result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		if p.keeper.IsKeygenResultSuccess(ctx, signerMsg, p.appKeys.GetSignerAddress().String()) {
			// This has been processed before.
			return nil, nil
		}

		log.Info("Saving keygen for ", signerMsg.Keygen.KeyType)

		// Save keygen to KVStore & private db
		p.keeper.SaveKeygen(ctx, signerMsg.Keygen)
		p.privateDb.SaveKeygen(signerMsg.Keygen)

		// Save result to KVStore & private db
		p.keeper.SaveKeygenResult(ctx, signerMsg)
		p.privateDb.SaveKeygenResult(signerMsg)

		// Add list the public key address to watch.
		p.addWatchAddress(signerMsg.Keygen)

		if !p.globalData.IsCatchingUp() {
			p.createPendingContracts(ctx, signerMsg.Keygen)
		}
	} else {
		// TODO: handle failure case
	}

	return nil, nil
}

func (p *Processor) addWatchAddress(msg *types.Keygen) {
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
			log.Verbose("adding watcher address ", msg.Address, " for chain ", chain)
			deyesClient.AddWatchAddresses(chain, []string{msg.Address})
		}
	}
}

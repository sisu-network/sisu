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

	// Update the address and pubkey of the keygen database.
	p.db.UpdateKeygenAddress(result.KeyType, result.Address, result.PubKeyBytes)
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

	if p.keeper.IsKeygenExisted(ctx, msg) {
		// This has been processed before.
		return nil, nil
	}

	if msg.Result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		// Save to KVStore
		p.keeper.SaveKeygen(ctx, msg)

		// Add list the public key address to watch.
		p.addWatchAddress(msg)

		if !p.globalData.IsCatchingUp() {
			p.createPendingContracts(ctx, msg)
		}
	} else {
		// TODO: handle failure case
	}

	return nil, nil
}

func (p *Processor) addWatchAddress(msg *types.KeygenResult) {
	// 2. Add the address to the watch list.
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Symbol
		deyesClient := p.deyesClients[chain]

		if libchain.GetKeyTypeForChain(chain) != msg.Keygen.KeyType {
			continue
		}

		if deyesClient == nil {
			log.Critical("Cannot find deyes client for chain", chain)
		} else {
			log.Verbose("adding watcher address ", msg.Keygen.Address, " for chain ", chain)
			deyesClient.AddWatchAddresses(chain, []string{msg.Keygen.Address})
		}
	}
}

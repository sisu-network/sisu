package tss

import (
	"fmt"

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

		// Update the address and pubkey of the keygen database.
		p.db.UpdateKeygenAddress(result.KeyType, result.Address, result.PubKeyBytes)

		// Update the result
		p.db.InsertKeygenResultSuccess(result.KeyType, result.KeygenIndex)
	}

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum)

	wrappedMsg := types.NewKeygenResultWithSigner(signer.String(), result.KeyType, resultEnum, result.PubKeyBytes, result.Address)
	p.txSubmit.SubmitMessage(wrappedMsg)
}

func (p *Processor) checkKeygenResult(ctx sdk.Context, wrappedMsg *types.KeygenResultWithSigner) error {
	keygenMsg := wrappedMsg.Keygen

	fmt.Println("Checking keygen result....")

	if wrappedMsg.Data.Result == types.KeygenResult_SUCCESS {
		// Check if we have this data in our private db.
		if p.db.GetKeygenResult(keygenMsg.KeyType, int(keygenMsg.Index)) != types.KeygenResult_SUCCESS {
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

	if p.keeper.IsKeygenResultSuccess(ctx, signerMsg) {
		// This has been processed before.
		fmt.Println("This has been processed before")
		return nil, nil
	}

	fmt.Println("Delivering keygen result")

	if msg.Result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		// Save to KVStore
		p.keeper.SaveKeygenResult(ctx, signerMsg)

		// Update the result
		p.db.InsertKeygenResultSuccess(signerMsg.Keygen.KeyType, int(signerMsg.Keygen.Index))

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

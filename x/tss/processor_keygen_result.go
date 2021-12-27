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
	resultEnum := types.KeygenResult_FAILURE
	if result.Success {
		resultEnum = types.KeygenResult_SUCCESS
	}

	signerMsg := types.NewKeygenResultWithSigner(
		p.appKeys.GetSignerAddress().String(),
		result.KeyType,
		result.KeygenIndex,
		resultEnum,
		result.PubKeyBytes,
		result.Address,
	)

	// Save the result to private db
	p.privateDb.SaveKeygenResult(signerMsg)

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum)

	p.txSubmit.SubmitMessage(signerMsg)
}

func (p *Processor) checkKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) error {
	fmt.Println("Checking keygen result....")

	if signerMsg.Data.Result == types.KeygenResult_SUCCESS {
		// Check if we have this data in our private db.
		if !p.privateDb.IsKeygenResultSuccess(signerMsg) {
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

	fmt.Println("Delivering keygen result")

	if msg.Result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		if p.keeper.IsKeygenResultSuccess(ctx, signerMsg) {
			// This has been processed before.
			fmt.Println("This has been processed before")
			return nil, nil
		}

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

package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dhtypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"

	libchain "github.com/sisu-network/lib/chain"
)

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result dhtypes.KeygenResult) {
	var resultEnum types.KeygenResult_Result
	switch result.Outcome {
	case dhtypes.OutcomeSuccess:
		resultEnum = types.KeygenResult_SUCCESS
	case dhtypes.OutcomeFailure:
		resultEnum = types.KeygenResult_FAILURE
	case dhtypes.OutcometNotSelected:
		resultEnum = types.KeygenResult_NOT_SELECTED
	}

	if resultEnum == types.KeygenResult_NOT_SELECTED {
		// No need to send result when this node is not selected.
		return
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
	p.publicDb.SaveKeygenResult(signerMsg)

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum)

	p.txSubmit.SubmitMessageAsync(signerMsg)
}

func (p *Processor) deliverKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) ([]byte, error) {
	if process, hash := p.shouldProcessMsg(ctx, signerMsg); process {
		p.doKeygenResult(ctx, signerMsg)
		p.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

func (p *Processor) doKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) ([]byte, error) {
	msg := signerMsg.Data

	log.Info("Delivering keygen result, result = ", msg.Result)

	result := p.getKeygenResult(ctx, signerMsg)

	// TODO: Get majority of the votes here.
	if result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		// Save result to KVStore & private db
		p.publicDb.SaveKeygenResult(signerMsg)

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

func (p *Processor) getKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) types.KeygenResult_Result {
	results := p.publicDb.GetAllKeygenResult(signerMsg.Keygen.KeyType, signerMsg.Keygen.Index)

	// Check the majority of the results
	successCount := 0
	for _, result := range results {
		if result.Data.Result == types.KeygenResult_SUCCESS {
			successCount += 1
		}
	}

	if successCount >= (len(results)+1)/2 {
		// TODO: Choose the address with most vote.
		// Save keygen Address
		log.Info("Saving keygen...")
		p.publicDb.SaveKeygen(signerMsg.Keygen)

		return types.KeygenResult_SUCCESS
	}

	return types.KeygenResult_FAILURE
}

func (p *Processor) addWatchAddress(msg *types.Keygen) {
	// 2. Add the address to the watch list.
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Symbol
		deyesClient := p.deyesClients[chain]

		if libchain.GetKeyTypeForChain(chain) != msg.KeyType {
			log.Info("!= msg.Keytype", chain, " ", msg.KeyType)
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

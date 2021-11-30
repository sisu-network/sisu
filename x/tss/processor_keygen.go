package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	dhTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	if p.globalData.IsCatchingUp() ||
		p.lastProposeBlockHeight != 0 && blockHeight-p.lastProposeBlockHeight <= ProposeBlockInterval {
		return
	}

	unavailableChains := make([]string, 0)
	for _, chainConfig := range p.config.SupportedChains {
		if !p.db.IsKeyExisted(chainConfig.Symbol) {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	// Broadcast a message.
	log.Info("Broadcasting TSS Keygen Proposal message. len(unavailableChains) = ", len(unavailableChains))
	signer := p.appKeys.GetSignerAddress()

	for _, chain := range unavailableChains {
		proposal := types.NewMsgKeygenProposal(
			signer.String(),
			chain,
			utils.GenerateRandomString(16),
			blockHeight,
		)
		log.Debug("Submitting proposal message for chain", chain)
		go func() {
			err := p.txSubmit.SubmitMessage(proposal)

			if err != nil {
				log.Error(err)
			}
		}()
	}

	p.lastProposeBlockHeight = blockHeight
}

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result dhTypes.KeygenResult) {
	// 1. Post result to the cosmos chain
	signer := p.appKeys.GetSignerAddress()

	resultEnum := types.KeygenResult_FAILURE
	if result.Success {
		resultEnum = types.KeygenResult_SUCCESS
	}

	msg := types.NewKeygenResult(signer.String(), result.Chain, resultEnum, result.PubKeyBytes, result.Address)
	p.txSubmit.SubmitMessage(msg)

	// 2. Add the address to the watch list.
	deyesClient := p.deyesClients[result.Chain]
	if deyesClient == nil {
		log.Critical("Cannot find deyes client for chain", result.Chain)
	} else {
		log.Verbose("adding watcher address", result.Address, "for chain", result.Chain)
		deyesClient.AddWatchAddresses(result.Chain, []string{result.Address})

		// Update the address and pubkey of the keygen database.
		p.db.UpdateKeygenAddress(result.Chain, result.Address, result.PubKeyBytes)
	}
}

func (p *Processor) CheckKeyGenProposal(msg *types.KeygenProposal) error {
	// TODO: Check if we see the same need to have keygen proposal here.
	return nil
}

func (p *Processor) DeliverKeyGenProposal(msg *types.KeygenProposal) ([]byte, error) {
	log.Info("Delivering keygen proposal")

	// TODO: Save data to KV store.
	if p.globalData.IsCatchingUp() {
		return nil, nil
	}

	if p.db.IsKeyExisted(msg.Chain) {
		log.Info("The keygen proposal has been processed")
		return nil, nil
	}

	err := p.db.CreateKeygen(msg.Chain)
	if err != nil {
		log.Error(err)
	}

	// Send a signal to Dheart to start keygen process.
	log.Info("Sending keygen request to Dheart. Chain =", msg.Chain)
	pubKeys := p.partyManager.GetActivePartyPubkeys()
	keygenId := GetKeygenId(msg.Chain, p.currentHeight, pubKeys)
	err = p.dheartClient.KeyGen(keygenId, msg.Chain, pubKeys)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Keygen request is sent successfully.")

	return []byte{}, nil
}

func (p *Processor) DeliverKeygenResult(ctx sdk.Context, msg *types.KeygenResult) ([]byte, error) {
	// TODO: Save data to KV store.
	if msg.Result == types.KeygenResult_SUCCESS {
		if status, _ := p.db.GetKeygenStatus(msg.Chain); status == db.StatusDeliveredToChain {
			log.Info("Keygen result has been processed for chain", msg.Chain)
			return nil, nil
		}

		// Update key address
		p.db.UpdateKeygenStatus(msg.Chain, db.StatusDeliveredToChain)

		// If this keygen is successful, prepare for contract deployment.
		// Save the pubkey to the keeper.
		p.keeper.SavePubKey(ctx, msg.Chain, msg.PubKeyBytes)

		// If this is a pubkey address of a ETH chain, save it to the store because we want to watch
		// transaction that funds the address (we will deploy contracts later).
		if chain.IsETHBasedChain(msg.Chain) {
			p.txOutputProducer.AddKeyAddress(ctx, msg.Chain, msg.Address)
		}

		// Check and see if we need to deploy some contracts. If we do, push them into the contract
		// queue for deployment later (after we receive some funding like ether to execute contract
		// deployment).
		p.txOutputProducer.SaveContractsToDeploy(msg.Chain)
	} else {
		// TODO: handle failure case
	}

	return nil, nil
}

package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
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

func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	if p.globalData.IsCatchingUp() ||
		p.lastProposeBlockHeight != 0 && blockHeight-p.lastProposeBlockHeight <= ProposeBlockInterval {
		return
	}

	// Check ECDSA only (for now)
	keyTypes := []string{libchain.KEY_TYPE_ECDSA}
	for _, keyType := range keyTypes {
		keygenEntity, err := p.db.GetKeyGen(keyType)
		if err != nil {
			log.Error("Cannot find keygen entity", err)
			continue
		}

		if keygenEntity != nil && keygenEntity.Status != "" {
			log.Info(keyType, "has been generated")
			continue
		}

		// Broadcast a message.
		signer := p.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenProposalWithSigner(
			signer.String(),
			keyType,
			utils.GenerateRandomString(16),
			blockHeight,
		)

		log.Info("Submitting proposal message for", keyType)
		go func() {
			err := p.txSubmit.SubmitMessage(proposal)

			if err != nil {
				log.Error(err)
			}
		}()
	}

	p.lastProposeBlockHeight = blockHeight
}

func (p *Processor) checkKeyGenProposal(ctx sdk.Context, wrapper *types.KeygenProposalWithSigner) error {
	msg := wrapper.Data

	if p.keeper.IsKeygenProposalExisted(ctx, msg) {
		log.Verbose("The keygen proposal has been processed")
		return ErrMessageHasBeenProcessed
	}

	return nil
}

func (p *Processor) deliverKeyGenProposal(ctx sdk.Context, wrapper *types.KeygenProposalWithSigner) ([]byte, error) {
	msg := wrapper.Data

	// TODO: Check if we have processed a keygen proposal recently.
	if p.keeper.IsKeygenProposalExisted(ctx, msg) {
		log.Verbose("The keygen proposal has been processed")
		return nil, nil
	}

	// Save this into Keeper
	p.keeper.SaveKeygenProposal(ctx, msg)

	// Save the keygen into the private db.
	blockHeight := p.currentHeight.Load().(int64)
	err := p.db.CreateKeygen(msg.KeyType, blockHeight)
	if err != nil {
		log.Error(err)
	}

	if p.globalData.IsCatchingUp() {
		return nil, nil
	}

	// Invoke TSS keygen in dheart
	p.doTss(msg, blockHeight)

	return []byte{}, nil
}

func (p *Processor) doTss(msg *types.KeygenProposal, blockHeight int64) {
	log.Info("Delivering keygen proposal")

	// Send a signal to Dheart to start keygen process.
	log.Info("Sending keygen request to Dheart. KeyType =", msg.KeyType)
	pubKeys := p.partyManager.GetActivePartyPubkeys()
	keygenId := GetKeygenId(msg.KeyType, blockHeight, pubKeys)

	err := p.dheartClient.KeyGen(keygenId, msg.KeyType, pubKeys)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Keygen request is sent successfully.")
}

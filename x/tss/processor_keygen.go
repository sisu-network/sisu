package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
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
	// TODO: We can replace this by sending command from client instead of running at the beginning
	// of each block.
	if p.globalData.IsCatchingUp() {
		return
	}

	// Check ECDSA only (for now)
	keyTypes := []string{libchain.KEY_TYPE_ECDSA}
	for _, keyType := range keyTypes {
		if p.keeper.IsKeygenExisted(ctx, keyType, 0) {
			continue
		}

		// Broadcast a message.
		signer := p.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenWithSigner(
			signer.String(),
			keyType,
			0,
		)

		// Create a new keygen entry in the db.
		p.privateDb.SaveKeygen(&types.Keygen{
			KeyType: keyType,
			Index:   0,
		})

		log.Info("Submitting proposal message for ", keyType)
		go func() {
			err := p.txSubmit.SubmitMessage(proposal)

			if err != nil {
				log.Error(err)
			}
		}()
	}
}

func (p *Processor) checkKeygen(ctx sdk.Context, wrapper *types.KeygenWithSigner) error {
	ok := p.privateDb.IsKeygenExisted(wrapper.Data.KeyType, int(wrapper.Data.Index))
	if !ok {
		return ErrCannotFindMessage
	}

	// TODO: Check if this is in the KVStore to avoid double processing.

	return nil
}

func (p *Processor) deliverKeygen(ctx sdk.Context, wrapper *types.KeygenWithSigner) ([]byte, error) {
	msg := wrapper.Data

	// TODO: Check if we have processed a keygen proposal recently.
	if p.keeper.IsKeygenExisted(ctx, msg.KeyType, int(msg.Index)) {
		log.Verbose("The keygen proposal has been processed")
		return nil, nil
	}

	log.Info("Delivering keygen....")

	// Save this into Keeper && private db.
	p.keeper.SaveKeygen(ctx, msg)
	p.privateDb.SaveKeygen(msg)

	if p.globalData.IsCatchingUp() {
		return nil, nil
	}

	// Invoke TSS keygen in dheart
	p.doTss(msg, ctx.BlockHeight())

	return []byte{}, nil
}

func (p *Processor) doTss(msg *types.Keygen, blockHeight int64) {
}

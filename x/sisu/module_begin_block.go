package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (a AppModule) beginBlock(ctx sdk.Context, blockHeight int64) {
	// Check keygen proposal
	if blockHeight > 1 {
		// We need to wait till block 2 for multistore of the app to be updated with latest account info
		// for signing.
		a.checkTssKeygen(ctx, blockHeight)
	}

	oldValue := a.globalData.IsCatchingUp()
	a.globalData.UpdateCatchingUp()
	newValue := a.globalData.IsCatchingUp()

	if oldValue && !newValue {
		log.Info("Setting Sisu readiness for dheart.")
		// This node has fully catched up with the blockchain, we need to inform dheart about this.
		a.mc.DheartClient().SetSisuReady(true)
		a.mc.DeyesClient().SetSisuReady(true)
	}
}

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/
func (a AppModule) checkTssKeygen(ctx sdk.Context, blockHeight int64) {
	// TODO: We can replace this by sending command from client instead of running at the beginning
	// of each block.
	if a.globalData.IsCatchingUp() || ctx.BlockHeight()%50 != 2 {
		return
	}

	keyTypes := []string{libchain.KEY_TYPE_ECDSA, libchain.KEY_TYPE_EDDSA}
	for _, keyType := range keyTypes {
		if a.keeper.IsKeygenExisted(ctx, keyType, 0) {
			continue
		}

		// Broadcast a message.
		signer := a.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenWithSigner(
			signer.String(),
			keyType,
			0,
		)

		log.Info("Submitting proposal message for ", keyType)
		a.txSubmit.SubmitMessageAsync(proposal)
	}
}

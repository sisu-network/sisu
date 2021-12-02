package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) OnPresignResult(result *htypes.PresignResult) {
	// Post the keysign result to cosmos chain.
	culprits := make([]*types.PartyID, len(result.Culprits))
	for i, culprit := range result.Culprits {
		culprits[i] = &types.PartyID{
			Id:      culprit.Id,
			Moniker: culprit.Moniker,
			Key:     culprit.Key,
			Index:   int32(culprit.Index),
		}
	}
	msg := types.NewPresignResult(
		p.appKeys.GetSignerAddress().String(),
		result.Chain,
		result.Success,
		result.PubKeyBytes,
		result.Address,
		culprits,
		)

	go p.txSubmit.SubmitMessage(msg)

	// Sends it to deyes for deployment.
	if result.Success {

		// TODO: handle success cases
	} else {
		// TODO: handle failure case here.
	}
}

func (p *Processor) CheckPresignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) DeliverPresignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.

	return nil, nil
}

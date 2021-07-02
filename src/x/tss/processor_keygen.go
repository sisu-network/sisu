package tss

import (
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) CheckKeyGenProposal(msg *types.KeygenProposal) error {
	// TODO: Check duplicated proposal here.
	return nil
}

func (p *Processor) DeliverKeyGenProposal(msg *types.KeygenProposal) ([]byte, error) {
	// 1. TODO: Check duplicated proposal here.

	// Just approve it for now.
	// 2. If this node supports the proposed chain and it's one of the top X validators, send an
	// approval vote to the keygen proposal.
	//    2a) Check this node is in the top N Validator
	//    2b) Check if this node supports chain X.
	supported := false
	for _, chainConfig := range p.config.SupportedChains {
		if chainConfig.Symbol == msg.ChainSymbol {
			supported = true
			break
		}
	}

	utils.LogDebug("Supported = ", supported)

	if !supported {
		// This is not supported by this current node
		return []byte{}, nil
	}

	// TODO: Save this proposal to KV store.
	utils.LogDebug("!p.appInfo.IsCatchingUp() = ", !p.appInfo.IsCatchingUp())

	if !p.appInfo.IsCatchingUp() {
		// Send vote message to everyone else
		signer := p.appKeys.GetSignerAddress()
		voteMsg := types.NewMsgKeygenProposalVote(signer.String(), msg.Id, types.KeygenProposalVote_APPROVE)

		utils.LogDebug("Sending this message...")

		go func() {
			err := p.txSubmit.SubmitMessage(voteMsg)
			if err != nil {
				utils.LogError(err)
			}
		}()
	}

	return []byte{}, nil
}

func (p *Processor) DeliverKeyGenProposalVote(msg *types.KeygenProposalVote) ([]byte, error) {
	return []byte{}, nil
}

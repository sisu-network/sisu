package tss

import (
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) CheckKeyGenProposal(msg *types.KeygenProposal) error {
	// TODO: Check duplicated proposal here.
	return nil
}

func (p *Processor) DeliverKeyGenProposal(msg *types.KeygenProposal) ([]byte, error) {
	// TODO: Check duplicated proposal here.
	// Just approve it for now.

	return []byte{}, nil
}

package service

import (
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type ChainPolling interface {
	// QueryRecentSolanBlock gets the recent solana transaction block hash and submit it to the Sisu
	// network.
	QueryRecentSolanBlock(chain string)
}

type defaultChainPolling struct {
	signer      string
	deyesClient external.DeyesClient
	txSubmit    common.TxSubmit
}

func NewChainPolling(signer string, deyesClient external.DeyesClient, txSubmit common.TxSubmit) ChainPolling {
	return &defaultChainPolling{
		signer:      signer,
		deyesClient: deyesClient,
		txSubmit:    txSubmit,
	}
}

func (p *defaultChainPolling) QueryRecentSolanBlock(chain string) {
	result, err := p.deyesClient.SolanaQueryRecentBlock(chain)
	if err != nil {
		log.Errorf("Failed to query solana recent block, err = ", err)
		return
	}

	// Broadcast result
	msg := types.NewUpdateSolanaRecentHashMsg(p.signer, chain, result.Hash, result.Height)
	p.txSubmit.SubmitMessageAsync(msg)
}

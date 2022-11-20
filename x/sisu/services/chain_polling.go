package services

import "github.com/sisu-network/sisu/x/sisu/external"

type ChainPolling interface {
	// QueryRecentSolanBlock gets the recent solana transaction block hash and submit it to the Sisu
	// network.
	QueryRecentSolanBlock()
}

type defaultChainPolling struct{}

func NewChainPolling(deyesClient external.DeyesClient) ChainPolling {
	return &defaultChainPolling{}
}

func (p *defaultChainPolling) QueryRecentSolanBlock() {

}

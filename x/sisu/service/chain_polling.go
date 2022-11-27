package service

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"go.uber.org/atomic"
)

var (
	PollingFrequency = 30 * 1000 // 30s
)

type ChainPolling interface {
	// QueryRecentSolanBlock gets the recent solana transaction block hash and submit it to the Sisu
	// network.
	Start(ctx sdk.Context, k keeper.Keeper)
	QueryRecentSolanBlock(chain string)
}

type defaultChainPolling struct {
	signer       string
	deyesClient  external.DeyesClient
	txSubmit     common.TxSubmit
	lastPollTime *atomic.Int64
}

func NewChainPolling(signer string, deyesClient external.DeyesClient, txSubmit common.TxSubmit) ChainPolling {
	return &defaultChainPolling{
		signer:       signer,
		deyesClient:  deyesClient,
		txSubmit:     txSubmit,
		lastPollTime: atomic.NewInt64(0),
	}
}

func (p *defaultChainPolling) Start(ctx sdk.Context, k keeper.Keeper) {
	solanaChain := ""
	// Start polling all solana chains
	params := k.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		if libchain.IsSolanaChain(chain) {
			solanaChain = chain
			break
		}
	}

	if solanaChain == "" {
		return
	}

	for {
		now := time.Now().UnixMilli()
		diff := now - int64(p.lastPollTime.Load())
		if now-int64(p.lastPollTime.Load()) < int64(PollingFrequency) {
			// We just poll recently, no need to poll again.
			sleepTime := PollingFrequency - int(diff)
			time.Sleep(time.Duration(sleepTime) * time.Microsecond)
			continue
		}

		p.QueryRecentSolanBlock(solanaChain)
	}
}

func (p *defaultChainPolling) QueryRecentSolanBlock(chain string) {
	now := time.Now().UnixMilli()
	p.lastPollTime.Store(now)

	result, err := p.deyesClient.SolanaQueryRecentBlock(chain)
	if err != nil {
		log.Errorf("Failed to query solana recent block, err = ", err)
		return
	}

	// Broadcast result
	msg := types.NewUpdateSolanaRecentHashMsg(p.signer, chain, result.Hash, result.Height)
	p.txSubmit.SubmitMessageAsync(msg)
}

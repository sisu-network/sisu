package tss

import (
	"fmt"
	"sort"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

type BlockSymbolPair struct {
	blockHeight int64
	chainSymbol string
}

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

	// Check if we have already processing this chain.
	found := false
	for _, pair := range p.keygenBlockPairs {
		if pair.chainSymbol == msg.ChainSymbol {
			found = true
			break
		}
	}

	fmt.Println("Found = ", found)

	if !found {
		// Add this chain to the processing queue. We will count votes in a few block later.
		p.keygenBlockPairs = append(p.keygenBlockPairs, BlockSymbolPair{
			blockHeight: p.currentHeight + int64(p.config.BlockProposalLength),
			chainSymbol: msg.ChainSymbol,
		})
		// Sort all pairs by block heights.
		sort.Slice(p.keygenBlockPairs, func(i, j int) bool {
			return p.keygenBlockPairs[i].blockHeight < p.keygenBlockPairs[j].blockHeight
		})
	}

	// TODO: Save this proposal to KV store.
	utils.LogDebug("!p.appInfo.IsCatchingUp() = ", !p.appInfo.IsCatchingUp())
	p.keygenVoteResult[msg.ChainSymbol] = make(map[string]bool)

	if !p.appInfo.IsCatchingUp() {
		// Send vote message to everyone else
		signer := p.appKeys.GetSignerAddress()
		voteMsg := types.NewMsgKeygenProposalVote(signer.String(), msg.ChainSymbol, types.KeygenProposalVote_APPROVE)

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
	fmt.Println("Signer = ", msg.Signer)

	voteResult := p.keygenVoteResult[msg.ChainSymbol]
	if voteResult == nil {
		voteResult = make(map[string]bool)
	}

	fmt.Println("msg.Vote = ", msg.Vote)
	fmt.Println("msg.Vote == types.KeygenProposalVote_APPROVE = ", msg.Vote == types.KeygenProposalVote_APPROVE)

	voteResult[msg.Signer] = msg.Vote == types.KeygenProposalVote_APPROVE
	p.keygenVoteResult[msg.ChainSymbol] = voteResult

	fmt.Println("len = ", len(p.keygenVoteResult[msg.ChainSymbol]))

	return []byte{}, nil
}

func (p *Processor) countKeygenVote() {
	chainSymbol := p.keygenBlockPairs[0].chainSymbol
	fmt.Println("Counting vote for ", chainSymbol)
	votesMap := p.keygenVoteResult[chainSymbol]

	fmt.Println("len(votesMap) =  ", len(votesMap))

	if len(votesMap) >= p.config.PoolSizeLowerBound {
		// n := utils.MinInt(len(votesMap), p.config.PoolSizeUpperBound)
		// TODO: Get top n validators from the map. For now, get all the validators.

		// 2. Send a signal to Tuktuk to start keygen process.
		utils.LogInfo("Sending keygen request to Tuktuk...")
		err := p.client.KeyGen(chainSymbol)
		if err != nil {
			utils.LogError(err)
			return
		}
		utils.LogInfo("Keygen request is sent successfully.")
	}
}

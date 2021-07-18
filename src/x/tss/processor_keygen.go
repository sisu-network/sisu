package tss

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tCommon "github.com/sisu-network/tuktuk/common"
	tTypes "github.com/sisu-network/tuktuk/types"
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
	voteResult := p.keygenVoteResult[msg.ChainSymbol]
	if voteResult == nil {
		voteResult = make(map[string]bool)
	}

	utils.LogDebug("msg = ", msg)

	voteResult[msg.Signer] = msg.Vote == types.KeygenProposalVote_APPROVE
	p.keygenVoteResult[msg.ChainSymbol] = voteResult

	return []byte{}, nil
}

func (p *Processor) countKeygenVote() {
	chainSymbol := p.keygenBlockPairs[0].chainSymbol
	votesMap := p.keygenVoteResult[chainSymbol]

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

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result tTypes.KeygenResult) {
	// Post result to the cosmos chain
	signer := p.appKeys.GetSignerAddress()

	resultEnum := types.KeygenResult_FAILURE
	if result.Success {
		resultEnum = types.KeygenResult_SUCCESS
	}

	msg := types.NewKeygenResult(signer.String(), result.Chain, resultEnum)
	p.txSubmit.SubmitMessage(msg)
}

func (p *Processor) DeliverKeygenResult(ctx sdk.Context, msg *types.KeygenResult) ([]byte, error) {
	// TODO: Accumulates results from others and check for bad actors

	// For now, only process self message sent from this node.
	if msg.Signer == p.appKeys.GetSignerAddress().String() {
		utils.LogDebug("Keygen: This is the same signer...")

		// Save this to KVStore
		chainsInfo, err := p.keeper.GetRecordedChainsOnSisu(ctx)
		if err != nil {
			utils.LogError(err)
			return nil, err
		}

		if chainsInfo.Chains == nil {
			chainsInfo.Chains = make(map[string]*types.ChainInfo)
		}

		// TODO: Add validators here.
		chainsInfo.Chains[msg.ChainSymbol] = &types.ChainInfo{
			Symbol: msg.ChainSymbol,
		}

		p.keeper.SetChainsInfo(ctx, chainsInfo)

		// TODO: Check if we need to deploy smart contracts.
		// Deploy smart contract if needed.
		p.BroadcastContractDeploymentMessage(msg.ChainSymbol)
	} else {
		utils.LogDebug("Keygen: message is from different signers.")
	}

	return nil, nil
}

func (p *Processor) BroadcastContractDeploymentMessage(chain string) {
	if tCommon.IsEthBasedChain(chain) {
		// Deploy smart contract for eth based chain. Add a few seconds delay for other nodes to update
		// latest results.
	}
}

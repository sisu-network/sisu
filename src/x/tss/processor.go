package tss

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
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

const (
	PROPOSE_BLOCK_INTERVAL = 1000
)

var (
	ERR_INVALID_MESSASGE_TYPE = fmt.Errorf("Invalid Message Type")
)

type Processor struct {
	keeper                 keeper.Keeper
	config                 config.TssConfig
	txSubmit               common.TxSubmit
	lastProposeBlockHeight int64
	appKeys                *common.AppKeys
	bridge                 *Bridge
	appInfo                *common.AppInfo
	currentHeight          int64

	// A map of chainSymbol -> map ()
	keygenVoteResult map[string]map[string]bool
	keygenBlockPairs []BlockSymbolPair
}

func NewProcessor(keeper keeper.Keeper,
	config config.TssConfig,
	appKeys *common.AppKeys,
	txSubmit common.TxSubmit,
	appInfo *common.AppInfo,
) *Processor {
	return &Processor{
		keeper:           keeper,
		appKeys:          appKeys,
		config:           config,
		txSubmit:         txSubmit,
		appInfo:          appInfo,
		keygenVoteResult: make(map[string]map[string]bool),
		// And array that stores block numbers where we should do final vote count.
		keygenBlockPairs: make([]BlockSymbolPair, 0),
	}
}

func (p *Processor) BeginBlock(blockHeight int64) {
	p.currentHeight = blockHeight
	// Check Vote result.
	for len(p.keygenBlockPairs) > 0 && !p.appInfo.IsCatchingUp() {
		utils.LogDebug("blockHeight = ", blockHeight)
		utils.LogDebug("p.keygenBlockPairs[0].blockHeight = ", p.keygenBlockPairs[0].blockHeight)

		if blockHeight < p.keygenBlockPairs[0].blockHeight {
			break
		}

		for len(p.keygenBlockPairs) > 0 && blockHeight >= p.keygenBlockPairs[0].blockHeight {
			chaimSymbol := p.keygenBlockPairs[0].chainSymbol

			// Now we count the votes
			p.countKeygenVote()

			// Remove the chain from processing queue.
			p.keygenBlockPairs = p.keygenBlockPairs[1:]
			delete(p.keygenVoteResult, chaimSymbol)
		}
	}
}

func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	recordedChains, err := p.keeper.GetRecordedChainsOnSisu(ctx)
	if err != nil {
		return
	}

	utils.LogInfo("recordedChains = ", recordedChains)

	unavailableChains := make([]string, 0)
	for _, chainConfig := range p.config.SupportedChains {
		if recordedChains[chainConfig.Symbol] == nil {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	if p.lastProposeBlockHeight == 0 || blockHeight-p.lastProposeBlockHeight > PROPOSE_BLOCK_INTERVAL {
		// Broadcast a message.
		utils.LogInfo("Broadcasting TSS Keygen Proposal message. len(unavailableChains) = ", len(unavailableChains))
		signer := p.appKeys.GetSignerAddress()

		for _, chain := range unavailableChains {
			// TODO: Add checking if a chain proposal has been submitted recently.
			proposal := types.NewMsgKeygenProposal(
				signer.String(),
				chain,
				utils.GenerateRandomString(16),
				blockHeight+PROPOSAL_BLOCK_LENGTH,
			)
			utils.LogDebug("Submitting proposal message for chain", chain)
			err := p.txSubmit.SubmitMessage(proposal)
			if err != nil {
				utils.LogError(err)
			}
		}

		p.lastProposeBlockHeight = blockHeight
	}
}

func (p *Processor) CheckTx(msgs []sdk.Msg) error {
	utils.LogDebug("TSSProcessor: checking tx. Message length = ", len(msgs))

	for _, msg := range msgs {
		if msg.Route() != types.ModuleName {
			return fmt.Errorf("Some message is not a TSS message")
		}

		if msg.Type() == types.MSG_TYPE_KEYGEN_PROPOSAL {
			typedMsg, ok := msg.(*types.KeygenProposal)
			if !ok {
				return ERR_INVALID_MESSASGE_TYPE
			}
			return p.CheckKeyGenProposal(typedMsg)
		}
	}

	return nil
}

func (p *Processor) DeliverTx(msg sdk.Msg) error {

	return nil
}

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

const (
	PROPOSE_BLOCK_INTERVAL = 1000
	PROPOSAL_BLOCK_LENGTH  = 5
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
}

func NewProcessor(keeper keeper.Keeper,
	config config.TssConfig,
	appKeys *common.AppKeys,
	txSubmit common.TxSubmit,
	appInfo *common.AppInfo,
) *Processor {
	return &Processor{
		keeper:   keeper,
		appKeys:  appKeys,
		config:   config,
		txSubmit: txSubmit,
		appInfo:  appInfo,
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

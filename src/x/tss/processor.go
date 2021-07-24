package tss

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/tssclients"
	"github.com/sisu-network/sisu/x/tss/types"
)

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
	appInfo                *common.AppInfo
	currentHeight          int64
	tuktukClient           *tssclients.Client
	logic                  *CrossChainLogic

	// This is a local database used for data specific to this node. For application state's data,
	// use KVStore.
	storage *TssStorage

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
		logic:            NewCrossChainLogic(),
	}
}

func (p *Processor) Init() {
	if p.config.Enable {
		p.connectToTuktuk()
		p.connectToDeyes()
	}

	var err error
	p.storage, err = NewTssStorage(p.config.Dir + "/processor.db")
	if err != nil {
		panic(err)
	}
}

// Connect to tuktuk server.
func (p *Processor) connectToTuktuk() {
	var err error
	url := fmt.Sprintf("http://%s:%d", p.config.Host, p.config.Port)
	utils.LogInfo("Connecting to tuktuk server at", url)

	p.tuktukClient, err = tssclients.Dial(url)

	if err != nil {
		utils.LogError(err)
		panic(err)
	}
	utils.LogInfo("Tuktuk server connected!")
}

func (p *Processor) connectToDeyes() {

}

func (p *Processor) BeginBlock(ctx sdk.Context, blockHeight int64) {
	p.currentHeight = blockHeight

	// Check keygen proposal
	p.CheckTssKeygen(ctx, blockHeight)

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

func (p *Processor) EndBlock(ctx sdk.Context) {
	// Check the list of transaction that have enough observations' attestation.
}

func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	if p.appInfo.IsCatchingUp() ||
		p.lastProposeBlockHeight != 0 && blockHeight-p.lastProposeBlockHeight <= PROPOSE_BLOCK_INTERVAL {
		return
	}

	chains, err := p.keeper.GetRecordedChainsOnSisu(ctx)
	if err != nil {
		return
	}

	utils.LogInfo("recordedChains = ", chains)

	unavailableChains := make([]string, 0)
	for _, chainConfig := range p.config.SupportedChains {
		// TODO: Remove this after testing.
		// if chains.Chains[chainConfig.Symbol] == nil {
		if true {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	// Broadcast a message.
	utils.LogInfo("Broadcasting TSS Keygen Proposal message. len(unavailableChains) = ", len(unavailableChains))
	signer := p.appKeys.GetSignerAddress()

	for _, chain := range unavailableChains {
		proposal := types.NewMsgKeygenProposal(
			signer.String(),
			chain,
			utils.GenerateRandomString(16),
			blockHeight+int64(p.config.BlockProposalLength),
		)
		utils.LogDebug("Submitting proposal message for chain", chain)
		go func() {
			err := p.txSubmit.SubmitMessage(proposal)
			if err != nil {
				utils.LogError(err)
			}
		}()
	}

	p.lastProposeBlockHeight = blockHeight
}

func (p *Processor) CheckTx(msgs []sdk.Msg) error {
	utils.LogDebug("TSSProcessor: checking tx. Message length = ", len(msgs))

	for _, msg := range msgs {
		if msg.Route() != types.ModuleName {
			return fmt.Errorf("Some message is not a TSS message")
		}

		utils.LogDebug("Msg type = ", msg.Type())

		switch msg.Type() {
		case types.MSG_TYPE_KEYGEN_PROPOSAL:
			typedMsg, ok := msg.(*types.KeygenProposal)
			if !ok {
				return ERR_INVALID_MESSASGE_TYPE
			}
			return p.CheckKeyGenProposal(typedMsg)

		case types.MSG_TYPE_KEYGEN_RESULT:
			// TODO: check this keygen result.
		}
	}

	return nil
}

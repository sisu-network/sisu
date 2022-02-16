package sisu

import (
	"fmt"
	"sync/atomic"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/mempool"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dhtypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	ProposeBlockInterval = 1000
)

var (
	ErrInvalidMessageType      = fmt.Errorf("Invalid Message Type")
	ErrMessageHasBeenProcessed = fmt.Errorf("Message has been processed")
	ErrCannotFindMessage       = fmt.Errorf("Cannot find the message in node's private db")
	ErrValueDoesNotMatch       = fmt.Errorf("Value does not match")
)

type BlockSymbolPair struct {
	blockHeight int64
	chain       string
}

// A major struct that processes complicated logic of TSS keysign and keygen. Read the documentation
// of keygen and keysign's flow before working on this.
type Processor struct {
	keeper                keeper.Keeper
	config                config.TssConfig
	tendermintPrivKey     crypto.PrivKey
	txSubmit              common.TxSubmit
	appKeys               common.AppKeys
	globalData            common.GlobalData
	partyManager          PartyManager
	txOutputProducer      TxOutputProducer
	lastContext           atomic.Value
	checkDuplicatedTxFunc mempool.PreCheckFunc
	txDecoder             sdk.TxDecoder

	// Public address of the key generated by TSS.
	keyAddress string

	// Dheart & Deyes client
	dheartClient tssclients.DheartClient
	deyesClient  tssclients.DeyesClient

	worldState       WorldState
	keygenVoteResult map[string]map[string]bool
	keygenBlockPairs []BlockSymbolPair

	publicDb  keeper.Storage
	privateDb keeper.Storage
}

func NewProcessor(k keeper.DefaultKeeper,
	publicDb keeper.Storage,
	privateDb keeper.Storage,
	config config.TssConfig,
	tendermintPrivKey crypto.PrivKey,
	appKeys *common.DefaultAppKeys,
	txDecoder sdk.TxDecoder,
	txSubmit common.TxSubmit,
	globalData common.GlobalData,
	dheartClient tssclients.DheartClient,
	deyesClient tssclients.DeyesClient,
	worldState WorldState,
) *Processor {
	p := &Processor{
		keeper:            &k,
		publicDb:          publicDb,
		privateDb:         privateDb,
		txDecoder:         txDecoder,
		appKeys:           appKeys,
		config:            config,
		tendermintPrivKey: tendermintPrivKey,
		txSubmit:          txSubmit,
		globalData:        globalData,
		partyManager:      NewPartyManager(globalData),
		keygenVoteResult:  make(map[string]map[string]bool),
		// And array that stores block numbers where we should do final vote count.
		keygenBlockPairs: make([]BlockSymbolPair, 0),
		dheartClient:     dheartClient,
		deyesClient:      deyesClient,
		worldState:       worldState,
	}

	return p
}

func (p *Processor) Init() {
	log.Info("Initializing TSS Processor...")

	p.txOutputProducer = NewTxOutputProducer(p.worldState, p.appKeys, p.publicDb, p.config)
}

func (p *Processor) BeginBlock(ctx sdk.Context, blockHeight int64) {
	p.setContext(ctx)

	// Check keygen proposal
	if blockHeight > 1 {
		// We need to wait till block 2 for multistore of the app to be updated with latest account info
		// for signing.
		p.CheckTssKeygen(ctx, blockHeight)
	}

	oldValue := p.globalData.IsCatchingUp()
	p.globalData.UpdateCatchingUp()
	newValue := p.globalData.IsCatchingUp()

	if oldValue && !newValue {
		log.Info("Setting Sisu readiness for dheart.")
		// This node has fully catched up with the blockchain, we need to inform dheart about this.
		p.dheartClient.SetSisuReady(true)
		p.deyesClient.SetSisuReady(true)
	}

	// TODO: Make keygen to be command instead of embedding inside the code.
	// Check Vote result.
	for len(p.keygenBlockPairs) > 0 && !p.globalData.IsCatchingUp() {
		log.Debug("blockHeight = ", blockHeight)
		if blockHeight < p.keygenBlockPairs[0].blockHeight {
			break
		}

		for len(p.keygenBlockPairs) > 0 && blockHeight >= p.keygenBlockPairs[0].blockHeight {
			// Remove the chain from processing queue.
			p.keygenBlockPairs = p.keygenBlockPairs[1:]
		}
	}

	// Calculate all token prices.
	p.calculateTokenPrices(ctx)
}

func (p *Processor) EndBlock(ctx sdk.Context) {
	if !p.globalData.IsCatchingUp() {
		// Inform dheart that we have reached end of block so that dheart could run presign works.
		height := ctx.BlockHeight()
		log.Verbose("End block reached, height = ", height)
		p.dheartClient.BlockEnd(height)
	}
}

func (p *Processor) setContext(ctx sdk.Context) {
	p.lastContext.Store(ctx)
}

// shouldProcessMsg counts how many validators have posted the same transaction to blockchain before
// processing.
//
// When adding new message type, remember to add its serialization in the GetTxRecordHash.
func (p *Processor) shouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	count := p.publicDb.SaveTxRecord(hash, signer)
	if count >= p.config.MajorityThreshold && !p.publicDb.IsTxRecordProcessed(hash) {
		return true, hash
	}

	return false, hash
}

///// Api Callback

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/
func (p *Processor) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	// TODO: We can replace this by sending command from client instead of running at the beginning
	// of each block.
	if p.globalData.IsCatchingUp() || ctx.BlockHeight()%50 != 2 {
		return
	}

	// Check ECDSA only (for now)
	keyTypes := []string{libchain.KEY_TYPE_ECDSA}
	for _, keyType := range keyTypes {
		if p.publicDb.IsKeygenExisted(keyType, 0) {
			continue
		}

		// Broadcast a message.
		signer := p.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenWithSigner(
			signer.String(),
			keyType,
			0,
		)

		log.Info("Submitting proposal message for ", keyType)
		p.txSubmit.SubmitMessageAsync(proposal)
	}
}

// Called after having key generation result from Sisu's api server.
func (p *Processor) OnKeygenResult(result dhtypes.KeygenResult) {
	var resultEnum types.KeygenResult_Result
	switch result.Outcome {
	case dhtypes.OutcomeSuccess:
		resultEnum = types.KeygenResult_SUCCESS
	case dhtypes.OutcomeFailure:
		resultEnum = types.KeygenResult_FAILURE
	case dhtypes.OutcometNotSelected:
		resultEnum = types.KeygenResult_NOT_SELECTED
	}

	if resultEnum == types.KeygenResult_NOT_SELECTED {
		// No need to send result when this node is not selected.
		return
	}

	signerMsg := types.NewKeygenResultWithSigner(
		p.appKeys.GetSignerAddress().String(),
		result.KeyType,
		result.KeygenIndex,
		resultEnum,
		result.PubKeyBytes,
		result.Address,
	)

	// Save the result to private db
	p.publicDb.SaveKeygenResult(signerMsg)

	log.Info("There is keygen result from dheart, resultEnum = ", resultEnum)

	p.txSubmit.SubmitMessageAsync(signerMsg)

	// Add list the public key address to watch.
	for _, chainConfig := range p.config.SupportedChains {
		chain := chainConfig.Id

		if libchain.GetKeyTypeForChain(chain) != result.KeyType {
			continue
		}

		log.Verbose("adding watcher address ", result.Address, " for chain ", chain)

		p.addWatchAddress(chain, result.Address)
	}
}

func (p *Processor) addWatchAddress(chain string, address string) {
	p.deyesClient.AddWatchAddresses(chain, []string{address})
}

package tss

import (
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/mempool"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/tssclients"
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
	deyesClients map[string]*tssclients.DeyesClient

	// A map of chain -> map ()
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
		deyesClients:     make(map[string]*tssclients.DeyesClient),
	}

	return p
}

func (p *Processor) Init() {
	log.Info("Initializing TSS Processor...")

	p.connectToDheart()
	p.connectToDeyes()

	p.txOutputProducer = NewTxOutputProducer(p.worldState, p.appKeys, p.publicDb, p.config)
}

// Connect to Dheart server and set private key for dheart. Note that this is the tendermint private
// key, not signer key in the keyring.
func (p *Processor) connectToDheart() {
	var err error
	url := fmt.Sprintf("http://%s:%d", p.config.DheartHost, p.config.DheartPort)
	log.Info("Connecting to Dheart server at", url)

	p.dheartClient, err = tssclients.DialDheart(url)
	if err != nil {
		log.Error("Failed to connect to Dheart. Err =", err)
		panic(err)
	}

	encryptedKey, err := p.appKeys.GetAesEncrypted(p.tendermintPrivKey.Bytes())
	if err != nil {
		log.Error("Failed to get encrypted private key. Err =", err)
		panic(err)
	}

	log.Info("p.tendermintPrivKey.Type() = ", p.tendermintPrivKey.Type())

	// Pass encrypted private key to dheart
	if err := p.dheartClient.SetPrivKey(hex.EncodeToString(encryptedKey), p.tendermintPrivKey.Type()); err != nil {
		panic(err)
	}

	log.Info("Dheart server connected!")
}

// Connecto to all deyes.
func (p *Processor) connectToDeyes() {
	for chain, chainConfig := range p.config.SupportedChains {
		log.Info("Deyes url = ", chainConfig.DeyesUrl)

		deyeClient, err := tssclients.DialDeyes(chainConfig.DeyesUrl)
		if err != nil {
			log.Error("Failed to connect to deyes", chain, ".Err =", err)
			panic(err)
		}

		if err := deyeClient.CheckHealth(); err != nil {
			panic(err)
		}

		p.deyesClients[chain] = deyeClient
	}

	p.worldState = NewWorldState(p.config, p.publicDb, p.deyesClients)
}

func (p *Processor) BeginBlock(ctx sdk.Context, blockHeight int64) {
	p.setContext(ctx)

	// Check keygen proposal
	if blockHeight > 1 {
		// We need to wait till block 2 for multistore of the app to be updated with latest account info
		// for signing.
		p.CheckTssKeygen(ctx, blockHeight)
	}

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

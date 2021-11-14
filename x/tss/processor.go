package tss

import (
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"github.com/sisu-network/tendermint/crypto"
	"github.com/sisu-network/tendermint/mempool"
	ttypes "github.com/sisu-network/tendermint/types"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/db"
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

// A major struct that processes complicated logic of TSS keysign and keygen. Read the documentation
// of keygen and keysign's flow before working on this.
type Processor struct {
	keeper                 keeper.Keeper
	config                 config.TssConfig
	tendermintPrivKey      crypto.PrivKey
	txSubmit               common.TxSubmit
	lastProposeBlockHeight int64
	appKeys                *common.AppKeys
	globalData             common.GlobalData
	currentHeight          int64
	partyManager           PartyManager
	txOutputProducer       TxOutputProducer
	lastContext            atomic.Value
	checkDuplicatedTxFunc  mempool.PreCheckFunc
	txDecoder              sdk.TxDecoder

	// Public address of the key generated by TSS.
	keyAddress string

	// Dheart & Deyes client
	dheartClient *tssclients.DheartClient
	deyesClients map[string]*tssclients.DeyesClient

	// A map of chain -> map ()
	worldState       WorldState
	keygenVoteResult map[string]map[string]bool
	keygenBlockPairs []BlockSymbolPair
	db               db.Database
}

func NewProcessor(keeper keeper.Keeper,
	config config.TssConfig,
	tendermintPrivKey crypto.PrivKey,
	appKeys *common.AppKeys,
	db db.Database,
	txDecoder sdk.TxDecoder,
	txSubmit common.TxSubmit,
	globalData common.GlobalData,
) *Processor {
	p := &Processor{
		keeper:            keeper,
		db:                db,
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

	if p.config.Enable {
		p.connectToDheart()
		p.connectToDeyes()
	}

	p.txOutputProducer = NewTxOutputProducer(p.worldState, p.keeper, p.appKeys, p.db, p.config)
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

	p.worldState = NewWorldState(p.config, p.db, p.deyesClients)
}

func (p *Processor) BeginBlock(ctx sdk.Context, blockHeight int64) {
	p.currentHeight = blockHeight
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
		log.Debug("p.keygenBlockPairs[0].blockHeight = ", p.keygenBlockPairs[0].blockHeight)

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
	// Do nothing
}

func (p *Processor) CheckTx(ctx sdk.Context, msgs []sdk.Msg) error {
	log.Debug("TSSProcessor: checking tx. Message length = ", len(msgs))

	for _, msg := range msgs {
		if msg.Route() != types.ModuleName {
			return fmt.Errorf("Some message is not a TSS message")
		}

		log.Debug("Checking tx: Msg type = ", msg.Type())

		switch msg.(type) {
		case *types.KeygenProposal:
			return p.CheckKeyGenProposal(msg.(*types.KeygenProposal))
		case *types.KeygenResult:
		case *types.ObservedTx:
			return p.CheckObservedTxs(ctx, msg.(*types.ObservedTx))
		case *types.TxOut:
			return p.CheckTxOut(ctx, msg.(*types.TxOut))
		case *types.KeysignResult:
			return p.CheckKeysignResult(ctx, msg.(*types.KeysignResult))
		}

		// switch msg.Type() {
		// case types.MSG_TYPE_KEYGEN_PROPOSAL:
		// 	return p.CheckKeyGenProposal(msg.(*types.KeygenProposal))

		// case types.MSG_TYPE_KEYGEN_RESULT:
		// 	// TODO: check this keygen result.
		// case types.MSG_TYPE_OBSERVED_TXS:

		// }
	}

	return nil
}

func (p *Processor) setContext(ctx sdk.Context) {
	p.lastContext.Store(ctx)
}

func (p *Processor) PreAddTxToMempoolFunc(txBytes ttypes.Tx) error {
	log.Verbose("checking new tx before adding into mempool....")

	tx, err := p.txDecoder(txBytes)
	if err != nil {
		log.Error("Failed to decode tx")
		return err
	}

	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		if msg.Route() != types.RouterKey {
			continue
		}

		switch msg.Type() {
		case types.MSG_TYPE_KEYGEN_PROPOSAL:
			proposalMsg := msg.(*types.KeygenProposal)
			// We dont get serialized data of the proposal msg because the serialized data contains
			// blockheight. Instead, we only check the chain of the proposal.
			key := fmt.Sprintf("KeygenProposal__%s", proposalMsg.Chain)
			if !p.db.MempoolTxExistedRange(key, p.currentHeight-PROPOSE_BLOCK_INTERVAL/2, p.currentHeight+PROPOSE_BLOCK_INTERVAL/2) {
				// Insert into the db
				p.db.InsertMempoolTxHash(key, p.currentHeight)
				return nil
			} else {
				err := fmt.Errorf("The keygen proposal has been inclued in a block for chain %s", proposalMsg.Chain)
				log.Verbose(err)
				return err
			}

		case types.MSG_TYPE_OBSERVED_TX:
			observedTx := msg.(*types.ObservedTx)
			hash := utils.KeccakHash32(string(observedTx.SerializeWithoutSigner()))

			if err := p.checkAndInsertMempoolTx(hash, "observed tx"); err != nil {
				return err
			}

		case types.MSG_TYPE_TX_OUT:
			txOut := msg.(*types.TxOut)
			hash := utils.KeccakHash32(string(txOut.SerializeWithoutSigner()))

			if err := p.checkAndInsertMempoolTx(hash, "tx out"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) checkAndInsertMempoolTx(hash, msgType string) error {
	if p.db.MempoolTxExisted(hash) {
		err := fmt.Errorf("%s has been added into the mempool! hash = %s", msgType, hash)
		log.Verbose(err)

		return err
	}

	log.Verbose("Inserting", msgType, "into the mempool table, hash =", hash)
	p.db.InsertMempoolTxHash(hash, p.currentHeight)

	return nil
}

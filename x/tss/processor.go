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
	currentHeight         atomic.Value
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

	db        db.Database
	privateDb keeper.PrivateDb
}

func NewProcessor(k keeper.DefaultKeeper,
	config config.TssConfig,
	tendermintPrivKey crypto.PrivKey,
	appKeys *common.DefaultAppKeys,
	db db.Database,
	dataDir string,
	txDecoder sdk.TxDecoder,
	txSubmit common.TxSubmit,
	globalData common.GlobalData,
) *Processor {
	p := &Processor{
		keeper:            &k,
		db:                db,
		privateDb:         keeper.NewPrivateDb(dataDir),
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

	p.currentHeight.Store(int64(0))

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
	p.currentHeight.Store(blockHeight)
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
		height := p.currentHeight.Load().(int64)
		log.Verbose("End block reached, height = ", height)
		p.dheartClient.BlockEnd(height)
	}
}

func (p *Processor) CheckTx(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		if msg.Route() != types.ModuleName {
			return fmt.Errorf("Some message is not a TSS message")
		}

		log.Info("Checking tx: Msg type = ", msg.Type())

		switch msg.(type) {
		case *types.KeygenWithSigner:
			return p.checkKeygen(ctx, msg.(*types.KeygenWithSigner))

		case *types.KeygenResultWithSigner:
			return p.checkKeygenResult(ctx, msg.(*types.KeygenResultWithSigner))

		case *types.TxInWithSigner:
			return p.checkTxIn(ctx, msg.(*types.TxInWithSigner))

		case *types.TxOutWithSigner:
			return p.checkTxOut(ctx, msg.(*types.TxOutWithSigner))

		case *types.KeysignResult:
			return p.checkKeysignResult(ctx, msg.(*types.KeysignResult))

		case *types.ContractsWithSigner:
			return p.checkContracts(ctx, msg.(*types.ContractsWithSigner))
		}
	}

	return nil
}

func (p *Processor) setContext(ctx sdk.Context) {
	p.lastContext.Store(ctx)
}

// PreAddTxToMempoolFunc checks if a tx has been included in a block. The hash of the tx is used to
// compare with other txs' hash. Only the first tx with such hash is included in the block. This is
// to avoid wasting space on Sisu's block due to duplicated tx submitted by multiple users.
func (p *Processor) PreAddTxToMempoolFunc(txBytes ttypes.Tx) error {
	log.Verbose("checking new tx before adding into mempool....")

	tx, err := p.txDecoder(txBytes)
	if err != nil {
		log.Error("Failed to decode tx")
		return err
	}

	msgs := tx.GetMsgs()
	log.Verbose("PreAddTxToMempoolFunc: msgs length = ", len(msgs))

	for _, msg := range msgs {
		if msg.Route() != types.RouterKey {
			continue
		}

		log.Verbose("PreAddTxToMempoolFunc: Msg type = ", msg.Type())

		switch msg.Type() {
		case types.MsgTypeKeygenWithSigner:
			keygenMsg := msg.(*types.KeygenWithSigner)

			key := fmt.Sprintf("keygen__%s__%d", keygenMsg.Data.KeyType, keygenMsg.Data.Index)
			if err := p.checkAndInsertMempoolTx(key, "keygen"); err != nil {
				return err
			}

		case types.MsgTypeKeygenResultWithSigner:
			resultMsg := msg.(*types.KeygenResultWithSigner)
			// Only do for success case
			if resultMsg.Data.Result == types.KeygenResult_SUCCESS {
				bz, err := resultMsg.Data.Marshal()
				if err != nil {
					return err
				}
				hash := utils.KeccakHash32(string(bz))
				key := fmt.Sprintf("keygenresult__%s__%d__%s", resultMsg.Keygen.KeyType, resultMsg.Keygen.Index, hash)
				if err := p.checkAndInsertMempoolTx(key, "keygen result"); err != nil {
					return err
				}
			}

		case types.MsgTypeTxInWithSigner:
			txIn := msg.(*types.TxInWithSigner).Data
			bz, err := txIn.Marshal()
			if err != nil {
				return err
			}

			hash := utils.KeccakHash32(string(bz))
			if err := p.checkAndInsertMempoolTx(hash, "tx in"); err != nil {
				return err
			}

		case types.MsgTypeTxOutWithSigner:
			txOut := msg.(*types.TxOutWithSigner).Data
			bz, err := txOut.Marshal()
			if err != nil {
				return err
			}

			hash := utils.KeccakHash32(string(bz))
			if err := p.checkAndInsertMempoolTx(hash, "tx out"); err != nil {
				return err
			}

		case types.MsgTypeContractsWithSigner:
			data := msg.(*types.ContractsWithSigner).Data

			bz, err := data.Marshal()
			if err == nil {
				hash := utils.KeccakHash32(string(bz))
				if err := p.checkAndInsertMempoolTx(hash, "contracts"); err != nil {
					return err
				}
			}

		}
	}

	fmt.Println("AAAA Done")

	return nil
}

func (p *Processor) checkAndInsertMempoolTx(hash, msgType string) error {
	if p.db.MempoolTxExisted(hash) {
		err := fmt.Errorf("%s has been added into the mempool! hash = %s", msgType, hash)
		log.Verbose(err)

		return err
	}

	log.Verbose("Inserting ", msgType, " into the mempool table, hash = ", hash)
	p.db.InsertMempoolTxHash(hash, p.currentHeight.Load().(int64))

	return nil
}

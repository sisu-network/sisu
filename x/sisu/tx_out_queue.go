package sisu

import (
	"fmt"
	"strconv"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	etypes "github.com/ethereum/go-ethereum/core/types"
	hTypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// TxOutQueue is an interface that batches outgoing transactions and sign them.
type TxOutQueue interface {
	Start()
	AddTxOut(txOut *types.TxOut)
	ProcessTxOuts(ctx sdk.Context)
}

type defaultTxOutQueue struct {
	keeper       keeper.Keeper
	globalData   common.GlobalData
	partyManager PartyManager
	txTracker    TxTracker
	dheartClient tssclients.DheartClient

	newTaskCh chan sdk.Context
	queue     []*types.TxOut
	lock      *sync.RWMutex
}

func NewTxOutQueue(keeper keeper.Keeper, globalData common.GlobalData, partyManager PartyManager,
	dheartClient tssclients.DheartClient, txTracker TxTracker) TxOutQueue {
	return &defaultTxOutQueue{
		keeper:       keeper,
		globalData:   globalData,
		partyManager: partyManager,
		dheartClient: dheartClient,
		txTracker:    txTracker,
		newTaskCh:    make(chan sdk.Context, 5),
		queue:        make([]*types.TxOut, 0),
		lock:         &sync.RWMutex{},
	}
}

func (q *defaultTxOutQueue) Start() {
	// Start the loop
	go q.loop()
	log.Info("TxOutQueue started")
}

func (q *defaultTxOutQueue) AddTxOut(txOut *types.TxOut) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, txOut)
}

func (q *defaultTxOutQueue) ProcessTxOuts(ctx sdk.Context) {

	q.newTaskCh <- ctx
}

func (q *defaultTxOutQueue) loop() {
	for {
		// Wait for new tx in to process
		ctx := <-q.newTaskCh

		// Read the queue
		q.lock.RLock()
		queue := q.queue
		q.queue = make([]*types.TxOut, 0)
		q.lock.RUnlock()

		if len(queue) == 0 {
			continue
		}

		if q.globalData.IsCatchingUp() {
			continue
		}

		for _, txOut := range queue {
			// Update the txOut to be delivered.
			q.txTracker.UpdateStatus(txOut.OutChain, txOut.OutHash, types.TxStatusDelivered)

			if libchain.IsETHBasedChain(txOut.OutChain) {
				q.signEthTx(ctx, txOut)
			}

			if libchain.IsCardanoChain(txOut.OutChain) {
				q.signCardanoTx(ctx, txOut)
			}
		}
	}
}

// signEthTx sends a TxOut to dheart for TSS signing.
func (q *defaultTxOutQueue) signEthTx(ctx sdk.Context, tx *types.TxOut) {
	log.Info("Delivering TXOUT for chain ", tx.OutChain, " tx hash = ", tx.OutHash)
	if tx.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		log.Info("This TxOut is a contract deployment")
	}

	ethTx := &etypes.Transaction{}
	if err := ethTx.UnmarshalBinary(tx.OutBytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
	}

	signer := libchain.GetEthChainSigner(tx.OutChain)
	if signer == nil {
		err := fmt.Errorf("cannot find signer for chain %s", tx.OutChain)
		log.Error(err)
	}

	hash := signer.Hash(ethTx)

	// 4. Send it to Dheart for signing.
	keysignReq := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_ECDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          q.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), tx.OutHash),
				InChain:     tx.InChain,
				OutChain:    tx.OutChain,
				OutHash:     tx.OutHash,
				BytesToSign: hash[:],
			},
		},
	}

	pubKeys := q.partyManager.GetActivePartyPubkeys()

	err := q.dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (q *defaultTxOutQueue) signCardanoTx(ctx sdk.Context, txOut *types.TxOut) {
	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx out: ", err)
		return
	}

	txHash, err := tx.Hash()
	if err != nil {
		log.Error("error when getting cardano tx hash: ", err)
		return
	}

	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          q.getKeysignRequestId(txOut.OutChain, ctx.BlockHeight(), txOut.OutHash),
				InChain:     txOut.InChain,
				OutChain:    txOut.OutChain,
				OutHash:     txOut.OutHash,
				BytesToSign: txHash[:],
			},
		},
	}

	pubKeys := q.partyManager.GetActivePartyPubkeys()
	err = q.dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (q *defaultTxOutQueue) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}

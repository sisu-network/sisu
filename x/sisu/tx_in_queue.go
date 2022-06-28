package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// TxInQueue is an interface that processes incoming transactions, produces corresponding transaction
// output.
type TxInQueue interface {
	Start()
	AddTxIn(height int64, txIn *types.TxIn)
	ProcessTxIns(ctx sdk.Context)
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit

	// queues stores list of TxIn by block height. We have to classify all incoming txs by their block
	// height so that all nodes will process the same list. If we have only single queue for all
	// TxIns, it's possible for different nodes to process different TxIn queues.
	queues     map[int64][]*types.TxIn
	newContext chan sdk.Context
	lock       *sync.RWMutex
}

func NewTxInQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	globalData common.GlobalData,
	txSubmit common.TxSubmit,
) TxInQueue {
	return &defaultTxInQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		globalData:       globalData,
		txSubmit:         txSubmit,
		newContext:       make(chan sdk.Context, 5),
		lock:             &sync.RWMutex{},
		queues:           make(map[int64][]*types.TxIn),
	}
}

func (q *defaultTxInQueue) Start() {
	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) AddTxIn(height int64, txIn *types.TxIn) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.queues[height] == nil {
		q.queues[height] = make([]*types.TxIn, 0, 10)
	}

	q.queues[height] = append(q.queues[height], txIn)
}

func (q *defaultTxInQueue) ProcessTxIns(ctx sdk.Context) {
	q.newContext <- ctx
}

func (q *defaultTxInQueue) loop() {
	for {
		// Wait for new tx in to process
		ctx := <-q.newContext
		q.processTxIns(ctx)
	}
}

func (q *defaultTxInQueue) processTxIns(ctx sdk.Context) {
	// Read the queue
	q.lock.RLock()
	queue := q.queues[ctx.BlockHeight()]
	q.lock.RUnlock()

	if len(queue) == 0 {
		return
	}

	// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
	// processor has.
	txOutWithSigners := q.txOutputProducer.GetTxOuts(ctx, ctx.BlockHeight(), queue)
	// Save this TxOut to database
	log.Verbose("len(txOut) = ", len(txOutWithSigners))
	if len(txOutWithSigners) > 0 {
		txOuts := make([]*types.TxOut, len(txOutWithSigners))
		for i, outWithSigner := range txOutWithSigners {
			txOut := outWithSigner.Data
			txOuts[i] = txOut

			// If this is a txOut deployment, mark the contract as being deployed.
			if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
				q.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
			}
		}
	}

	log.Info("Broadcasting txout....")

	for _, txOutWithSigner := range txOutWithSigners {
		q.txSubmit.SubmitMessageAsync(txOutWithSigner)
	}
}

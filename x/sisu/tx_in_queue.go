package sisu

import (
	"sync"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// TxInQueue is an interface that processes incoming transactions, produces corresponding transaction
// output.
type TxInQueue interface {
	Start()
	AddTxIn(txIn *types.TxIn)
	ProcessTxIns()
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit

	newTaskCh chan bool
	queue     []*types.TxIn
	lock      *sync.RWMutex
}

func NewTxInQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	globalData common.GlobalData,
	txSubmit common.TxSubmit,
	txTracker TxTracker,
) TxInQueue {
	return &defaultTxInQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		globalData:       globalData,
		txSubmit:         txSubmit,
		newTaskCh:        make(chan bool, 5),
		queue:            make([]*types.TxIn, 0),
		lock:             &sync.RWMutex{},
	}
}

func (q *defaultTxInQueue) Start() {
	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) AddTxIn(txIn *types.TxIn) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, txIn)
}

func (q *defaultTxInQueue) ProcessTxIns() {
	q.newTaskCh <- true
}

func (q *defaultTxInQueue) loop() {
	for {
		// Wait for new tx in to process
		<-q.newTaskCh

		// Read the queue
		q.lock.RLock()
		queue := q.queue
		q.queue = make([]*types.TxIn, 0)
		q.lock.RUnlock()

		if len(queue) == 0 {
			continue
		}

		ctx := q.globalData.GetReadOnlyContext()

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

		// If this node is not catching up, broadcast the tx.
		if !q.globalData.IsCatchingUp() && len(txOutWithSigners) > 0 {
			log.Info("Broadcasting txout....")

			for _, txOutWithSigner := range txOutWithSigners {
				q.txSubmit.SubmitMessageAsync(txOutWithSigner)
			}
		}
	}
}

package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	MaxPendingTxCacheSize = 1000
)

type TransferRequest struct {
	ctx sdk.Context
}

type TransferQueue interface {
	Start(ctx sdk.Context)
	Stop()
	ProcessTransfers(ctx sdk.Context)
}

type defaultTransferQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	txSubmit         common.TxSubmit
	stopCh           chan bool
	submittedTxs     *lru.Cache

	newRequestCh chan TransferRequest
	lock         *sync.RWMutex
}

func NewTransferQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
) TransferQueue {
	cache, err := lru.New(MaxPendingTxCacheSize)
	if err != nil {
		panic(err)
	}

	return &defaultTransferQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		txSubmit:         txSubmit,
		newRequestCh:     make(chan TransferRequest, 10),
		lock:             &sync.RWMutex{},
		stopCh:           make(chan bool),
		submittedTxs:     cache,
	}
}

func (q *defaultTransferQueue) Start(ctx sdk.Context) {
	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTransferQueue) Stop() {
	q.stopCh <- true
}

func (q *defaultTransferQueue) ProcessTransfers(ctx sdk.Context) {
	q.newRequestCh <- TransferRequest{
		ctx: ctx,
	}
}

func (q *defaultTransferQueue) loop() {
	for {
		select {
		case request := <-q.newRequestCh:
			// Wait for new tx in to process
			q.processBatch(request.ctx)
		case <-q.stopCh:
			return
		}
	}
}

func (q *defaultTransferQueue) processBatch(ctx sdk.Context) {
	params := q.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		queue := q.keeper.GetTransferQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		log.Debug("Queue length = ", len(queue))

		remaining := make([]*types.Transfer, 0)
		batchSize := params.GetMaxTransferOutBatch(chain)
		for i := 0; i < len(queue); i += batchSize {
			if _, ok := q.submittedTxs.Get(queue[i].Id); ok {
				log.Warn("Tx with id ", queue[i].Id, " is already submitted")
				continue
			}

			txOutMsgs, err := q.txOutputProducer.GetTxOuts(ctx, chain, queue[i:i+batchSize])
			if err != nil {
				log.Error(err)
				remaining = append(remaining, queue[i:i+batchSize]...)
				continue
			}

			for j := 0; j < batchSize; j++ {
				log.Verbose("Adding to submited txs ", queue[i+j].Id)
				q.submittedTxs.Add(queue[i+j].Id, true)
			}

			if len(txOutMsgs) > 0 {
				log.Infof("Broadcasting txout with length %d on chain %s", len(txOutMsgs), chain)
				for _, txOutMsg := range txOutMsgs {
					q.txSubmit.SubmitMessageAsync(txOutMsg)
				}
			}
		}
	}
}

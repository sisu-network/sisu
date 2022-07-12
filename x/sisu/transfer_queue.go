package sisu

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type TransferRequest struct {
	ctx sdk.Context
}

type TransferQueue interface {
	Start(ctx sdk.Context)
	Stop()
	ProcessTransfers(ctx sdk.Context)
	ClearInMemoryPendingTransfers(chain string)
}

type defaultTransferQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	txSubmit         common.TxSubmit
	stopCh           chan bool

	// In-memory list that stores all newly added transfer in a block (grouped by chain)
	chainsWithSubmission map[string]bool
	newRequestCh         chan TransferRequest
	lock                 *sync.RWMutex
}

func NewTransferQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
) TransferQueue {
	return &defaultTransferQueue{
		keeper:               keeper,
		txOutputProducer:     txOutputProducer,
		txSubmit:             txSubmit,
		newRequestCh:         make(chan TransferRequest, 10),
		lock:                 &sync.RWMutex{},
		chainsWithSubmission: make(map[string]bool),
		stopCh:               make(chan bool),
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

func (q *defaultTransferQueue) ClearInMemoryPendingTransfers(chain string) {
	q.chainsWithSubmission[chain] = false
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

		// Check if this chain has some pending tx.
		pendings := q.keeper.GetPendingTransfers(ctx, chain)
		if len(pendings) > 0 {
			continue
		}
		if q.chainsWithSubmission[chain] {
			continue
		}

		batchSize := params.GetMaxTransferOutBatch(chain)
		txOutMsgs, err := q.txOutputProducer.GetTxOuts(ctx, chain, queue[:batchSize])
		if err != nil {
			log.Error(err)
			continue
		}

		fmt.Println("len(txOutMsgs) = ", len(txOutMsgs), " on chain ", chain)

		if len(txOutMsgs) > 0 {
			q.chainsWithSubmission[chain] = true
			log.Info("Broadcasting txout....")
			for _, txOutMsg := range txOutMsgs {
				q.txSubmit.SubmitMessageAsync(txOutMsg)
			}
		}
	}
}

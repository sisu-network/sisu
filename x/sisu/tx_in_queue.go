package sisu

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type TransferBatchRequest struct {
	ctx     sdk.Context
	batches []*types.TransferBatch
}

type TransferQueue interface {
	Start(ctx sdk.Context)
	ProcessTransfers(ctx sdk.Context)
	ClearInMemoryPendingTransfers(chain string)
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit

	// In-memory list that stores all newly added transfer in a block (grouped by chain)
	chainsWithSubmission map[string]bool
	lastCheckPoints      map[string]*types.GatewayCheckPoint
	newRequestCh         chan TransferBatchRequest
	lock                 *sync.RWMutex
}

func NewTxInQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	globalData common.GlobalData,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
) TransferQueue {
	return &defaultTxInQueue{
		keeper:               keeper,
		txOutputProducer:     txOutputProducer,
		globalData:           globalData,
		txSubmit:             txSubmit,
		newRequestCh:         make(chan TransferBatchRequest, 10),
		lock:                 &sync.RWMutex{},
		chainsWithSubmission: make(map[string]bool),
		lastCheckPoints:      make(map[string]*types.GatewayCheckPoint),
	}
}

func (q *defaultTxInQueue) Start(ctx sdk.Context) {
	// Load all last checkpoints
	q.lastCheckPoints = q.keeper.GetAllGatewayCheckPoints(ctx)

	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) ProcessTransfers(ctx sdk.Context) {
	q.newRequestCh <- TransferBatchRequest{
		ctx: ctx,
	}
}

func (q *defaultTxInQueue) ClearInMemoryPendingTransfers(chain string) {
	q.chainsWithSubmission[chain] = false
}

func (q *defaultTxInQueue) loop() {
	for {
		// Wait for new tx in to process
		request := <-q.newRequestCh
		q.processBatch(request)
	}
}

func (q *defaultTxInQueue) processBatch(request TransferBatchRequest) {
	ctx := request.ctx

	params := q.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		queue := q.keeper.GetTransferQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		// Check if this chain has some pending tx.
		pendings := q.keeper.GetPendingTransfers(ctx, chain)
		if len(pendings) > 0 {
			fmt.Println("PEnding 1 is not empty")
			continue
		}
		if q.chainsWithSubmission[chain] {
			fmt.Println("PEnding 2 is not empty")
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

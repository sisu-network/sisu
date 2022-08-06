package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
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
	appKeys          common.AppKeys

	newRequestCh chan TransferRequest
	lock         *sync.RWMutex
}

func NewTransferQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
	appKeys common.AppKeys,
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
		appKeys:          appKeys,
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

		batchSize := utils.MinInt(params.GetMaxTransferOutBatch(chain), len(queue))
		batch := queue[0:batchSize]

		if _, ok := q.submittedTxs.Get(queue[0].Id); ok {
			log.Warn("Tx with id ", queue[0].Id, " is already submitted")
			continue
		}

		txOutMsgs, err := q.txOutputProducer.GetTxOuts(ctx, chain, batch)
		if err != nil {
			log.Error("Failed to get txOut on chain ", chain, ", err = ", err)
			// Submit transfer failure transaction if this is an ETH based chain
			ids := q.getTransferIds(batch)
			msg := types.NewTransferFailureMsg(q.appKeys.GetSignerAddress().String(), &types.TransferFailure{
				Ids:     ids,
				Chain:   chain,
				Message: err.Error(),
			})
			q.txSubmit.SubmitMessageAsync(msg)
			continue
		}

		for j := 0; j < batchSize; j++ {
			log.Verbose("Adding to submited txs ", queue[j].Id)
			q.submittedTxs.Add(queue[j].Id, true)
		}

		if len(txOutMsgs) > 0 {
			log.Infof("Broadcasting txout with length %d on chain %s", len(txOutMsgs), chain)
			for _, txOutMsg := range txOutMsgs {
				q.txSubmit.SubmitMessageAsync(txOutMsg)
			}
		}
	}
}

func (q *defaultTransferQueue) getTransferIds(batch []*types.Transfer) []string {
	ids := make([]string, len(batch))

	for i, transfer := range batch {
		ids[i] = transfer.Id
	}

	return ids
}

package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
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
	appKeys          common.AppKeys
	privateDb        keeper.PrivateDb
	newRequestCh     chan TransferRequest
	valsManager      ValidatorManager
	lock             *sync.RWMutex
	globalData       common.GlobalData
}

func NewTransferQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
	appKeys common.AppKeys,
	privateDb keeper.PrivateDb,
	valsManager ValidatorManager,
	globalData common.GlobalData,
) TransferQueue {
	return &defaultTransferQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		txSubmit:         txSubmit,
		newRequestCh:     make(chan TransferRequest, 10),
		lock:             &sync.RWMutex{},
		stopCh:           make(chan bool),
		appKeys:          appKeys,
		privateDb:        privateDb,
		valsManager:      valsManager,
		globalData:       globalData,
	}
}

func (q *defaultTransferQueue) Start(ctx sdk.Context) {
	// Start the loop
	go q.loop()
	log.Info("TransferQueue started")
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
	if q.globalData.IsCatchingUp() {
		// This app is still catching up with block. Do nothing here.
		return
	}

	params := q.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		if q.privateDb.GetHoldProcessing(types.TransferHoldKey, chain) {
			log.Verbose("Another transfer is being processed")
			continue
		}

		queue := q.keeper.GetTransferQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		// Check if the this node is the assigned node for the first transfer in the queue.
		transfer := queue[0]
		assignedNode := q.valsManager.GetAssignedValidator(ctx, transfer.Id)
		if assignedNode.AccAddress != q.appKeys.GetSignerAddress().String() {
			continue
		}

		log.Verbosef("Assigned node for transfer %s is %s", transfer.Id, assignedNode.AccAddress)

		batchSize := utils.MinInt(params.GetMaxTransferOutBatch(chain), len(queue))
		batch := queue[0:batchSize]

		txOutMsgs, err := q.txOutputProducer.GetTxOuts(ctx, chain, batch)
		if err != nil {
			log.Error("Failed to get txOut on chain ", chain, ", err = ", err)

			ids := q.getTransferIds(batch)
			msg := types.NewTransferFailureMsg(q.appKeys.GetSignerAddress().String(), &types.TransferFailure{
				Ids:     ids,
				Chain:   chain,
				Message: err.Error(),
			})
			q.txSubmit.SubmitMessageAsync(msg)

			continue
		}

		if len(txOutMsgs) > 0 {
			log.Infof("Broadcasting txout with length %d on chain %s", len(txOutMsgs), chain)
			for _, txOutMsg := range txOutMsgs {
				q.txSubmit.SubmitMessageAsync(txOutMsg)
			}

			q.privateDb.SetHoldProcessing(types.TransferHoldKey, chain, true)
		}
	}
}

func (q *defaultTransferQueue) getTransferIds(batch []*types.TransferDetails) []string {
	ids := make([]string, len(batch))

	for i, transfer := range batch {
		ids[i] = transfer.Id
	}

	return ids
}

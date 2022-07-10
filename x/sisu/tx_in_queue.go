package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type TxInRequest struct {
	ctx       sdk.Context
	transfers map[string][]*types.TransferOutData
}

// TxInQueue is an interface that processes incoming transactions, produces corresponding transaction
// output.
type TxInQueue interface {
	Start(ctx sdk.Context)
	ClearTransferQueue(ctx sdk.Context)
	AddTxIn(ctx sdk.Context, txIn *types.TxsIn)
	GetTransfer(ctx sdk.Context) map[string][]*types.Transfer
	ProcessTxIns(ctx sdk.Context)
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit

	// In-memory list that stores all newly added transfer in a block (grouped by chain)
	queues          map[string][]*types.Transfer
	pendingChains   map[string]bool
	lastCheckPoints map[string]*types.GatewayCheckPoint
	newRequestCh    chan TxInRequest
	lock            *sync.RWMutex
}

func NewTxInQueue(
	keeper keeper.Keeper,
	txOutputProducer TxOutputProducer,
	globalData common.GlobalData,
	txSubmit common.TxSubmit,
	tssConfig config.TssConfig,
) TxInQueue {
	return &defaultTxInQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		globalData:       globalData,
		txSubmit:         txSubmit,
		newRequestCh:     make(chan TxInRequest, 10),
		lock:             &sync.RWMutex{},
		queues:           make(map[string][]*types.Transfer),
		pendingChains:    make(map[string]bool),
		lastCheckPoints:  make(map[string]*types.GatewayCheckPoint),
	}
}

func (q *defaultTxInQueue) Start(ctx sdk.Context) {
	// Load all last checkpoints
	q.lastCheckPoints = q.keeper.GetAllGatewayCheckPoints(ctx)

	q.ClearTransferQueue(ctx)

	// Load pending transaction into memory
	params := q.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		transfers := q.keeper.GetPendingTransfers(ctx, chain)
		if len(transfers) > 0 {
			q.pendingChains[chain] = true
		}
	}

	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) ClearTransferQueue(ctx sdk.Context) {
	params := q.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		q.queues[chain] = make([]*types.Transfer, 0)
	}
}

func (q *defaultTxInQueue) AddTxIn(ctx sdk.Context, txsInput *types.TxsIn) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for _, request := range txsInput.Requests {
		if q.queues[request.ToChain] == nil {
			log.Warn("Unsupported chain ", txsInput.Chain)
			continue
		}

		transfer := &types.Transfer{
			Recipient: request.Recipient,
			Token:     request.Token,
			Amount:    request.Amount,
		}

		q.queues[request.ToChain] = append(q.queues[request.ToChain], transfer)
	}
}

func (q *defaultTxInQueue) GetTransfer(ctx sdk.Context) map[string][]*types.Transfer {
	copy := make(map[string][]*types.Transfer)
	for k, v := range q.queues {
		copy[k] = v
	}
	return copy
}

func (q *defaultTxInQueue) ProcessTxIns(ctx sdk.Context) {
	q.newRequestCh <- TxInRequest{
		ctx: ctx,
	}
}

func (q *defaultTxInQueue) loop() {
	for {
		// Wait for new tx in to process
		ctx := <-q.newRequestCh
		q.processTxIns(ctx)
	}
}

func (q *defaultTxInQueue) processTxIns(request TxInRequest) {
	ctx := request.ctx
	params := q.keeper.GetParams(ctx)

	// Read the queue from the context.
	for _, chain := range params.SupportedChains {
		q.lock.RLock()
		pending := q.pendingChains[chain]
		q.lock.RUnlock()

		if pending {
			continue
		}

		pendings := q.keeper.GetPendingTransfers(ctx, chain)
		if len(pendings) > 0 {
			continue
		}

		queue := q.keeper.GetTransferQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		batchSize := params.GetMaxTransferOutBatch(chain)
		transfers := queue[:batchSize]
		txOutMsgs, err := q.txOutputProducer.GetTxOuts(ctx, chain, transfers)
		if err != nil {
			log.Error(err)
			continue
		}

		q.lock.Lock()
		q.pendingChains[chain] = true
		q.lock.Unlock()

		// // If this is a txOut deployment, mark the contract as being deployed.
		// if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		// 	q.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
		// }

		if len(txOutMsgs) > 0 {
			log.Info("Broadcasting txout....")
			for _, txOutMsg := range txOutMsgs {
				q.txSubmit.SubmitMessageAsync(txOutMsg)
			}
		}
	}
}

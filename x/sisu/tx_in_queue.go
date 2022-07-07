package sisu

import (
	"fmt"
	"math/big"
	"sort"
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
	AddTxIn(ctx sdk.Context, txIn *types.TxsIn)
	RemoveTransfers(ctx sdk.Context, transfers []*types.TransferOutData)
	ProcessTxIns(ctx sdk.Context)
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit
	supportedChains  []string

	// queues stores list of TxIn by chain.
	queues          map[string][]*types.TransferOutData
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
	supportedChains := make([]string, 0)
	for chain := range tssConfig.SupportedChains {
		supportedChains = append(supportedChains, chain)
	}
	sort.Strings(supportedChains)

	return &defaultTxInQueue{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		globalData:       globalData,
		txSubmit:         txSubmit,
		newRequestCh:     make(chan TxInRequest, 5),
		lock:             &sync.RWMutex{},
		queues:           make(map[string][]*types.TransferOutData),
		lastCheckPoints:  make(map[string]*types.GatewayCheckPoint),
		supportedChains:  supportedChains,
	}
}

func (q *defaultTxInQueue) Start(ctx sdk.Context) {
	// Load all last checkpoints
	q.lastCheckPoints = q.keeper.GetAllGatewayCheckPoints(ctx)

	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) AddTxIn(ctx sdk.Context, txsInput *types.TxsIn) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.queues[txsInput.Chain] == nil {
		q.queues[txsInput.Chain] = make([]*types.TransferOutData, 0, len(txsInput.Requests))
	}

	for _, request := range txsInput.Requests {
		// Get token from keeper.
		tokens := q.keeper.GetTokens(ctx, []string{request.Token})
		token := tokens[request.Token]
		if token == nil {
			log.Warn("AddTxIn: cannot find token ", request.Token)
			continue
		}

		amount, ok := new(big.Int).SetString(request.Amount, 10)
		if !ok {
			log.Error("Cannot set string for big integer")
			continue
		}

		transfer := &types.TransferOutData{
			BlockHeight: txsInput.Height,
			DestChain:   request.ToChain,
			Recipient:   request.Recipient,
			Token:       token,
			Amount:      amount,
			InHash:      request.Hash,
		}

		q.queues[txsInput.Chain] = append(q.queues[txsInput.Chain], transfer)
	}
}

func (q *defaultTxInQueue) RemoveTransfers(ctx sdk.Context, transfers []*types.TransferOutData) {
	fmt.Println("Removing transfers....")
	copy := make(map[string][]*types.TransferOutData)
	q.lock.Lock()
	for k, v := range q.queues {
		copy[k] = v
	}

	for chain, v := range copy {
		newArr := make([]*types.TransferOutData, 0)
		for _, curTransfer := range v {
			found := false
			for _, transfer := range transfers {
				if curTransfer.InChain == transfer.InChain && curTransfer.InHash == transfer.InHash {
					fmt.Println("Removing transfer with chain and hash", transfer.InChain, transfer.InHash)
					found = true
					break
				}
			}

			if !found {
				newArr = append(newArr, curTransfer)
			}
		}

		q.queues[chain] = newArr
	}
	q.lock.Unlock()
}

func (q *defaultTxInQueue) ProcessTxIns(ctx sdk.Context) {
	copy := make(map[string][]*types.TransferOutData)
	q.lock.Lock()
	for k, v := range q.queues {
		copy[k] = v
	}

	q.queues = make(map[string][]*types.TransferOutData)
	q.lock.Unlock()

	q.newRequestCh <- TxInRequest{
		ctx:       ctx,
		transfers: copy,
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
	transfersByChains := make([][]*types.TransferOutData, 0)
	// Read the queue
	q.lock.RLock()
	for _, chain := range q.supportedChains {
		if len(request.transfers[chain]) != 0 {
			transfersByChains = append(transfersByChains, request.transfers[chain])
		}
	}
	q.lock.RUnlock()

	for _, transfers := range transfersByChains {
		destChain := transfers[0].DestChain
		newCheckPoint := q.keeper.GetGatewayCheckPoint(ctx, destChain)
		lastCheckPoint := q.lastCheckPoints[destChain]
		if newCheckPoint != nil && lastCheckPoint != nil && newCheckPoint.BlockHeight == lastCheckPoint.BlockHeight {
			log.Verbose("No new checkpoint, skip processing tx in for chain ", destChain)
			continue
		}
		q.lastCheckPoints[destChain] = newCheckPoint

		fmt.Println("Creating txout....")

		// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
		// processor has.
		txOutMsgs, notProcessed := q.txOutputProducer.GetTxOuts(ctx, destChain, transfers)
		fmt.Println("txOutMsgs, notProcessed = ", len(txOutMsgs), len(notProcessed))

		if len(notProcessed) > 0 {
			// Add the not processed txs back to the queue to be processed in the next block
			q.lock.Lock()
			q.queues[destChain] = append(notProcessed, q.queues[destChain]...)
			q.lock.Unlock()
		}

		// Save this TxOut to database
		log.Verbose("len(txOut) = ", len(txOutMsgs))
		if len(txOutMsgs) > 0 {
			txOuts := make([]*types.TxOut, len(txOutMsgs))
			for i, outWithSigner := range txOutMsgs {
				txOut := outWithSigner.Data
				txOuts[i] = txOut

				// If this is a txOut deployment, mark the contract as being deployed.
				if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
					q.keeper.UpdateContractsStatus(ctx, txOut.OutChain, txOut.ContractHash, string(types.TxOutStatusSigning))
				}
			}

			log.Info("Broadcasting txout....")
			for _, txOutWithSigner := range txOutMsgs {
				q.txSubmit.SubmitMessageAsync(txOutWithSigner)
			}
		}
	}
}

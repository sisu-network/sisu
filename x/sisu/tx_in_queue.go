package sisu

import (
	"fmt"
	"math/big"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type TxInRequest struct {
	Ctx sdk.Context

	HasNewCarnadoBlock bool
	CarnadoBlockHeight int64
}

// TxInQueue is an interface that processes incoming transactions, produces corresponding transaction
// output.
type TxInQueue interface {
	Start()
	AddTxIn(ctx sdk.Context, txIn *types.TxsIn)
	ProcessTxIns(request TxInRequest)
}

type defaultTxInQueue struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	txSubmit         common.TxSubmit

	// queues stores list of TxIn by block height. We have to classify all incoming txs by their block
	// height so that all nodes will process the same list. If we have only single queue for all
	// TxIns, it's possible for different nodes to process different TxIn queues.
	queues     map[int64][]*transferOutData
	newContext chan TxInRequest
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
		newContext:       make(chan TxInRequest, 5),
		lock:             &sync.RWMutex{},
		queues:           make(map[int64][]*transferOutData),
	}
}

func (q *defaultTxInQueue) Start() {
	// Start the loop
	go q.loop()
	log.Info("TxInQueue started")
}

func (q *defaultTxInQueue) AddTxIn(ctx sdk.Context, txsInput *types.TxsIn) {
	q.lock.Lock()
	defer q.lock.Unlock()

	fmt.Println("AAAA 00000 0000000")

	height := ctx.BlockHeight()

	if q.queues[height] == nil {
		q.queues[height] = make([]*transferOutData, 0, len(txsInput.Requests))
	}

	fmt.Println("AAAA 00000 1111111")

	fmt.Println("AAAA 00000")

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

		q.queues[height] = append(q.queues[height], &transferOutData{
			blockHeight: txsInput.Height,
			destChain:   request.ToChain,
			recipient:   request.Recipient,
			token:       token,
			amount:      amount,
		})
	}

	fmt.Println("AAAA 11111")
}

func (q *defaultTxInQueue) ProcessTxIns(request TxInRequest) {
	q.newContext <- request
}

func (q *defaultTxInQueue) loop() {
	for {
		// Wait for new tx in to process
		ctx := <-q.newContext
		q.processTxIns(ctx)
	}
}

func (q *defaultTxInQueue) processTxIns(request TxInRequest) {
	ctx := request.Ctx

	// Read the queue
	q.lock.RLock()
	queue := q.queues[ctx.BlockHeight()]
	delete(q.queues, ctx.BlockHeight())
	q.lock.RUnlock()

	if len(queue) == 0 {
		return
	}

	// Creates and broadcast TxOuts. This has to be deterministic based on all the data that the
	// processor has.
	txOutWithSigners, notProcessed := q.txOutputProducer.GetTxOuts(ctx, queue)

	if len(notProcessed) > 0 {
		// Add the not processed txs back to the queue to be processed in the next block
		q.lock.Lock()
		queue = q.queues[ctx.BlockHeight()+1]
		if queue == nil {
			queue = notProcessed
		} else {
			// prepend to the beginning of the array
			queue = append(notProcessed, queue...)
		}
		q.lock.Unlock()
	}

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

		log.Info("Broadcasting txout....")
		for _, txOutWithSigner := range txOutWithSigners {
			q.txSubmit.SubmitMessageAsync(txOutWithSigner)
		}
	}
}

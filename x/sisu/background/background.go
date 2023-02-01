package background

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type Background interface {
	Start()
	Update(ctx sdk.Context)
}

type defaultBackground struct {
	keeper           keeper.Keeper
	txOutputProducer chains.TxOutputProducer
	txSubmit         components.TxSubmit
	appKeys          components.AppKeys
	privateDb        keeper.PrivateDb
	newRequestCh     chan TransferRequest
	valsManager      components.ValidatorManager
	globalData       components.GlobalData
	dheartCli        external.DheartClient
	partyManager     components.PartyManager
	stopCh           chan bool
}

func NewBackground(
	keeper keeper.Keeper,
	txOutputProducer chains.TxOutputProducer,
	txSubmit components.TxSubmit,
	appKeys components.AppKeys,
	privateDb keeper.PrivateDb,
	valsManager components.ValidatorManager,
	globalData components.GlobalData,
	dheartCli external.DheartClient,
	partyManager components.PartyManager,
) Background {
	return &defaultBackground{
		keeper:           keeper,
		txOutputProducer: txOutputProducer,
		txSubmit:         txSubmit,
		newRequestCh:     make(chan TransferRequest, 10),
		stopCh:           make(chan bool),
		appKeys:          appKeys,
		privateDb:        privateDb,
		valsManager:      valsManager,
		globalData:       globalData,
		dheartCli:        dheartCli,
		partyManager:     partyManager,
	}
}

func (b *defaultBackground) Start() {
	// Start the loop
	go b.loop()
	log.Info("Backround started")
}

func (q *defaultBackground) loop() {
	for {
		select {
		case request := <-q.newRequestCh:
			// Wait for new tx in to process
			q.Process(request.ctx)
		case <-q.stopCh:
			return
		}
	}
}

func (b *defaultBackground) Stop() {
	b.stopCh <- true
}

func (q *defaultBackground) Update(ctx sdk.Context) {
	q.newRequestCh <- TransferRequest{
		ctx: ctx,
	}
}

func (b *defaultBackground) Process(ctx sdk.Context) {
	params := b.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		// Process admin commands queue.
		cmdQ := b.keeper.GetCommandQueue(ctx, chain)
		if len(cmdQ) > 0 {
			// Admin command has higher priority

		} else {
			// Process transfer queue
			b.processTranfserQueue(ctx, chain, params)
		}
	}

	b.processTxOut(ctx, params)
}

func (b *defaultBackground) processCmdQueue(ctx sdk.Context, chain string, cmd *types.Command) {
	switch cmd.Type.(type) {
	case *types.Command_PauseResume:
	}
}

func (b *defaultBackground) processTranfserQueue(ctx sdk.Context, chain string, params *types.Params) {
	if b.globalData.IsCatchingUp() {
		// This app is still catching up with block. Do nothing here.
		return
	}

	queue := b.keeper.GetTransferQueue(ctx, chain)
	if len(queue) == 0 {
		return
	}

	// Check if the this node is the assigned node for the first transfer in the queue.
	transfer := queue[0]
	assignedNode := b.valsManager.GetAssignedValidator(ctx, transfer.Id)
	if assignedNode.AccAddress != b.appKeys.GetSignerAddress().String() {
		return
	}

	log.Verbosef("Assigned node for transfer %s is %s", transfer.Id, assignedNode.AccAddress)

	batchSize := utils.MinInt(params.GetMaxTransferOutBatch(chain), len(queue))
	batch := queue[0:batchSize]

	txOutMsgs, err := b.txOutputProducer.GetTxOuts(ctx, chain, batch)
	if err != nil {
		log.Error("Failed to get txOut on chain ", chain, ", err = ", err)

		ids := b.getTransferIds(batch)
		msg := types.NewTransferFailureMsg(b.appKeys.GetSignerAddress().String(), &types.TransferFailure{
			Ids:     ids,
			Chain:   chain,
			Message: err.Error(),
		})
		b.txSubmit.SubmitMessageAsync(msg)

		return
	}

	if len(txOutMsgs) > 0 {
		log.Infof("Broadcasting txout with length %d on chain %s", len(txOutMsgs), chain)
		for _, txOutMsg := range txOutMsgs {
			b.txSubmit.SubmitMessageAsync(txOutMsg)
		}

		b.privateDb.SetHoldProcessing(types.TransferHoldKey, chain, true)
	}
}

func (b *defaultBackground) getTransferIds(batch []*types.TransferDetails) []string {
	ids := make([]string, len(batch))

	for i, transfer := range batch {
		ids[i] = transfer.Id
	}

	return ids
}

func (b *defaultBackground) processTxOut(ctx sdk.Context, params *types.Params) {
	for _, chain := range params.SupportedChains {
		if b.privateDb.GetHoldProcessing(types.TxOutHoldKey, chain) {
			log.Verbosef("Another TxOut is being processed on chain %s", chain)
			continue
		}

		queue := b.keeper.GetTxOutQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		b.privateDb.SetHoldProcessing(types.TxOutHoldKey, chain, true)

		txOut := queue[0]
		if !b.globalData.IsCatchingUp() {
			SignTxOut(ctx, b.dheartCli, b.partyManager, txOut)
		}
	}
}

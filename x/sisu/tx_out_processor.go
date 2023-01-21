package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type TxOutProcessor interface {
	Start()
	ProcessTxOut(ctx sdk.Context)
	Stop()
}

type defaultTxOutProcessor struct {
	keeper       keeper.Keeper
	privateDb    keeper.PrivateDb
	newRequestCh chan sdk.Context
	txOutSigner  *txOutSigner
	globalData   components.GlobalData
	stopCh       chan bool
}

func NewTxOutProcessor(keeper keeper.Keeper, privateDb keeper.PrivateDb, txOutSigner *txOutSigner,
	globalData components.GlobalData) TxOutProcessor {
	return &defaultTxOutProcessor{
		keeper:       keeper,
		privateDb:    privateDb,
		newRequestCh: make(chan sdk.Context, 10),
		txOutSigner:  txOutSigner,
		globalData:   globalData,
		stopCh:       make(chan bool),
	}
}

func (d *defaultTxOutProcessor) Start() {
	go d.loop()
}

func (d *defaultTxOutProcessor) Stop() {
	d.stopCh <- true
}

func (d *defaultTxOutProcessor) loop() {
	for {
		select {
		case ctx := <-d.newRequestCh:
			d.processTxOut(ctx)

		case <-d.stopCh:
			return
		}
	}
}

func (d *defaultTxOutProcessor) ProcessTxOut(ctx sdk.Context) {
	d.newRequestCh <- ctx
}

func (d *defaultTxOutProcessor) processTxOut(ctx sdk.Context) {
	params := d.keeper.GetParams(ctx)
	for _, chain := range params.SupportedChains {
		if d.privateDb.GetHoldProcessing(types.TxOutHoldKey, chain) {
			log.Verbose("Another TxOut is being processed")
			continue
		}

		queue := d.keeper.GetTxOutQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		d.privateDb.SetHoldProcessing(types.TxOutHoldKey, chain, true)

		txOut := queue[0]
		if !d.globalData.IsCatchingUp() {
			d.txOutSigner.signTxOut(ctx, txOut)
		}
	}
}

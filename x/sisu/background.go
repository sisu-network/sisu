package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type BackgroundInput struct {
	Context           sdk.Context
	NewConfirmedTxIns []*types.TxIn
}

type Background interface {
	Start()
	ProcessData(input BackgroundInput)
}

type defaultBackground struct {
	inputCh       chan BackgroundInput
	keeper        keeper.Keeper
	privateDb     keeper.PrivateDb
	valsManager   ValidatorManager
	bridgeManager chains.BridgeManager
}

func NewBackground(k keeper.Keeper, privateDb keeper.PrivateDb, valsManager ValidatorManager,
	bridgeManager chains.BridgeManager) Background {
	return &defaultBackground{
		inputCh:       make(chan BackgroundInput),
		keeper:        k,
		privateDb:     privateDb,
		valsManager:   valsManager,
		bridgeManager: bridgeManager,
	}
}

func (b *defaultBackground) Start() {
	go b.loop()
}

func (b *defaultBackground) ProcessData(input BackgroundInput) {
	b.inputCh <- input
}

func (b *defaultBackground) loop() {
	for {
		select {
		case input := <-b.inputCh:
			b.process(input)
		}
	}
}

func (b *defaultBackground) process(input BackgroundInput) {
}

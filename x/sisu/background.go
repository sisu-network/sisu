package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
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
	ctx := input.Context

	newTransfers := make(map[string][]*types.TransferDetails)

	for _, txIn := range input.NewConfirmedTxIns {
		confirmedTxIn := b.keeper.GetConfirmedTxIn(ctx, txIn.Id)
		if confirmedTxIn == nil {
			log.Errorf("Critical error: cannot find the confirmed TxIn in the keeper when processing new Confirmed TxIn")
			continue
		}

		details := b.keeper.GetTxInDetails(ctx, txIn.Id)
		if details == nil {
			log.Errorf("Critical error: cannot find the TxInDetails in the keeper when processing new Confirmed TxIn")
			continue
		}

		// Parse the details
		chain := details.Data.FromChain
		bridge := b.bridgeManager.GetBridge(ctx, chain)
		if bridge == nil {
			log.Errorf("Cannot find bridge for chain %s", chain)
			continue
		}

		transfers, err := bridge.ParseIncomginTx(ctx, chain, details.Data.Serialize)
		if err != nil {
			log.Errorf("Failed to parse incoming transaction for chain %s, err = %s", chain, err)
			continue
		}

		if len(transfers) == 0 {
			continue
		}

		if newTransfers[chain] == nil {
			newTransfers[chain] = make([]*types.TransferDetails, 0)
		}

		newTransfers[chain] = append(newTransfers[chain], transfers...)
	}
}

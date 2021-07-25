package tss

import (
	eTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	tTypes "github.com/sisu-network/tuktuk/types"
)

type ApiHandler struct {
	processor *Processor
	keeper    *keeper.Keeper
}

func NewApi(processor *Processor, keeper *keeper.Keeper) *ApiHandler {
	return &ApiHandler{
		processor: processor,
	}
}

func (a *ApiHandler) Version() string {
	return "1.0"
}

// Empty function for checking health only.
func (api *ApiHandler) CheckHealth() {
}

func (a *ApiHandler) KeygenResult(result tTypes.KeygenResult) bool {
	utils.LogInfo("There is a TSS Result")

	a.processor.OnKeygenResult(result)
	return true
}

// This is a API endpoint to receive transactions with To address we are interested in.
func (a *ApiHandler) PostObservedTxs(txs *eTypes.Txs) {
	utils.LogDebug("There is new list of transactions from deyes")

	// There is a new transaction that we are interested in.
	a.processor.ProcessObservedTxs(txs)
}

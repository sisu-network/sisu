package tss

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/keeper"
)

type ApiHandler struct {
	processor *Processor
	keeper    *keeper.DefaultKeeper
}

func NewApi(processor *Processor, keeper *keeper.DefaultKeeper) *ApiHandler {
	return &ApiHandler{
		processor: processor,
	}
}

func (a *ApiHandler) Version() string {
	return "1.0"
}

// Empty function for checking health only.
func (a *ApiHandler) CheckHealth() {
}

// This is a API endpoint to receive transactions with To address we are interested in.
func (a *ApiHandler) PostObservedTxs(txs *etypes.Txs) {
	log.Debug("There is new list of transactions from deyes from chain ", txs.Chain)

	for _, tx := range txs.Arr {
		ethTx := &ethtypes.Transaction{}

		err := ethTx.UnmarshalBinary(tx.Serialized)
		if err != nil {
			log.Error("Cannot unmarshall transaction ", err)
		}
	}

	// There is a new transaction that we are interested in.
	a.processor.OnTxIns(txs)
}

func (a *ApiHandler) PostDeploymentResult(result *etypes.DispatchedTxResult) {
	go a.processor.OnTxDeploymentResult(result)
}

package sisu

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
)

type ApiHandler struct {
	processor *Processor
}

func NewApi(processor *Processor) *ApiHandler {
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

func (a *ApiHandler) KeygenResult(result htypes.KeygenResult) bool {
	log.Info("There is a Keygen Result")

	a.processor.OnKeygenResult(result)
	return true
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

func (a *ApiHandler) KeysignResult(result *htypes.KeysignResult) {
	log.Info("There is keysign result")
	go a.processor.OnKeysignResult(result)
}

func (a *ApiHandler) PostDeploymentResult(result *etypes.DispatchedTxResult) {
	go a.processor.OnTxDeploymentResult(result)
}

func (a *ApiHandler) UpdateGasPrice(request *etypes.GasPriceRequest) {
	log.Info("Received update gas price request")
	go a.processor.OnUpdateGasPriceRequest(request)
}

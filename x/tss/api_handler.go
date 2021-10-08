package tss

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
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

func (a *ApiHandler) KeygenResult(result htypes.KeygenResult) bool {
	utils.LogInfo("There is a Keygen Result")

	a.processor.OnKeygenResult(result)
	return true
}

// This is a API endpoint to receive transactions with To address we are interested in.
func (a *ApiHandler) PostObservedTxs(txs *etypes.Txs) {
	utils.LogDebug("There is new list of transactions from deyes")

	for _, tx := range txs.Arr {
		ethTx := &ethtypes.Transaction{}

		err := ethTx.UnmarshalBinary(tx.Serialized)
		if err != nil {
			utils.LogError("Cannot unmarshall transaction ", err)
		}
	}

	// There is a new transaction that we are interested in.
	a.processor.OnObservedTxs(txs)
}

func (a *ApiHandler) KeysignResult(result *htypes.KeysignResult) {
	utils.LogInfo("There is keysign result")
	go a.processor.OnKeysignResult(result)
}

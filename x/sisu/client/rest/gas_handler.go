package rest

import (
	"encoding/json"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"math/big"
	"net/http"

	"github.com/sisu-network/lib/log"
)

type gasCostRequest struct {
	Chain   string `json:"chain,omitempty"`
	TokenId string `json:"token_id,omitempty"`
}

type gasCostResponse struct {
	Chain   string `json:"chain,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	GasCost int64  `json:"gas_cost,omitempty"`
}

func (a *ExternalHandler) newGasCostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &gasCostRequest{}
		if err := ReadRESTReq(r, req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("Come here 1")
		gasPrice, err := a.worldState.GetGasPrice(req.Chain)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("Come here 2")
		// TODO: correct gasLimit here
		gasLimit := big.NewInt(8_000_000)
		tokenPrice, err := a.worldState.GetTokenPrice(req.TokenId)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("Come here 3")
		nativeTokenPrice, err := a.worldState.GetNativeTokenPriceForChain(req.Chain)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		gasCost, err := helper.GetGasCostInToken(gasLimit, gasPrice, big.NewInt(tokenPrice), big.NewInt(nativeTokenPrice))
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		output, err := json.Marshal(&gasCostResponse{
			Chain:   req.Chain,
			TokenId: req.TokenId,
			GasCost: gasCost.Int64(),
		})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Come here 4")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(output)
	}
}

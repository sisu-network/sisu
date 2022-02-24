package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sisu-network/lib/log"
)

type gasCostResponse struct {
	Chain   string `json:"chain,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	GasCost int64  `json:"gas_cost,omitempty"`
}

func (a *ExternalHandler) newGasCostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryStr := r.URL.Query()
		tokenId := queryStr.Get("token_id")
		if len(tokenId) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("missing token_id")); err != nil {
				log.Warn(err)
			}
			return
		}
		chainId := queryStr.Get("chain")
		if len(chainId) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("missing chain")); err != nil {
				log.Warn(err)
			}
			return
		}

		gasCost, err := a.worldState.GetGasCostInToken(tokenId, chainId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		output, err := json.Marshal(&gasCostResponse{
			Chain:   chainId,
			TokenId: tokenId,
			GasCost: gasCost,
		})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(output)
	}
}

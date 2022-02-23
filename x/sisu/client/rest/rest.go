package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/sisu-network/sisu/x/sisu/world"
)

type ExternalHandler struct {
	worldState world.WorldState
}

func NewExternalHandler(worldState world.WorldState) *ExternalHandler {
	return &ExternalHandler{
		worldState: worldState,
	}
}

func (e *ExternalHandler) RegisterRoutes(_ client.Context, r *mux.Router) {
	r.HandleFunc("/getGasFeeInToken", e.newGasCostHandler()).Methods(http.MethodGet)
}

// RegisterRoutes registers sisu-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
}

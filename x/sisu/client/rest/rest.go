package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type ExternalHandler struct {
	keeper      keeper.Keeper
	globalData  components.GlobalData
	deyesClient external.DeyesClient
}

func NewExternalHandler(keeper keeper.Keeper, globalData components.GlobalData,
	deyesClient external.DeyesClient) *ExternalHandler {
	return &ExternalHandler{
		keeper:      keeper,
		globalData:  globalData,
		deyesClient: deyesClient,
	}
}

func (e *ExternalHandler) RegisterRoutes(_ client.Context, r *mux.Router) {
	r.HandleFunc("/getGasFeeInToken", e.newGasCostHandler()).Methods(http.MethodGet)
	r.HandleFunc("/getPubKeys", e.newPubkeyHandler()).Methods(http.MethodGet)
	r.Use(customCORSHeader())
}

func customCORSHeader() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, req)
		})
	}
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

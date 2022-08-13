package rest

import (
	"fmt"
	"net/http"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

var (
	getPoolPath = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Get amount of coins from swap pool
	r.HandleFunc(
		"/swap/pool",
		getPoolHandlerFn(clientCtx),
	).Methods("GET")
}

// HTTP request handler to query a delegator unbonding delegations
func getPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := ctx.Query(getPoolPath)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		ctx = ctx.WithHeight(height)
		rest.PostProcessResponse(w, ctx, res)
	}
}

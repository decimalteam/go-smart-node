package rest

import (
	"fmt"
	"net/http"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Get wallet
	r.HandleFunc(
		"/wallet/{address}",
		getHandlerFn(clientCtx, types.QueryWallet, "address"),
	).Methods("GET")

	// Get wallets by owner
	r.HandleFunc(
		"/wallets/{owner}",
		getHandlerFn(clientCtx, types.QueryWallets, "owner"),
	).Methods("GET")

	// Get transaction by id
	r.HandleFunc(
		"/transaction/{tx_id}",
		getHandlerFn(clientCtx, types.QueryTransaction, "tx_id"),
	).Methods("GET")

	// Get transactions by wallet
	r.HandleFunc(
		"/transactions/{wallet}",
		getHandlerFn(clientCtx, types.QueryTransactions, "wallet"),
	).Methods("GET")

}

func getHandlerFn(clientCtx client.Context, subpath, paramName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := vars[paramName]

		res, height, err := clientCtx.Query(fmt.Sprintf(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, subpath, param)))
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

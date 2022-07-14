package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/gorilla/mux"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router, queryRoute string) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r, queryRoute)
	registerTxRoutes(clientCtx, r, queryRoute)
}

package rest

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	allCoinsPath  = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCoins)
	allChecksPath = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChecks)
	paramsPath    = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Get all coins from a store
	r.HandleFunc(
		"/coin/all_coins",
		allCoinsHandlerFn(clientCtx),
	).Methods("GET")

	// Get coin info by denom
	r.HandleFunc(
		"/coin/coin/{coinDenom}",
		getCoinByDenomHandlerFn(clientCtx),
	).Methods("GET")

	// Get all checks from a store
	r.HandleFunc(
		"/coin/all_checks",
		allChecksHandlerFn(clientCtx),
	).Methods("GET")

	// Get check by hash
	r.HandleFunc(
		"/coin/check/{checkHash}",
		getCheckByHashHandlerFn(clientCtx),
	).Methods("GET")

	// Get the current coin parameter values
	r.HandleFunc(
		"/coin/parameters",
		paramsHandlerFn(clientCtx),
	).Methods("GET")
}

// HTTP request handler to query a delegator unbonding delegations
func allCoinsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.Query(allCoinsPath)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getCoinByDenomHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coinDenom := vars["coinDenom"]

		res, height, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryCoin, coinDenom))
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func allChecksHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.Query(allChecksPath)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getCheckByHashHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		checkHash := vars["checkHash"]

		checkBytes, err := base64.URLEncoding.DecodeString(checkHash)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.NewQueryCheckParams(checkBytes)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheck), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query the coin params values
func paramsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(paramsPath, nil)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

package rest

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"net/http"

	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router, queryRoute string) {
	// Get the total supply of a collection
	r.HandleFunc(
		"/nft/supply/{denom}", getSupply(clientCtx, queryRoute),
	).Methods("GET")

	// Get the collections of NFTs owned by an address
	r.HandleFunc(
		"/nft/owner/{delegatorAddr}", getOwner(clientCtx, queryRoute),
	).Methods("GET")

	// Get the NFTs owned by an address from a given collection
	r.HandleFunc(
		"/nft/owner/{delegatorAddr}/collection/{denom}", getOwnerByDenom(clientCtx, queryRoute),
	).Methods("GET")

	// Get all the NFT from a given collection
	r.HandleFunc(
		"/nft/collection/{denom}", getCollection(clientCtx, queryRoute),
	).Methods("GET")

	// Query all denoms
	r.HandleFunc(
		"/nft/denoms", getDenoms(clientCtx, queryRoute),
	).Methods("GET")

	// Query a single NFT
	r.HandleFunc(
		"/nft/collection/{denom}/nft/{id}", getNFT(clientCtx, queryRoute),
	).Methods("GET")
}

func getSupply(clientCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := mux.Vars(r)["denom"]

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.Query(fmt.Sprintf("custom/%s/supply/%s", queryRoute, denom))
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getOwner(clientCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address, err := sdk.AccAddressFromBech32(mux.Vars(r)["delegatorAddr"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryBalanceParams(address, "")
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/owner", queryRoute), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)

	}
}

func getOwnerByDenom(clientCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars["denom"]
		address, err := sdk.AccAddressFromBech32(vars["delegatorAddr"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.NewQueryBalanceParams(address, denom)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/ownerByDenom", queryRoute), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getCollection(clientCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := mux.Vars(r)["denom"]

		params := types.NewQueryCollectionParams(denom)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/collection", queryRoute), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getDenoms(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/denoms", queryRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getNFT(clientCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars["denom"]
		id := vars["id"]

		params := types.NewQueryNFTParams(denom, id)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/nft", queryRoute), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

package api

import sdk "github.com/cosmos/cosmos-sdk/types"

type NFT struct {
	ID        string
	Owners    TokenOwners
	Creator   string
	TokenURI  string
	Reserve   sdk.Int
	AllowMint bool
}

type TokenOwners []TokenOwner

type TokenOwner struct {
	Address     string
	SubTokenIDs []uint64
}

func (api *API) NFTCollections() ([]string, error) {
	type responseDenomsType struct {
		Result []string `json:"result"`
	}
	// request
	res, err := api.rest.R().Get("/nft/denoms")
	if err = processConnectionError(res, err); err != nil {
		return []string{}, err
	}
	// json decode
	response := responseDenomsType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return len(response.Result) > 0, false
	})
	if err != nil {
		return []string{}, err
	}
	// process result
	return response.Result, nil
}

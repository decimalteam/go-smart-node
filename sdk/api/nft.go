package api

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFT struct {
	ID        string
	Owners    []TokenOwner
	Creator   string
	TokenURI  string
	Denom     string
	Reserve   sdk.Int
	AllowMint bool
}

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

func (api *API) NFTsByDenom(denom string) ([]string, error) {
	type responseType struct {
		Result struct {
			Denom string   `json:"denom"`
			Nfts  []string `json:"nfts"`
		} `json:"result"`
	}
	// request
	res, err := api.rest.R().Get(fmt.Sprintf("/nft/collection/%s", denom))
	if err = processConnectionError(res, err); err != nil {
		return []string{}, err
	}
	// json decode
	response := responseType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return response.Result.Denom == denom, false
	})
	if err != nil {
		return []string{}, err
	}
	// process result
	return response.Result.Nfts, nil
}

func (api *API) NFT(denom string, id string) (*NFT, error) {
	type responseType struct {
		Result struct {
			Id       string `json:"id"`
			Creator  string `json:"creator"`
			TokenUri string `json:"token_uri"`
			Reserve  string `json:"reserve"`
			Owners   []struct {
				Address     string   `json:"address"`
				SubTokenIds []string `json:"sub_token_ids"`
			}
		} `json:"result"`
	}
	// request
	res, err := api.rest.R().Get(fmt.Sprintf("/nft/collection/%s/nft/%s", denom, id))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	response := responseType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return response.Result.Id == id, false
	})
	if err != nil {
		return nil, err
	}
	// process result
	var ok bool
	result := NFT{}
	result.ID = response.Result.Id
	result.Denom = denom
	result.Creator = response.Result.Creator
	result.TokenURI = response.Result.TokenUri
	result.Reserve, ok = sdk.NewIntFromString(response.Result.Reserve)
	if !ok {
		return nil, fmt.Errorf("cannot parse sdk.Int from %s", response.Result.Reserve)
	}
	result.Owners = make([]TokenOwner, len(response.Result.Owners))
	for i, o := range response.Result.Owners {
		subtokens := make([]uint64, len(o.SubTokenIds))
		for j := range o.SubTokenIds {
			subtokens[j], err = strconv.ParseUint(o.SubTokenIds[j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot parse subtoken owner: %s, sub token id: %s", o.Address, o.SubTokenIds[j])
			}
		}
		result.Owners[i] = TokenOwner{Address: o.Address, SubTokenIDs: subtokens}
	}

	return &result, nil
}

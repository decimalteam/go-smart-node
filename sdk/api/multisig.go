package api

import (
	"fmt"
	"strconv"
)

type MultisigWallet struct {
	Address      string
	Owners       []string
	Weights      []uint64
	Threshold    uint64
	LegacyOwners []string
}

func (api *API) MultisigWalletsByOwner(owner string) ([]MultisigWallet, error) {
	type responseType struct {
		Result []struct {
			Address      string   `json:"address"`
			Owners       []string `json:"owners"`
			Weights      []string `json:"weights"`
			LegacyOwners []string `json:"legacy_owners"`
			Threshold    string   `json:"threshold"`
		} `json:"result"`
		Height string `json:"height"`
	}
	// request
	res, err := api.rest.R().Get(fmt.Sprintf("/multisig/wallets/%s", owner))
	if err = processConnectionError(res, err); err != nil {
		return []MultisigWallet{}, err
	}
	// json decode
	response := responseType{}
	err = universalJSONDecode(res.Body(), &response, nil, func() (bool, bool) {
		return response.Height > "", false
	})
	if err != nil {
		return []MultisigWallet{}, err
	}
	// process result
	var result = make([]MultisigWallet, len(response.Result))
	for i, w := range response.Result {
		result[i].Address = w.Address
		result[i].Owners = w.Owners
		result[i].Threshold, err = strconv.ParseUint(w.Threshold, 10, 64)
		if err != nil {
			return []MultisigWallet{}, err
		}
		result[i].Weights = make([]uint64, len(w.Weights))
		for j := range w.Weights {
			result[i].Weights[j], err = strconv.ParseUint(w.Weights[j], 10, 64)
			if err != nil {
				return []MultisigWallet{}, err
			}
		}
	}
	return result, nil

}

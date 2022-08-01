package api

import (
	"errors"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddressResult contains API response fields.
type AddressResult struct {
	ID      uint64            `json:"id"`
	Address string            `json:"address"`
	Type    string            `json:"type"`
	Nonce   string            `json:"nonce"`
	Balance map[string]string `json:"balance"`
	Txes    uint64            `json:"txes"`
}

// Address requests full information about specified address
func (api *API) Address(address string) (*AddressResult, error) {
	type respDirectAddress struct {
		Result struct {
			Value struct {
				AccountNumber string `json:"account_number"`
				Address       string `json:"address"`
				Sequence      string `json:"sequence"`
			} `json:"base_account"`
		} `json:"result"`
	}
	// request
	res, err := api.rest.R().Get(fmt.Sprintf("/auth/accounts/%s", address))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := respDirectAddress{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result.Value.AccountNumber > "", respErr.StatusCode != 0
	})
	//empty account (no transactions), it's normal
	if errors.Is(err, ErrMissingLogic) {
		return &AddressResult{
			ID:      0,
			Address: address,
			Nonce:   "0",
		}, nil
	}
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	// process result
	accNumber, _ := strconv.ParseUint(respValue.Result.Value.AccountNumber, 10, 64)

	return &AddressResult{
		ID:      accNumber,
		Address: respValue.Result.Value.Address,
		Nonce:   respValue.Result.Value.Sequence,
	}, nil
}

// AccountNumberAndSequence requests account number and current sequence (nonce) of specified address.
func (api *API) AccountNumberAndSequence(address string) (uint64, uint64, error) {
	adrRes, err := api.Address(address)
	if err != nil {
		return 0, 0, err
	}
	seq, _ := strconv.ParseUint(adrRes.Nonce, 10, 64)
	return adrRes.ID, seq, nil

}

// AddressBalance returns balance for address
func (api *API) AddressBalance(address string) (sdk.Coins, error) {
	type directBalanceResult struct {
		Height string `json:"height"`
		Result []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"result"`
	}
	// request
	res, err := api.rest.R().Get("/bank/balances/" + address)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := directBalanceResult{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Height > "", respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	// process result
	var result sdk.Coins
	for _, rawCoin := range respValue.Result {
		amount, ok := sdk.NewIntFromString(rawCoin.Amount)
		if !ok {
			return nil, fmt.Errorf("not ok Amount='%s'", rawCoin.Amount)
		}
		result = result.Add(sdk.NewCoin(rawCoin.Denom, amount))
	}
	return result, nil
}

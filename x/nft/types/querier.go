package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryCollectionParams defines the params for queries:
type QueryCollectionParams struct {
	Denom string `json:"denom"`
}

// NewQueryCollectionParams creates a new instance of QuerySupplyParams
func NewQueryCollectionParams(denom string) QueryCollectionParams {
	return QueryCollectionParams{Denom: denom}
}

// Bytes exports the Denom as bytes
func (q QueryCollectionParams) Bytes() []byte {
	return []byte(q.Denom)
}

// QueryBalanceParams params for query 'custom/nfts/balance'
type QueryBalanceParams struct {
	Owner sdk.AccAddress
	Denom string // optional
}

// NewQueryBalanceParams creates a new instance of QuerySupplyParams
func NewQueryBalanceParams(owner sdk.AccAddress, denom ...string) QueryBalanceParams {
	if len(denom) > 0 {
		return QueryBalanceParams{
			Owner: owner,
			Denom: denom[0],
		}
	}
	return QueryBalanceParams{Owner: owner}
}

// QueryNFTParams params for query 'custom/nfts/nft'
type QueryNFTParams struct {
	Denom   string `json:"denom"`
	TokenID string `json:"token_id"`
}

// NewQueryNFTParams creates a new instance of QueryNFTParams
func NewQueryNFTParams(denom, id string) QueryNFTParams {
	return QueryNFTParams{
		Denom:   denom,
		TokenID: id,
	}
}

// QuerySubTokensParams params for query 'custom/nfts/sub_tokens'
type QuerySubTokensParams struct {
	Denom       string   `json:"denom"`
	TokenID     string   `json:"token_id"`
	SubTokenIDs []uint64 `json:"sub_token_ids"`
}

// NewQuerySubTokensParams creates a new instance of QuerySubTokensParams
func NewQuerySubTokensParams(denom, id string, subTokenIDs []uint64) QuerySubTokensParams {
	return QuerySubTokensParams{
		Denom:       denom,
		TokenID:     id,
		SubTokenIDs: subTokenIDs,
	}
}

type ResponseSubTokens []ResponseSubToken

type ResponseSubToken struct {
	ID      uint64
	Reserve sdk.Coin
}

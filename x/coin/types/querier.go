package types

import "github.com/cosmos/cosmos-sdk/types/query"

const (
	QueryParams = "params" // 'custom/coin/params'
	QueryCoin   = "coin"   // 'custom/coin/coin/{denom}'
	QueryCoins  = "coins"  // 'custom/coin/coins'
	QueryCheck  = "check"  // 'custom/coin/check' with req
	QueryChecks = "checks" // 'custom/coin/checks'
)

// QueryCheckParams params for query 'custom/coin/check'
type QueryCheckParams struct {
	Hash []byte `json:"hash"`
}

// NewQueryCheckParams creates a new instance of QueryCheckParams
func NewQueryCheckParams(hash []byte) QueryCheckParams {
	return QueryCheckParams{
		Hash: hash,
	}
}

func NewQueryCoinRequest(symbol string) *QueryCoinRequest {
	return &QueryCoinRequest{
		Symbol: symbol,
	}
}

func NewQueryCheckRequest(hash []byte) *QueryCheckRequest {
	return &QueryCheckRequest{
		Hash: hash,
	}
}

func NewQueryCoinsRequest(pagination *query.PageRequest) *QueryCoinsRequest {
	return &QueryCoinsRequest{
		Pagination: pagination,
	}
}

func NewQueryChecksRequest(pagination *query.PageRequest) *QueryChecksRequest {
	return &QueryChecksRequest{
		Pagination: pagination,
	}
}

func NewQueryParamsRequest() *QueryParamsRequest { return &QueryParamsRequest{} }

package types

import "github.com/cosmos/cosmos-sdk/types/query"

const (
	QueryCoins  = "coins"  // 'custom/coin/coins'
	QueryCoin   = "coin"   // 'custom/coin/coin/{denom}'
	QueryChecks = "checks" // 'custom/coin/checks'
	QueryCheck  = "check"  // 'custom/coin/check' with req
	QueryParams = "params" // 'custom/coin/params'
)

func NewQueryCoinsRequest(pagination *query.PageRequest) *QueryCoinsRequest {
	return &QueryCoinsRequest{
		Pagination: pagination,
	}
}

func NewQueryCoinRequest(denom string) *QueryCoinRequest {
	return &QueryCoinRequest{
		Denom: denom,
	}
}

func NewQueryCoinByEVMRequest(drc20_address string) *QueryCoinByEVMRequest {
	return &QueryCoinByEVMRequest{
		Drc20Address: drc20_address,
	}
}

func NewQueryChecksRequest(pagination *query.PageRequest) *QueryChecksRequest {
	return &QueryChecksRequest{
		Pagination: pagination,
	}
}

func NewQueryCheckRequest(hash []byte) *QueryCheckRequest {
	return &QueryCheckRequest{
		Hash: hash,
	}
}

func NewQueryParamsRequest() *QueryParamsRequest { return &QueryParamsRequest{} }

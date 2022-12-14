package api

import (
	"context"

	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type NFTCollection = nfttypes.Collection
type NFTToken = nfttypes.Token

type SubToken = nfttypes.SubToken

// Returns all NFT collections (denoms)
func (api *API) NFTCollections() ([]NFTCollection, error) {
	client := nfttypes.NewQueryClient(api.grpcClient)
	colls := make([]nfttypes.Collection, 0)
	req := &nfttypes.QueryCollectionsRequest{
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.Collections(
			context.Background(),
			req,
		)
		if err != nil {
			return []NFTCollection{}, err
		}
		if len(res.Collections) == 0 {
			break
		}
		colls = append(colls, res.Collections...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return colls, nil
}

// Returns NFT IDs from collection
func (api *API) NFTCollection(creator, denom string) (NFTCollection, error) {
	client := nfttypes.NewQueryClient(api.grpcClient)
	res, err := client.Collection(
		context.Background(),
		&nfttypes.QueryCollectionRequest{
			Creator: creator,
			Denom:   denom,
		},
	)
	if err != nil {
		return NFTCollection{}, err
	}
	return res.Collection, nil
}

func (api *API) NFTToken(id string) (NFTToken, error) {
	client := nfttypes.NewQueryClient(api.grpcClient)
	res, err := client.Token(
		context.Background(),
		&nfttypes.QueryTokenRequest{
			TokenId: id,
		},
	)
	if err != nil {
		return NFTToken{}, err
	}
	return res.Token, nil
}

func (api *API) NFTSubToken(tokenID string, subTokenID string) (SubToken, error) {
	client := nfttypes.NewQueryClient(api.grpcClient)
	res, err := client.SubToken(
		context.Background(),
		&nfttypes.QuerySubTokenRequest{
			TokenId:    tokenID,
			SubTokenId: subTokenID,
		},
	)
	if err != nil {
		return SubToken{}, err
	}
	return res.SubToken, nil
}

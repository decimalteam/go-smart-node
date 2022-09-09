package api

import (
	"context"

	nftTypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

type NFTCollection = nftTypes.Collection
type NFT struct {
	nftTypes.BaseNFT
	Denom string
}

type SubToken = nftTypes.SubToken

// Returns all NFT collections (denoms)
func (api *API) NFTCollections() ([]string, error) {
	client := nftTypes.NewQueryClient(api.grpcClient)
	res, err := client.QueryDenoms(
		context.Background(),
		&nftTypes.QueryDenomsRequest{},
	)
	if err != nil {
		return []string{}, err
	}
	return res.Denoms, nil
}

// Returns NFT IDs from collection
func (api *API) NFTCollection(denom string) (NFTCollection, error) {
	client := nftTypes.NewQueryClient(api.grpcClient)
	res, err := client.QueryCollection(
		context.Background(),
		&nftTypes.QueryCollectionRequest{
			Denom: denom,
		},
	)
	if err != nil {
		return NFTCollection{}, err
	}
	return res.Collection, nil
}

func (api *API) NFT(denom string, id string) (NFT, error) {
	client := nftTypes.NewQueryClient(api.grpcClient)
	res, err := client.QueryNFT(
		context.Background(),
		&nftTypes.QueryNFTRequest{
			Denom:   denom,
			TokenId: id,
		},
	)
	if err != nil {
		return NFT{}, err
	}
	return NFT{res.NFT, denom}, nil
}

func (api *API) NFTSubTokens(denom string, tokenID string, ids []uint64) ([]SubToken, error) {
	client := nftTypes.NewQueryClient(api.grpcClient)
	res, err := client.QuerySubTokens(
		context.Background(),
		&nftTypes.QuerySubTokensRequest{
			Denom:   denom,
			TokenID: tokenID,
		},
	)
	if err != nil {
		return []SubToken{}, err
	}
	return res.SubTokens, nil
}

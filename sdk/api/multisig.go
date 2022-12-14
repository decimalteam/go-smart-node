package api

import (
	"context"

	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type MultisigWallet = multisigtypes.Wallet
type MultisigTransaction = multisigtypes.Transaction

type MultisigTransactionInfo struct {
	Transaction MultisigTransaction
	Signers     []string
	Completed   bool
}

func (api *API) MultisigWalletsByOwner(owner string) ([]MultisigWallet, error) {
	client := multisigtypes.NewQueryClient(api.grpcClient)
	wallets := make([]MultisigWallet, 0)
	req := &multisigtypes.QueryWalletsRequest{
		Owner:      owner,
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.Wallets(
			context.Background(),
			req,
		)
		if err != nil {
			return []MultisigWallet{}, err
		}
		wallets = append(wallets, res.Wallets...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return wallets, nil
}

func (api *API) MultisigWalletByAddress(address string) (MultisigWallet, error) {
	client := multisigtypes.NewQueryClient(api.grpcClient)
	res, err := client.Wallet(
		context.Background(),
		&multisigtypes.QueryWalletRequest{
			Wallet: address,
		},
	)
	if err != nil {
		return MultisigWallet{}, err
	}
	return res.Wallet, nil
}

func (api *API) MultisigTransactionsByWallet(address string) ([]MultisigTransaction, error) {
	client := multisigtypes.NewQueryClient(api.grpcClient)
	txs := make([]MultisigTransaction, 0)
	req := &multisigtypes.QueryTransactionsRequest{
		Wallet:     address,
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.Transactions(
			context.Background(),
			req,
		)
		if err != nil {
			return []MultisigTransaction{}, err
		}
		txs = append(txs, res.Transactions...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return txs, nil
}

func (api *API) MultisigTransactionsByID(txID string) (MultisigTransactionInfo, error) {
	client := multisigtypes.NewQueryClient(api.grpcClient)
	res, err := client.Transaction(
		context.Background(),
		&multisigtypes.QueryTransactionRequest{
			Id: txID,
		},
	)
	if err != nil {
		return MultisigTransactionInfo{}, err
	}
	return MultisigTransactionInfo{
		Transaction: res.Transaction,
		Signers:     res.Signers,
		Completed:   res.Completed,
	}, nil
}

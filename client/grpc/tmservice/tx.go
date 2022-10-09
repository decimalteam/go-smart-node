package tmservice

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
)

func getTxSearch(ctx context.Context, clientCtx client.Context, query string, prove bool, page, perPage int64, orderBy string) ([]*TxResult, int, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, 0, err
	}

	intPage := int(page)
	intPerPage := int(perPage)

	res, err := node.TxSearch(ctx, query, prove, &intPage, &intPerPage, orderBy)

	protoResults := make([]*TxResult, len(res.Txs))
	for i, v := range res.Txs {
		proof := v.Proof.ToProto()
		txResult := v.TxResult
		protoResults[i] = &TxResult{
			Hash:     v.Hash,
			Height:   v.Height,
			Index:    v.Index,
			TxResult: &txResult,
			Tx:       v.Tx,
			TxProof:  &proof,
		}
	}

	return protoResults, res.TotalCount, nil
}

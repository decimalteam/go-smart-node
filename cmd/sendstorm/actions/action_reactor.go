package actions

import (
	"fmt"
	"math/rand"
)

type ActionReactor struct {
	wags []*WeightedAG
	wsum int64
}

type WeightedAG struct {
	Weight int64
	AG     ActionGenerator
}

// TODO: parameters for generator
func (ar *ActionReactor) Add(generatorName string, weight int64) error {
	var wag *WeightedAG = &WeightedAG{Weight: weight}
	switch generatorName {
	// coin
	case "CreateCoin":
		wag.AG = NewCreateCoinGenerator(3, 9, 100, 1000, 1000, 2000, 1000000, 2000000)
	case "SendCoin":
		wag.AG = NewSendCoinGenerator(500, 20000)
	case "BuyCoin":
		wag.AG = NewBuyCoinGenerator(500, 20000, "del")
	case "SellCoin":
		wag.AG = NewSellCoinGenerator(500, 20000, "del")
	case "MultiSendCoin":
		wag.AG = NewMultiSendCoinGenerator(500, 20000, 1, 10)
	case "UpdateCoin":
		wag.AG = NewUpdateCoinGenerator(1, 100, "del")
	case "BurnCoin":
		wag.AG = NewBurnCoinGenerator(1, 1000)
	case "RedeemCheck":
		wag.AG = NewRedeemCheckGenerator(1, 1000)
	// nft
	case "MintNFT":
		wag.AG = NewMintNFTGenerator(1, 100, 100, 1000, 1, 10)
	case "TransferNFT":
		wag.AG = NewTransferNFTGenerator()
	case "EditNFT":
		wag.AG = NewEditNFTGenerator(1, 100)
	case "BurnNFT":
		wag.AG = NewBurnNFTGenerator()
	case "UpdateReserveNFT":
		wag.AG = NewUpdateReserveNFTGenerator(1, 100)
	// multisig
	case "CreateMultisigWallet":
		wag.AG = NewCreateMultisigWalletGenerator()
	case "DepositMultisigWallet":
		wag.AG = NewDepositMultisigWalletGenerator(100, 10000)
	case "CreateMultisigTransaction":
		wag.AG = NewCreateMultisigTransactionGenerator(100, 10000)
	case "SignMultisigTransaction":
		wag.AG = NewSignMultisigTransactionGenerator()
	// validator
	case "CreateValidator":
		wag.AG = NewCreateValidatorGenerator(100, 1000)
	case "EditValidator":
		wag.AG = NewEditValidatorGenerator()
	case "Delegate":
		wag.AG = NewDelegateGenerator(1, 100)
	case "DelegateNFT":
		wag.AG = NewDelegateNFTGenerator()
	case "Undelegate":
		wag.AG = NewUndelegateGenerator()
	case "UndelegateNFT":
		wag.AG = NewUndelegateNFTGenerator()
	case "Redelegate":
		wag.AG = NewRedelegateGenerator()
	case "RedelegateNFT":
		wag.AG = NewRedelegateNFTGenerator()
	}
	if wag.AG == nil {
		return fmt.Errorf("%s: unknown generator name", generatorName)
	}
	ar.wsum += weight
	ar.wags = append(ar.wags, wag)
	return nil
}

// choose generator and generate action
func (ar *ActionReactor) Generate() Action {
	w := rand.Int63n(ar.wsum)
	for _, wag := range ar.wags {
		if w < wag.Weight {
			return wag.AG.Generate()
		}
		w -= wag.Weight
	}
	// we can not be here, this is for stub
	return ar.wags[0].AG.Generate()
}

func (ar *ActionReactor) Update(ui UpdateInfo) {
	for _, wag := range ar.wags {
		wag.AG.Update(ui)
	}
}

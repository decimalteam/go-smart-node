package main

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

//TODO: parameters for generator
func (ar *ActionReactor) Add(generatorName string, weight int64) error {
	var wag *WeightedAG = nil
	switch generatorName {
	case "CreateCoin":
		{
			wag = &WeightedAG{
				AG:     NewCreateCoinGenerator(3, 9, 100, 1000, 1000, 2000, 1000000, 2000000),
				Weight: weight,
			}
		}
	case "SendCoin":
		{
			wag = &WeightedAG{
				AG:     NewSendCoinGenerator(500, 20000),
				Weight: weight,
			}
		}
	case "BuyCoin":
		{
			wag = &WeightedAG{
				AG:     NewBuyCoinGenerator(500, 20000, "del"),
				Weight: weight,
			}
		}
	case "SellCoin":
		{
			wag = &WeightedAG{
				AG:     NewSellCoinGenerator(500, 20000, "del"),
				Weight: weight,
			}
		}
	case "MultiSendCoin":
		{
			wag = &WeightedAG{
				AG:     NewMultiSendCoinGenerator(500, 20000, 1, 10),
				Weight: weight,
			}
		}
	case "UpdateCoin":
		{
			wag = &WeightedAG{
				AG:     NewUpdateCoinGenerator(1, 100, "del"),
				Weight: weight,
			}
		}
	case "MintNFT":
		{
			wag = &WeightedAG{
				AG:     NewMintNFTGenerator(1, 100, 100, 1000, 1, 10),
				Weight: weight,
			}
		}
	}
	if wag == nil {
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

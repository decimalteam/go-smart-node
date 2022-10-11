package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewStakeCoin creates a new Stake instance for stake in the coin.
func NewStakeCoin(stake sdk.Coin) Stake {
	return Stake{
		Type:        StakeType_Coin,
		ID:          stake.Denom,
		Stake:       stake,
		SubTokenIDs: []int64{},
	}
}

// NewStakeNFT creates a new Stake instance for stake in the NFT token.
func NewStakeNFT(tokenID string, subTokenIDs []int64, reserve sdk.Coin) Stake {
	return Stake{
		Type:        StakeType_NFT,
		ID:          tokenID,
		Stake:       reserve,
		SubTokenIDs: subTokenIDs,
	}
}

func (s Stake) GetType() StakeType {
	return s.Type
}

func (s Stake) GetID() string {
	return s.ID
}

func (s Stake) GetStake() sdk.Coin {
	return s.Stake
}

func (s Stake) GetSubTokenIDs() []int64 {
	return s.SubTokenIDs
}

package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewStakeCoin creates a new Stake instance for stake in the coin.
func NewStakeCoin(stake sdk.Coin) Stake {
	return Stake{
		Type:        StakeType_Coin,
		ID:          stake.Denom,
		Stake:       stake,
		SubTokenIDs: []uint32{},
	}
}

// NewStakeNFT creates a new Stake instance for stake in the NFT token.
func NewStakeNFT(tokenID string, subTokenIDs []uint32, reserve sdk.Coin) Stake {
	return Stake{
		Type:        StakeType_NFT,
		ID:          tokenID,
		Stake:       reserve,
		SubTokenIDs: subTokenIDs,
	}
}

// GetType returns the stake type.
func (s Stake) GetType() StakeType {
	return s.Type
}

// GetID returns the stake ID.
// For stake in Coin: contains coin denom value.
// For stake in NFT: contains NFT token ID value.
func (s Stake) GetID() string {
	return s.ID
}

// GetStake returns amount of the coin staked.
// For stake in Coin: contains actually staked coin.
// For stake in NFT: contains total reserve of staked NFT sub-tokens.
func (s Stake) GetStake() sdk.Coin {
	return s.Stake
}

// GetSubTokenIDs returns the list of staked NFT sub-token IDs.
func (s Stake) GetSubTokenIDs() []uint32 {
	return s.SubTokenIDs
}

func (s Stake) AddSubTokens(newSubTokens []uint32) ([]uint32, error) {
	if s.Type != StakeType_NFT {
		return nil, fmt.Errorf("delegation type is coin")
	}
	existSubTokens := make(map[uint32]bool)
	for _, v := range s.SubTokenIDs {
		existSubTokens[v] = false
	}

	for _, v := range newSubTokens {
		if _, ok := existSubTokens[v]; ok {
			return nil, fmt.Errorf("subtoken exists in delegations")
		}
		s.SubTokenIDs = append(s.SubTokenIDs, v)
	}

	return newSubTokens, nil
}

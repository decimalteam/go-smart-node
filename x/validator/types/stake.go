package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewStakeCoin creates a new Stake instance for stake in the coin.
func NewStakeCoin(stake sdk.Coin) Stake {
	return Stake{
		Type:        StakeType_Coin,
		ID:          stake.Denom,
		Stake:       stake,
		SubTokenIDs: nil,
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
		return nil, errors.WrongStakeType
	}
	var result = make([]uint32, len(s.SubTokenIDs)+len(newSubTokens))
	existSubTokens := make(map[uint32]bool)
	for i, v := range s.SubTokenIDs {
		if existSubTokens[v] {
			return nil, errors.SubTokenExistsInStake
		}
		existSubTokens[v] = true
		result[i] = v
	}
	for i, v := range newSubTokens {
		if existSubTokens[v] {
			return nil, errors.SubTokenExistsInStake
		}
		existSubTokens[v] = true
		result[i+len(s.SubTokenIDs)] = v
	}

	return result, nil
}

func (s Stake) IsEmpty() bool {
	return s.Stake.IsNil() || s.Stake.IsZero() || (s.Type == StakeType_NFT && len(s.SubTokenIDs) == 0)
}

func (s Stake) Add(a Stake) (Stake, error) {
	if s.Type != a.Type {
		return Stake{}, errors.WrongStakeType
	}
	if s.ID != a.ID {
		return Stake{}, errors.WrongStakeID
	}
	var result Stake
	var err error
	result.Type = s.Type
	result.ID = s.ID
	result.Stake = s.Stake.Add(a.GetStake())
	if s.Type == StakeType_NFT {
		result.SubTokenIDs, err = s.AddSubTokens(a.SubTokenIDs)
		if err != nil {
			return Stake{}, err
		}
	}
	return result, nil
}

func (s Stake) Equal(a *Stake) bool {
	if s.Type != a.Type {
		return false
	}
	if s.ID != a.ID {
		return false
	}
	if !s.Stake.Equal(a.Stake) {
		return false
	}
	if len(s.SubTokenIDs) != len(a.SubTokenIDs) {
		return false
	}
	for _, id1 := range s.SubTokenIDs {
		isHere := false
		for _, id2 := range a.SubTokenIDs {
			if id1 == id2 {
				isHere = true
				break
			}
		}
		if !isHere {
			return false
		}
	}
	return true
}

// return true if set has subset of IDs
func SetHasSubset(set []uint32, subset []uint32) bool {
	for _, id1 := range subset {
		inSet := false
		for _, id2 := range set {
			if id1 == id2 {
				inSet = true
				break
			}
		}
		if !inSet {
			return false
		}
	}
	return true
}

// return true if sets have same elements
func SetHasIntersection(set []uint32, otherset []uint32) bool {
	for _, id1 := range otherset {
		for _, id2 := range set {
			if id1 == id2 {
				return true
			}
		}
	}
	return false
}

// returns elements of (set-subset)
func SetSubstract(set []uint32, subset []uint32) []uint32 {
	var result = []uint32{}
	for _, id1 := range set {
		substract := true
		for _, id2 := range subset {
			if id1 == id2 {
				substract = false
				break
			}
		}
		if substract {
			result = append(result, id1)
		}
	}
	return result
}

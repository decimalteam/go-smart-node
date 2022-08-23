package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id string, creator, tokenURI string, reserve sdk.Coin, allowMint bool) BaseNFT {
	return BaseNFT{
		ID:        id,
		TokenURI:  strings.TrimSpace(tokenURI),
		Creator:   creator,
		Reserve:   reserve,
		AllowMint: allowMint,
	}
}

// GetID returns the ID of the token
func (bnft BaseNFT) GetID() string {
	return bnft.ID
}

func (bnft BaseNFT) GetOwners() TokenOwners {
	return bnft.Owners
}

// GetTokenURI returns the path to optional extra properties
func (bnft BaseNFT) GetTokenURI() string {
	return bnft.TokenURI
}

func (bnft BaseNFT) GetCreator() string {
	return bnft.Creator
}

// EditMetadata edits metadata of an nft
func (bnft BaseNFT) EditMetadata(tokenURI string) BaseNFT {
	bnft.TokenURI = tokenURI
	return bnft
}

func (bnft BaseNFT) SetOwners(owners TokenOwners) BaseNFT {
	bnft.Owners = owners
	return bnft
}

func (bnft BaseNFT) GetReserve() sdk.Coin {
	return bnft.Reserve
}

func (bnft BaseNFT) GetAllowMint() bool {
	return bnft.AllowMint
}

func (bnft BaseNFT) GenSubTokenIDs(quantity uint64) SortedUintArray {
	var prevSubTokenID uint64 = 0
	for _, o := range bnft.GetOwners() {
		max := o.GetSubTokenIDs().Max()
		if max > prevSubTokenID {
			prevSubTokenID = max
		}
	}

	subTokenIDs := make(SortedUintArray, quantity)

	for i := uint64(0); i < quantity; i++ {
		prevSubTokenID++
		subTokenIDs[i] = prevSubTokenID
	}

	return subTokenIDs
}

func (bnft BaseNFT) AddOwnerSubTokenIDs(ownerAddress string, subTokenIDs SortedUintArray) BaseNFT {
	owner := bnft.GetOwners().GetOwner(ownerAddress)
	if owner == nil {
		owner = NewTokenOwner(ownerAddress, subTokenIDs)

		bnft = bnft.SetOwners(bnft.GetOwners().SetOwner(owner))
	} else {
		for _, id := range subTokenIDs {
			owner = owner.SetSubTokenID(id)
		}

		bnft = bnft.SetOwners(bnft.GetOwners().SetOwner(owner))
	}

	return bnft
}

func (bnft BaseNFT) String() string {
	return fmt.Sprintf(`ID:				%s
Owners:			%s
TokenURI:		%s`,
		bnft.ID,
		bnft.Owners.String(),
		bnft.TokenURI,
	)
}

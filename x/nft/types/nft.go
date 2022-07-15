package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/exported"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
	"strings"
)

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id string, creator, owner string, tokenURI string, reserve sdk.Int, subTokenIDs []int64, allowMint bool) BaseNFT {
	return BaseNFT{
		ID: id,
		Owners: TokenOwners{TokenOwner{
			Address:     owner,
			SubTokenIDs: subTokenIDs,
		}},
		TokenURI:  strings.TrimSpace(tokenURI),
		Creator:   creator,
		Reserve:   reserve,
		AllowMint: allowMint,
	}
}

// GetID returns the ID of the token
func (bnft BaseNFT) GetID() string { return bnft.ID }

func (bnft BaseNFT) GetOwners() exported.TokenOwners { return &bnft.Owners }

// GetTokenURI returns the path to optional extra properties
func (bnft BaseNFT) GetTokenURI() string { return bnft.TokenURI }

func (bnft BaseNFT) GetCreator() string {
	return bnft.Creator
}

// EditMetadata edits metadata of an nft
func (bnft BaseNFT) EditMetadata(tokenURI string) exported.NFT {
	bnft.TokenURI = tokenURI
	return &bnft
}

func (bnft BaseNFT) SetOwners(owners exported.TokenOwners) exported.NFT {
	bnft.Owners = owners.(TokenOwners)
	return &bnft
}

func (bnft BaseNFT) GetReserve() sdk.Int {
	return bnft.Reserve
}

func (bnft BaseNFT) GetAllowMint() bool {
	return bnft.AllowMint
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

// ----------------------------------------------------------

// NFTs define a list of NFT
type NFTs []BaseNFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...BaseNFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}

	return NFTs(nfts).Sort()
}

// Append appends two sets of NFTs
func (nfts NFTs) Append(nftsB ...BaseNFT) NFTs {
	return append(nfts, nftsB...).Sort()
}

// Find returns the searched collection from the set
func (nfts NFTs) Find(id string) (nft BaseNFT, found bool) {
	index := nfts.find(id)
	if index == -1 {
		return nft, false
	}
	return nfts[index], true
}

// Update removes and replaces an NFT from the set
func (nfts NFTs) Update(id string, nft exported.NFT) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(append(nfts[:index], *nft.(*BaseNFT)), nfts[index+1:]...), true
}

// Remove removes an NFT from the set of NFTs
func (nfts NFTs) Remove(id string) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(nfts[:index], nfts[index+1:]...), true
}

// String follows stringer interface
func (nfts NFTs) String() string {
	if len(nfts) == 0 {
		return ""
	}

	out := ""
	for _, nft := range nfts {
		out += fmt.Sprintf("%v\n", nft.String())
	}
	return out[:len(out)-1]
}

// Empty returns true if there are no NFTs and false otherwise.
func (nfts NFTs) Empty() bool {
	return len(nfts) == 0
}

func (nfts NFTs) find(id string) int {
	return FindUtil(nfts, id)
}

// Findable and Sort interfaces
func (nfts NFTs) ElAtIndex(index int) string { return nfts[index].GetID() }
func (nfts NFTs) Len() int                   { return len(nfts) }
func (nfts NFTs) Less(i, j int) bool         { return strings.Compare(nfts[i].GetID(), nfts[j].GetID()) == -1 }
func (nfts NFTs) Swap(i, j int)              { nfts[i], nfts[j] = nfts[j], nfts[i] }

var _ sort.Interface = NFTs{}

// Sort is a helper function to sort the set of coins in place
func (nfts NFTs) Sort() NFTs {
	sort.Sort(nfts)
	return nfts
}

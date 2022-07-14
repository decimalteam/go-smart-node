package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/exported"
	"fmt"
	"sort"
	"strings"
)

// NewIDCollection creates a new IDCollection instance
func NewIDCollection(denom string, ids []string) IDCollection {
	return IDCollection{
		Denom: strings.TrimSpace(denom),
		IDs:   SortedStringArray(ids).Sort(),
	}
}

func (m IDCollection) Supply() int {
	return len(m.IDs)
}

// AddID adds an ID to the idCollection
func (m IDCollection) AddID(id string) IDCollection {
	m.IDs = append(m.IDs, id).Sort()
	return m
}

// DeleteID deletes an ID from an ID Collection
func (m IDCollection) DeleteID(id string) (IDCollection, error) {
	index := FindUtil(m.IDs, id)
	if index == -1 {
		return m, ErrUnknownNFT(m.Denom, id)
	}

	m.IDs = append(m.IDs[:index], m.IDs[index+1:]...)

	return m, nil
}

// NewCollection creates a new NFT Collection
func NewCollection(denom string, nfts NFTs) Collection {
	return Collection{
		Denom: strings.TrimSpace(denom),
		NFTs:  nfts,
	}
}

// EmptyCollection returns an empty collection
func EmptyCollection() Collection {
	return NewCollection("", NewNFTs())
}

// GetNFT gets a NFT from the collection
func (collection Collection) GetNFT(id string) (nft exported.NFT, err error) {
	nft, found := collection.NFTs.Find(id)
	if found {
		return nft, nil
	}
	return nil, ErrUnknownNFT(collection.Denom, id)
}

// ContainsNFT returns whether or not a Collection contains an NFT
func (collection Collection) ContainsNFT(id string) bool {
	_, err := collection.GetNFT(id)
	return err == nil
}

// AddNFT adds an NFT to the collection
func (collection Collection) AddNFT(nft exported.NFT) (Collection, error) {
	id := nft.GetID()
	exists := collection.ContainsNFT(id)
	if exists {
		collNFT, err := collection.GetNFT(id)
		if err != nil {
			return collection, ErrUnknownNFT(collection.Denom, id)
		}

		ownerAddress := nft.GetOwners().GetOwners()[0].GetAddress()

		subTokenIDs := nft.GetOwners().GetOwners()[0].GetSubTokenIDs()
		owner := collNFT.GetOwners().GetOwner(ownerAddress)
		if owner == nil {
			owner = NewTokenOwner(ownerAddress, subTokenIDs)

			owners, err := collNFT.GetOwners().SetOwner(owner)
			if err != nil {
				return collection, err
			}

			collNFT = collNFT.SetOwners(owners)
		} else {
			for _, id := range subTokenIDs {
				owner = owner.SetSubTokenID(id)
			}

			owners, err := collNFT.GetOwners().SetOwner(owner)
			if err != nil {
				return collection, err
			}

			collNFT = collNFT.SetOwners(owners)
		}

		updatedNFTs, found := collection.NFTs.Update(id, collNFT)
		if !found {
			return collection, ErrUnknownNFT(collection.Denom, id)
		}

		collection.NFTs = updatedNFTs
	} else {
		collection.NFTs = collection.NFTs.Append(nft.(BaseNFT))
	}

	return collection, nil
}

// UpdateNFT updates an NFT from a collection
func (collection Collection) UpdateNFT(nft exported.NFT) (Collection, error) {
	nfts, ok := collection.NFTs.Update(nft.GetID(), nft)
	if !ok {
		return collection, ErrUnknownNFT(collection.Denom, nft.GetID())
	}

	collection.NFTs = nfts
	return collection, nil
}

// DeleteNFT deletes an NFT from a collection
func (collection Collection) DeleteNFT(nft exported.NFT) (Collection, error) {
	nfts, ok := collection.NFTs.Remove(nft.GetID())
	if !ok {
		return collection, ErrUnknownNFT(collection.Denom, nft.GetID())
	}
	collection.NFTs = nfts
	return collection, nil
}

// Supply gets the total supply of NFTs of a collection
func (collection Collection) Supply() int {
	return len(collection.NFTs)
}

// String follows stringer interface
func (collection Collection) String() string {
	return fmt.Sprintf(`Denom: 				%s
NFTs:

%s`,
		collection.Denom,
		collection.NFTs.String(),
	)
}

//------------------------------------------------------------

// Collections define an array of Collection
type Collections []Collection

// NewCollections creates a new set of NFTs
func NewCollections(collections ...Collection) Collections {
	if len(collections) == 0 {
		return Collections{}
	}
	return Collections(collections).Sort()
}

// Append appends two sets of Collections
func (collections Collections) Append(collectionsB ...Collection) Collections {
	return append(collections, collectionsB...).Sort()
}

// Find returns the searched collection from the set
func (collections Collections) Find(denom string) (Collection, bool) {
	index := collections.find(denom)
	if index == -1 {
		return Collection{}, false
	}
	return collections[index], true
}

// Remove removes a collection from the set of collections
func (collections Collections) Remove(denom string) (Collections, bool) {
	index := collections.find(denom)
	if index == -1 {
		return collections, false
	}
	collections[len(collections)-1], collections[index] = collections[index], collections[len(collections)-1]
	return collections[:len(collections)-1], true
}

// String follows stringer interface
func (collections Collections) String() string {
	if len(collections) == 0 {
		return ""
	}

	out := ""
	for _, collection := range collections {
		out += fmt.Sprintf("%v\n", collection.String())
	}
	return out[:len(out)-1]
}

// Empty returns true if there are no collections and false otherwise.
func (collections Collections) Empty() bool {
	return len(collections) == 0
}

func (collections Collections) find(denom string) (idx int) {
	return FindUtil(collections, denom)
}

//-----------------------------------------------------------------------------
// Sort & Findable interfaces

func (collections Collections) ElAtIndex(index int) string { return collections[index].Denom }
func (collections Collections) Len() int                   { return len(collections) }
func (collections Collections) Less(i, j int) bool {
	return strings.Compare(collections[i].Denom, collections[j].Denom) == -1
}
func (collections Collections) Swap(i, j int) {
	collections[i], collections[j] = collections[j], collections[i]
}

var _ sort.Interface = Collections{}

// Sort is a helper function to sort the set of coins inplace
func (collections Collections) Sort() Collections {
	sort.Sort(collections)
	return collections
}

func expNftToBaseNft(nft exported.NFT) BaseNFT { return *nft.(*BaseNFT) }

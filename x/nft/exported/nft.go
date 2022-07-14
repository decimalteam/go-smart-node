package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO abcd зачем эти интерфейсы? 1. Их никто не использует 2. В коде путаница - BaseNFT вперемешку с proto BaseNFT

// NFT non fungible token interface
type NFT interface {
	GetID() string
	GetOwners() TokenOwners
	SetOwners(owners TokenOwners) NFT
	GetCreator() string
	GetTokenURI() string
	EditMetadata(tokenURI string) NFT
	GetReserve() sdk.Int
	GetAllowMint() bool
	String() string
}

type TokenOwner interface {
	GetAddress() string
	GetSubTokenIDs() []int64
	SetSubTokenID(id int64) TokenOwner
	SortSubTokensFix() TokenOwner
	RemoveSubTokenID(id int64) TokenOwner
	String() string
}

type TokenOwners interface {
	GetOwners() []TokenOwner
	SetOwner(owner TokenOwner) (TokenOwners, error)
	GetOwner(address string) TokenOwner
	String() string
}

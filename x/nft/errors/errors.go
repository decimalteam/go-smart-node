package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "nft"

var (
	Internal                      = errors.New(codespace, 101, "internal error")
	InvalidCollection             = errors.New(codespace, 102, "invalid NFT collection")
	UnknownCollection             = errors.New(codespace, 103, "unknown NFT collection")
	InvalidNFT                    = errors.New(codespace, 104, "invalid NFT")
	UnknownNFT                    = errors.New(codespace, 105, "unknown NFT")
	InvalidQuantity               = errors.New(codespace, 106, "invalid NFT quantity")
	InvalidReserve                = errors.New(codespace, 107, "invalid NFT reserve")
	NotCreatorBurn                = errors.New(codespace, 108, "only the creator can burn token")
	NotCreatorMint                = errors.New(codespace, 109, "only the creator can mint a token")
	InvalidDenom                  = errors.New(codespace, 110, "invalid denom name")
	InvalidTokenID                = errors.New(codespace, 111, "Invalid token id")
	NotUniqueSubTokenIDs          = errors.New(codespace, 112, "not unique SubTokenIDs")
	NotUniqueTokenURI             = errors.New(codespace, 113, "not unique tokenURI")
	NotUniqueTokenID              = errors.New(codespace, 114, "not unique tokenID")
	OwnerDoesNotOwnSubTokenID     = errors.New(codespace, 115, "owner does not own sub tokenID")
	InvalidSender                 = errors.New(codespace, 116, "invalid sender address")
	InvalidRecipientAddress       = errors.New(codespace, 117, "invalid recipient address")
	ForbiddenToTransferToYourself = errors.New(codespace, 118, "forbidden to transfer to yourself")
	NotCreatorUpdateReserve       = errors.New(codespace, 119, "only the creator can update reserve a token")
	InsufficientFunds             = errors.New(codespace, 120, "insufficient funds are required")
	NotSetValueLowerNow           = errors.New(codespace, 121, "invalid new reserve")
	WrongReserveCoinDenom         = errors.New(codespace, 122, "wrong reserve coin denom")
	EmptyTokenURI                 = errors.New(codespace, 123, "empty tokenURI")
	SubTokenDoesNotExists         = errors.New(codespace, 124, "sub token with received ID does not exist")
	NotAllowedMint                = errors.New(codespace, 125, "minting is disabled for the token")
	NotCreatorUpdate              = errors.New(codespace, 126, "only the creator can update a token")
	SameTokenURI                  = errors.New(codespace, 127, "same token URI")
	// invariants

	NftSupply             = errors.New(codespace, 150, "NFTs count not equal to total nft supply")
	InvalidSubTokensLen   = errors.New(codespace, 151, "invalid sub tokens len for nft")
	UnknownSubTokenForNFT = errors.New(codespace, 152, "unknown sub token id for nft")
)

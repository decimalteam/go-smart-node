package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "nft"

var (
	Internal                      = errors.New(codespace, 1, "internal error")
	InvalidCollection             = errors.New(codespace, 2, "invalid NFT collection")
	UnknownCollection             = errors.New(codespace, 3, "unknown NFT collection")
	InvalidNFT                    = errors.New(codespace, 4, "invalid NFT")
	UnknownNFT                    = errors.New(codespace, 5, "unknown NFT")
	InvalidQuantity               = errors.New(codespace, 6, "invalid NFT quantity")
	InvalidReserve                = errors.New(codespace, 7, "invalid NFT reserve")
	NotAllowedBurn                = errors.New(codespace, 8, "only the creator can burn token")
	NotAllowedMint                = errors.New(codespace, 9, "only the creator can mint a token")
	InvalidDenom                  = errors.New(codespace, 10, "invalid denom name")
	InvalidTokenID                = errors.New(codespace, 11, "Invalid token id")
	NotUniqueSubTokenIDs          = errors.New(codespace, 12, "not unique SubTokenIDs")
	NotUniqueTokenURI             = errors.New(codespace, 13, "not unique tokenURI")
	NotUniqueTokenID              = errors.New(codespace, 14, "not unique tokenID")
	OwnerDoesNotOwnSubTokenID     = errors.New(codespace, 15, "owner does not own sub tokenID")
	InvalidSender                 = errors.New(codespace, 16, "invalid sender address")
	InvalidRecipientAddress       = errors.New(codespace, 17, "invalid recipient address")
	ForbiddenToTransferToYourself = errors.New(codespace, 18, "forbidden to transfer to yourself")
	NotAllowedUpdateReserve       = errors.New(codespace, 19, "only the creator can update reserve a token")
	InsufficientFunds             = errors.New(codespace, 20, "insufficient funds are required")
	NotSetValueLowerNow           = errors.New(codespace, 21, "invalid new reserve")
	WrongReserveCoinDenom         = errors.New(codespace, 22, "wrong reserve coin denom")
	EmptyTokenURI                 = errors.New(codespace, 23, "empty tokenURI")
	SubTokenDoesNotExists         = errors.New(codespace, 24, "sub token with received ID does not exist")
	// invariants

	NftSupply             = errors.New(codespace, 25, "NFTs count not equal to total nft supply")
	InvalidSubTokensLen   = errors.New(codespace, 26, "invalid sub tokens len for nft")
	UnknownSubTokenForNFT = errors.New(codespace, 27, "unknown sub token id for nft")
)

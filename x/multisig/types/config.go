package types

// MultisigTransactionIDPrefix is prefix for multisig transaction ID.
const MultisigTransactionIDPrefix = "d0mstx"

// Multisignature wallet limitations.
const (
	MinOwnerCount = 2
	MaxOwnerCount = 16
	MinWeight     = 1
	MaxWeight     = 1024
)

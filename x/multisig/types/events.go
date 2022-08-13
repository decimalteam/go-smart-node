package types

const (
	// Event attributes
	AttributeValueCategory = ModuleName

	// Common
	AttributeKeySender      = "sender"
	AttributeKeyWallet      = "wallet"
	AttributeKeyTransaction = "transaction"

	// CreateWallet
	AttributeKeyOwners    = "owners"
	AttributeKeyWeights   = "weights"
	AttributeKeyThreshold = "threshold"

	// CreateTransaction
	AttributeKeyReceiver = "receiver"
	AttributeKeyCoins    = "coins"

	// SignTransaction
	AttributeKeySignerWeight  = "signer_weight"
	AttributeKeyConfirmations = "confirmations"
	AttributeKeyConfirmed     = "confirmed"
)

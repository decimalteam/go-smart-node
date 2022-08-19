package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (rec *LegacyRecord) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(rec.Address, DecimalPrefix)
	if err != nil {
		return errors.InvalidLegacyBech32Address
	}
	// record must be not empty
	if len(rec.Coins) == 0 && len(rec.Nfts) == 0 && len(rec.Wallets) == 0 {
		return errors.NoInfoForLegacyAddress
	}

	for _, coin := range rec.Coins {
		if coin.Amount.IsZero() || coin.Amount.IsNegative() {
			return errors.OneOfLegacyAddresыСoinsBalanceIsNegativeOrZero
		}
	}
	// wallets addresses must be valid bech32 addresses
	for _, w := range rec.Wallets {
		_, err := sdk.GetFromBech32(w, DecimalPrefix)
		if err != nil {
			return errors.WalletAddressIsNotValidBech32
		}
	}
	return nil
}

/*
// ValidateBasic runs stateless checks on the message.
func (msg MsgReturnLegacyBalance) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}
	// Validate receiver
	if _, err := sdk.AccAddressFromBech32(msg.Receiver); err != nil {
		return ErrInvalidReceiverAddress(msg.Receiver)
	}
	// Validate public key
	if len(msg.PublicKeyBytes) != ethsecp256k1.PubKeySize {
		return ErrInvalidPublicKeyLength(len(msg.PublicKeyBytes))
	}
	// Validate receiver and public key
	address, err := bech32.ConvertAndEncode(config.Bech32Prefix, ethsecp256k1.PubKey{Key: msg.PublicKeyBytes}.Address())
	if err != nil {
		return ErrCannnotGetAddressFromPublicKey(err.Error())
	}
	if address != msg.Receiver {
		return ErrNoMatchReceiverAndPKey(msg.Receiver, address)
	}
	// Validate old address
	_, err = commonTypes.GetLegacyAddressFromPubKey(msg.PublicKeyBytes)
	if err != nil {
		return ErrCannnotGetAddressFromPublicKey(err.Error())
	}
	return nil
}
*/

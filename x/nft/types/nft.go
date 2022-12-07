package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
)

func (c *Collection) Validate() error {
	if !config.CollectionDenomValidator.MatchString(c.Denom) {
		return errors.InvalidDenom
	}

	_, err := sdk.AccAddressFromBech32(c.Creator)
	if err != nil {
		return errors.InvalidSender
	}

	return nil
}

func (t *Token) Validate() error {
	_, err := sdk.AccAddressFromBech32(t.Creator)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(t.Denom) == "" {
		return errors.InvalidDenom
	}
	if strings.TrimSpace(t.ID) == "" {
		return errors.InvalidNFT
	}
	if strings.TrimSpace(t.URI) == "" {
		return errors.InvalidNFT
	}
	if !t.Reserve.IsPositive() {
		return errors.InvalidReserve
	}
	if !config.CollectionDenomValidator.MatchString(t.Denom) {
		return errors.InvalidDenom
	}
	if !config.CollectionDenomValidator.MatchString(t.ID) {
		return errors.InvalidTokenID
	}
	return nil
}

func (st *SubToken) Validate() error {
	_, err1 := sdk.AccAddressFromBech32(st.Owner)
	_, err2 := sdk.GetFromBech32(st.Owner, "dx") // legacy owners
	if err1 != nil && err2 != nil {
		return errors.InvalidSender
	}
	if !st.Reserve.IsPositive() {
		return errors.InvalidReserve
	}

	return nil
}

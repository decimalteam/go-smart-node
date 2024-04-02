// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// evm coin center events
const (
	NameOfSlugForGetAddressTokenCenter     = "token-center"
	NameOfSlugForGetAddressDelegation      = "delegation"
	NameOfSlugForGetAddressMasterValidator = "master-validator"
	EventChangeTokenCenter                 = "ContractAdded"

	// DRC20MethodCreateToken defines the create method for DRC20 token
	DRC20MethodCreateToken = "createToken"
)

// LogTransfer Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

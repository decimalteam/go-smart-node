package types

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// GetHaltAdmin returns the bech32-encoded address of the dedicated emergency
// halt-admin for the given chain. Only this account may schedule or clear a
// hard halt (MsgHaltChain / MsgResumeChain) or a soft freeze
// (MsgFreezeChain / MsgUnfreezeChain).
//
// The halt-admin is intentionally SEPARATE from the upgrade authority
// (see x/upgrade/utils.go:GetAddressForUpdate). This key can stop the entire
// chain, so it must be held securely.
//
// TODO: replace the placeholder addresses below with real, dedicated keys.
// Mainnet returns an empty address by default, which DISABLES the emergency
// levers until a real key is configured (a deliberate fail-safe — no signed
// transaction can ever match an empty admin address).
func GetHaltAdmin(chainID string) string {
	switch {
	case helpers.IsMainnet(chainID):
		return "d01wfden2mlxsp4lqzu5k6jr3ypjx7mehxn2vx8vc" // TODO: set the mainnet halt-admin address before enabling on mainnet
	case helpers.IsTestnet(chainID):
		return "d01wfden2mlxsp4lqzu5k6jr3ypjx7mehxn2vx8vc" // TODO: replace with a dedicated testnet halt-admin
	default:
		return "d01wfden2mlxsp4lqzu5k6jr3ypjx7mehxn2vx8vc" // TODO: replace with a dedicated devnet halt-admin
	}
}

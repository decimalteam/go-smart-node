package config

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/decimalteam/ethermint/types"
)

const (
	// AppName defines the application name.
	AppName = "decimal"

	// AppBinName defines the application binary name.
	AppBinName = "dscd"
)

// Full list of registered chain IDs can be found here: https://chainlist.org/
const (
	// MainnetChainID defines EVM chain ID used for DSC mainnet network.
	MainnetChainID = 2020

	// TestnetChainID defines EVM chain ID used for DSC testnet network.
	TestnetChainID = 202020
)

const (
	// Bech32Prefix defines the Bech32 prefix used for EthAccounts
	Bech32Prefix = "d0"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

const (
	// BaseDenom defines to the default denomination used in Decimal (staking, EVM, governance, etc.)
	// TODO: Load it from
	BaseDenom = "del"
)

var (
	DataPath = fmt.Sprintf("%s/.decimal/daemon/data", os.Getenv("HOME"))
	// NOTE: THIS IS GLOBAL OBJECT TO STORE UPDATES HEIGHTS TO CALCULATE GRACE PERIODS
	// PATH TO FILE WILL BE REDEFINED IN app/app.go:NewDSC
	UpdatesInfo = &UpdatesInfoStruct{} // NewUpdatesInfo(filepath.Join(DataPath, UpdatesName))
	UpdatesName = "updates.json"
	GracePeriod = int64(600 * 24 * 1) // 600 - average blocks per hour, grace period to use inside inGracePeriod
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetPurpose(sdk.Purpose)
	config.SetCoinType(ethermint.Bip44CoinType)
}

// RegisterBaseDenom registers the base denomination to the SDK.
func RegisterBaseDenom() {
	if err := sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, ethermint.BaseDenomUnit)); err != nil {
		panic(err)
	}
}

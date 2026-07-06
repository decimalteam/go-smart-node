// Package redenom implements the one-time DEL redenomination: it divides every
// DEL-denominated amount in both Cosmos module state and EVM contract storage by
// a fixed divisor (1000), keeping the 18-decimal base unit. It is invoked from an
// in-place upgrade handler (app/upgrades.go) and from the offline `dscd
// redenom-dryrun` command, which both call Redenominate so production and the
// dry-run exercise identical logic.
//
// Scaling rule: an amount is scaled iff it is denominated in the base coin "del".
// Custom-coin supplies/balances are left untouched; a custom coin's DEL *reserve*
// is scaled (its bonding-curve price in DEL therefore drops 1000x, which is correct
// because one new DEL is worth 1000 old DEL).
//
// Integer division floors, so per-item amounts lose sub-divisor dust. To keep the
// (morally enforced) pool==sum-of-stakes and supply==sum-of-balances relations
// exact, module-account pool balances and the del supply are DERIVED as exact sums
// of the already-scaled authoritative records rather than scaled independently.
package redenom

import (
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	coinkeeper "bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	evmkeeper "github.com/decimalteam/ethermint/x/evm/keeper"

	feekeeper "bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	legacykeeper "bitbucket.org/decimalteam/go-smart-node/x/legacy/keeper"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
)

// DefaultDivisor is the redenomination factor (1000 old DEL -> 1 new DEL).
var DefaultDivisor = sdkmath.NewInt(1000)

// Keepers bundles every keeper Redenominate needs. Pointer keepers (coin, nft,
// legacy, evm) have pointer receivers; validator/gov/account/bank are values.
type Keepers struct {
	Bank      bankkeeper.Keeper
	Coin      *coinkeeper.Keeper
	Validator validatorkeeper.Keeper
	NFT       *nftkeeper.Keeper
	Legacy    *legacykeeper.Keeper
	Gov       govkeeper.Keeper
	Account   authkeeper.AccountKeeper
	EVM       *evmkeeper.Keeper
	Fee       *feekeeper.Keeper // x/fee (customfee): fiat oracle price, from which the EVM base fee is derived
}

// StoreKeys carries the raw store keys used for direct store manipulation (the
// bank balance/supply rewrite and the EVM storage rewrite go below the public
// keeper APIs, which expose no supply-consistent SetBalance in cosmos-sdk v0.46).
type StoreKeys struct {
	Bank store.StoreKey
	EVM  store.StoreKey
}

// Report summarizes a redenomination run for logging and dry-run assertions.
type Report struct {
	Divisor sdkmath.Int

	OldSupplyDel  sdkmath.Int
	NewSupplyDel  sdkmath.Int
	SupplyRemoved sdkmath.Int // OldSupplyDel - NewSupplyDel (~999/1000 of old: the intended ÷1000 reduction)
	RoundingDust  sdkmath.Int // floor(OldSupplyDel/divisor) - NewSupplyDel (true precision loss from per-item flooring)

	// Derived module-account del targets (exact sums of scaled records).
	BondedTarget       sdkmath.Int
	NotBondedTarget    sdkmath.Int
	CoinReserveTarget  sdkmath.Int // del held by the coin module (sum of scaled custom reserves)
	ReservedPoolTarget sdkmath.Int // del held by the nft reserved_pool (sum of scaled subtoken del reserves)

	// Counts.
	BankAccountsScaled  int
	DelegationsScaled   int
	UndelegationsScaled int
	HoldsScaled         int
	RedelegationsScaled int
	ValidatorsRepowered int
	CustomReservesScaled int
	NFTReservesScaled   int
	LegacyRecordsScaled int
	FeePricesScaled     int // base-denom fiat oracle prices multiplied by div (fee layer: EVM base fee + Cosmos tx fees)

	// Validators that were Bonded but whose new power floors to 0 (rollout gate).
	ValidatorsZeroed []string

	// EVM.
	EVMCoinStakeSlots   int
	EVMNFTReserveSlots  int
	EVMValidatorReserveSlots int
	EVMAutoUnbondSlots  int
	EVMFrozenSlots      int
	EVMCheckSlotsCleared int

	// Checks void+refund.
	ChecksVoided      int
	ChecksRefundedDel sdkmath.Int

	// WDEL (wrapped DEL) ledger scaling: staked DEL is custodied as WDEL, so its ERC20
	// ledger must be divided in lockstep to stay backed 1:1 by the (scaled) native del.
	WDELHoldersScaled  int
	WDELTotalSupplyOld sdkmath.Int
	WDELTotalSupplyNew sdkmath.Int

	// EVM NFT collection _reserve scaling (DecimalNFTBeaconProxy DRC721/DRC1155).
	EVMNFTCollectionsScaled int
	EVMNFTCollectionsFailed int         // rejected: floor(Σreserve/div) > scaled native (slot mis-id); left UNSCALED
	EVMNFTFailedDel         sdkmath.Int // total native del in the rejected collections
}

func newReport(div sdkmath.Int) Report {
	z := sdkmath.ZeroInt()
	return Report{
		Divisor: div, OldSupplyDel: z, NewSupplyDel: z, SupplyRemoved: z, RoundingDust: z,
		BondedTarget: z, NotBondedTarget: z, CoinReserveTarget: z, ReservedPoolTarget: z,
		ChecksRefundedDel: z, WDELTotalSupplyOld: z, WDELTotalSupplyNew: z, EVMNFTFailedDel: z,
	}
}

// divFloor returns floor(x / d). All DEL amounts are non-negative, and Int.Quo
// truncates toward zero, so this is a true floor here.
func divFloor(x, d sdkmath.Int) sdkmath.Int {
	if x.IsNil() {
		return sdkmath.ZeroInt()
	}
	return x.Quo(d)
}

// Redenominate divides every DEL amount in Cosmos and EVM state by divisor.
//
// Order is load-bearing (see scaleCosmos): authoritative records are scaled first,
// then module-account pool balances and the del supply are reconciled as exact
// sums, then validator power is recomputed from the scaled stakes + scaled custom
// reserves, then EVM contract storage is rewritten and outstanding DEL checks are
// voided and refunded.
func Redenominate(ctx sdk.Context, k Keepers, sk StoreKeys, divisor sdkmath.Int) (Report, error) {
	if divisor.IsNil() || !divisor.IsPositive() || divisor.LTE(sdkmath.OneInt()) {
		return Report{}, fmt.Errorf("redenom: divisor must be > 1, got %s", divisor)
	}
	rep := newReport(divisor)

	base := k.Validator.BaseDenom(ctx)
	if base == "" {
		return rep, fmt.Errorf("redenom: empty base denom")
	}

	rep.OldSupplyDel = k.Bank.GetSupply(ctx, base).Amount

	if err := scaleCosmos(ctx, k, sk, divisor, base, &rep); err != nil {
		return rep, fmt.Errorf("redenom: cosmos scaling: %w", err)
	}

	if err := scaleEVM(ctx, k, sk, divisor, base, &rep); err != nil {
		return rep, fmt.Errorf("redenom: evm scaling: %w", err)
	}

	rep.NewSupplyDel = k.Bank.GetSupply(ctx, base).Amount
	rep.SupplyRemoved = rep.OldSupplyDel.Sub(rep.NewSupplyDel)
	rep.RoundingDust = rep.OldSupplyDel.Quo(divisor).Sub(rep.NewSupplyDel)

	ctx.Logger().Info("redenomination complete",
		"divisor", divisor.String(),
		"oldSupply", rep.OldSupplyDel.String(),
		"newSupply", rep.NewSupplyDel.String(),
		"supplyRemoved", rep.SupplyRemoved.String(),
		"roundingDust", rep.RoundingDust.String(),
		"bankAccountsScaled", rep.BankAccountsScaled,
		"delegations", rep.DelegationsScaled,
		"stakeHoldsScaled", rep.HoldsScaled,
		"validatorsRepowered", rep.ValidatorsRepowered,
		"validatorsZeroed", len(rep.ValidatorsZeroed),
		"customReserves", rep.CustomReservesScaled,
		"nftReserves", rep.NFTReservesScaled,
		"feePricesScaled", rep.FeePricesScaled,
		"evmCoinStakeSlots", rep.EVMCoinStakeSlots,
		"evmNftReserveSlots", rep.EVMNFTReserveSlots,
		"evmValidatorReserveSlots", rep.EVMValidatorReserveSlots,
		"checksVoided", rep.ChecksVoided,
		"checksRefundedDel", rep.ChecksRefundedDel.String(),
	)
	if len(rep.ValidatorsZeroed) > 0 {
		sort.Strings(rep.ValidatorsZeroed)
		ctx.Logger().Error("redenomination: bonded validators floored to zero power",
			"validators", rep.ValidatorsZeroed)
	}
	return rep, nil
}

package app

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// coinHoldStakeID converts bech32 validator and delegator addresses to their EVM
// common.Address representations and computes the contract stakeId for a coin hold.
// It is exported (lower-case package-level) so tests can call it directly and lock
// the validator/delegator argument ordering against future regressions.
func coinHoldStakeID(validatorBech32, delegatorBech32 string, token common.Address, endTime int64) (common.Hash, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorBech32)
	if err != nil {
		return common.Hash{}, fmt.Errorf("validator address %s: %w", validatorBech32, err)
	}
	delAddr, err := sdk.AccAddressFromBech32(delegatorBech32)
	if err != nil {
		return common.Hash{}, fmt.Errorf("delegator address %s: %w", delegatorBech32, err)
	}
	return computeStakeID(
		common.BytesToAddress(valAddr.Bytes()),
		common.BytesToAddress(delAddr.Bytes()),
		token,
		big.NewInt(0),
		big.NewInt(endTime),
	), nil
}

// delegationStorageBase is the ERC-7201 base slot of DelegationStorage
// (DECIMAL_DELEGATION_COMMON_STORAGE_LOCATION in the Solidity contracts).
var delegationStorageBase = func() *big.Int {
	v, _ := new(big.Int).SetString("c1dae510251b57b62087f142cc8746564a134308f8810580b66b08a29b6def00", 16)
	return v
}()

// stakesMappingSlot is the slot of the `_stakes` mapping (DelegationStorage field index 1).
var stakesMappingSlot = new(big.Int).Add(delegationStorageBase, big.NewInt(1))

// holdStartTimeFieldOffset is the slot offset of `holdStartTime` within the Stake struct
// (validator0, delegator1, token2, amount3, tokenId4, tokenType5, holdTimestamp6, holdStartTime7).
const holdStartTimeFieldOffset = 7

// computeStakeID reproduces the contract's
// keccak256(abi.encodePacked(validator, delegator, token, tokenId, holdTimestamp)).
func computeStakeID(validator, delegator, token common.Address, tokenId, holdTimestamp *big.Int) common.Hash {
	packed := make([]byte, 0, 124)
	packed = append(packed, validator.Bytes()...)
	packed = append(packed, delegator.Bytes()...)
	packed = append(packed, token.Bytes()...)
	packed = append(packed, common.LeftPadBytes(tokenId.Bytes(), 32)...)
	packed = append(packed, common.LeftPadBytes(holdTimestamp.Bytes(), 32)...)
	return common.BytesToHash(crypto.Keccak256(packed))
}

// holdStartTimeSlot returns the storage slot of `_stakes[stakeID].holdStartTime`.
func holdStartTimeSlot(stakeID common.Hash) common.Hash {
	start := new(big.Int).SetBytes(crypto.Keccak256(
		append(stakeID.Bytes(), common.BigToHash(stakesMappingSlot).Bytes()...),
	))
	return common.BigToHash(new(big.Int).Add(start, big.NewInt(holdStartTimeFieldOffset)))
}

// selectHoldStart implements Decision B/C: among holds sharing one contract stakeId
// (same HoldEndTime), use the earliest start; replace any 0 with the floor (upgrade
// block time) so no hold reads as >=1yr old by accident.
func selectHoldStart(starts []int64, floor int64) int64 {
	best := int64(0)
	for _, s := range starts {
		if s == 0 {
			s = floor
		}
		if best == 0 || s < best {
			best = s
		}
	}
	if best == 0 {
		return floor
	}
	return best
}

// BackfillHoldStartTimes reads every native hold's start time and writes it into the
// delegation contract's _stakes[stakeId].holdStartTime slot (Decision: sync node -> contract).
func BackfillHoldStartTimes(ctx sdk.Context, app *DSC) (int, error) {
	floor := ctx.BlockTime().Unix()

	// denom -> token EVM address (one KV scan), with the base-denom -> WDEL overlay.
	denomToToken := map[string]common.Address{}
	app.CoinKeeper.IterateCoinDRC(ctx, func(denom, drc20 string) bool {
		if drc20 != "" {
			denomToToken[denom] = common.HexToAddress(drc20)
		}
		return false
	})
	if wdelHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressWDEL); err == nil {
		if a := common.HexToAddress(wdelHex); a != (common.Address{}) {
			denomToToken[app.CoinKeeper.GetBaseDenom(ctx)] = a
		}
	}

	delegationHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return 0, fmt.Errorf("resolve delegation contract: %w", err)
	}
	delegationAddr := common.HexToAddress(delegationHex)
	if delegationAddr == (common.Address{}) {
		return 0, fmt.Errorf("delegation contract address not resolved")
	}

	// Group hold starts per stakeId so multi-hold same-end buckets collapse to the earliest start.
	type key struct {
		val, del string
		token    common.Address
		endTime  int64
	}
	startsByKey := map[key][]int64{}
	app.ValidatorKeeper.IterateAllDelegations(ctx, func(d validatortypes.Delegation) bool {
		denom := d.Stake.Stake.Denom
		token, ok := denomToToken[denom]
		if !ok {
			return false // coin never bridged to EVM -> no contract stake
		}
		for _, h := range d.Stake.GetHolds() {
			if h.HoldEndTime == 0 {
				continue
			}
			k := key{d.Validator, d.Delegator, token, h.HoldEndTime}
			startsByKey[k] = append(startsByKey[k], h.HoldStartTime)
		}
		return false
	})

	written := 0
	for k, starts := range startsByKey {
		stakeID, err := coinHoldStakeID(k.val, k.del, k.token, k.endTime)
		if err != nil {
			return written, fmt.Errorf("backfill stakeId: %w", err)
		}
		start := selectHoldStart(starts, floor)
		slot := holdStartTimeSlot(stakeID)
		val := common.BigToHash(big.NewInt(start))

		app.EvmKeeper.SetState(ctx, delegationAddr, slot, val.Bytes())
		if got := app.EvmKeeper.GetState(ctx, delegationAddr, slot); got != val {
			return written, fmt.Errorf("slot read-back mismatch for stake %s: got %s want %s", stakeID.Hex(), got.Hex(), val.Hex())
		}
		written++
	}
	return written, nil
}

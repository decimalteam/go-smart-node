package redenom

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// scaleCosmos divides every "del" amount in Cosmos module state by div. Order is
// load-bearing:
//  1. scale the authoritative validator records (delegations/undelegations/redelegations);
//  2. rewrite all "del" bank balances (floor) and set supply = sum of new balances;
//  3. scale coin reserves/volumes (base "del" volume is set to the new bank supply
//     so the two stay exactly equal) — needed before the power recompute, which
//     prices custom-coin stakes off coin.Reserve;
//  4. scale NFT and legacy reserves;
//  5. recompute every validator's consensus power from the scaled stakes + scaled
//     custom-coin prices (mirrors PayRewards) and scale accrued rewards; the
//     validator EndBlocker's BlockValidatorUpdates emits the diff to Tendermint, so
//     this function does NOT call ApplyAndReturnValidatorSetUpdates itself;
//  6. scale DEL-denominated threshold params (coin BaseVolume, nft MinReserveAmount,
//     gov MinDeposit) so minimums stay economically equivalent.
func scaleCosmos(ctx sdk.Context, k Keepers, sk StoreKeys, div sdkmath.Int, base string, rep *Report) error {
	scaleValidatorStakes(ctx, k, div, base, rep)

	newSupply := scaleBankDel(ctx, k, sk, div, base, rep)

	if err := scaleCoins(ctx, k, div, base, newSupply, rep); err != nil {
		return err
	}
	if err := scaleNFTReserves(ctx, k, div, base, rep); err != nil {
		return err
	}
	scaleLegacy(ctx, k, div, base, rep)

	if err := repowerValidators(ctx, k, div, base, rep); err != nil {
		return err
	}
	if err := scaleThresholds(ctx, k, div, base); err != nil {
		return err
	}
	return nil
}

// scaleValidatorStakes scales the "del" amount of every delegation, undelegation
// entry and redelegation entry. Custom-coin stakes are left untouched. NFT stakes
// whose reserve coin is "del" are scaled too (their Stake.Stake holds the DEL reserve).
func scaleValidatorStakes(ctx sdk.Context, k Keepers, div sdkmath.Int, base string, rep *Report) {
	for _, del := range k.Validator.GetAllDelegations(ctx) {
		if del.Stake.Stake.Denom == base {
			del.Stake.Stake.Amount = divFloor(del.Stake.Stake.Amount, div)
			k.Validator.SetDelegation(ctx, del)
			rep.DelegationsScaled++
		}
	}

	var ubds []validatortypes.Undelegation
	k.Validator.IterateUndelegations(ctx, func(_ int64, ubd validatortypes.Undelegation) bool {
		ubds = append(ubds, ubd)
		return false
	})
	for _, ubd := range ubds {
		changed := false
		for i := range ubd.Entries {
			if ubd.Entries[i].Stake.Stake.Denom == base {
				ubd.Entries[i].Stake.Stake.Amount = divFloor(ubd.Entries[i].Stake.Stake.Amount, div)
				changed = true
			}
		}
		if changed {
			k.Validator.SetUndelegation(ctx, ubd)
			rep.UndelegationsScaled++
		}
	}

	var reds []validatortypes.Redelegation
	k.Validator.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		reds = append(reds, red)
		return false
	})
	for _, red := range reds {
		changed := false
		for i := range red.Entries {
			if red.Entries[i].Stake.Stake.Denom == base {
				red.Entries[i].Stake.Stake.Amount = divFloor(red.Entries[i].Stake.Stake.Amount, div)
				changed = true
			}
		}
		if changed {
			k.Validator.SetRedelegation(ctx, red)
			rep.RedelegationsScaled++
		}
	}
}

// scaleBankDel rewrites every account's "del" balance to floor(old/div) and sets the
// "del" supply to the exact sum of the new balances. It writes directly to the bank
// store because cosmos-sdk v0.46 exposes no supply-consistent public SetBalance; the
// key encoding here is byte-identical to x/bank/keeper setBalance/setSupply. Module
// pool accounts are floored like everyone else: each pool holds a sum of per-item
// amounts and floor(sum) >= sum(floor), so a pool can never end up under-funding the
// (also floored) records it backs. Returns the new total "del" supply.
func scaleBankDel(ctx sdk.Context, k Keepers, sk StoreKeys, div sdkmath.Int, base string, rep *Report) sdkmath.Int {
	type bal struct {
		addr sdk.AccAddress
		amt  sdkmath.Int
	}
	var dels []bal
	k.Bank.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == base {
			dels = append(dels, bal{addr: addr, amt: coin.Amount})
		}
		return false
	})

	bankStore := ctx.KVStore(sk.Bank)
	newSupply := sdkmath.ZeroInt()
	for _, b := range dels {
		newAmt := divFloor(b.amt, div)
		setBankBalanceDel(bankStore, b.addr, base, newAmt)
		newSupply = newSupply.Add(newAmt)
		rep.BankAccountsScaled++
	}
	setBankSupply(bankStore, base, newSupply)
	return newSupply
}

// scaleCoins scales coin volumes/reserves. The base "del" coin's volume is set to the
// exact new bank supply (keeping coin-module supply == bank supply); its reserve, if
// any, is floored. Each custom coin keeps its volume/limits and has only its
// DEL-denominated reserve floored.
func scaleCoins(ctx sdk.Context, k Keepers, div sdkmath.Int, base string, newSupply sdkmath.Int, rep *Report) error {
	for _, c := range k.Coin.GetCoins(ctx) {
		if c.Denom == base {
			if err := k.Coin.UpdateCoinVR(ctx, c.Denom, newSupply, divFloor(c.Reserve, div)); err != nil {
				return err
			}
			continue
		}
		if err := k.Coin.UpdateCoinVR(ctx, c.Denom, c.Volume, divFloor(c.Reserve, div)); err != nil {
			return err
		}
		rep.CustomReservesScaled++
	}
	return nil
}

// scaleNFTReserves scales DEL-denominated token and sub-token reserves. Reserves
// denominated in a custom coin are left untouched.
func scaleNFTReserves(ctx sdk.Context, k Keepers, div sdkmath.Int, base string, rep *Report) error {
	for _, col := range k.NFT.GetCollections(ctx) {
		creator, err := sdk.AccAddressFromBech32(col.Creator)
		if err != nil {
			return err
		}
		for _, token := range k.NFT.GetTokens(ctx, creator, col.Denom) {
			if token.Reserve.Denom == base {
				token.Reserve.Amount = divFloor(token.Reserve.Amount, div)
				k.NFT.SetTokenReserveForMigration(ctx, token.ID, token.Reserve)
			}
			for _, st := range k.NFT.GetSubTokens(ctx, token.ID) {
				if st.Reserve != nil && st.Reserve.Denom == base {
					st.Reserve.Amount = divFloor(st.Reserve.Amount, div)
					k.NFT.SetSubToken(ctx, token.ID, st)
					rep.NFTReservesScaled++
				}
			}
		}
	}
	return nil
}

// scaleLegacy scales DEL coins held in legacy records (their bank backing lives in
// the LegacyCoinPool module account, covered by scaleBankDel).
func scaleLegacy(ctx sdk.Context, k Keepers, div sdkmath.Int, base string, rep *Report) {
	for _, rec := range k.Legacy.GetLegacyRecords(ctx) {
		changed := false
		for i := range rec.Coins {
			if rec.Coins[i].Denom == base {
				rec.Coins[i].Amount = divFloor(rec.Coins[i].Amount, div)
				changed = true
			}
		}
		if changed {
			k.Legacy.SetLegacyRecord(ctx, rec)
			rep.LegacyRecordsScaled++
		}
	}
}

// repowerValidators recomputes each validator's consensus power from the already
// scaled stakes and scaled custom-coin prices, and scales accrued rewards. It mirrors
// the recompute in PayRewards (x/validator/keeper/reward.go) but pays nothing out and
// does NOT apply the validator-set updates — the EndBlocker's BlockValidatorUpdates
// does that, emitting the power diff to Tendermint in the same block.
func repowerValidators(ctx sdk.Context, k Keepers, div sdkmath.Int, base string, rep *Report) error {
	ccs := k.Validator.GetAllCustomCoinsStaked(ctx)
	prices := k.Validator.CalculateCustomCoinPrices(ctx, ccs)
	delsByVal := k.Validator.GetAllDelegationsByValidator(ctx)

	for _, val := range k.Validator.GetAllValidators(ctx) {
		op := val.GetOperator()

		rs, err := k.Validator.GetValidatorRS(ctx, op)
		if err != nil {
			// validator without a rewards record — nothing authoritative to update
			continue
		}
		rs.Rewards = divFloor(rs.Rewards, div)
		rs.TotalRewards = divFloor(rs.TotalRewards, div)

		dels := delsByVal[op.String()]
		total, err := k.Validator.CalculateTotalPowerWithDelegationsAndPrices(ctx, op, validatortypes.Delegations(dels), prices)
		if err != nil {
			return err
		}
		newPower := validatorkeeper.TokensToConsensusPower(total)
		wasBonded := val.Status == validatortypes.BondStatus_Bonded

		// Reconcile the consensus power index for EVERY validator currently in it, not
		// just bonded ones. The index key encodes the validator's power
		// (PotentialConsensusPower == RS.Stake, which GetAllValidators overlaid into
		// val.Stake), so HasValidatorByPowerIndex(val)/DeleteValidatorByPowerIndex(val)
		// match the live entry keyed by the OLD (pre-scale) power.
		//
		// SetOnline indexes an Online validator while its Status is still Unbonded (it
		// only becomes Bonded at the next EndBlocker), so NON-bonded validators can be in
		// the index. The previous version re-keyed only bonded validators and forced
		// RS.Stake=0 for the rest — which left such an Online+Unbonded "candidate" with a
		// stale ~div×-too-large index key AND RS.Stake=0, i.e. {Unbonded, Online, Stake==0}.
		// That state matches no case in the next-block EndBlocker's transition switch
		// (val_state_change.go) and hits its default panic, recovered into a truncated,
		// inconsistent validator set. Re-keying every indexed validator at its scaled power
		// keeps the invariant index-key-power == RS.Stake and the switch well-defined.
		inIndex := k.Validator.HasValidatorByPowerIndex(ctx, val)
		if inIndex {
			k.Validator.DeleteValidatorByPowerIndex(ctx, val)
		}

		// Offline/idle validators (not bonded and not in the index) keep stake 0 so they
		// stay in the EndBlocker-valid {Unbonded, !Online, Stake==0} state, matching the
		// per-block recompute in PayRewards (reward.go:316-318).
		if !wasBonded && !inIndex {
			newPower = 0
		}
		rs.Stake = newPower
		k.Validator.SetValidatorRS(ctx, op, rs)

		if inIndex {
			// Re-add to the index (at the scaled power) ONLY for validators that
			// legitimately belong there — a bonded validator or an online candidate.
			// The pre-scale entry was already deleted above; anything else stays
			// de-indexed, which is safe because the EndBlocker switch only walks indexed
			// validators. The resulting (status, online, stake) must be a state that
			// switch can classify, never {Unbonded, Online, Stake==0}.
			switch {
			case wasBonded:
				// Bonded: stays indexed at its scaled power. Power 0 is the valid
				// {Bonded, Stake==0} state (EndBlocker unbonds it next block); flag it for
				// the rollout gate.
				val.Stake = newPower
				k.Validator.SetValidatorByPowerIndex(ctx, val)
				if newPower == 0 {
					rep.ValidatorsZeroed = append(rep.ValidatorsZeroed, val.OperatorAddress)
				}
			case val.Online && newPower > 0:
				// Online+Unbonded candidate with positive scaled power: stays a valid
				// candidate ({Unbonded, Online, Stake>0}), promoted normally next block.
				val.Stake = newPower
				k.Validator.SetValidatorByPowerIndex(ctx, val)
			default:
				// Online+Unbonded candidate whose scaled power floored to 0 (the
				// would-be {Unbonded, Online, Stake==0} default-panic state), or any
				// anomalous indexed non-bonded validator. Keep it OUT of the index; take an
				// online one offline so it settles into the valid idle state
				// {Unbonded, !Online, Stake==0}, and can re-online once it again has power.
				if val.Online {
					val.Online = false
					val.Stake = 0
					k.Validator.SetValidator(ctx, val)
					ctx.Logger().Error("redenomination: online candidate validator floored to 0 power — taken offline",
						"validator", val.OperatorAddress)
				}
			}
		}
		rep.ValidatorsRepowered++
	}
	return nil
}

// scaleThresholds scales DEL-denominated minimum params so they stay economically
// equivalent after redenomination. (The coin MinCoinReserve config constant is a
// build-time value scaled in the upgrade binary, not here.)
func scaleThresholds(ctx sdk.Context, k Keepers, div sdkmath.Int, base string) error {
	cp := k.Coin.GetParams(ctx)
	cp.BaseVolume = divFloor(cp.BaseVolume, div)
	k.Coin.SetParams(ctx, cp)

	np := k.NFT.GetParams(ctx)
	np.MinReserveAmount = divFloor(np.MinReserveAmount, div)
	k.NFT.SetParams(ctx, np)

	dp := k.Gov.GetDepositParams(ctx)
	for i := range dp.MinDeposit {
		if dp.MinDeposit[i].Denom == base {
			dp.MinDeposit[i].Amount = divFloor(dp.MinDeposit[i].Amount, div)
		}
	}
	k.Gov.SetDepositParams(ctx, dp)
	return nil
}

// setBankBalanceDel writes addr's balance of denom directly to the bank store,
// replicating x/bank/keeper (BaseSendKeeper).setBalance byte-for-byte: the account
// balance store plus the denom->address reverse index, deleting both on a zero
// balance (bank invariants prohibit persisted zero balances).
func setBankBalanceDel(store sdk.KVStore, addr sdk.AccAddress, denom string, amt sdkmath.Int) {
	accountStore := prefix.NewStore(store, banktypes.CreateAccountBalancesPrefix(addr))
	denomPrefixStore := prefix.NewStore(store, banktypes.CreateDenomAddressPrefix(denom))

	if amt.IsZero() {
		accountStore.Delete([]byte(denom))
		denomPrefixStore.Delete(address.MustLengthPrefix(addr))
		return
	}
	b, err := amt.Marshal()
	if err != nil {
		panic(err)
	}
	accountStore.Set([]byte(denom), b)
	denomAddrKey := address.MustLengthPrefix(addr)
	if !denomPrefixStore.Has(denomAddrKey) {
		denomPrefixStore.Set(denomAddrKey, []byte{0})
	}
}

// setBankSupply writes the total supply of denom directly to the bank store,
// replicating x/bank/keeper (BaseKeeper).setSupply (zero supply is deleted).
func setBankSupply(store sdk.KVStore, denom string, amt sdkmath.Int) {
	supplyStore := prefix.NewStore(store, banktypes.SupplyKey)
	if amt.IsZero() {
		supplyStore.Delete([]byte(denom))
		return
	}
	b, err := amt.Marshal()
	if err != nil {
		panic(err)
	}
	supplyStore.Set([]byte(denom), b)
}

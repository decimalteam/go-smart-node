package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	types "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
// of it, updating unbonding delegations & redelegations appropriately
//
// CONTRACT:
//
//	slashFactor is non-negative
//
// CONTRACT:
//
//	Infraction was committed equal to or less than an unbonding period in the past,
//	so all unbonding delegations and redelegations from that height are stored
//
// CONTRACT:
//
//	Slash will not slash unbonded validators (for the above reason)
//
// CONTRACT:
//
//	Infraction was committed at the current height or at a past height,
//	not at a height in the future
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) sdkmath.Int {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(
			"WARNING: ignored attempt to slash a nonexistent validator; we recommend you investigate immediately",
			"validator", consAddr.String(),
		)
		return sdk.NewInt(0)
	}

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()

	// call the before-modification hook
	k.BeforeValidatorModified(ctx, operatorAddress)

	switch {
	case infractionHeight > ctx.BlockHeight():
		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():
		// Special-case slash at current height for efficiency - we don't need to
		// look through unbonding delegations or redelegations.
		logger.Info(
			"slashing at current height; not scanning unbonding delegations & redelegations",
			"height", infractionHeight,
		)

	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding delegations from slashed validator
		undelegations := k.GetUndelegationsFromValidator(ctx, operatorAddress)
		for _, undelegation := range undelegations {
			amountSlashed := k.SlashUndelegation(ctx, undelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}
		}

		// Iterate through redelegations from slashed source validator
		redelegations := k.GetRedelegationsFromSrcValidator(ctx, operatorAddress)
		for _, redelegation := range redelegations {
			amountSlashed := k.SlashRedelegation(ctx, validator, redelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}
		}
	}

	// slash delegations
	delegations := k.GetValidatorDelegations(ctx, operatorAddress)
	for _, delegation := range delegations {
		switch delegation.Stake.Type {
		case types.StakeType_Coin:
			calcSlashCoinStake(delegation.Stake, slashFactor)
		case types.StakeType_NFT:
			token, found := k.nftKeeper.GetToken(ctx, delegation.Stake.ID)
			if !found {
				panic(fmt.Errorf("can't find token %s in slashing time", delegation.Stake.ID))
			}
			var subtokens []nfttypes.SubToken
			for _, id := range delegation.Stake.SubTokenIDs {
				sub, found := k.nftKeeper.GetSubToken(ctx, token.ID, id)
				if !found {
					panic(fmt.Errorf("can't find subtoken %d for token %s in slashing time", id, delegation.Stake.ID))
				}
				if sub.Reserve == nil {
					sub.Reserve = &token.Reserve
				}
				subtokens = append(subtokens, sub)
			}
			calcSlashNFTStake(delegation.Stake, subtokens, slashFactor)
		}
	}

	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(remainingSlashAmount, validator.Tokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.

	// we need to calculate the *effective* slash fraction for distribution
	if validator.Tokens.IsPositive() {
		effectiveFraction := sdk.NewDecFromInt(tokensToBurn).QuoRoundUp(sdk.NewDecFromInt(validator.Tokens))
		// possible if power has changed
		if effectiveFraction.GT(sdk.OneDec()) {
			effectiveFraction = sdk.OneDec()
		}
		// call the before-slashed hook
		k.BeforeValidatorSlashed(ctx, operatorAddress, effectiveFraction)
	}

	// Deduct from validator's bonded tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	validator = k.RemoveValidatorTokens(ctx, validator, tokensToBurn)

	switch validator.GetStatus() {
	case types.Bonded:
		if err := k.burnBondedTokens(ctx, tokensToBurn); err != nil {
			panic(err)
		}
	case types.Unbonding, types.Unbonded:
		if err := k.burnNotBondedTokens(ctx, tokensToBurn); err != nil {
			panic(err)
		}
	default:
		panic("invalid validator status")
	}

	logger.Info(
		"validator slashed by slash factor",
		"validator", validator.GetOperator().String(),
		"slash_factor", slashFactor.String(),
		"burned", tokensToBurn,
	)
	return tokensToBurn
}

// jail a validator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.jailValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info("validator jailed", "validator", consAddr)
}

// unjail a validator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.unjailValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info("validator un-jailed", "validator", consAddr)
}

// slash an unbonding delegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
func (k Keeper) SlashUndelegation(ctx sdk.Context, undelegation types.Undelegation,
	infractionHeight int64, slashFactor sdk.Dec,
) (totalSlashAmount sdkmath.Int) {
	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	burnedAmount := sdk.ZeroInt()

	// perform slashing on all entries within the unbonding delegation
	for i, entry := range undelegation.Entries {
		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding delegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Don't slash more tokens than held
		// Possible since the unbonding delegation may already
		// have been slashed, and slash amounts are calculated
		// according to stake held at time of infraction
		unbondingSlashAmount := sdk.MinInt(slashAmount, entry.Balance)

		// Update unbonding delegation if necessary
		if unbondingSlashAmount.IsZero() {
			continue
		}

		burnedAmount = burnedAmount.Add(unbondingSlashAmount)
		entry.Balance = entry.Balance.Sub(unbondingSlashAmount)
		undelegation.Entries[i] = entry
		k.SetUndelegation(ctx, undelegation)
	}

	if err := k.burnNotBondedTokens(ctx, burnedAmount); err != nil {
		panic(err)
	}

	return totalSlashAmount
}

// slash a redelegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
// NOTE this is only slashing for prior infractions from the source validator
func (k Keeper) SlashRedelegation(ctx sdk.Context, srcValidator types.Validator, redelegation types.Redelegation,
	infractionHeight int64, slashFactor sdk.Dec,
) (totalSlashAmount sdkmath.Int) {
	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	bondedBurnedAmount, notBondedBurnedAmount := sdk.ZeroInt(), sdk.ZeroInt()

	// perform slashing on all entries within the redelegation
	for _, entry := range redelegation.Entries {
		// If redelegation started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Redelegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Unbond from target validator
		sharesToUnbond := slashFactor.Mul(entry.SharesDst)
		if sharesToUnbond.IsZero() {
			continue
		}

		valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress := sdk.MustAccAddressFromBech32(redelegation.DelegatorAddress)

		delegation, found := k.GetDelegation(ctx, delegatorAddress, valDstAddr)
		if !found {
			// If deleted, delegation has zero shares, and we can't unbond any more
			continue
		}

		if sharesToUnbond.GT(delegation.Shares) {
			sharesToUnbond = delegation.Shares
		}

		tokensToBurn, err := k.Unbond(ctx, delegatorAddress, valDstAddr, sharesToUnbond)
		if err != nil {
			panic(fmt.Errorf("error unbonding delegator: %v", err))
		}

		dstValidator, found := k.GetValidator(ctx, valDstAddr)
		if !found {
			panic("destination validator not found")
		}

		// tokens of a redelegation currently live in the destination validator
		// therefor we must burn tokens from the destination-validator's bonding status
		switch {
		case dstValidator.IsBonded():
			bondedBurnedAmount = bondedBurnedAmount.Add(tokensToBurn)
		case dstValidator.IsUnbonded() || dstValidator.IsUnbonding():
			notBondedBurnedAmount = notBondedBurnedAmount.Add(tokensToBurn)
		default:
			panic("unknown validator status")
		}
	}

	if err := k.burnBondedTokens(ctx, bondedBurnedAmount); err != nil {
		panic(err)
	}

	if err := k.burnNotBondedTokens(ctx, notBondedBurnedAmount); err != nil {
		panic(err)
	}

	return totalSlashAmount
}

// slash delegations records, make changes in nft, burn coins from pools
func (k Keeper) SlashDelegations(ctx sdk.Context, operatorAddress sdk.ValAddress, slashFactor sdk.Dec) {
	var coinsToBurn, nftCoinsToBurn sdk.Coins
	delegations := k.GetValidatorDelegations(ctx, operatorAddress)
	var slashEvents = make(map[string]types.DelegatorSlash)
	for _, delegation := range delegations {
		newStake, slashCoin, slashNFT, nftChanges := k.calcSlashStake(ctx, delegation.Stake, slashFactor)
		// accumulate coins to burn
		coinsToBurn = coinsToBurn.Add(slashCoin.Slash)
		for _, sub := range slashNFT.SubTokens {
			nftCoinsToBurn = nftCoinsToBurn.Add(sub.Slash)
		}
		// make changes in stake and nft
		k.SetDelegation(ctx, types.Delegation{
			Delegator: delegation.Delegator,
			Validator: delegation.Validator,
			Stake:     newStake,
		})
		for _, sub := range nftChanges.subtokens {
			k.nftKeeper.SetSubToken(ctx, nftChanges.tokenID, sub)
		}
		// accumulate events
		ev, ok := slashEvents[delegation.Delegator]
		if !ok {
			ev = types.DelegatorSlash{
				Delegator: delegation.Delegator,
			}
		}
		if slashCoin.Slash.Denom != "" {
			ev.Coins = append(ev.Coins, slashCoin)
		}
		if slashNFT.ID != "" {
			ev.NFTs = append(ev.NFTs, slashNFT)
		}
		slashEvents[delegation.Delegator] = ev
	}
	// TODO: burn coins from pools
	if !coinsToBurn.IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, types.BondedPoolName, coinsToBurn)
	}
	if !nftCoinsToBurn.IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, nfttypes.ReservedPool, coinsToBurn)
	}
	// ...
	// emit result event
	resultEvent := &types.EventSlash{
		Validator: operatorAddress.String(),
	}
	for _, ev := range slashEvents {
		resultEvent.Delegators = append(resultEvent.Delegators, ev)
	}
	events.EmitTypedEvent(ctx, resultEvent)
}

type nftSubtokenChanges struct {
	tokenID   string
	subtokens []nfttypes.SubToken
}

// calcSlashStake prepare changes. Only calculate, no changes
func (k Keeper) calcSlashStake(ctx sdk.Context, stake types.Stake, slashFactor sdk.Dec) (types.Stake,
	types.SlashCoin, types.SlashNFT, nftSubtokenChanges) {
	var newStake types.Stake
	var slashCoin types.SlashCoin
	var slashNFT types.SlashNFT
	var newSubTokens []nfttypes.SubToken
	var nftChanges nftSubtokenChanges
	var err error
	switch stake.Type {
	case types.StakeType_Coin:
		newStake, slashCoin, err = calcSlashCoinStake(stake, slashFactor)
		if err != nil {
			panic(err)
		}
	case types.StakeType_NFT:
		token, found := k.nftKeeper.GetToken(ctx, stake.ID)
		if !found {
			panic(fmt.Errorf("can't find token %s in slashing time", stake.ID))
		}
		var subtokens []nfttypes.SubToken
		for _, id := range stake.SubTokenIDs {
			sub, found := k.nftKeeper.GetSubToken(ctx, token.ID, id)
			if !found {
				panic(fmt.Errorf("can't find subtoken %d for token %s in slashing time", id, stake.ID))
			}
			if sub.Reserve == nil {
				sub.Reserve = &token.Reserve
			}
			subtokens = append(subtokens, sub)
		}
		newStake, slashNFT, newSubTokens, err = calcSlashNFTStake(stake, subtokens, slashFactor)
		if err != nil {
			panic(err)
		}
		nftChanges = nftSubtokenChanges{tokenID: token.ID, subtokens: newSubTokens}
	}
	return newStake, slashCoin, slashNFT, nftChanges
}

// calcSlashCoinStake return modified stake and SlashCoin for event. Only calculate, no changes
func calcSlashCoinStake(stake types.Stake, fraction sdk.Dec) (types.Stake, types.SlashCoin, error) {
	if stake.Type != types.StakeType_Coin {
		return types.Stake{}, types.SlashCoin{}, fmt.Errorf("call slashCoinStake not for coin stake")
	}
	slashAmount := sdk.NewDecFromInt(stake.Stake.Amount).Mul(fraction).RoundInt()
	newAmount := stake.Stake.Amount.Sub(slashAmount)
	newStake := types.Stake{
		Type:  types.StakeType_Coin,
		ID:    stake.ID,
		Stake: sdk.NewCoin(stake.Stake.Denom, newAmount),
	}
	slashCoin := types.SlashCoin{
		Slash: sdk.NewCoin(stake.Stake.Denom, slashAmount),
	}
	return newStake, slashCoin, nil
}

// slashCoinStake return modified stake and subtokens, SlashNFT for event. Only calculate, no changes
// stake.SubTokenIDs and subtokens order must be same, subtokens reserve must be not nil filled
func calcSlashNFTStake(stake types.Stake, subtokens []nfttypes.SubToken, fraction sdk.Dec) (types.Stake, types.SlashNFT, []nfttypes.SubToken, error) {
	if stake.Type != types.StakeType_NFT {
		return types.Stake{}, types.SlashNFT{}, []nfttypes.SubToken{}, fmt.Errorf("call slashCoinStake not for coin stake")
	}
	if len(stake.SubTokenIDs) != len(subtokens) {
		return types.Stake{}, types.SlashNFT{}, []nfttypes.SubToken{}, fmt.Errorf("lengths of stake subtokens and subtokens not equal")
	}

	totalSlash := sdk.ZeroInt()
	totalStake := sdk.ZeroInt()

	var newSubTokens []nfttypes.SubToken
	var slashSubTokens []types.SlashNFTSubToken

	for i := range stake.SubTokenIDs {
		if stake.SubTokenIDs[i] != subtokens[i].ID {
			return types.Stake{}, types.SlashNFT{}, []nfttypes.SubToken{}, fmt.Errorf("subtokens id not equal")
		}
		if subtokens[i].Reserve == nil {
			return types.Stake{}, types.SlashNFT{}, []nfttypes.SubToken{}, fmt.Errorf("subtokens reserve is nil")
		}
		slashAmount := sdk.NewDecFromInt(subtokens[i].Reserve.Amount).Mul(fraction).RoundInt()
		newAmount := subtokens[i].Reserve.Amount.Sub(slashAmount)
		newReserve := sdk.NewCoin(subtokens[i].Reserve.Denom, newAmount)

		newSubTokens = append(newSubTokens, nfttypes.SubToken{
			ID:      subtokens[i].ID,
			Owner:   subtokens[i].Owner,
			Reserve: &newReserve,
		})
		slashSubTokens = append(slashSubTokens, types.SlashNFTSubToken{
			ID:    subtokens[i].ID,
			Slash: sdk.NewCoin(subtokens[i].Reserve.Denom, slashAmount),
		})

		totalSlash = totalSlash.Add(slashAmount)
		totalStake = totalStake.Add(newAmount)
	}

	newStake := types.Stake{
		Type:        types.StakeType_NFT,
		ID:          stake.ID,
		Stake:       sdk.NewCoin(stake.Stake.Denom, totalStake),
		SubTokenIDs: stake.SubTokenIDs,
	}
	slashNFT := types.SlashNFT{
		ID:        stake.ID,
		SubTokens: slashSubTokens,
	}

	return newStake, slashNFT, newSubTokens, nil
}

package keeper

import (
	"fmt"

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
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	if infractionHeight > ctx.BlockHeight() {
		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))
	}

	validator, found := k.GetValidatorByConsAddrDecimal(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(
			"WARNING: ignored attempt to slash a nonexistent validator; we recommend you investigate immediately",
			"validator", consAddr.String(),
		)
		return
	}

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()
	valStatuses := make(map[string]types.BondStatus)
	for _, v := range k.GetAllValidators(ctx) {
		valStatuses[v.OperatorAddress] = v.Status
	}

	// call the before-modification hook
	k.BeforeValidatorModified(ctx, operatorAddress)

	//////////////////////////////////////////////////
	// 1. precalculation for correct burn of slashed coins
	accum := NewSlashesAccumulator(k, ctx, slashFactor, NewDecreasingFactors())
	for _, delegation := range k.GetValidatorDelegations(ctx, operatorAddress) {
		accum.AddDelegation(delegation, validator.Status, true)
	}
	if infractionHeight < ctx.BlockHeight() {
		for _, undelegation := range k.GetUndelegationsFromValidator(ctx, operatorAddress) {
			accum.AddUndelegation(undelegation, infractionHeight, true)
		}
		for _, redelegation := range k.GetRedelegationsFromSrcValidator(ctx, operatorAddress) {
			accum.AddRedelegation(redelegation, infractionHeight, valStatuses, true)
		}
	}

	// precalculation to check future coins burns
	var factors = NewDecreasingFactors()
	for _, coin := range accum.GetAllCoinsToBurn() {
		if coin.Denom == k.coinKeeper.GetBaseDenom(ctx) {
			factors.SetFactor(coin.Denom, sdk.OneDec())
			continue
		}
		f, err := k.coinKeeper.GetDecreasingFactor(ctx, coin)
		if err != nil {
			panic(fmt.Errorf("error in GetDecreasingFactor %s: %s", coin.Denom, err.Error()))
		}
		factors.SetFactor(coin.Denom, f)
	}

	//////////////////////////////////////////////////
	// 2. new accumulator to prepare changes
	accum = NewSlashesAccumulator(k, ctx, slashFactor, factors)
	for _, delegation := range k.GetValidatorDelegations(ctx, operatorAddress) {
		accum.AddDelegation(delegation, validator.Status, false)
	}
	if infractionHeight < ctx.BlockHeight() {
		for _, undelegation := range k.GetUndelegationsFromValidator(ctx, operatorAddress) {
			accum.AddUndelegation(undelegation, infractionHeight, false)
		}
		for _, redelegation := range k.GetRedelegationsFromSrcValidator(ctx, operatorAddress) {
			accum.AddRedelegation(redelegation, infractionHeight, valStatuses, false)
		}
	}

	//////////////////////////////////////////////////
	// 3. do changes
	for _, delegation := range accum.GetDelegations() {
		k.SetDelegation(ctx, delegation)
	}
	for _, undelegation := range accum.GetUndelegations() {
		k.SetUndelegation(ctx, undelegation)
	}
	for _, redelegation := range accum.GetRedelegations() {
		k.SetRedelegation(ctx, redelegation)
	}
	for _, chng := range accum.GetNFTChanges() {
		for _, sub := range chng.subtokens {
			k.nftKeeper.SetSubToken(ctx, chng.tokenID, sub)
		}
	}

	//////////////////////////////////////////////////
	// 4. burn coins
	if !accum.GetCoinsToBurnBonded().IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, types.BondedPoolName, accum.GetCoinsToBurnBonded())
		if err != nil {
			panic(fmt.Errorf("error in burn in bonded pool: %s", err.Error()))
		}
	}
	if !accum.GetCoinsToBurnUnbonded().IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, types.NotBondedPoolName, accum.GetCoinsToBurnUnbonded())
		if err != nil {
			panic(fmt.Errorf("error in burn in not_bonded pool: %s", err.Error()))
		}
	}
	if !accum.GetCoinsToBurnNFT().IsZero() {
		err := k.coinKeeper.BurnPoolCoins(ctx, nfttypes.ReservedPool, accum.GetCoinsToBurnNFT())
		if err != nil {
			panic(fmt.Errorf("error in burn for nft reserved pool: %s", err.Error()))
		}
	}

	//////////////////////////////////////////////////
	// 5. change stakes of custom coins
	stakeDecreasing := accum.GetCoinsToBurnBonded().Add(accum.GetCoinsToBurnUnbonded()...).Add(accum.GetCoinsToBurnNFT()...)
	for _, coin := range stakeDecreasing {
		k.SubCustomCoinStaked(ctx, coin)
	}

	//////////////////////////////////////////////////
	// 6. emit event
	ev := accum.GetEvent(validator.OperatorAddress)
	events.EmitTypedEvent(ctx, &ev)

	logger.Info(
		"validator slashed by slash factor",
		"validator", validator.GetOperator().String(),
		"slash_factor", slashFactor.String(),
		"burned", accum.GetAllCoinsToBurn(),
	)
}

// jail a validator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}

	validator.Jailed = true
	validator.Online = false
	k.SetValidator(ctx, validator)
	// Jailed validator will be processe in ApplyAndReturnValidatorSetUpdates
	//k.DeleteValidatorByPowerIndex(ctx, validator)
	// Deleting of start height reset MissedBlockCounter
	// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
	k.DeleteStartHeight(ctx, consAddr)

	k.Logger(ctx).Info("validator jailed", "validator", consAddr)
}

// unjail a validator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
	}
	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.Logger(ctx).Info("validator un-jailed", "validator", consAddr)
}

// structure accumulates changes during slash process
type slashesAccumulator struct {
	newDelegations          []types.Delegation
	newUndelegations        []types.Undelegation
	newRedelegations        []types.Redelegation
	delegationSlashEvents   map[string]types.DelegatorSlash
	undelegationSlashEvents map[undelegationKey]types.UndelegateSlash
	redelegationSlashEvents map[redelegationKey]types.RedelegateSlash
	nftChanges              []nftSubtokenChanges
	coinsToBurnBonded       sdk.Coins
	coinsToBurnUnbonded     sdk.Coins
	nftCoinsToBurn          sdk.Coins
	keeper                  Keeper
	ctx                     sdk.Context
	slashFactor             sdk.Dec
	factors                 *decreasingFactors
}

type undelegationKey struct {
	delegator string
	validator string
}

type redelegationKey struct {
	delegator    string
	validatorSrc string
	validatorDst string
}

func NewSlashesAccumulator(k Keeper, ctx sdk.Context, slashFactor sdk.Dec, factors *decreasingFactors) *slashesAccumulator {
	return &slashesAccumulator{
		coinsToBurnBonded:       sdk.NewCoins(),
		coinsToBurnUnbonded:     sdk.NewCoins(),
		nftCoinsToBurn:          sdk.NewCoins(),
		keeper:                  k,
		ctx:                     ctx,
		slashFactor:             slashFactor,
		factors:                 factors,
		delegationSlashEvents:   make(map[string]types.DelegatorSlash),
		undelegationSlashEvents: make(map[undelegationKey]types.UndelegateSlash),
		redelegationSlashEvents: make(map[redelegationKey]types.RedelegateSlash),
	}
}

func (sa *slashesAccumulator) GetAllCoinsToBurn() sdk.Coins {
	return sa.coinsToBurnBonded.Add(sa.coinsToBurnUnbonded...).Add(sa.nftCoinsToBurn...)
}

func (sa *slashesAccumulator) GetCoinsToBurnBonded() sdk.Coins {
	return sa.coinsToBurnBonded
}

func (sa *slashesAccumulator) GetCoinsToBurnUnbonded() sdk.Coins {
	return sa.coinsToBurnUnbonded
}

func (sa *slashesAccumulator) GetCoinsToBurnNFT() sdk.Coins {
	return sa.nftCoinsToBurn
}

func (sa *slashesAccumulator) GetDelegations() []types.Delegation {
	return sa.newDelegations
}

func (sa *slashesAccumulator) GetUndelegations() []types.Undelegation {
	return sa.newUndelegations
}

func (sa *slashesAccumulator) GetRedelegations() []types.Redelegation {
	return sa.newRedelegations
}

func (sa *slashesAccumulator) GetNFTChanges() []nftSubtokenChanges {
	return sa.nftChanges
}

func (sa *slashesAccumulator) GetEvent(operatorAddress string) types.EventValidatorSlash {
	var result types.EventValidatorSlash
	result.Validator = operatorAddress
	for _, ev := range sa.delegationSlashEvents {
		result.Delegators = append(result.Delegators, ev)
	}
	for _, ev := range sa.undelegationSlashEvents {
		result.Undelegations = append(result.Undelegations, ev)
	}
	for _, ev := range sa.redelegationSlashEvents {
		result.Redelegations = append(result.Redelegations, ev)
	}
	return result
}

// AddDelegation process delegation. if simulate = true, only calculates coins for burning
func (sa *slashesAccumulator) AddDelegation(delegation types.Delegation, validatorStatus types.BondStatus, simulate bool) {
	newStake, slashCoin, slashNFT, nftChanges := sa.keeper.calcSlashStake(sa.ctx, delegation.Stake, sa.slashFactor, sa.factors)

	switch delegation.Stake.Type {
	case types.StakeType_Coin:
		switch validatorStatus {
		case types.BondStatus_Bonded:
			sa.coinsToBurnBonded = sa.coinsToBurnBonded.Add(slashCoin.Slash)
		case types.BondStatus_Unbonded, types.BondStatus_Unbonding:
			sa.coinsToBurnUnbonded = sa.coinsToBurnUnbonded.Add(slashCoin.Slash)
		}
	case types.StakeType_NFT:
		for _, sub := range slashNFT.SubTokens {
			sa.nftCoinsToBurn = sa.nftCoinsToBurn.Add(sub.Slash)
		}
	}

	if simulate {
		return
	}

	// add changes
	sa.newDelegations = append(sa.newDelegations, types.Delegation{
		Delegator: delegation.Delegator,
		Validator: delegation.Validator,
		Stake:     newStake,
	})
	if nftChanges.tokenID != "" {
		sa.nftChanges = append(sa.nftChanges, nftChanges)
	}
	// accumulate events
	ev, ok := sa.delegationSlashEvents[delegation.Delegator]
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
	sa.delegationSlashEvents[delegation.Delegator] = ev

}

// AddUndelegation process undelegation record. if simulate = true, only calculates coins for burning
func (sa *slashesAccumulator) AddUndelegation(undelegation types.Undelegation, infractionHeight int64, simulate bool) {
	var doChanges bool
	var newUndelegation types.Undelegation
	now := sa.ctx.BlockHeader().Time
	newUndelegation.Delegator = undelegation.Delegator
	newUndelegation.Validator = undelegation.Validator
	for _, entry := range undelegation.Entries {
		if entry.CreationHeight < infractionHeight || entry.IsMature(now) {
			newUndelegation.Entries = append(newUndelegation.Entries, entry)
			continue
		}
		doChanges = true
		newStake, slashCoin, slashNFT, nftChanges := sa.keeper.calcSlashStake(sa.ctx, entry.Stake, sa.slashFactor, sa.factors)

		switch entry.Stake.Type {
		case types.StakeType_Coin:
			sa.coinsToBurnUnbonded = sa.coinsToBurnUnbonded.Add(slashCoin.Slash)
		case types.StakeType_NFT:
			for _, sub := range slashNFT.SubTokens {
				sa.nftCoinsToBurn = sa.nftCoinsToBurn.Add(sub.Slash)
			}
		}
		if simulate {
			continue
		}

		// add changes
		newUndelegation.Entries = append(newUndelegation.Entries, types.UndelegationEntry{
			CreationHeight: entry.CreationHeight,
			CompletionTime: entry.CompletionTime,
			Stake:          newStake,
		})
		if nftChanges.tokenID != "" {
			sa.nftChanges = append(sa.nftChanges, nftChanges)
		}
		// accumulate events
		key := undelegationKey{undelegation.Delegator, undelegation.Validator}
		ev, ok := sa.undelegationSlashEvents[key]
		if !ok {
			ev = types.UndelegateSlash{
				Delegator: undelegation.Delegator,
				Validator: undelegation.Validator,
			}
		}
		if slashCoin.Slash.Denom != "" {
			ev.Coins = append(ev.Coins, slashCoin)
		}
		if slashNFT.ID != "" {
			ev.NFTs = append(ev.NFTs, slashNFT)
		}
		sa.undelegationSlashEvents[key] = ev
	}
	if !simulate && doChanges {
		sa.newUndelegations = append(sa.newUndelegations, newUndelegation)
	}
}

// AddRedelegation process redelegation record. if simulate = true, only calculates coins for burning
func (sa *slashesAccumulator) AddRedelegation(redelegation types.Redelegation, infractionHeight int64,
	validatorStatuses map[string]types.BondStatus, simulate bool) {

	var doChanges bool
	var newRedelegation types.Redelegation
	now := sa.ctx.BlockHeader().Time
	newRedelegation.Delegator = redelegation.Delegator
	newRedelegation.ValidatorSrc = redelegation.ValidatorSrc
	newRedelegation.ValidatorDst = redelegation.ValidatorDst
	for _, entry := range redelegation.Entries {
		if entry.CreationHeight < infractionHeight || entry.IsMature(now) {
			newRedelegation.Entries = append(newRedelegation.Entries, entry)
			continue
		}
		doChanges = true
		newStake, slashCoin, slashNFT, nftChanges := sa.keeper.calcSlashStake(sa.ctx, entry.Stake, sa.slashFactor, sa.factors)

		switch entry.Stake.Type {

		case types.StakeType_Coin:
			sa.coinsToBurnUnbonded = sa.coinsToBurnUnbonded.Add(slashCoin.Slash)

		case types.StakeType_NFT:
			for _, sub := range slashNFT.SubTokens {
				sa.nftCoinsToBurn = sa.nftCoinsToBurn.Add(sub.Slash)
			}
		}
		if simulate {
			continue
		}

		// add changes
		newRedelegation.Entries = append(newRedelegation.Entries, types.RedelegationEntry{
			CreationHeight: entry.CreationHeight,
			CompletionTime: entry.CompletionTime,
			Stake:          newStake,
		})
		if nftChanges.tokenID != "" {
			sa.nftChanges = append(sa.nftChanges, nftChanges)
		}
		// accumulate events
		key := redelegationKey{redelegation.Delegator, redelegation.ValidatorSrc, redelegation.ValidatorDst}
		ev, ok := sa.redelegationSlashEvents[key]
		if !ok {
			ev = types.RedelegateSlash{
				Delegator:    redelegation.Delegator,
				ValidatorSrc: redelegation.ValidatorSrc,
				ValidatorDst: redelegation.ValidatorDst,
			}
		}
		if slashCoin.Slash.Denom != "" {
			ev.Coins = append(ev.Coins, slashCoin)
		}
		if slashNFT.ID != "" {
			ev.NFTs = append(ev.NFTs, slashNFT)
		}
		sa.redelegationSlashEvents[key] = ev
	}
	if !simulate && doChanges {
		sa.newRedelegations = append(sa.newRedelegations, newRedelegation)
	}
}

type nftSubtokenChanges struct {
	tokenID   string
	subtokens []nfttypes.SubToken
}

// calcSlashStake prepare changes. Only calculate, no changes
func (k Keeper) calcSlashStake(ctx sdk.Context, stake types.Stake, slashFactor sdk.Dec, factors *decreasingFactors) (types.Stake,
	types.SlashCoin, types.SlashNFT, nftSubtokenChanges) {
	var newStake types.Stake
	var slashCoin types.SlashCoin
	var slashNFT types.SlashNFT
	var newSubTokens []nfttypes.SubToken
	var nftChanges nftSubtokenChanges
	var err error
	switch stake.Type {
	case types.StakeType_Coin:
		newStake, slashCoin, err = calcSlashCoinStake(stake, slashFactor, factors)
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
		newStake, slashNFT, newSubTokens, err = calcSlashNFTStake(stake, subtokens, slashFactor, factors)
		if err != nil {
			panic(err)
		}
		nftChanges = nftSubtokenChanges{tokenID: token.ID, subtokens: newSubTokens}
	}
	return newStake, slashCoin, slashNFT, nftChanges
}

// calcSlashCoinStake return modified stake and SlashCoin for event. Only calculate, no changes
func calcSlashCoinStake(stake types.Stake, slashFactor sdk.Dec, factors *decreasingFactors) (types.Stake, types.SlashCoin, error) {
	if stake.Type != types.StakeType_Coin {
		return types.Stake{}, types.SlashCoin{}, fmt.Errorf("call slashCoinStake not for coin stake")
	}
	slashAmount := sdk.NewDecFromInt(stake.Stake.Amount).Mul(slashFactor).Mul(factors.Factor(stake.Stake.Denom)).RoundInt()
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
func calcSlashNFTStake(stake types.Stake, subtokens []nfttypes.SubToken, slashFactor sdk.Dec, factors *decreasingFactors) (types.Stake, types.SlashNFT, []nfttypes.SubToken, error) {
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
		slashAmount := sdk.NewDecFromInt(subtokens[i].Reserve.Amount).Mul(slashFactor).RoundInt()
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

// helper type
// Decreasing factors
// map of coin denom -> sdk.Dec. By default return 1.0 for any coin

type decreasingFactors struct {
	factors map[string]sdk.Dec
}

func NewDecreasingFactors() *decreasingFactors {
	return &decreasingFactors{
		factors: make(map[string]sdk.Dec),
	}
}

func (df *decreasingFactors) Factor(denom string) sdk.Dec {
	v, ok := df.factors[denom]
	if !ok {
		return sdk.OneDec()
	}
	return v
}

func (df *decreasingFactors) SetFactor(denom string, factor sdk.Dec) {
	df.factors[denom] = factor
}

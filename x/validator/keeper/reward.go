package keeper

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var daoAccount = "d01pk2rurh73er88p032qrd6kq5xmu53thjqc22mu"
var developAccount = "d01gsa4w0cuyjqwt9j7qtc32m6n0lkyxfan9s2ghh"

var DAOCommission = sdk.NewDec(5).QuoInt64(100)
var DevelopCommission = sdk.NewDec(5).QuoInt64(100)

func (k Keeper) PayRewards(ctx sdk.Context) error {
    ctx.Logger().Info("Starting PayRewards function")
    e := types.EventPayRewards{}

    validators := k.GetAllValidators(ctx)
    ctx.Logger().Info("Got all validators", "count", len(validators))

    delByValidator := k.GetAllDelegationsByValidator(ctx)
    ctx.Logger().Info("Got all delegations by validator", "count", len(delByValidator))

    customCoinStaked := k.GetAllCustomCoinsStaked(ctx)
    ctx.Logger().Info("Got all custom coins staked", "count", len(customCoinStaked))

    customCoinPrices := k.CalculateCustomCoinPrices(ctx, customCoinStaked)
    ctx.Logger().Info("Calculated custom coin prices", "count", len(customCoinPrices))

    ctx.Logger().Debug("custom coin staked", "is", customCoinStaked)
    ctx.Logger().Debug("custom prices", "is", customCoinPrices)

    allRewards := sdk.NewInt(0)
    allDelegationSum := sdk.NewInt(0)
    allHoldBigOneYearsSum := sdk.NewInt(0)

    coin, err := k.coinKeeper.GetCoin(ctx, k.BaseDenom(ctx))
    if err != nil {
        ctx.Logger().Error("Failed to get base coin", "error", err)
        return err
    }
    allEmmision := coin.LimitVolume
    ctx.Logger().Info("Got base coin and emission", "baseDenom", k.BaseDenom(ctx), "emission", allEmmision)

    for i, val := range validators {
        ctx.Logger().Info("Processing validator", "index", i, "address", val.GetOperator())
        if val.Rewards.IsZero() {
            ctx.Logger().Info("Validator has zero rewards, skipping", "address", val.GetOperator())
            continue
        }
        validator := val.GetOperator()

        totalStake, err := k.CalculateTotalPowerWithDelegationsAndPrices(ctx, val.GetOperator(), delByValidator[validator.String()], customCoinPrices)
        if err != nil {
            ctx.Logger().Error("Failed to calculate total power", "error", err)
            return err
        }
        ctx.Logger().Info("Calculated total stake for validator", "address", validator, "totalStake", totalStake)

        allDelegationSum = allDelegationSum.Add(totalStake)

        rewards := val.Rewards
        allRewards = allRewards.Add(rewards)
        ctx.Logger().Info("Added rewards to total", "validator", validator, "rewards", rewards, "allRewards", allRewards)

        for _, del := range delByValidator[validator.String()] {
            delStake := del.GetStake().GetStake()

            if delStake.Denom != k.BaseDenom(ctx) {
                delCoinPrice, ok := customCoinPrices[delStake.Denom]
                if !ok {
                    ctx.Logger().Error("Price not found for custom coin", "denom", delStake.Denom)
                    return fmt.Errorf("not found price for custom coin %s, base denom is %s, validator is %s, delegator is %s", delStake.Denom, k.BaseDenom(ctx), validator.String(), del.Delegator)
                }
                baseAmount := sdk.NewDecFromInt(sdk.NewInt(0)).TruncateInt()
                var deleteHolds []*types.StakeHold
                for _, hold := range del.GetStake().GetHolds() {
                    ctx.Logger().Info("Processing hold", "startTime", hold.HoldStartTime, "endTime", hold.HoldEndTime)
                    dateStart := time.Unix(hold.HoldStartTime, 0)
                    dateEnd := time.Unix(hold.HoldEndTime, 0)
                    if hold.HoldStartTime == 0 {
                        deleteHolds = append(deleteHolds, hold)
                    }
                    difference := dateEnd.Sub(dateStart)
                    if (difference.Hours() / 24 / 365) >= 1 {
                        baseAmount = baseAmount.Add(sdk.NewDecFromInt(hold.Amount).Mul(delCoinPrice).TruncateInt())
                        ctx.Logger().Info("Added to base amount", "amount", baseAmount)
                    }
                }
                allHoldBigOneYearsSum = allHoldBigOneYearsSum.Add(baseAmount)
            }
        }
    }

    ctx.Logger().Info("Calculating percentForHold")
    percentForHold := sdk.NewDecFromInt(allDelegationSum).Quo(sdk.NewDecFromInt(allEmmision)).Mul(sdk.NewDec(100)).TruncateInt()
    percentForHold = sdk.NewInt(100).Sub(percentForHold)
    if percentForHold.IsNegative() {
        percentForHold = sdk.NewInt(0)
    }
    ctx.Logger().Info("Calculated percentForHold", "percent", percentForHold)

    ctx.Logger().Info("Calculating sumRewardForHold")
    sumRewardForHold := sdk.NewDecFromInt(allRewards).Mul(sdk.NewDecFromInt(percentForHold).QuoInt64(100)).TruncateInt()
    ctx.Logger().Info("Calculated sumRewardForHold", "sum", sumRewardForHold)

    for _, val := range validators {
        ctx.Logger().Info("Processing validator for rewards", "address", val.GetOperator())
        if val.Rewards.IsZero() {
            ctx.Logger().Info("Validator has zero rewards, skipping", "address", val.GetOperator())
            continue
        }
        validator := val.GetOperator()
        rewards := val.Rewards
        accumRewards := sdk.NewDecFromInt(rewards).Mul(sdk.NewDecFromInt(percentForHold).QuoInt64(100)).TruncateInt()

        daoWallet := sdk.MustAccAddressFromBech32(daoAccount)
        developWallet := sdk.MustAccAddressFromBech32(developAccount)

        ctx.Logger().Info("Calculating and sending DAO commission")
        daoVal := sdk.NewDecFromInt(rewards).Mul(DAOCommission).TruncateInt()
        err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, daoWallet, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), daoVal)))
        if err != nil {
            ctx.Logger().Error("Failed to send DAO commission", "error", err)
            return err
        }

        ctx.Logger().Info("Calculating and sending Develop commission")
        developVal := sdk.NewDecFromInt(rewards).Mul(DevelopCommission).TruncateInt()
        err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, developWallet, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), developVal)))
        if err != nil {
            ctx.Logger().Error("Failed to send Develop commission", "error", err)
            return err
        }

        rewards = rewards.Sub(daoVal)
        rewards = rewards.Sub(developVal)

        ctx.Logger().Info("Calculating and sending validator commission")
        valComission := sdk.NewDecFromInt(rewards).Mul(val.Commission).TruncateInt()
		var valRewardAddress sdk.AccAddress

        // RewardAddress may be legacy address with 'dx' prefix
        valRewardAddress, err = sdk.GetFromBech32(val.RewardAddress, config.Bech32PrefixAccAddr)
        if err != nil {
            ctx.Logger().Info("Failed to parse reward address with current prefix, trying legacy prefix", "address", val.RewardAddress)
            valRewardAddress, err = sdk.GetFromBech32(val.RewardAddress, "dx")
            if err != nil {
                ctx.Logger().Error("Failed to parse reward address with legacy prefix", "address", val.RewardAddress, "error", err)
                return err
            }
        }
        ctx.Logger().Info("Parsed validator reward address", "address", valRewardAddress.String())

        err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, valRewardAddress, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), valComission)))
        if err != nil {
            ctx.Logger().Error("Failed to send validator commission", "error", err)
            return err
        }
        ctx.Logger().Info("Sent validator commission", "amount", valComission, "address", valRewardAddress.String())

        rewards = rewards.Sub(valComission)

        // event
        valEvent := types.ValidatorReward{
            Validator:   validator.String(),
            Dao:         daoVal,
            Develop:     developVal,
            Commission:  valComission,
            Accumulated: accumRewards,
            Delegators:  nil,
        }

        totalStake, err := k.CalculateTotalPowerWithDelegationsAndPrices(ctx, val.GetOperator(), delByValidator[validator.String()], customCoinPrices)
        if err != nil {
            ctx.Logger().Error("Failed to calculate total power for validator", "error", err)
            return err
        }
        ctx.Logger().Info("Calculated total stake for validator", "validator", validator.String(), "totalStake", totalStake)

        remainder := rewards
        for _, del := range delByValidator[validator.String()] {
            reward := sdk.NewIntFromBigInt(rewards.BigInt())
            // calculate share
            delStake := del.GetStake().GetStake()

            baseAmount := delStake.Amount
            sumHold := sdk.NewDecFromInt(sdk.NewInt(0)).TruncateInt()
            if delStake.Denom != k.BaseDenom(ctx) {
                delCoinPrice, ok := customCoinPrices[delStake.Denom]
                if !ok {
                    ctx.Logger().Error("Price not found for custom coin", "denom", delStake.Denom)
                    return fmt.Errorf("not found price for custom coin %s, base denom is %s, validator is %s, delegator is %s", delStake.Denom, k.BaseDenom(ctx), validator.String(), del.Delegator)
                }
                baseAmount = sdk.NewDecFromInt(delStake.Amount).Mul(delCoinPrice).TruncateInt()
                for _, hold := range del.GetStake().GetHolds() {
                    dateStart := time.Unix(hold.HoldStartTime, 0)
                    dateEnd := time.Unix(hold.HoldEndTime, 0)
                    difference := dateEnd.Sub(dateStart)
                    if (difference.Hours() / 24 / 365) >= 1 {
                        sumHold = sumHold.Add(sdk.NewDecFromInt(hold.Amount).Mul(delCoinPrice).TruncateInt())
                    }
                }
            } else {
                for _, hold := range del.GetStake().GetHolds() {
                    dateStart := time.Unix(hold.HoldStartTime, 0)
                    dateEnd := time.Unix(hold.HoldEndTime, 0)
                    difference := dateEnd.Sub(dateStart)
                    if (difference.Hours() / 24 / 365) >= 1 {
                        sumHold = sumHold.Add(sdk.NewDecFromInt(hold.Amount).TruncateInt())
                    }
                }
            }

            reward = reward.Mul(baseAmount).Quo(totalStake)
            if reward.LT(sdk.NewInt(1)) {
                continue
            }
            // pay reward
            err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, del.GetDelegator(), sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), reward)))
            if err != nil {
                ctx.Logger().Error("Failed to send reward to delegator", "error", err)
                continue
            }
            remainder = remainder.Sub(reward)
            // event
            if del.GetStake().GetType() == types.StakeType_Coin {
                // rewards coins
                delEvent := types.DelegatorReward{
                    Delegator: del.Delegator,
                    Coins: []types.StakeReward{
                        {
                            ID:       k.BaseDenom(ctx),
                            Reward:   reward,
                            RewardID: del.GetStake().GetID(),
                        },
                    },
                    NFTs: nil,
                }
                valEvent.Delegators = append(valEvent.Delegators, delEvent)
            }
            if del.GetStake().GetType() == types.StakeType_NFT {
                // rewards nft
                nftEvent := types.DelegatorReward{
                    Delegator: del.Delegator,
                    Coins:     nil,
                    NFTs: []types.StakeReward{
                        {
                            ID:       del.GetStake().GetID(),
                            Reward:   reward,
                            RewardID: del.GetStake().GetID(),
                        },
                    },
                }
                valEvent.Delegators = append(valEvent.Delegators, nftEvent)
            }
            if sumHold.LT(sdk.NewInt(1)) {
                continue
            }
            if allHoldBigOneYearsSum.LT(sdk.NewInt(1)) {
                continue
            }
            if sumRewardForHold.LT(sdk.NewInt(1)) {
                continue
            }
			rewardHold := sumRewardForHold.Mul(sumHold).Quo(allHoldBigOneYearsSum)
            if rewardHold.LT(sdk.NewInt(1)) {
                continue
            }
            // pay reward
            err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, del.GetDelegator(), sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), rewardHold)))
            if err != nil {
                ctx.Logger().Error("Failed to send hold reward to delegator", "error", err, "delegator", del.Delegator)
                continue
            }
            // event
            if del.GetStake().GetType() == types.StakeType_Coin {
                // rewards coins
                delEvent := types.DelegatorReward{
                    Delegator: del.Delegator,
                    Coins: []types.StakeReward{
                        {
                            ID:       k.BaseDenom(ctx),
                            Reward:   rewardHold,
                            RewardID: del.GetStake().GetID(),
                        },
                    },
                    NFTs: nil,
                }
                valEvent.DelegatorHolds = append(valEvent.DelegatorHolds, delEvent)
            }
            if del.GetStake().GetType() == types.StakeType_NFT {
                // rewards nft
                ctx.Logger().Info("NFT hold reward not implemented", "delegator", del.Delegator)
                //nftEvent := types.DelegatorReward{
                //    Delegator: del.Delegator,
                //    Coins:     nil,
                //    NFTs: []types.StakeReward{
                //        {
                //            ID:       del.GetStake().GetID(),
                //            Reward:   reward,
                //            RewardID: del.GetStake().GetID(),
                //        },
                //    },
                //}
                //valEvent.Delegators = append(valEvent.Delegators, nftEvent)
            }
        }
        // update validator rewards
        valRewards, err := k.GetValidatorRS(ctx, validator)
        if err != nil {
            ctx.Logger().Error("Failed to get validator rewards", "error", err, "validator", validator)
            return err
        }
        valRewards.Rewards = sdk.ZeroInt()
        valRewards.Stake = TokensToConsensusPower(totalStake)
        if val.Status != types.BondStatus_Bonded {
            valRewards.Stake = 0
        }
        k.SetValidatorRS(ctx, validator, valRewards)
        ctx.Logger().Info("Updated validator rewards", "validator", validator, "stake", valRewards.Stake)

        if val.Status == types.BondStatus_Bonded {
            k.DeleteValidatorByPowerIndex(ctx, val)
            val.Stake = TokensToConsensusPower(totalStake)
            k.SetValidatorByPowerIndex(ctx, val)
            ctx.Logger().Info("Updated bonded validator power index", "validator", validator, "stake", val.Stake)
        }

        e.Validators = append(e.Validators, valEvent)
        ctx.Logger().Info("Added validator event", "validator", validator)
    }

    err = events.EmitTypedEvent(ctx, &e)
    if err != nil {
        ctx.Logger().Error("Failed to emit PayRewards event", "error", err)
        return err
    }
    ctx.Logger().Info("Emitted PayRewards event")

    ctx.Logger().Info("Finished PayRewards function")
    return nil
}

const DAOAddress1 = "d01tlrmpaxn6v2xzwxftmd77tunuxctnqy7apu9x6"
const DAOAddress2 = "d015v26l75pzlc7u02p9u7xg80uynkj2lck2ehtjq"
const DAOAddress3 = "d016tfxmfsfftnaukvwuagc4svs8qsp6pzh3s2uky"

func (k Keeper) getDAO(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("d01pk2rurh73er88p032qrd6kq5xmu53thjqc22mu")
	if err != nil {
		return nil, err
	}

	wallet, err := k.multisigKeeper.GetWallet(ctx, address.String())
	if err != nil {
		return nil, err
	}

	if wallet.Address != "" {
		return address, nil
	}

	owner1, err := sdk.AccAddressFromBech32(DAOAddress1)
	if err != nil {
		return nil, err
	}
	owner2, err := sdk.AccAddressFromBech32(DAOAddress2)
	if err != nil {
		return nil, err
	}
	owner3, err := sdk.AccAddressFromBech32(DAOAddress3)
	if err != nil {
		return nil, err
	}

	owners := []string{
		owner1.String(), owner2.String(), owner3.String(),
	}

	wallet = multisig.Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   []uint32{1, 1, 1},
		Threshold: 3}

	k.multisigKeeper.SetWallet(ctx, wallet)
	return address, nil
}

const DevelopAddress1 = "d01zq5slrn7988ml6pldn8dhu9r5aaphgkdr9juqa"
const DevelopAddress2 = "d01l8qhz9pc4nct0j342r3l6380lg5dpxs2a5eawj"
const DevelopAddress3 = "d01864lut35xymux2xwnlvdwnjww7q7xy25guzz3f"

func (k Keeper) getDevelop(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("d01gsa4w0cuyjqwt9j7qtc32m6n0lkyxfan9s2ghh")
	if err != nil {
		return nil, err
	}

	wallet, err := k.multisigKeeper.GetWallet(ctx, address.String())
	if err != nil {
		return nil, err
	}
	if wallet.Address != "" {
		return address, nil
	}

	owner1, err := sdk.AccAddressFromBech32(DevelopAddress1)
	if err != nil {
		return nil, err
	}
	owner2, err := sdk.AccAddressFromBech32(DevelopAddress2)
	if err != nil {
		return nil, err
	}
	owner3, err := sdk.AccAddressFromBech32(DevelopAddress3)
	if err != nil {
		return nil, err
	}

	owners := []string{
		owner1.String(), owner2.String(), owner3.String(),
	}

	wallet = multisig.Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   []uint32{1, 1, 1},
		Threshold: 3}

	k.multisigKeeper.SetWallet(ctx, wallet)
	return address, nil
}

func (k Keeper) CalculateCustomCoinPrices(ctx sdk.Context, ccs map[string]sdkmath.Int) map[string]sdk.Dec {
	ctx.Logger().Debug("custom coins staked", "is", ccs)
	prices := make(map[string]sdk.Dec)
	for denom, staked := range ccs {
		coin, err := k.coinKeeper.GetCoin(ctx, denom)
		if err != nil {
			panic(err)
		}

		allPrice := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(coin.CRR), staked)
		prices[denom] = sdk.NewDecFromInt(allPrice).Quo(sdk.NewDecFromInt(staked))
	}

	return prices
}

func (k Keeper) calculateCustomCoinPrice(ctx sdk.Context, denom string, staked sdkmath.Int) sdk.Dec {
	coin, err := k.coinKeeper.GetCoin(ctx, denom)
	if err != nil {
		panic(err)
	}

	allPrice := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(coin.CRR), staked)
	return sdk.NewDecFromInt(allPrice).Quo(sdk.NewDecFromInt(staked))
}

func (k Keeper) CalculateTotalPowerWithDelegationsAndPrices(ctx sdk.Context, validator sdk.ValAddress, delegations types.Delegations, ccp map[string]sdk.Dec) (sdkmath.Int, error) {
	stakeInBaseCoin := sdk.ZeroInt()
	for _, del := range delegations {
		delStake := del.GetStake().GetStake()

		baseAmount := delStake.Amount
		if delStake.Denom != k.BaseDenom(ctx) {
			delCoinPrice, ok := ccp[delStake.Denom]
			if !ok {
				return stakeInBaseCoin, fmt.Errorf("not found price for custom coin %s, base denom is %s, validator is %s, delegator is %s", delStake.Denom, k.BaseDenom(ctx), validator.String(), del.Delegator)
			}
			baseAmount = sdk.NewDecFromInt(delStake.Amount).Mul(delCoinPrice).TruncateInt()
		}

		stakeInBaseCoin = stakeInBaseCoin.Add(baseAmount)
	}

	return stakeInBaseCoin, nil
}

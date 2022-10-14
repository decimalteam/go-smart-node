package keeper

import (
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var DAOCommission = sdk.NewDec(5).QuoInt64(100)
var DevelopCommission = sdk.NewDec(5).QuoInt64(100)

func (k Keeper) PayRewards(ctx sdk.Context, totalStake sdkmath.Int) error {
	events := types.EventPayRewards{}
	nbPool := k.GetNotBondedPool(ctx).String()

	validators := k.GetAllValidators(ctx)
	delByValidator := k.GetAllDelegationsByValidator(ctx)
	for _, val := range validators {
		if val.Rewards.IsZero() {
			continue
		}
		validator := val.GetOperator()
		rewards := val.Rewards
		accumRewards := rewards

		daoWallet, err := k.getDAO(ctx)
		if err != nil {
			return err
		}
		developWallet, err := k.getDevelop(ctx)
		if err != nil {
			return err
		}

		// dao commission
		daoVal := sdk.NewDecFromInt(rewards).Mul(DAOCommission).TruncateInt()
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, nbPool, daoWallet, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), daoVal)))
		if err != nil {
			return err
		}

		// develop commission
		developVal := sdk.NewDecFromInt(rewards).Mul(DevelopCommission).TruncateInt()
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, nbPool, developWallet, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), developVal)))
		if err != nil {
			return err
		}

		rewards = rewards.Sub(daoVal)
		rewards = rewards.Sub(developVal)

		// validator commission
		valComission := sdk.NewDecFromInt(rewards).Mul(val.Commission).TruncateInt()
		valRewardAddress := sdk.MustAccAddressFromBech32(val.RewardAddress)
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, nbPool, valRewardAddress, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), valComission)))
		if err != nil {
			return err
		}

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

		remainder := rewards
		for _, del := range delByValidator[validator.String()] {
			reward := sdk.NewIntFromBigInt(rewards.BigInt())
			// calculate share
			delStake := del.GetStake().GetStake()
			if del.Stake.SubTokenIDs != nil && len(del.Stake.SubTokenIDs) != 0 {
				delStake = k.getSumSubTokensReserve(ctx, del.GetStake().GetID(), del.GetStake().GetSubTokenIDs())
			}

			defAmount := k.ToBaseCoin(ctx, delStake)

			reward = reward.Mul(defAmount.Amount).Quo(totalStake)
			if reward.LT(sdk.NewInt(1)) {
				continue
			}

			// pay reward
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, nbPool, del.GetDelegator(), sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), reward)))
			if err != nil {
				continue
			}
			remainder.Sub(reward)
			// event
			delEvent := types.DelegatorReward{
				Delegator: del.Delegator,
				Coins: []types.StakeReward{
					{
						ID:     k.BaseDenom(ctx),
						Reward: reward,
					},
				},
				NFTs: nil,
			}
			valEvent.Delegators = append(valEvent.Delegators, delEvent)
		}
		// update validator rewards
		valRewards, err := k.GetValidatorRewards(ctx, validator)
		if err != nil {
			return err
		}
		valRewards.Rewards = sdk.ZeroInt()
		k.SetValidatorRewards(ctx, validator, valRewards)

		events.Validators = append(events.Validators, valEvent)
	}

	err := ctx.EventManager().EmitTypedEvents(&events)
	if err != nil {
		return err
	}

	return nil
}

const DAOAddress1 = "dx18tay9ayumxjun9sexlq4t3nvt7zts5typnyjdr"
const DAOAddress2 = "dx1w54s4wq8atjmmu4snv0tt72qpvtg38megw5ngn"
const DAOAddress3 = "dx19ws36j00axpk0ytumc20l9wyv0ae26zygk2z0f"

func (k Keeper) getDAO(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("dx1pk2rurh73er88p032qrd6kq5xmu53thjylflsr")
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

const DevelopAddress1 = "dx1fpjhs2wlaz6dd95d0lmxj5tfrmncwg437jh0y3"
const DevelopAddress2 = "dx1lfleqkc39pt2jkyhr7m845x207kh5d9av3423z"
const DevelopAddress3 = "dx1f46tyn4wmnvuxfj9cu5yn6vn939spfzt3yhxey"

func (k Keeper) getDevelop(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("dx1gsa4w0cuyjqwt9j7qtc32m6n0lkyxfanphfaug")
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
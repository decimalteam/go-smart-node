package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// test slashesAccumulator coins factor calculation (1 part of slash)
func TestFactorCalculation(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	coins := sdk.NewCoins(sdk.NewCoin(dsc.ValidatorKeeper.BaseDenom(ctx), helpers.EtherToWei(sdkmath.NewInt(1_000_000))))
	addrDels, addrVals := generateAddresses(dsc, ctx, 10, coins)

	slashFactor := sdk.MustNewDecFromStr("0.5") // 50% for easy checking

	msgserver := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	// heights:
	// 1: create, delegate
	// 2: non-slashing undel,redel
	// 3: slashing undel,redel

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx := sdk.WrapSDKContext(ctx)
	//////////////////////////////
	// 1. create coins, validators, nfts

	// coin for factor 1.0: crr = 80, vol=2000, res=2000, collector=100, burn=100
	_, err := dsc.CoinKeeper.CreateCoin(goCtx, cointypes.NewMsgCreateCoin(
		addrDels[0],
		"factor10",
		"factor10",
		80,
		helpers.EtherToWei(sdkmath.NewInt(2000)),
		helpers.EtherToWei(sdkmath.NewInt(2000)),
		helpers.EtherToWei(sdkmath.NewInt(10_000)),
		"-",
	))
	require.NoError(t, err)
	// coin for factor 0.5: crr = 100, vol=1000, res=2000, collector=400, burn=200
	_, err = dsc.CoinKeeper.CreateCoin(goCtx, cointypes.NewMsgCreateCoin(
		addrDels[0],
		"factor05",
		"factor05",
		100,
		helpers.EtherToWei(sdkmath.NewInt(1000)),
		helpers.EtherToWei(sdkmath.NewInt(2000)),
		helpers.EtherToWei(sdkmath.NewInt(10_000)),
		"-",
	))
	require.NoError(t, err)
	// coin for factor 0.0: crr = 10, vol=2000, res=2000, collector=200, burn=400
	_, err = dsc.CoinKeeper.CreateCoin(goCtx, cointypes.NewMsgCreateCoin(
		addrDels[0],
		"factor00",
		"factor00",
		10,
		helpers.EtherToWei(sdkmath.NewInt(2000)),
		helpers.EtherToWei(sdkmath.NewInt(2000)),
		helpers.EtherToWei(sdkmath.NewInt(10_000)),
		"-",
	))
	require.NoError(t, err)

	// mint nft to check nft slashes
	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		addrDels[0],
		"abc",
		"token1",
		"http://localhost",
		false,
		addrDels[0],
		2,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200))), // 50% = 100
	))
	require.NoError(t, err)
	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		addrDels[0],
		"abc",
		"token2",
		"http://127.0.0.1",
		false,
		addrDels[0],
		2,
		sdk.NewCoin("factor10", helpers.EtherToWei(sdkmath.NewInt(100))), // 50% = 50
	))
	require.NoError(t, err)
	// mint nft for undelegation/redelegation check
	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		addrDels[0],
		"abc",
		"token_redelegation",
		"http://example.org",
		false,
		addrDels[0],
		2,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(140))), // 50% = 70
	))
	require.NoError(t, err)
	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		addrDels[0],
		"abc",
		"token_undelegation",
		"http://example.com",
		false,
		addrDels[0],
		2,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(160))), // 50% = 80
	))
	require.NoError(t, err)

	// create validator
	msgCreate, err := types.NewMsgCreateValidator(
		addrVals[0],
		addrDels[1],
		PKs[0],
		types.Description{Moniker: "monik"},
		sdk.ZeroDec(),
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1000))), //200 for delegation, 200+200 for unbondings, 200+200 for redelegations
	)
	require.NoError(t, err)
	_, err = msgserver.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	val, found := dsc.ValidatorKeeper.GetValidator(ctx, addrVals[0])
	require.True(t, found)
	val.Online = true
	val.Status = types.BondStatus_Bonded
	dsc.ValidatorKeeper.SetValidator(ctx, val)
	dsc.ValidatorKeeper.SetValidatorRS(ctx, addrVals[0], types.ValidatorRS{
		Rewards:      sdk.ZeroInt(),
		TotalRewards: sdk.ZeroInt(),
		Stake:        1,
	})

	// second validator to check redelegations
	msgCreate, err = types.NewMsgCreateValidator(
		addrVals[1],
		addrDels[2],
		PKs[1],
		types.Description{Moniker: "monik2"},
		sdk.ZeroDec(),
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	)
	require.NoError(t, err)
	_, err = msgserver.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	dsc.ValidatorKeeper.SetValidatorRS(ctx, addrVals[1], types.ValidatorRS{
		Rewards:      sdk.ZeroInt(),
		TotalRewards: sdk.ZeroInt(),
		Stake:        1,
	})

	//////////////////////////////
	// 2. prepare delegations, fee_collector, unbondings
	_, err = msgserver.Delegate(goCtx, types.NewMsgDelegate(
		addrDels[0],
		addrVals[0],
		sdk.NewCoin("factor10", helpers.EtherToWei(sdkmath.NewInt(100))), // 50% = 50
	))
	require.NoError(t, err)
	_, err = msgserver.Delegate(goCtx, types.NewMsgDelegate(
		addrDels[0],
		addrVals[0],
		sdk.NewCoin("factor05", helpers.EtherToWei(sdkmath.NewInt(400))), // 50% = 200
	))
	require.NoError(t, err)
	_, err = msgserver.Delegate(goCtx, types.NewMsgDelegate(
		addrDels[0],
		addrVals[0],
		sdk.NewCoin("factor00", helpers.EtherToWei(sdkmath.NewInt(800))), // 50% = 400
	))
	require.NoError(t, err)
	_, err = msgserver.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		addrDels[0],
		addrVals[0],
		"token1",
		[]uint32{1}, // 200 base coin, 50% = 100
	))
	require.NoError(t, err)
	_, err = msgserver.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		addrDels[0],
		addrVals[0],
		"token2",
		[]uint32{1}, // 100 of factor10, 50% = 50
	))
	require.NoError(t, err)
	_, err = msgserver.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		addrDels[0],
		addrVals[0],
		"token_redelegation",
		[]uint32{1}, // 140 base coin, 50% = 70
	))
	require.NoError(t, err)
	_, err = msgserver.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		addrDels[0],
		addrVals[0],
		"token_undelegation",
		[]uint32{1}, // 160 base coin, 50% = 80
	))
	require.NoError(t, err)

	// Now total stake for slashing is:
	// 1000 of base coin (from CreateValidator)+200 from nft token1 + 140 token_redelegation + 160 token_undelegation -(200+200 redel/undel)
	// 50% = 550
	// 100 of factor10 +100 from nft token2
	// 50% = 100
	// 400 of factor05
	// 50% = 200
	// 800 of factor00
	// 50% = 400

	// prepare fee_collector: 100 factor10, 400 factor05, 200 factor00
	err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, addrDels[0], authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewCoin("factor10", helpers.EtherToWei(sdkmath.NewInt(100)))))
	require.NoError(t, err)
	err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, addrDels[0], authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewCoin("factor05", helpers.EtherToWei(sdkmath.NewInt(400)))))
	require.NoError(t, err)
	err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, addrDels[0], authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewCoin("factor00", helpers.EtherToWei(sdkmath.NewInt(200)))))
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)
	// create non-slashing unbonding, redelegation
	_, err = msgserver.Undelegate(goCtx, types.NewMsgUndelegate(addrDels[0], addrVals[0],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200))),
	))
	require.NoError(t, err)
	_, err = msgserver.Redelegate(goCtx, types.NewMsgRedelegate(addrDels[0], addrVals[0], addrVals[1],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200))),
	))
	require.NoError(t, err)
	//

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 2)
	goCtx = sdk.WrapSDKContext(ctx)
	// create slashing undelegation, redelegation
	_, err = msgserver.Undelegate(goCtx, types.NewMsgUndelegate(addrDels[0], addrVals[0],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200))),
	))
	require.NoError(t, err)
	_, err = msgserver.Redelegate(goCtx, types.NewMsgRedelegate(addrDels[0], addrVals[0], addrVals[1],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200))),
	))
	require.NoError(t, err)
	_, err = msgserver.UndelegateNFT(goCtx, types.NewMsgUndelegateNFT(addrDels[0], addrVals[0],
		"token_undelegation", []uint32{1},
	))
	require.NoError(t, err)
	_, err = msgserver.RedelegateNFT(goCtx, types.NewMsgRedelegateNFT(addrDels[0], addrVals[0], addrVals[1],
		"token_redelegation", []uint32{1},
	))
	require.NoError(t, err)

	_, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], addrVals[0], "token_undelegation")
	require.False(t, found)

	//////////////////////////////
	// 2. calculate: this is part of Slash() to check correct calculation
	operatorAddress := addrVals[0]
	k := dsc.ValidatorKeeper
	infractionHeight := ctx.BlockHeight() - 1
	validator, found := k.GetValidator(ctx, operatorAddress)
	require.True(t, found)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// ------------- START PART OF Slash()
	valStatuses := make(map[string]types.BondStatus)
	for _, v := range k.GetAllValidators(ctx) {
		valStatuses[v.OperatorAddress] = v.Status
	}

	accum := keeper.NewSlashesAccumulator(k, ctx, slashFactor, keeper.NewDecreasingFactors())
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
	var factors = keeper.NewDecreasingFactors()
	for _, coin := range accum.GetAllCoinsToBurn() {
		if coin.Denom == dsc.CoinKeeper.GetBaseDenom(ctx) {
			factors.SetFactor(coin.Denom, sdk.OneDec())
			continue
		}
		f, err := dsc.CoinKeeper.GetDecreasingFactor(ctx, coin)
		if err != nil {
			panic(fmt.Errorf("error in GetDecreasingFactor %s: %s", coin.Denom, err.Error()))
		}
		factors.SetFactor(coin.Denom, f)
	}
	// ------------- END PART OF Slash()

	//////////////////////////////
	// 3. check calculation factors
	require.True(t, factors.Factor("factor10").Equal(sdk.MustNewDecFromStr("1.0")))
	require.True(t, factors.Factor("factor05").Equal(sdk.MustNewDecFromStr("0.5")))
	require.True(t, factors.Factor("factor00").Equal(sdk.MustNewDecFromStr("0.0")))
	coinsToBurn := accum.GetAllCoinsToBurn()
	require.True(t, coinsToBurn.IsEqual(
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(550))),
			sdk.NewCoin("factor10", helpers.EtherToWei(sdkmath.NewInt(100))),
			sdk.NewCoin("factor05", helpers.EtherToWei(sdkmath.NewInt(200))),
			sdk.NewCoin("factor00", helpers.EtherToWei(sdkmath.NewInt(400))),
		),
	))
	// correct coinsToBurn for next checks
	coinsToBurn = coinsToBurn.Sub(
		sdk.NewCoin("factor05", sdk.NewDecFromInt(coinsToBurn.AmountOf("factor05")).Mul(factors.Factor("factor05")).RoundInt()),
		sdk.NewCoin("factor00", coinsToBurn.AmountOf("factor00")),
	)

	//////////////////////////////
	// 4. make real slash
	startBondedBalance := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	startNotBondedBalance := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	startNFTPool := dsc.BankKeeper.GetAllBalances(ctx, dsc.AccountKeeper.GetModuleAddress(nfttypes.ReservedPool))

	consAddr, err := validator.GetConsAddr()
	require.NoError(t, err)

	// SLASH, SLASH, SLASH
	dsc.ValidatorKeeper.Slash(ctx, consAddr, infractionHeight, 0, slashFactor)

	//////////////////////////////
	// 5. check pools
	endBondedBalance := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	endNotBondedBalance := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	endNFTPool := dsc.BankKeeper.GetAllBalances(ctx, dsc.AccountKeeper.GetModuleAddress(nfttypes.ReservedPool))

	start := startBondedBalance.Add(startNotBondedBalance...).Add(startNFTPool...)
	endPlusBurn := endBondedBalance.Add(endNotBondedBalance...).Add(endNFTPool...).Add(coinsToBurn...)

	require.True(t, start.AmountOf(cmdcfg.BaseDenom).Equal(endPlusBurn.AmountOf(cmdcfg.BaseDenom)))
	require.True(t, start.IsEqual(endPlusBurn), "start =%s, end = %s", start, endPlusBurn)

	//////////////////////////////
	// 5. check delegations, un/redelegations and nfts
	del, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], addrVals[0], cmdcfg.BaseDenom)
	require.True(t, found)
	require.True(t, del.Stake.Stake.IsEqual(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(100)))))
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], addrVals[0], "token1")
	require.True(t, found)
	require.True(t, del.Stake.Stake.IsEqual(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(100)))))
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], addrVals[0], "token2")
	require.True(t, found)
	require.True(t, del.Stake.Stake.IsEqual(sdk.NewCoin("factor10", helpers.EtherToWei(sdk.NewInt(50)))))
	// check undelegation
	undel, found := dsc.ValidatorKeeper.GetUndelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	for _, ent := range undel.Entries {
		switch ent.CreationHeight {
		case 2:
			// non-slashed
			switch ent.Stake.ID {
			case cmdcfg.BaseDenom:
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200)))),
				)
			}
		case 3:
			// slashed
			switch ent.Stake.ID {
			case cmdcfg.BaseDenom:
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))),
				)
			case "token_undelegation":
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(80)))),
				)
				subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, "token_undelegation", 1)
				require.True(t, found)
				require.True(t,
					subtoken.Reserve.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(80)))),
				)
			}
		}
	}

	// check redelegation
	redel, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)
	for _, ent := range redel.Entries {
		switch ent.CreationHeight {
		case 2:
			// non-slashed
			switch ent.Stake.ID {
			case cmdcfg.BaseDenom:
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(200)))),
				)
			}
		case 3:
			// slashed
			switch ent.Stake.ID {
			case cmdcfg.BaseDenom:
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))),
				)
			case "token_undelegation":
				require.True(t,
					ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(70)))),
				)
				subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, "token_redelegation", 1)
				require.True(t, found)
				require.True(t,
					subtoken.Reserve.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(70)))),
				)
			}
		}
	}
}

func TestJailUnjail(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))
	consAdr := sdk.GetConsAddress(PKs[0])

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	//
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[0]))
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	powerVals, _, _ := dsc.ValidatorKeeper.GetAllValidatorsByPowerIndex(ctx)
	require.Len(t, powerVals, 2) // genesisn validator + 1 created

	dsc.ValidatorKeeper.Jail(ctx, consAdr)
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	powerVals, _, _ = dsc.ValidatorKeeper.GetAllValidatorsByPowerIndex(ctx)
	require.Len(t, powerVals, 1)
}

func TestEndBlockAfterJail(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))
	consAdr := sdk.GetConsAddress(PKs[0])

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	//
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[0]))
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	// begin
	dsc.ValidatorKeeper.Jail(ctx, consAdr)
	// in-block
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[0]))
	require.NoError(t, err)
	// end
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	val, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	require.Equal(t, types.BondStatus_Unbonded, val.Status)
	require.False(t, val.Online)
	require.True(t, val.Jailed)
}

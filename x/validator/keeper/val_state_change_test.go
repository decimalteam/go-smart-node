package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestStateOnlineOffline(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)
	nbPool := dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress()
	bPool := dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress()

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	genesisVals := dsc.ValidatorKeeper.GetValidators(ctx, 10)
	require.Len(t, genesisVals, 1)
	genesisVal := genesisVals[0]
	require.True(t, genesisVal.ConsensusPower() > 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB := dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsZero())
	startBalanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)
	balanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)

	//
	goCtx := sdk.WrapSDKContext(ctx)
	tokenID := "aaaaaaa"
	_, err := dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		accs[1],
		"abcdef",
		tokenID,
		"URL",
		false,
		accs[1],
		2,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.NoError(t, err)

	////////////////////////////////////////////////
	// 1. create second validator
	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[1], accs[1], PKs[1], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	// delegate NFT
	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(accs[1], vals[1], tokenID, []uint32{1}))
	require.NoError(t, err)
	subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, nbPool.String())

	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// new validator is not online, there is no changes in tendermint validators and powers
	require.Len(t, updates, 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	// check balance
	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
	require.True(t, balanceB.IsEqual(startBalanceB))
	// check nft
	subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, nbPool.String())

	////////////////////////////////////////////////
	// 2. increment block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	////////////////////////////////////////////////
	// 3. set second validator online
	msgOnline := types.NewMsgSetOnline(vals[1])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOnline(goCtx, msgOnline)
	require.NoError(t, err)
	// last validators must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)
	// check start height
	require.Equal(t, ctx.BlockHeight(), dsc.ValidatorKeeper.GetStartHeight(ctx, sdk.GetConsAddress(PKs[1])))

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(100+100)) // see MsgCreateValidator stake+NFT stake (MintToken)
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)
	newValidator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
	require.True(t, found)
	require.Equal(t, newValidator.ConsensusPower()+genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	// check pool
	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsZero())
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
	require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake)))
	// check nft
	subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, bPool.String())

	////////////////////////////////////////////////
	// 4. increment block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	////////////////////////////////////////////////
	// 5. set second validator offline
	msgOffline := types.NewMsgSetOffline(vals[1])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOffline(goCtx, msgOffline)
	require.NoError(t, err)
	// last validator must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)
	require.Equal(t, int64(-1), dsc.ValidatorKeeper.GetStartHeight(ctx, sdk.GetConsAddress(PKs[1])))

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(0))                  // 0 mean 'remove from validators'
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1) // genesis + new
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	// check pool
	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
	require.True(t, balanceB.IsEqual(startBalanceB))
	// check nft
	subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, nbPool.String())
	// check second subtoken
	subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, accs[1].String())

}

func TestApplyAndReturnValidatorSetUpdates(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)
	nbPool := dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress()
	bPool := dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress()

	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.RedelegationTime = time.Hour
	dsc.ValidatorKeeper.SetParams(ctx, params)

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	genesisVals := dsc.ValidatorKeeper.GetValidators(ctx, 10)
	require.Len(t, genesisVals, 1)
	genesisVal := genesisVals[0]
	require.True(t, genesisVal.ConsensusPower() > 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB := dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsZero())
	startBalanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)
	balanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)

	// create custom coin
	ccDenom := "custom"
	initVolume := keeper.TokensFromConsensusPower(100000000000)
	initReserve := keeper.TokensFromConsensusPower(1000)
	limitVolume := keeper.TokensFromConsensusPower(100000000000000000)
	crr := uint64(50)

	_, err := dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, ""))
	require.NoError(t, err)
	// ----------------------------

	// create nfts
	tokenID := "nft_denom"
	subTokenReserve := sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(100))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(accs[1], "collection", tokenID, "uri", true, accs[1], 5, subTokenReserve))
	require.NoError(t, err)
	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(accs[1], "collection", tokenID, "uri", true, accs[0], 2, subTokenReserve))
	require.NoError(t, err)

	// ----------------------------
	goCtx := sdk.WrapSDKContext(ctx)

	////////////////////////////////////////////////
	// 1. create two validators

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	// first
	{
		msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)

		// delegate NFT
		_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(accs[1], vals[0], tokenID, []uint32{1, 2}))
		require.NoError(t, err)
	}
	// second
	{
		msgCreate, err := types.NewMsgCreateValidator(vals[1], accs[1], PKs[1], types.Description{Moniker: "monik1"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)

		// delegate NFT
		_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(accs[0], vals[1], tokenID, []uint32{6}))
		require.NoError(t, err)
	}
	// third: validator need to receive redelegation
	{
		msgCreate, err := types.NewMsgCreateValidator(vals[2], accs[2], PKs[2], types.Description{Moniker: "monik2"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)
	}
	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// new validator is not online, there is a changes in tendermint validators and powers with 0 powers
	require.Len(t, updates, 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	// check balance
	{
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake).Add(creatorStake).Add(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB))
	}
	// check nft
	{
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
	}

	////////////////////////////////////////////////
	// 2. increment block

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)

	//////////////////////////////////////////////////
	//// 3. set validators online

	{
		// first val
		msgOnline := types.NewMsgSetOnline(vals[0])
		_, err = msgsrv.SetOnline(goCtx, msgOnline)
		require.NoError(t, err)

		// second val
		msgOnline = types.NewMsgSetOnline(vals[1])
		_, err = msgsrv.SetOnline(goCtx, msgOnline)
		require.NoError(t, err)

		// third val
		msgOnline = types.NewMsgSetOnline(vals[2])
		_, err = msgsrv.SetOnline(goCtx, msgOnline)
		require.NoError(t, err)
	}

	// last validators must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// check updates
	{
		require.Len(t, updates, 3)
		require.Equal(t, updates[0].Power, int64(100+200)) // see MsgCreateValidator stake+NFT stake (MintToken)
		require.Equal(t, updates[1].Power, int64(100+100)) // see MsgCreateValidator stake+NFT stake (MintToken)
		require.Equal(t, updates[2].Power, int64(100))     // see MsgCreateValidator stake
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 4)

		newValidator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, newValidator1.Status)
		newValidator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, newValidator2.Status)
		newValidator3, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[2])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, newValidator2.Status)

		totalPower := newValidator1.ConsensusPower() + newValidator2.ConsensusPower() + newValidator3.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}

	// check pools
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins()))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake).Add(creatorStake).Add(creatorStake)))

		// check nft
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
	}

	////////////////////////////////////////////////
	// 4. increment block

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)

	////////////////////////////////////////////////
	// 5. first val to offline

	{
		// first val
		msgOffline := types.NewMsgSetOffline(vals[0])
		_, err = msgsrv.SetOffline(goCtx, msgOffline)
		require.NoError(t, err)
	}

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// check updates
	{
		require.Len(t, updates, 1)
		require.Equal(t, updates[0].Power, int64(0))
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 3)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator2.Status)
		validator3, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[2])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator3.Status)

		totalPower := validator2.ConsensusPower() + validator3.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake).Add(creatorStake)))

		// check nft
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
	}

	////////////////////////////////////////////////
	// 6. increment block

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)

	////////////////////////////////////////////////
	// 7. undelegate power from second val
	var completionTime time.Time
	{
		// second val
		ubd := types.NewMsgUndelegate(accs[1], vals[1], creatorStake)
		_, err = msgsrv.Undelegate(goCtx, ubd)
		require.NoError(t, err)

		ubdNft := types.NewMsgUndelegateNFT(accs[0], vals[1], tokenID, []uint32{6})
		resp, err := msgsrv.UndelegateNFT(goCtx, ubdNft)
		completionTime = resp.CompletionTime
		require.NoError(t, err)
	}
	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// check updates
	{
		require.Len(t, updates, 1)
		require.Equal(t, updates[0].Power, int64(0))
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonding, validator2.Status)
		validator3, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[2])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator3.Status)

		totalPower := validator3.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake).Add(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake)))

		// check nft
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
	}

	////////////////////////////////////////////////
	// 8. increment block

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(completionTime.Add(time.Hour * 2))
	goCtx = sdk.WrapSDKContext(ctx)

	////////////////////////////////////////////////
	// 9. unbond all mature validators

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 0)
	// check validators
	{
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator2.Status)
		validator3, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[2])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator3.Status)

		totalPower := validator3.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake)))

		// check nft
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, nbPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, accs[0].String())
	}

	////////////////////////////////////////////////
	// 10. increment block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(completionTime.Add(time.Hour * 2))
	goCtx = sdk.WrapSDKContext(ctx)

	////////////////////////////////////////////////
	// 11. redelegate

	// redelegate first -> third
	{
		// first val
		red := types.NewMsgRedelegate(accs[0], vals[0], vals[2], creatorStake)
		_, err = msgsrv.Redelegate(goCtx, red)
		require.NoError(t, err)

		redNft := types.NewMsgRedelegateNFT(accs[1], vals[0], vals[2], tokenID, []uint32{1, 2})
		_, err := msgsrv.RedelegateNFT(goCtx, redNft)
		require.NoError(t, err)
	}

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 0)

	////////////////////////////////////////////////
	// 11. after redelegate
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(completionTime.Add(time.Hour * 10))
	goCtx = sdk.WrapSDKContext(ctx)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 0)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(completionTime.Add(time.Hour * 11))
	goCtx = sdk.WrapSDKContext(ctx)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 1)
	// check validators
	{
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonding, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator2.Status)
		validator3, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[2])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator3.Status)

		totalPower := validator2.ConsensusPower() + validator3.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins()))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		// validator 3 stake + completed redelegation
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake).Add(creatorStake)))

		// check nft
		subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 2)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, bPool.String())
		subtoken, found = dsc.NFTKeeper.GetSubToken(ctx, tokenID, 6)
		require.True(t, found)
		require.Equal(t, subtoken.Owner, accs[0].String())
	}
}

func TestCheckDelegations(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	defaultParams := dsc.ValidatorKeeper.GetParams(ctx)
	defaultParams.MaxDelegations = 4
	dsc.ValidatorKeeper.SetParams(ctx, defaultParams)

	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	genesisVals := dsc.ValidatorKeeper.GetValidators(ctx, 10)
	require.Len(t, genesisVals, 1)
	genesisVal := genesisVals[0]
	require.True(t, genesisVal.ConsensusPower() > 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())

	// create custom coin
	ccDenom := "custom"
	initVolume := keeper.TokensFromConsensusPower(100000000000)
	initReserve := keeper.TokensFromConsensusPower(1000)
	limitVolume := keeper.TokensFromConsensusPower(100000000000000000)
	crr := uint64(50)

	_, err := dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, ""))
	require.NoError(t, err)
	// ----------------------------

	// create custom coin
	ccDenom2 := "custom2"

	initVolume2 := keeper.TokensFromConsensusPower(1000000)
	initReserve2 := keeper.TokensFromConsensusPower(1000)
	limitVolume2 := keeper.TokensFromConsensusPower(100000000000000000)

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[1], ccDenom2, "da", crr, initVolume2, initReserve2, limitVolume2, ""))
	require.NoError(t, err)
	// ----------------------------

	// create custom coin
	ccDenom3 := "custom3"

	initVolume3 := keeper.TokensFromConsensusPower(100000000)
	initReserve3 := keeper.TokensFromConsensusPower(1000)
	limitVolume3 := keeper.TokensFromConsensusPower(100000000000000000)

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[2], ccDenom3, "d", crr, initVolume3, initReserve3, limitVolume3, ""))
	require.NoError(t, err)
	// ----------------------------

	goCtx := sdk.WrapSDKContext(ctx)
	valK := dsc.ValidatorKeeper
	////////////////////////////////////////////////
	// 1. create two validators

	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	// first
	{
		msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)
	}
	// delegates
	stake1 := types.NewStakeCoin(sdk.NewCoin(ccDenom, helpers.EtherToWei(sdkmath.NewInt(1000))))
	stake2 := types.NewStakeCoin(sdk.NewCoin(ccDenom2, helpers.EtherToWei(sdkmath.NewInt(1200))))
	stake3 := types.NewStakeCoin(sdk.NewCoin(ccDenom3, helpers.EtherToWei(sdkmath.NewInt(1400))))
	stake4 := types.NewStakeCoin(creatorStake)
	stake5 := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000000))))
	{
		val, _ := valK.GetValidator(ctx, vals[0])

		{
			err = valK.Delegate(ctx, accs[0], val, stake1)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[1], val, stake2)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[2], val, stake3)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[2], val, stake4)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[1], val, stake5)
			require.NoError(t, err)
		}
	}

	{
		val, _ := valK.GetValidator(ctx, vals[0])

		dels := valK.GetAllDelegationsByValidator(ctx)
		require.Len(t, dels[val.GetOperator().String()], 6)

		valK.CheckDelegations(ctx, val)

		dels = valK.GetAllDelegationsByValidator(ctx)
		require.Len(t, dels[val.GetOperator().String()], 4)

		updatedRS, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		minus := keeper.TokensToConsensusPower(valK.ToBaseCoin(ctx, stake1.Stake).Amount)
		minus += keeper.TokensToConsensusPower(valK.ToBaseCoin(ctx, stake3.Stake).Amount)
		require.Equal(t, val.Stake-minus, updatedRS.Stake)
	}

}

func TestPoolBreakingByCancelRedelegation(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	startBondedBalance := dsc.BankKeeper.GetBalance(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress(), cmdcfg.BaseDenom)

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	stake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1000)))
	halfstake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(500)))
	cancelstake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))

	// 1. create validators
	goCtx := sdk.WrapSDKContext(ctx)

	msgCreate, err := types.NewMsgCreateValidator(
		vals[1],
		accs[1],
		PKs[0],
		types.Description{Moniker: "monik1"},
		sdk.ZeroDec(),
		stake,
	)
	require.NoError(t, err)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)

	msgCreate, err = types.NewMsgCreateValidator(
		vals[2],
		accs[2],
		PKs[1],
		types.Description{Moniker: "monik2"},
		sdk.ZeroDec(),
		stake,
	)
	require.NoError(t, err)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)

	// 2. set online
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[1]))
	require.NoError(t, err)
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[2]))
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// 3. add redelegation
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)

	_, err = msgsrv.Redelegate(goCtx, types.NewMsgRedelegate(
		accs[1],
		vals[1],
		vals[2],
		halfstake,
	))
	require.NoError(t, err)

	height := ctx.BlockHeight()
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// 4. set one validator offline + cancel redelegation
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOffline(goCtx, types.NewMsgSetOffline(vals[1]))
	require.NoError(t, err)

	_, err = msgsrv.CancelRedelegation(goCtx, types.NewMsgCancelRedelegation(
		accs[1],
		vals[1],
		vals[2],
		height,
		cancelstake,
	))
	require.NoError(t, err)
	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// check pool balances
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	bondedBalance := dsc.BankKeeper.GetBalance(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress(), cmdcfg.BaseDenom)
	notBondedBalance := dsc.BankKeeper.GetBalance(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress(), cmdcfg.BaseDenom)

	// stake of vals[1] (1000) (online)
	require.True(t, bondedBalance.IsEqual(startBondedBalance.Add(stake)), "bonded: %s, expect: %s", bondedBalance, stake)
	// stake of vals[0] (500+100 cancel) + 400 in redelegation
	require.True(t, notBondedBalance.IsEqual(stake), "not bonded: %s, expect: %s", notBondedBalance, stake)
}

func TestValidatorsCandidates(t *testing.T) {
	// candidates - bonded online validators with small stake
	// they are out of top X validators, they are out of tendermint validators

	const maxValidators = 10

	powerByPos := func(i int) int64 {
		return int64(maxValidators*2 + 1 - i)
	}

	_, dsc, ctx := createTestInput(t)
	genesisValidators := dsc.ValidatorKeeper.GetAllValidators(ctx)

	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.MaxValidators = maxValidators
	dsc.ValidatorKeeper.SetParams(ctx, params)

	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	///////////////////////////////
	// 0. start: validator count is maxValidators * 2
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(10000000)))
	accs, vals := generateAddresses(dsc, ctx, maxValidators*2, sdk.NewCoins(balance))

	goCtx := sdk.WrapSDKContext(ctx)
	for i := 0; i < maxValidators*2; i++ {
		msgCreate, err := types.NewMsgCreateValidator(
			vals[i],
			accs[i],
			PKs[i],
			types.Description{Moniker: "monik"},
			sdk.ZeroDec(),
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(powerByPos(i)))),
		)
		require.NoError(t, err)
		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)
		_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[i]))
		require.NoError(t, err)
	}
	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, maxValidators-len(genesisValidators))

	// check powers and statuses
	expectTotalPower := int64(0)
	for _, gv := range genesisValidators {
		expectTotalPower += gv.Stake
	}
	for i := 0; i < maxValidators-len(genesisValidators); i++ {
		expectTotalPower += powerByPos(i)
	}
	totalPower := dsc.ValidatorKeeper.GetLastTotalPower(ctx)
	require.Equal(t, expectTotalPower, totalPower.Int64())

	for i := 0; i < maxValidators*2; i++ {
		val, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[i])
		require.True(t, found)
		require.True(t, val.Online)
		require.True(t, val.IsBonded())
		require.True(t, val.Stake > 0)
	}

	///////////////////////////////
	// 1. turn off top validator
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	goCtx = sdk.WrapSDKContext(ctx)
	_, err := msgsrv.SetOffline(goCtx, types.NewMsgSetOffline(vals[0]))
	require.NoError(t, err)

	// 1 validator out
	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 1)

	// 1 validator from candidates in
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))
	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 1)

	// check powers and statuses
	expectTotalPower = int64(0)
	for _, gv := range genesisValidators {
		expectTotalPower += gv.Stake
	}
	for i := 0; i < maxValidators-len(genesisValidators); i++ {
		expectTotalPower += powerByPos(i + 1)
	}
	totalPower = dsc.ValidatorKeeper.GetLastTotalPower(ctx)
	require.Equal(t, expectTotalPower, totalPower.Int64())

	val, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	require.False(t, val.Online)
	require.True(t, val.IsUnbonded())
	require.True(t, val.Stake == 0)
}

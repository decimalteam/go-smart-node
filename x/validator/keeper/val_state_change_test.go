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
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)

	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)
	// delegate NFT
	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(accs[1], vals[0], tokenID, []uint32{1}))
	require.NoError(t, err)
	subtoken, found := dsc.NFTKeeper.GetSubToken(ctx, tokenID, 1)
	require.True(t, found)
	require.Equal(t, subtoken.Owner, nbPool.String())

	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// new validator is not online, there is not changes in tendermint validators and powers
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
	msgOnline := types.NewMsgSetOnline(vals[0])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOnline(goCtx, msgOnline)
	require.NoError(t, err)
	// last validators must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(100+100)) // see MsgCreateValidator stake+NFT stake (MintToken)
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)
	newValidator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
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
	msgOffline := types.NewMsgSetOffline(vals[0])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOffline(goCtx, msgOffline)
	require.NoError(t, err)
	// last validator must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(0)) // 0 mean 'remove from validators'
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)
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
	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// new validator is not online, there is not changes in tendermint validators and powers
	require.Len(t, updates, 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	// check balance
	{
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake).Add(creatorStake)))
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
	}

	// last validators must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// check updates
	{
		require.Len(t, updates, 2)
		require.Equal(t, updates[0].Power, int64(100+200)) // see MsgCreateValidator stake+NFT stake (MintToken)
		require.Equal(t, updates[1].Power, int64(100+100)) // see MsgCreateValidator stake+NFT stake (MintToken)
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 3)

		newValidator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, newValidator1.Status)
		newValidator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, newValidator2.Status)

		totalPower := newValidator1.ConsensusPower() + newValidator2.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}

	// check pools
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins()))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
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
		msgOnline := types.NewMsgSetOffline(vals[0])
		_, err = msgsrv.SetOffline(goCtx, msgOnline)
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
		require.Equal(t, types.BondStatus_Unbonding, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator2.Status)

		totalPower := validator2.ConsensusPower() + genesisVal.ConsensusPower()
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
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonding, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonding, validator2.Status)

		totalPower := genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake).Add(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB))

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
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator2.Status)

		totalPower := genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB))

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

	// redelegate first -> second
	{
		// first val
		red := types.NewMsgRedelegate(accs[0], vals[0], vals[1], creatorStake)
		_, err = msgsrv.Redelegate(goCtx, red)
		require.NoError(t, err)

		redNft := types.NewMsgRedelegateNFT(accs[1], vals[0], vals[1], tokenID, []uint32{1, 2})
		_, err := msgsrv.RedelegateNFT(goCtx, redNft)
		require.NoError(t, err)
	}

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 0)

	////////////////////////////////////////////////
	// 11. after redelegate
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(completionTime.Add(time.Hour * 2))
	goCtx = sdk.WrapSDKContext(ctx)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	require.Len(t, updates, 1)
	// check validators
	{
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

		validator1, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Unbonded, validator1.Status)
		validator2, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[1])
		require.True(t, found)
		require.Equal(t, types.BondStatus_Bonded, validator2.Status)

		totalPower := validator2.ConsensusPower() + genesisVal.ConsensusPower()
		require.Equal(t, dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64(), totalPower)
	}
	// check pool
	{
		// check balance
		balanceNB = dsc.BankKeeper.GetAllBalances(ctx, nbPool)
		require.True(t, balanceNB.IsEqual(sdk.NewCoins()))
		balanceB = dsc.BankKeeper.GetAllBalances(ctx, bPool)
		require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake)))

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
	defaultParams.MaxDelegations = 3
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

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[1], ccDenom2, "da", crr, initVolume, initReserve, limitVolume, ""))
	require.NoError(t, err)
	// ----------------------------

	// create custom coin
	ccDenom3 := "custom3"

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[2], ccDenom3, "d", crr, initVolume, initReserve, limitVolume, ""))
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
	{
		val, _ := valK.GetValidator(ctx, vals[0])

		{
			err = valK.Delegate(ctx, accs[0], val, stake1)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[1], val, stake2)
			require.NoError(t, err)
			err = valK.Delegate(ctx, accs[2], val, stake3)
			require.NoError(t, err)
		}
	}

	{
		val, _ := valK.GetValidator(ctx, vals[0])

		dels := valK.GetAllDelegationsByValidator(ctx)
		require.Len(t, dels[val.GetOperator().String()], 4)

		valK.CheckDelegations(ctx, val, dels[val.GetOperator().String()])

		dels = valK.GetAllDelegationsByValidator(ctx)
		require.Len(t, dels[val.GetOperator().String()], 3)

		updatedRS, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		minus := keeper.TokensToConsensusPower(valK.ToBaseCoin(ctx, stake1.Stake).Amount)
		require.Equal(t, val.Stake-minus, updatedRS.Stake)
	}

}

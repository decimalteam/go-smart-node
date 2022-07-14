package keeper

//var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
//
//const (
//	DefaultBondDenom = "del"
//)
//
//// Hogpodge of all sorts of input required for testing.
//// `initPower` is converted to an amount of tokens.
//// If `initPower` is 0, no addrs get created.
//func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, Keeper) {
//	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
//	keyParams := sdk.NewKVStoreKey(params.StoreKey)
//	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
//	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
//	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
//
//	db := dbm.NewMemDB()
//	ms := store.NewCommitMultiStore(db)
//	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
//	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
//	err := ms.LoadLatestVersion()
//	require.Nil(t, err)
//
//	_config := sdk.GetConfig()
//	_config.SetCoinType(60)
//	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
//	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
//	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
//	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
//
//	ctx := sdk.NewContext(ms, abcitypes.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
//	ctx = ctx.WithConsensusParams(
//		&abcitypes.ConsensusParams{
//			Validator: &abcitypes.ValidatorParams{
//				PubKeyTypes: []string{types3.ABCIPubKeyTypeEd25519},
//			},
//		},
//	)
//	cdc := MakeTestCodec()
//
//	initCoins := sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(100000000000)))
//
//	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)
//
//	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
//
//	initTokens := sdk.NewInt(initPower).Mul(PowerReduction)
//
//	accountKeeper := auth.NewAccountKeeper(
//		cdc,    // amino codec
//		keyAcc, // target store
//		pk.Subspace(auth.DefaultParamspace),
//		auth.ProtoBaseAccount, // prototype
//	)
//
//	bk := bank.NewBaseKeeper(
//		accountKeeper,
//		pk.Subspace(bank.DefaultParamspace),
//		nil,
//	)
//
//	maccPerms := map[string][]string{
//		auth.FeeCollectorName: nil,
//		nftTypes.ReservedPool: {supply.Burner},
//	}
//	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)
//
//	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, supplyKeeper, config.GetDefaultConfig(config.ChainID))
//
//	coinConfig := config.GetDefaultConfig("decimal")
//	coinKeeper.SetCoin(ctx, coin.Coin{
//		Title:  coinConfig.TitleBaseCoin,
//		Symbol: coinConfig.SymbolBaseCoin,
//		Volume: coinConfig.InitialVolumeBaseCoin,
//	})
//	nftkeeper := NewKeeper(cdc, keyCoin, supplyKeeper, DefaultBondDenom)
//
//	totalSupply := sdk.NewCoins(sdk.NewCoin(DefaultBondDenom, initTokens.MulRaw(int64(len(nftTypes.Addrs)))))
//	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
//	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
//
//	// fill all the Addrs with some coins, set the loose pool tokens simultaneously
//	for _, addr := range nftTypes.Addrs {
//		_, err := bk.AddCoins(ctx, addr, initCoins)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	return ctx, nftkeeper
//}
//
//// MakeTestCodec creates a codec for testing
//func MakeTestCodec() *codec.Codec {
//	var cdc = codec.New()
//
//	// Register Msgs
//	cdc.RegisterInterface((*exported.NFT)(nil), nil)
//	cdc.RegisterInterface((*exported.TokenOwners)(nil), nil)
//	cdc.RegisterInterface((*exported.TokenOwner)(nil), nil)
//	cdc.RegisterConcrete(&nftTypes.BaseNFT{}, "nft/BaseNFT", nil)
//	cdc.RegisterConcrete(&nftTypes.IDCollection{}, "nft/IDCollection", nil)
//	cdc.RegisterConcrete(&nftTypes.Collection{}, "nft/Collection", nil)
//	cdc.RegisterConcrete(&nftTypes.Owner{}, "nft/Owner", nil)
//	cdc.RegisterConcrete(&nftTypes.TokenOwner{}, "nft/TokenOwner", nil)
//	cdc.RegisterConcrete(&nftTypes.TokenOwners{}, "nft/TokenOwners", nil)
//	cdc.RegisterConcrete(nftTypes.MsgTransferNFT{}, "nft/msg_transfer", nil)
//	cdc.RegisterConcrete(nftTypes.MsgEditNFTMetadata{}, "nft/msg_edit_metadata", nil)
//	cdc.RegisterConcrete(nftTypes.MsgMintNFT{}, "nft/msg_mint", nil)
//	cdc.RegisterConcrete(nftTypes.MsgBurnNFT{}, "nft/msg_burn", nil)
//
//	// Register AppAccount
//	cdc.RegisterInterface((*exported2.Account)(nil), nil)
//	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/coin/base_account", nil)
//	supply.RegisterCodec(cdc)
//	codec.RegisterCrypto(cdc)
//
//	return cdc
//}

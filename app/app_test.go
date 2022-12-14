package app

import (
	"testing"
)

func TestDSCExport(t *testing.T) {
	//db := dbm.NewMemDB()
	//simapp.SetupOptions{
	//	Logger:             nil,
	//	DB:                 nil,
	//	InvCheckPeriod:     0,
	//	HomePath:           "",
	//	SkipUpgradeHeights: nil,
	//	EncConfig:          simappparams.EncodingConfig{},
	//	AppOpts:            types.AppOptions(),
	//}
	//app := NewDSC(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	//
	//gs := GenesisStateWithSingleValidator(t, app)
	//stateBytes, err := json.MarshalIndent(gs, "", "  ")
	//require.NoError(t, err)
	//
	//// Initialize the chain
	//app.InitChain(
	//	abci.RequestInitChain{
	//		ChainId:       MainnetChainIDPrefix + "-1",
	//		Validators:    []abci.ValidatorUpdate{},
	//		AppStateBytes: stateBytes,
	//	},
	//)
	//app.Commit()
	//
	//// Making a new app object with the db, so that initchain hasn't been called
	//app2 := NewDSC(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	//_, err = app2.ExportAppStateAndValidators(false, []string{})
	//require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

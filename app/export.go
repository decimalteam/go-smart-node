package app

import (
	"encoding/json"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	validator "bitbucket.org/decimalteam/go-smart-node/x/validator"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *DSC) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string,
) (servertypes.ExportedApp, error) {
	// Creates context with current height and checks txs for ctx to be usable by start of next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0

		//if err := app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs); err != nil {
		//	return servertypes.ExportedApp{}, err
		//}
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := validator.WriteValidators(ctx, app.ValidatorKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

package keeper

import (
	v2 "bitbucket.org/decimalteam/go-smart-node/x/swap/migrations/v2"
	v3 "bitbucket.org/decimalteam/go-smart-node/x/swap/migrations/v3"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var _ module.MigrationHandler = Migrator{}.Migrate1to2

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

// Migrate1to2 migrates from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.UpdateParams(ctx, &m.keeper.paramSpace)
}

// Migrate2to3 migrates from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.UpdateParams(ctx, &m.keeper.paramSpace)
}

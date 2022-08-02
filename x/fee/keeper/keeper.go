package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.BinaryCodec // The amino codec for binary encoding/decoding.

	bankKeeper keeper.Keeper

	BaseDenom *string
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey, bankKeeper keeper.Keeper, baseDenom string) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		bankKeeper: bankKeeper,
		BaseDenom:  &baseDenom,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SavePrice(
	ctx sdk.Context,
	price float64,
) error {
	fmt.Println(price)

	return nil
}

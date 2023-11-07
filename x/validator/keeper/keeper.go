package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

//var _ types.ValidatorSet = Keeper{}
//var _ types.DelegationSet = Keeper{}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine.
type Keeper struct {
	storeKey       storetypes.StoreKey
	cdc            codec.BinaryCodec
	authKeeper     types.AccountKeeper
	bankKeeper     types.BankKeeper
	nftKeeper      types.NFTKeeper
	coinKeeper     types.CoinKeeper
	multisigKeeper types.MultisigKeeper
	hooks          types.StakingHooks
	paramstore     paramtypes.Subspace
}

// NewKeeper creates new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NFTKeeper,
	ck types.CoinKeeper,
	mk types.MultisigKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := ak.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		authKeeper:     ak,
		bankKeeper:     bk,
		nftKeeper:      nk,
		coinKeeper:     ck,
		multisigKeeper: mk,
		paramstore:     ps,
		hooks:          nil,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetLastTotalPower loads the last total validators power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastTotalPowerKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)
	return ip.Int
}

// SetLastTotalPower sets the last total validators power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdkmath.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: power})
	store.Set(types.GetLastTotalPowerKey(), bz)
}

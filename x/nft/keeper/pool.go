package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// func (k *Keeper) GetReservedPool(ctx sdk.Context) exported.ModuleAccountI {
// 	return k.supplyKeeper.GetModuleAccount(ctx, types.ReservedPool)
// }

func (k *Keeper) ReserveTokens(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ReservedPool, amount)
}

func (k *Keeper) ReturnTokensTo(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ReservedPool, address, amount)
}

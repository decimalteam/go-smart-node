package v3

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
)

func UpdateParams(ctx sdk.Context, paramstore *paramtypes.Subspace) error {
	ctx.Logger().Info("Upgrade is successful: congrats people!")
	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}

	newTimeOut := time.Hour * 30
	newTimeIn := time.Hour * 13

	paramstore.Set(ctx, types.KeyLockedTimeOut, newTimeOut)
	paramstore.Set(ctx, types.KeyLockedTimeIn, newTimeIn)
	return nil
}

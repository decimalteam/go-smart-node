package v3

import (
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"time"
)

func UpdateParams(ctx sdk.Context, paramstore *paramtypes.Subspace) error {
	fmt.Println("hiiiiii2")
	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}

	newTimeOut := time.Hour * 30
	newTimeIn := time.Hour * 15

	paramstore.Set(ctx, types.KeyLockedTimeOut, newTimeOut)
	paramstore.Set(ctx, types.KeyLockedTimeIn, newTimeIn)
	return nil
}

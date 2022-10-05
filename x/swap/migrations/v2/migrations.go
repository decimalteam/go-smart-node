package v2

import (
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"time"
)

func UpdateParams(ctx sdk.Context, paramstore *paramtypes.Subspace) error {
	fmt.Println("hiiiiii")
	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}

	newTimeOut := time.Hour * 25
	newTimeIn := time.Hour * 13

	paramstore.Set(ctx, types.KeyLockedTimeOut, newTimeOut)
	paramstore.Set(ctx, types.KeyLockedTimeIn, newTimeIn)
	return nil
}

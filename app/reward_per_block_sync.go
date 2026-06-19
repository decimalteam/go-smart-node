package app

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// SyncRewardPerBlockOverride reads the per-block reward override from the
// master-validator contract (getRewardPerBlock) and mirrors it into validator
// module state (sync contract -> node). This seeds any value set on the contract
// before the node logic shipped; afterwards the EVM hook keeps the two in sync.
// It returns the synced amount and whether the override is enabled.
func SyncRewardPerBlockOverride(ctx sdk.Context, app *DSC) (sdkmath.Int, bool, error) {
	validatorHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressMasterValidator)
	if err != nil {
		return sdk.ZeroInt(), false, fmt.Errorf("resolve master-validator contract: %w", err)
	}
	contractAddr := common.HexToAddress(validatorHex)
	if contractAddr == (common.Address{}) {
		return sdk.ZeroInt(), false, fmt.Errorf("master-validator contract address not resolved")
	}

	validatorAbi, err := validator.ValidatorMetaData.GetAbi()
	if err != nil {
		return sdk.ZeroInt(), false, fmt.Errorf("load master-validator abi: %w", err)
	}

	const methodCall = "getRewardPerBlock"
	res, err := app.EvmKeeper.CallEVM(ctx, *validatorAbi, common.Address(validatortypes.ModuleAddress), contractAddr, false, methodCall)
	if err != nil {
		return sdk.ZeroInt(), false, fmt.Errorf("call %s: %w", methodCall, err)
	}

	out, err := validatorAbi.Unpack(methodCall, res.Ret)
	if err != nil {
		return sdk.ZeroInt(), false, fmt.Errorf("unpack %s: %w", methodCall, err)
	}
	if len(out) != 2 {
		return sdk.ZeroInt(), false, fmt.Errorf("unexpected %s return arity: %d", methodCall, len(out))
	}

	rewardPerBlock, ok := out[0].(*big.Int)
	if !ok {
		return sdk.ZeroInt(), false, fmt.Errorf("unexpected rewardPerBlock type %T", out[0])
	}
	enabled, ok := out[1].(bool)
	if !ok {
		return sdk.ZeroInt(), false, fmt.Errorf("unexpected enabled type %T", out[1])
	}

	amount := sdkmath.NewIntFromBigInt(rewardPerBlock)
	if enabled {
		app.ValidatorKeeper.SetRewardPerBlockOverride(ctx, amount)
	} else {
		app.ValidatorKeeper.ClearRewardPerBlockOverride(ctx)
	}

	return amount, enabled, nil
}

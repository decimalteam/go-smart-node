package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// EmergencyStatus queries the current emergency hard-halt and soft-freeze state.
func (k Querier) EmergencyStatus(c context.Context, req *types.QueryEmergencyStatusRequest) (*types.QueryEmergencyStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Missing records are reported as inactive/unfrozen (zero values).
	halt, _ := k.GetHaltInfo(ctx)
	freeze, _ := k.GetFreezeInfo(ctx)

	return &types.QueryEmergencyStatusResponse{
		Halt:   halt,
		Freeze: freeze,
	}, nil
}

package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (p Plan) String() string {
	due := p.DueAt()
	downloadedAt := p.DownloadAt()
	return fmt.Sprintf(`Upgrade Plan
  Name: %s
  %s
  %s
  Info: %s.`, p.Name, due, downloadedAt, p.Info)
}

// ValidateBasic does basic validation of a Plan
func (p Plan) ValidateBasic() error {
	if !p.Time.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("time-based upgrades have been deprecated in the SDK")
	}
	if p.UpgradedClientState != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("upgrade logic for IBC has been moved to the IBC module")
	}
	if len(p.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	if p.Height <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "height must be greater than 0")
	}

	return nil
}

// ShouldExecute returns true if the Plan is ready to execute given the current context
func (p Plan) ShouldExecute(ctx sdk.Context) bool {
	if p.Height > 0 {
		return p.Height <= ctx.BlockHeight()
	}
	return false
}

// DueAt is a string representation of when this plan is due to be executed
func (p Plan) DueAt() string {
	return fmt.Sprintf("height: %d", p.Height)
}

// DownloadAt is a string representation of when this plan binary file is downloaded
func (p Plan) DownloadAt() string {
	return fmt.Sprintf("download at: %d", p.ToDownload)
}

func (p Plan) Mapping() map[string][]string {
	var mapping map[string][]string
	err := json.Unmarshal([]byte(p.Info), &mapping)
	if err != nil {
		return nil
	}
	return mapping
}

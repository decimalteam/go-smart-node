package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc references the global validator module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to modules/validator and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())

// RegisterInterfaces registers concrete implementations of specific interfaces.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateValidator{},
		&MsgEditValidator{},
		&MsgSetOnline{},
		&MsgSetOffline{},
		&MsgDelegate{},
		&MsgDelegateNFT{},
		&MsgRedelegate{},
		&MsgRedelegateNFT{},
		&MsgUndelegate{},
		&MsgUndelegateNFT{},
		&MsgCancelRedelegation{},
		&MsgCancelRedelegationNFT{},
		&MsgCancelUndelegation{},
		&MsgCancelUndelegationNFT{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

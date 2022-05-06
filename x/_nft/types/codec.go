package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc references the global coin module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to modules/coin and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())

// RegisterInterfaces registers concrete implementations of specific inferfaces.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgBurnNFT{},
		&MsgMintNFT{},
		&MsgEditNFTMetadata{},
		&MsgTransferNFT{},
		&MsgUpdateReserveNFT{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

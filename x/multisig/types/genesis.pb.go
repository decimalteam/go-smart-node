// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: decimal/multisig/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the module's genesis state.
type GenesisState struct {
	// wallets defines all registered multisig wallets.
	Wallets []Wallet `protobuf:"bytes,1,rep,name=wallets,proto3" json:"wallets"`
	// transactions defines all registered multisig transactions.
	Transactions []GenesisTransaction `protobuf:"bytes,2,rep,name=transactions,proto3" json:"transactions"`
	// params defines all the module's parameters.
	Params Params `protobuf:"bytes,3,opt,name=params,proto3" json:"params"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_6405f2144b1d3117, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetWallets() []Wallet {
	if m != nil {
		return m.Wallets
	}
	return nil
}

func (m *GenesisState) GetTransactions() []GenesisTransaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// GenesisTransaction defines old multisig transaction (only send coins).
type GenesisTransaction struct {
	Id        string                                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Wallet    string                                   `protobuf:"bytes,2,opt,name=wallet,proto3" json:"wallet,omitempty"`
	Receiver  string                                   `protobuf:"bytes,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Coins     github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,4,rep,name=coins,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"coins"`
	Signers   []string                                 `protobuf:"bytes,5,rep,name=signers,proto3" json:"signers,omitempty"`
	CreatedAt int64                                    `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (m *GenesisTransaction) Reset()         { *m = GenesisTransaction{} }
func (m *GenesisTransaction) String() string { return proto.CompactTextString(m) }
func (*GenesisTransaction) ProtoMessage()    {}
func (*GenesisTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_6405f2144b1d3117, []int{1}
}
func (m *GenesisTransaction) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisTransaction.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisTransaction.Merge(m, src)
}
func (m *GenesisTransaction) XXX_Size() int {
	return m.Size()
}
func (m *GenesisTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisTransaction proto.InternalMessageInfo

func (m *GenesisTransaction) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GenesisTransaction) GetWallet() string {
	if m != nil {
		return m.Wallet
	}
	return ""
}

func (m *GenesisTransaction) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *GenesisTransaction) GetCoins() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Coins
	}
	return nil
}

func (m *GenesisTransaction) GetSigners() []string {
	if m != nil {
		return m.Signers
	}
	return nil
}

func (m *GenesisTransaction) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "decimal.multisig.v1.GenesisState")
	proto.RegisterType((*GenesisTransaction)(nil), "decimal.multisig.v1.GenesisTransaction")
}

func init() { proto.RegisterFile("decimal/multisig/v1/genesis.proto", fileDescriptor_6405f2144b1d3117) }

var fileDescriptor_6405f2144b1d3117 = []byte{
	// 469 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0x63, 0xa7, 0x4d, 0xe9, 0x52, 0x71, 0x58, 0x7a, 0xd8, 0x16, 0xe1, 0x9a, 0x5c, 0xf0,
	0x25, 0xbb, 0x0d, 0x70, 0x00, 0x71, 0x6a, 0x38, 0x70, 0x05, 0x17, 0x09, 0x89, 0x4b, 0xb5, 0xb6,
	0x47, 0x66, 0xd5, 0xd8, 0x1b, 0xed, 0x4c, 0x02, 0xbc, 0x05, 0x8f, 0x81, 0x38, 0xf3, 0x10, 0x3d,
	0x56, 0x9c, 0xe0, 0x02, 0x28, 0x79, 0x11, 0x64, 0x7b, 0xd3, 0x0a, 0x11, 0x38, 0xd9, 0xb3, 0xf3,
	0xcd, 0x3f, 0xff, 0xee, 0x0c, 0xbb, 0x57, 0x40, 0x6e, 0x2a, 0x3d, 0x55, 0xd5, 0x7c, 0x4a, 0x06,
	0x4d, 0xa9, 0x16, 0x63, 0x55, 0x42, 0x0d, 0x68, 0x50, 0xce, 0x9c, 0x25, 0xcb, 0x6f, 0x7b, 0x44,
	0xae, 0x11, 0xb9, 0x18, 0x1f, 0xee, 0x97, 0xb6, 0xb4, 0x6d, 0x5e, 0x35, 0x7f, 0x1d, 0x7a, 0x78,
	0x90, 0x5b, 0xac, 0x2c, 0x9e, 0x75, 0x89, 0x2e, 0xf0, 0xa9, 0xa8, 0x8b, 0x54, 0xa6, 0x11, 0xd4,
	0x62, 0x9c, 0x01, 0xe9, 0xb1, 0xca, 0xad, 0xa9, 0x7d, 0x7e, 0xb8, 0xc9, 0xc8, 0x55, 0xc7, 0x8e,
	0x89, 0x37, 0x31, 0x33, 0xed, 0x74, 0xe5, 0xbb, 0x0c, 0xbf, 0x07, 0x6c, 0xef, 0x79, 0xe7, 0xfe,
	0x94, 0x34, 0x01, 0x7f, 0xca, 0x76, 0xde, 0xe9, 0xe9, 0x14, 0x08, 0x45, 0x10, 0xf7, 0x93, 0x9b,
	0x0f, 0xee, 0xc8, 0x0d, 0xd7, 0x91, 0xaf, 0x5b, 0x66, 0xb2, 0x75, 0xf1, 0xe3, 0xa8, 0x97, 0xae,
	0x2b, 0xf8, 0x4b, 0xb6, 0x47, 0x4e, 0xd7, 0xa8, 0x73, 0x32, 0xb6, 0x46, 0x11, 0xb6, 0x0a, 0xf7,
	0x37, 0x2a, 0xf8, 0xae, 0xaf, 0xae, 0x79, 0xaf, 0xf6, 0x87, 0x04, 0x7f, 0xc2, 0x06, 0x9d, 0x61,
	0xd1, 0x8f, 0x83, 0x7f, 0xda, 0x79, 0xd1, 0x22, 0x5e, 0xc0, 0x17, 0x0c, 0x3f, 0x85, 0x8c, 0xff,
	0xdd, 0x85, 0xdf, 0x62, 0xa1, 0x29, 0x44, 0x10, 0x07, 0xc9, 0x6e, 0x1a, 0x9a, 0x82, 0x1f, 0xb3,
	0x41, 0xe7, 0x5f, 0x84, 0xcd, 0xd9, 0x44, 0x7c, 0xfd, 0x32, 0xda, 0xf7, 0xa3, 0x38, 0x29, 0x0a,
	0x07, 0x88, 0xa7, 0xe4, 0x4c, 0x5d, 0xa6, 0x9e, 0xe3, 0x8f, 0xd8, 0x0d, 0x07, 0x39, 0x98, 0x05,
	0xb8, 0xd6, 0xd5, 0xff, 0x6a, 0xae, 0x48, 0xae, 0xd9, 0x76, 0x33, 0x3e, 0x14, 0x5b, 0xed, 0xab,
	0x1c, 0x48, 0xcf, 0x37, 0x03, 0x96, 0x7e, 0xc0, 0xf2, 0x99, 0x35, 0xf5, 0xe4, 0xb8, 0xb9, 0xc6,
	0xe7, 0x9f, 0x47, 0x49, 0x69, 0xe8, 0xed, 0x3c, 0x93, 0xb9, 0xad, 0xfc, 0x6e, 0xf8, 0xcf, 0x08,
	0x8b, 0x73, 0x45, 0x1f, 0x66, 0x80, 0x6d, 0x01, 0xa6, 0x9d, 0x32, 0x17, 0x6c, 0x07, 0x4d, 0x59,
	0x83, 0x43, 0xb1, 0x1d, 0xf7, 0x93, 0xdd, 0x74, 0x1d, 0xf2, 0xbb, 0x8c, 0xe5, 0x0e, 0x34, 0x41,
	0x71, 0xa6, 0x49, 0x0c, 0xe2, 0x20, 0xe9, 0xa7, 0xbb, 0xfe, 0xe4, 0x84, 0x26, 0xe9, 0xc5, 0x32,
	0x0a, 0x2e, 0x97, 0x51, 0xf0, 0x6b, 0x19, 0x05, 0x1f, 0x57, 0x51, 0xef, 0x72, 0x15, 0xf5, 0xbe,
	0xad, 0xa2, 0xde, 0x9b, 0xc7, 0x99, 0xa1, 0x6c, 0x9e, 0x9f, 0x03, 0x49, 0xeb, 0x4a, 0xe5, 0x1f,
	0x9f, 0x40, 0x57, 0xaa, 0xb4, 0x23, 0xac, 0xb4, 0xa3, 0x51, 0x6d, 0x0b, 0x50, 0xef, 0xaf, 0x97,
	0xac, 0x75, 0x96, 0x0d, 0xda, 0x0d, 0x7b, 0xf8, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x4e, 0x1f,
	0xd9, 0x32, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Transactions) > 0 {
		for iNdEx := len(m.Transactions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Transactions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Wallets) > 0 {
		for iNdEx := len(m.Wallets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Wallets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *GenesisTransaction) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisTransaction) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisTransaction) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CreatedAt != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Signers) > 0 {
		for iNdEx := len(m.Signers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Signers[iNdEx])
			copy(dAtA[i:], m.Signers[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.Signers[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Coins) > 0 {
		for iNdEx := len(m.Coins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Coins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Wallet) > 0 {
		i -= len(m.Wallet)
		copy(dAtA[i:], m.Wallet)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Wallet)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Wallets) > 0 {
		for _, e := range m.Wallets {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Transactions) > 0 {
		for _, e := range m.Transactions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *GenesisTransaction) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Wallet)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Coins) > 0 {
		for _, e := range m.Coins {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Signers) > 0 {
		for _, s := range m.Signers {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.CreatedAt != 0 {
		n += 1 + sovGenesis(uint64(m.CreatedAt))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Wallets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Wallets = append(m.Wallets, Wallet{})
			if err := m.Wallets[len(m.Wallets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transactions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Transactions = append(m.Transactions, GenesisTransaction{})
			if err := m.Transactions[len(m.Transactions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GenesisTransaction) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisTransaction: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisTransaction: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Wallet", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Wallet = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Coins = append(m.Coins, types.Coin{})
			if err := m.Coins[len(m.Coins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signers", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signers = append(m.Signers, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibcswap/interchainswap/market.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
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

type PoolSide int32

const (
	PoolSide_NATIVE PoolSide = 0
	PoolSide_REMOTE PoolSide = 1
)

var PoolSide_name = map[int32]string{
	0: "NATIVE",
	1: "REMOTE",
}

var PoolSide_value = map[string]int32{
	"NATIVE": 0,
	"REMOTE": 1,
}

func (x PoolSide) String() string {
	return proto.EnumName(PoolSide_name, int32(x))
}

func (PoolSide) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_11407e7be1a65d42, []int{0}
}

type PoolStatus int32

const (
	PoolStatus_POOL_STATUS_INITIAL PoolStatus = 0
	PoolStatus_POOL_STATUS_READY   PoolStatus = 1
)

var PoolStatus_name = map[int32]string{
	0: "POOL_STATUS_INITIAL",
	1: "POOL_STATUS_READY",
}

var PoolStatus_value = map[string]int32{
	"POOL_STATUS_INITIAL": 0,
	"POOL_STATUS_READY":   1,
}

func (x PoolStatus) String() string {
	return proto.EnumName(PoolStatus_name, int32(x))
}

func (PoolStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_11407e7be1a65d42, []int{1}
}

type PoolAsset struct {
	Side    PoolSide    `protobuf:"varint,1,opt,name=side,proto3,enum=ibcswap.v4.interchainswap.PoolSide" json:"side,omitempty"`
	Balance *types.Coin `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Weight  uint32      `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Decimal uint32      `protobuf:"varint,4,opt,name=decimal,proto3" json:"decimal,omitempty"`
}

func (m *PoolAsset) Reset()         { *m = PoolAsset{} }
func (m *PoolAsset) String() string { return proto.CompactTextString(m) }
func (*PoolAsset) ProtoMessage()    {}
func (*PoolAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_11407e7be1a65d42, []int{0}
}
func (m *PoolAsset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolAsset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolAsset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolAsset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolAsset.Merge(m, src)
}
func (m *PoolAsset) XXX_Size() int {
	return m.Size()
}
func (m *PoolAsset) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolAsset.DiscardUnknown(m)
}

var xxx_messageInfo_PoolAsset proto.InternalMessageInfo

func (m *PoolAsset) GetSide() PoolSide {
	if m != nil {
		return m.Side
	}
	return PoolSide_NATIVE
}

func (m *PoolAsset) GetBalance() *types.Coin {
	if m != nil {
		return m.Balance
	}
	return nil
}

func (m *PoolAsset) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *PoolAsset) GetDecimal() uint32 {
	if m != nil {
		return m.Decimal
	}
	return 0
}

type InterchainLiquidityPool struct {
	PoolId                string       `protobuf:"bytes,1,opt,name=poolId,proto3" json:"poolId,omitempty"`
	Assets                []*PoolAsset `protobuf:"bytes,2,rep,name=assets,proto3" json:"assets,omitempty"`
	Supply                *types.Coin  `protobuf:"bytes,3,opt,name=supply,proto3" json:"supply,omitempty"`
	Status                PoolStatus   `protobuf:"varint,4,opt,name=status,proto3,enum=ibcswap.v4.interchainswap.PoolStatus" json:"status,omitempty"`
	EncounterPartyPort    string       `protobuf:"bytes,5,opt,name=encounterPartyPort,proto3" json:"encounterPartyPort,omitempty"`
	EncounterPartyChannel string       `protobuf:"bytes,6,opt,name=encounterPartyChannel,proto3" json:"encounterPartyChannel,omitempty"`
}

func (m *InterchainLiquidityPool) Reset()         { *m = InterchainLiquidityPool{} }
func (m *InterchainLiquidityPool) String() string { return proto.CompactTextString(m) }
func (*InterchainLiquidityPool) ProtoMessage()    {}
func (*InterchainLiquidityPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_11407e7be1a65d42, []int{1}
}
func (m *InterchainLiquidityPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InterchainLiquidityPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InterchainLiquidityPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InterchainLiquidityPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterchainLiquidityPool.Merge(m, src)
}
func (m *InterchainLiquidityPool) XXX_Size() int {
	return m.Size()
}
func (m *InterchainLiquidityPool) XXX_DiscardUnknown() {
	xxx_messageInfo_InterchainLiquidityPool.DiscardUnknown(m)
}

var xxx_messageInfo_InterchainLiquidityPool proto.InternalMessageInfo

func (m *InterchainLiquidityPool) GetPoolId() string {
	if m != nil {
		return m.PoolId
	}
	return ""
}

func (m *InterchainLiquidityPool) GetAssets() []*PoolAsset {
	if m != nil {
		return m.Assets
	}
	return nil
}

func (m *InterchainLiquidityPool) GetSupply() *types.Coin {
	if m != nil {
		return m.Supply
	}
	return nil
}

func (m *InterchainLiquidityPool) GetStatus() PoolStatus {
	if m != nil {
		return m.Status
	}
	return PoolStatus_POOL_STATUS_INITIAL
}

func (m *InterchainLiquidityPool) GetEncounterPartyPort() string {
	if m != nil {
		return m.EncounterPartyPort
	}
	return ""
}

func (m *InterchainLiquidityPool) GetEncounterPartyChannel() string {
	if m != nil {
		return m.EncounterPartyChannel
	}
	return ""
}

type InterchainMarketMaker struct {
	PoolId  string                   `protobuf:"bytes,1,opt,name=poolId,proto3" json:"poolId,omitempty"`
	Pool    *InterchainLiquidityPool `protobuf:"bytes,2,opt,name=pool,proto3" json:"pool,omitempty"`
	FeeRate uint64                   `protobuf:"varint,3,opt,name=feeRate,proto3" json:"feeRate,omitempty"`
}

func (m *InterchainMarketMaker) Reset()         { *m = InterchainMarketMaker{} }
func (m *InterchainMarketMaker) String() string { return proto.CompactTextString(m) }
func (*InterchainMarketMaker) ProtoMessage()    {}
func (*InterchainMarketMaker) Descriptor() ([]byte, []int) {
	return fileDescriptor_11407e7be1a65d42, []int{2}
}
func (m *InterchainMarketMaker) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InterchainMarketMaker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InterchainMarketMaker.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InterchainMarketMaker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterchainMarketMaker.Merge(m, src)
}
func (m *InterchainMarketMaker) XXX_Size() int {
	return m.Size()
}
func (m *InterchainMarketMaker) XXX_DiscardUnknown() {
	xxx_messageInfo_InterchainMarketMaker.DiscardUnknown(m)
}

var xxx_messageInfo_InterchainMarketMaker proto.InternalMessageInfo

func (m *InterchainMarketMaker) GetPoolId() string {
	if m != nil {
		return m.PoolId
	}
	return ""
}

func (m *InterchainMarketMaker) GetPool() *InterchainLiquidityPool {
	if m != nil {
		return m.Pool
	}
	return nil
}

func (m *InterchainMarketMaker) GetFeeRate() uint64 {
	if m != nil {
		return m.FeeRate
	}
	return 0
}

func init() {
	proto.RegisterEnum("ibcswap.v4.interchainswap.PoolSide", PoolSide_name, PoolSide_value)
	proto.RegisterEnum("ibcswap.v4.interchainswap.PoolStatus", PoolStatus_name, PoolStatus_value)
	proto.RegisterType((*PoolAsset)(nil), "ibcswap.v4.interchainswap.PoolAsset")
	proto.RegisterType((*InterchainLiquidityPool)(nil), "ibcswap.v4.interchainswap.InterchainLiquidityPool")
	proto.RegisterType((*InterchainMarketMaker)(nil), "ibcswap.v4.interchainswap.InterchainMarketMaker")
}

func init() {
	proto.RegisterFile("ibcswap/interchainswap/market.proto", fileDescriptor_11407e7be1a65d42)
}

var fileDescriptor_11407e7be1a65d42 = []byte{
	// 519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xc1, 0x6e, 0xd3, 0x4c,
	0x10, 0xce, 0xb6, 0xf9, 0xdd, 0xbf, 0x53, 0x51, 0x85, 0x45, 0xa5, 0x2e, 0x07, 0x2b, 0x4a, 0x41,
	0x8a, 0x7a, 0x58, 0x2b, 0x69, 0x24, 0x84, 0x54, 0x0e, 0xa1, 0x04, 0xc9, 0x52, 0xd2, 0x44, 0x9b,
	0x80, 0x04, 0x97, 0x6a, 0xed, 0x2c, 0xcd, 0xaa, 0x8e, 0xd7, 0x78, 0x37, 0x29, 0x79, 0x0b, 0x78,
	0x0e, 0x1e, 0x04, 0x8e, 0x3d, 0x72, 0x44, 0xc9, 0x8b, 0xa0, 0xdd, 0x3a, 0x94, 0xa2, 0xa6, 0xbd,
	0x79, 0x3c, 0xdf, 0xb7, 0x33, 0xdf, 0x37, 0x33, 0xb0, 0x2f, 0xc2, 0x48, 0x5d, 0xb0, 0xd4, 0x17,
	0x89, 0xe6, 0x59, 0x34, 0x62, 0x22, 0xb1, 0xe1, 0x98, 0x65, 0xe7, 0x5c, 0x93, 0x34, 0x93, 0x5a,
	0xe2, 0xbd, 0x1c, 0x44, 0xa6, 0x0d, 0x72, 0x13, 0xf7, 0xc4, 0x8b, 0xa4, 0x1a, 0x4b, 0xe5, 0x87,
	0x4c, 0x71, 0x7f, 0x5a, 0x0b, 0xb9, 0x66, 0x35, 0x3f, 0x92, 0x22, 0xb9, 0xa2, 0x56, 0xbe, 0x21,
	0xd8, 0xec, 0x49, 0x19, 0x37, 0x95, 0xe2, 0x1a, 0x3f, 0x87, 0xa2, 0x12, 0x43, 0xee, 0xa2, 0x32,
	0xaa, 0x6e, 0xd7, 0xf7, 0xc9, 0xca, 0x77, 0x89, 0xe1, 0xf4, 0xc5, 0x90, 0x53, 0x4b, 0xc0, 0x87,
	0xb0, 0x11, 0xb2, 0x98, 0x25, 0x11, 0x77, 0xd7, 0xca, 0xa8, 0xba, 0x55, 0xdf, 0x23, 0x57, 0x85,
	0x89, 0x29, 0x4c, 0xf2, 0xc2, 0xe4, 0x58, 0x8a, 0x84, 0x2e, 0x91, 0xf8, 0x31, 0x38, 0x17, 0x5c,
	0x9c, 0x8d, 0xb4, 0xbb, 0x5e, 0x46, 0xd5, 0x07, 0x34, 0x8f, 0xb0, 0x0b, 0x1b, 0x43, 0x1e, 0x89,
	0x31, 0x8b, 0xdd, 0xa2, 0x4d, 0x2c, 0xc3, 0xca, 0xf7, 0x35, 0xd8, 0x0d, 0xfe, 0x34, 0xd2, 0x16,
	0x9f, 0x26, 0x62, 0x28, 0xf4, 0xcc, 0x34, 0x63, 0x5e, 0x4b, 0xa5, 0x8c, 0x83, 0xa1, 0xed, 0x7e,
	0x93, 0xe6, 0x11, 0x3e, 0x02, 0x87, 0x19, 0x71, 0xca, 0x5d, 0x2b, 0xaf, 0x57, 0xb7, 0xea, 0x4f,
	0xef, 0x51, 0x65, 0x9d, 0xa0, 0x39, 0x07, 0xd7, 0xc0, 0x51, 0x93, 0x34, 0x8d, 0x67, 0xb6, 0xc7,
	0x3b, 0x75, 0xe5, 0x40, 0xfc, 0x12, 0x1c, 0xa5, 0x99, 0x9e, 0x28, 0xdb, 0xfd, 0x76, 0xfd, 0xd9,
	0x7d, 0x36, 0x5a, 0x30, 0xcd, 0x49, 0x98, 0x00, 0xe6, 0x49, 0x24, 0x27, 0x06, 0xd8, 0x63, 0x99,
	0x51, 0x97, 0x69, 0xf7, 0x3f, 0xab, 0xe9, 0x96, 0x0c, 0x6e, 0xc0, 0xce, 0xcd, 0xbf, 0xc7, 0x23,
	0x96, 0x24, 0x3c, 0x76, 0x1d, 0x4b, 0xb9, 0x3d, 0x59, 0xf9, 0x8a, 0x60, 0xe7, 0xda, 0xc9, 0x8e,
	0xdd, 0xa6, 0x0e, 0x3b, 0xe7, 0xd9, 0x4a, 0x1f, 0xdf, 0x40, 0xd1, 0x7c, 0xe5, 0xf3, 0xad, 0xdf,
	0x21, 0x6a, 0xc5, 0x84, 0xa8, 0xe5, 0x9b, 0xe9, 0x7e, 0xe4, 0x9c, 0x32, 0xcd, 0xad, 0xa5, 0x45,
	0xba, 0x0c, 0x0f, 0x2a, 0xf0, 0xff, 0x72, 0xad, 0x30, 0x80, 0x73, 0xd2, 0x1c, 0x04, 0xef, 0x5a,
	0xa5, 0x82, 0xf9, 0xa6, 0xad, 0x4e, 0x77, 0xd0, 0x2a, 0xa1, 0x83, 0x23, 0x80, 0x6b, 0xcf, 0xf0,
	0x2e, 0x3c, 0xea, 0x75, 0xbb, 0xed, 0xd3, 0xfe, 0xa0, 0x39, 0x78, 0xdb, 0x3f, 0x0d, 0x4e, 0x82,
	0x41, 0xd0, 0x6c, 0x97, 0x0a, 0x78, 0x07, 0x1e, 0xfe, 0x9d, 0xa0, 0xad, 0xe6, 0xeb, 0xf7, 0x25,
	0xf4, 0xaa, 0xff, 0x63, 0xee, 0xa1, 0xcb, 0xb9, 0x87, 0x7e, 0xcd, 0x3d, 0xf4, 0x65, 0xe1, 0x15,
	0x2e, 0x17, 0x5e, 0xe1, 0xe7, 0xc2, 0x2b, 0x7c, 0x78, 0x71, 0x26, 0xf4, 0x68, 0x12, 0x92, 0x48,
	0x8e, 0x7d, 0xb3, 0xd1, 0xf6, 0x3a, 0x22, 0x19, 0xfb, 0xcb, 0xfb, 0x9b, 0x36, 0xfc, 0xcf, 0xff,
	0x1e, 0xa1, 0x9e, 0xa5, 0x5c, 0x85, 0x8e, 0xc5, 0x1e, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x18,
	0xf3, 0xf6, 0xe4, 0xab, 0x03, 0x00, 0x00,
}

func (m *PoolAsset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolAsset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolAsset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Decimal != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.Decimal))
		i--
		dAtA[i] = 0x20
	}
	if m.Weight != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.Weight))
		i--
		dAtA[i] = 0x18
	}
	if m.Balance != nil {
		{
			size, err := m.Balance.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Side != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.Side))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *InterchainLiquidityPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterchainLiquidityPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InterchainLiquidityPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EncounterPartyChannel) > 0 {
		i -= len(m.EncounterPartyChannel)
		copy(dAtA[i:], m.EncounterPartyChannel)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.EncounterPartyChannel)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.EncounterPartyPort) > 0 {
		i -= len(m.EncounterPartyPort)
		copy(dAtA[i:], m.EncounterPartyPort)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.EncounterPartyPort)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Status != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x20
	}
	if m.Supply != nil {
		{
			size, err := m.Supply.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Assets) > 0 {
		for iNdEx := len(m.Assets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Assets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarket(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.PoolId) > 0 {
		i -= len(m.PoolId)
		copy(dAtA[i:], m.PoolId)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.PoolId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *InterchainMarketMaker) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterchainMarketMaker) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InterchainMarketMaker) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.FeeRate != 0 {
		i = encodeVarintMarket(dAtA, i, uint64(m.FeeRate))
		i--
		dAtA[i] = 0x18
	}
	if m.Pool != nil {
		{
			size, err := m.Pool.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarket(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.PoolId) > 0 {
		i -= len(m.PoolId)
		copy(dAtA[i:], m.PoolId)
		i = encodeVarintMarket(dAtA, i, uint64(len(m.PoolId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMarket(dAtA []byte, offset int, v uint64) int {
	offset -= sovMarket(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PoolAsset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Side != 0 {
		n += 1 + sovMarket(uint64(m.Side))
	}
	if m.Balance != nil {
		l = m.Balance.Size()
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.Weight != 0 {
		n += 1 + sovMarket(uint64(m.Weight))
	}
	if m.Decimal != 0 {
		n += 1 + sovMarket(uint64(m.Decimal))
	}
	return n
}

func (m *InterchainLiquidityPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PoolId)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	if len(m.Assets) > 0 {
		for _, e := range m.Assets {
			l = e.Size()
			n += 1 + l + sovMarket(uint64(l))
		}
	}
	if m.Supply != nil {
		l = m.Supply.Size()
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovMarket(uint64(m.Status))
	}
	l = len(m.EncounterPartyPort)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	l = len(m.EncounterPartyChannel)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	return n
}

func (m *InterchainMarketMaker) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PoolId)
	if l > 0 {
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.Pool != nil {
		l = m.Pool.Size()
		n += 1 + l + sovMarket(uint64(l))
	}
	if m.FeeRate != 0 {
		n += 1 + sovMarket(uint64(m.FeeRate))
	}
	return n
}

func sovMarket(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMarket(x uint64) (n int) {
	return sovMarket(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PoolAsset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: PoolAsset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolAsset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Side", wireType)
			}
			m.Side = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Side |= PoolSide(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Balance == nil {
				m.Balance = &types.Coin{}
			}
			if err := m.Balance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			m.Weight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Weight |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimal", wireType)
			}
			m.Decimal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimal |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func (m *InterchainLiquidityPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: InterchainLiquidityPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterchainLiquidityPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Assets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Assets = append(m.Assets, &PoolAsset{})
			if err := m.Assets[len(m.Assets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Supply", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Supply == nil {
				m.Supply = &types.Coin{}
			}
			if err := m.Supply.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= PoolStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncounterPartyPort", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncounterPartyPort = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncounterPartyChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncounterPartyChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func (m *InterchainMarketMaker) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarket
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
			return fmt.Errorf("proto: InterchainMarketMaker: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterchainMarketMaker: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pool", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
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
				return ErrInvalidLengthMarket
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pool == nil {
				m.Pool = &InterchainLiquidityPool{}
			}
			if err := m.Pool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType)
			}
			m.FeeRate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeRate |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarket
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
func skipMarket(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMarket
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
					return 0, ErrIntOverflowMarket
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
					return 0, ErrIntOverflowMarket
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
				return 0, ErrInvalidLengthMarket
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMarket
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMarket
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMarket        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMarket          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMarket = fmt.Errorf("proto: unexpected end of group")
)
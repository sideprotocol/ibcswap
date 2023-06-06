// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibc/applications/atomic_swap/v1/query.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type QueryOrdersRequest struct {
}

func (m *QueryOrdersRequest) Reset()         { *m = QueryOrdersRequest{} }
func (m *QueryOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*QueryOrdersRequest) ProtoMessage()    {}
func (*QueryOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{0}
}
func (m *QueryOrdersRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOrdersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOrdersRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOrdersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOrdersRequest.Merge(m, src)
}
func (m *QueryOrdersRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryOrdersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOrdersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOrdersRequest proto.InternalMessageInfo

type QueryOrdersResponse struct {
	Orders []*Order `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
}

func (m *QueryOrdersResponse) Reset()         { *m = QueryOrdersResponse{} }
func (m *QueryOrdersResponse) String() string { return proto.CompactTextString(m) }
func (*QueryOrdersResponse) ProtoMessage()    {}
func (*QueryOrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{1}
}
func (m *QueryOrdersResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOrdersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOrdersResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOrdersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOrdersResponse.Merge(m, src)
}
func (m *QueryOrdersResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryOrdersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOrdersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOrdersResponse proto.InternalMessageInfo

func (m *QueryOrdersResponse) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{2}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params defines the parameters of the module.
	Params *Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{3}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

// QueryEscrowAddressRequest is the request type for the EscrowAddress RPC method.
type QueryEscrowAddressRequest struct {
	// unique port identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// unique channel identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (m *QueryEscrowAddressRequest) Reset()         { *m = QueryEscrowAddressRequest{} }
func (m *QueryEscrowAddressRequest) String() string { return proto.CompactTextString(m) }
func (*QueryEscrowAddressRequest) ProtoMessage()    {}
func (*QueryEscrowAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{4}
}
func (m *QueryEscrowAddressRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEscrowAddressRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEscrowAddressRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEscrowAddressRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEscrowAddressRequest.Merge(m, src)
}
func (m *QueryEscrowAddressRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryEscrowAddressRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEscrowAddressRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEscrowAddressRequest proto.InternalMessageInfo

func (m *QueryEscrowAddressRequest) GetPortId() string {
	if m != nil {
		return m.PortId
	}
	return ""
}

func (m *QueryEscrowAddressRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

// QueryEscrowAddressResponse is the response type of the EscrowAddress RPC method.
type QueryEscrowAddressResponse struct {
	// the escrow account address
	EscrowAddress string `protobuf:"bytes,1,opt,name=escrow_address,json=escrowAddress,proto3" json:"escrow_address,omitempty"`
}

func (m *QueryEscrowAddressResponse) Reset()         { *m = QueryEscrowAddressResponse{} }
func (m *QueryEscrowAddressResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEscrowAddressResponse) ProtoMessage()    {}
func (*QueryEscrowAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd3ffdce48e373f, []int{5}
}
func (m *QueryEscrowAddressResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEscrowAddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEscrowAddressResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEscrowAddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEscrowAddressResponse.Merge(m, src)
}
func (m *QueryEscrowAddressResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEscrowAddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEscrowAddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEscrowAddressResponse proto.InternalMessageInfo

func (m *QueryEscrowAddressResponse) GetEscrowAddress() string {
	if m != nil {
		return m.EscrowAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryOrdersRequest)(nil), "ibc.applications.atomic_swap.v1.QueryOrdersRequest")
	proto.RegisterType((*QueryOrdersResponse)(nil), "ibc.applications.atomic_swap.v1.QueryOrdersResponse")
	proto.RegisterType((*QueryParamsRequest)(nil), "ibc.applications.atomic_swap.v1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "ibc.applications.atomic_swap.v1.QueryParamsResponse")
	proto.RegisterType((*QueryEscrowAddressRequest)(nil), "ibc.applications.atomic_swap.v1.QueryEscrowAddressRequest")
	proto.RegisterType((*QueryEscrowAddressResponse)(nil), "ibc.applications.atomic_swap.v1.QueryEscrowAddressResponse")
}

func init() {
	proto.RegisterFile("ibc/applications/atomic_swap/v1/query.proto", fileDescriptor_fcd3ffdce48e373f)
}

var fileDescriptor_fcd3ffdce48e373f = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x6b, 0xd4, 0x40,
	0x18, 0xde, 0x54, 0x5d, 0xe9, 0x94, 0x7a, 0x18, 0x0b, 0xd6, 0xa0, 0x71, 0x09, 0x58, 0x8b, 0xd2,
	0x4c, 0xb3, 0x15, 0x0f, 0x0a, 0x8a, 0x4a, 0x0f, 0x05, 0xf1, 0x63, 0xfd, 0x38, 0x78, 0x59, 0x26,
	0xc9, 0x90, 0x0e, 0x24, 0x79, 0xa7, 0x99, 0xd9, 0x2d, 0xa5, 0xf4, 0xe2, 0x2f, 0x10, 0x3c, 0x7a,
	0xf5, 0xc7, 0x78, 0x2c, 0x78, 0xf1, 0x24, 0xb2, 0xeb, 0x0f, 0x91, 0xcc, 0x4c, 0xea, 0x06, 0x2b,
	0x59, 0x7b, 0x4b, 0xde, 0x8f, 0xe7, 0x79, 0xe6, 0x79, 0xdf, 0x17, 0xdd, 0xe1, 0x51, 0x4c, 0xa8,
	0x10, 0x19, 0x8f, 0xa9, 0xe2, 0x50, 0x48, 0x42, 0x15, 0xe4, 0x3c, 0x1e, 0xca, 0x7d, 0x2a, 0xc8,
	0x38, 0x24, 0x7b, 0x23, 0x56, 0x1e, 0x04, 0xa2, 0x04, 0x05, 0xf8, 0x06, 0x8f, 0xe2, 0x60, 0xb6,
	0x38, 0x98, 0x29, 0x0e, 0xc6, 0xa1, 0xbb, 0x92, 0x42, 0x0a, 0xba, 0x96, 0x54, 0x5f, 0xa6, 0xcd,
	0xbd, 0x1d, 0x83, 0xcc, 0x41, 0x92, 0x88, 0x4a, 0x66, 0xf0, 0xc8, 0x38, 0x8c, 0x98, 0xa2, 0x21,
	0x11, 0x34, 0xe5, 0x85, 0xc6, 0xaa, 0x6b, 0xdb, 0xf4, 0x68, 0x2a, 0x53, 0x7b, 0x2d, 0x05, 0x48,
	0x33, 0x46, 0xa8, 0xe0, 0x84, 0x16, 0x05, 0x28, 0x2b, 0x4a, 0x67, 0xfd, 0x15, 0x84, 0x5f, 0x55,
	0x5c, 0x2f, 0xca, 0x84, 0x95, 0x72, 0xc0, 0xf6, 0x46, 0x4c, 0x2a, 0xff, 0x2d, 0xba, 0xdc, 0x88,
	0x4a, 0x01, 0x85, 0x64, 0xf8, 0x21, 0xea, 0x82, 0x8e, 0xac, 0x3a, 0xbd, 0x73, 0xeb, 0x4b, 0xfd,
	0xb5, 0xa0, 0xe5, 0xa9, 0x81, 0x06, 0x18, 0xd8, 0xae, 0x13, 0xb2, 0x97, 0xb4, 0xa4, 0xf9, 0x09,
	0xd9, 0x3b, 0x4b, 0x56, 0x47, 0x2d, 0xd9, 0x23, 0xd4, 0x15, 0x3a, 0xb2, 0xea, 0xf4, 0x9c, 0xf5,
	0xa5, 0xfe, 0xad, 0x56, 0x32, 0x0b, 0x60, 0xdb, 0xfc, 0xd7, 0xe8, 0xaa, 0xc6, 0xdd, 0x96, 0x71,
	0x09, 0xfb, 0x8f, 0x93, 0xa4, 0x64, 0xb2, 0x26, 0xc5, 0x57, 0xd0, 0x45, 0x01, 0xa5, 0x1a, 0xf2,
	0x44, 0xc3, 0x2f, 0x0e, 0xba, 0xd5, 0xef, 0x4e, 0x82, 0xaf, 0x23, 0x14, 0xef, 0xd2, 0xa2, 0x60,
	0x59, 0x95, 0x5b, 0xd0, 0xb9, 0x45, 0x1b, 0xd9, 0x49, 0xfc, 0xa7, 0xc8, 0x3d, 0x0d, 0xd4, 0x6a,
	0xbe, 0x89, 0x2e, 0x31, 0x9d, 0x18, 0x52, 0x93, 0xb1, 0xe0, 0xcb, 0x6c, 0xb6, 0xbc, 0xff, 0xe5,
	0x3c, 0xba, 0xa0, 0x51, 0xf0, 0x67, 0x07, 0x75, 0x8d, 0x6c, 0xbc, 0xd5, 0xfa, 0xbe, 0xbf, 0xbd,
	0x73, 0xef, 0xfe, 0x5f, 0x93, 0x91, 0xe9, 0xaf, 0x7d, 0xf8, 0xf6, 0xeb, 0xd3, 0x42, 0x0f, 0x7b,
	0xc4, 0xee, 0x51, 0xbd, 0x3f, 0xf5, 0xfa, 0x18, 0x07, 0xf1, 0x0f, 0x07, 0x2d, 0x37, 0x1e, 0x8a,
	0xef, 0xcf, 0xc7, 0x77, 0x9a, 0xe5, 0xee, 0x83, 0x33, 0xf5, 0x5a, 0xc9, 0x6f, 0xb4, 0xe4, 0xe7,
	0xf8, 0xd9, 0xbf, 0x24, 0xdb, 0x11, 0x49, 0x72, 0xf8, 0x67, 0x7c, 0x47, 0xa4, 0x1a, 0xaa, 0x24,
	0x87, 0x76, 0xd4, 0x47, 0xa4, 0x39, 0x1d, 0x6d, 0xbf, 0xd9, 0xf1, 0x79, 0xed, 0x6f, 0xdc, 0xc9,
	0xbc, 0xf6, 0x37, 0xcf, 0xa8, 0xdd, 0x7e, 0x73, 0x2e, 0x4f, 0x86, 0x5f, 0x27, 0x9e, 0x73, 0x3c,
	0xf1, 0x9c, 0x9f, 0x13, 0xcf, 0xf9, 0x38, 0xf5, 0x3a, 0xc7, 0x53, 0xaf, 0xf3, 0x7d, 0xea, 0x75,
	0xde, 0x6f, 0xa7, 0x5c, 0xed, 0x8e, 0xa2, 0x20, 0x86, 0x9c, 0x48, 0x9e, 0x30, 0x7d, 0xcb, 0x31,
	0x64, 0x15, 0xa0, 0x01, 0xb9, 0x47, 0x72, 0x48, 0x46, 0x19, 0x93, 0x06, 0x3f, 0xdc, 0xdc, 0xdc,
	0x30, 0x1c, 0x1b, 0x3a, 0xaf, 0x0e, 0x04, 0x93, 0x51, 0x57, 0xf7, 0x6d, 0xfd, 0x0e, 0x00, 0x00,
	0xff, 0xff, 0x4c, 0xa8, 0x40, 0x37, 0xdf, 0x04, 0x00, 0x00,
}

func (m *QueryOrdersRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOrdersRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOrdersRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryOrdersResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOrdersResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOrdersResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Orders) > 0 {
		for iNdEx := len(m.Orders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Orders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryEscrowAddressRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEscrowAddressRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEscrowAddressRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChannelId) > 0 {
		i -= len(m.ChannelId)
		copy(dAtA[i:], m.ChannelId)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ChannelId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PortId) > 0 {
		i -= len(m.PortId)
		copy(dAtA[i:], m.PortId)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.PortId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryEscrowAddressResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEscrowAddressResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEscrowAddressResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EscrowAddress) > 0 {
		i -= len(m.EscrowAddress)
		copy(dAtA[i:], m.EscrowAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.EscrowAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryOrdersRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryOrdersResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Orders) > 0 {
		for _, e := range m.Orders {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryEscrowAddressRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PortId)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.ChannelId)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryEscrowAddressResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EscrowAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryOrdersRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryOrdersRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOrdersRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryOrdersResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryOrdersResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOrdersResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Orders", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Orders = append(m.Orders, &Order{})
			if err := m.Orders[len(m.Orders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryEscrowAddressRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEscrowAddressRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEscrowAddressRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PortId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PortId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryEscrowAddressResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEscrowAddressResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEscrowAddressResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EscrowAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EscrowAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)

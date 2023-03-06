// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: ibc/applications/atomic_swap/v1/tx.proto

package types

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	MakeSwap(ctx context.Context, in *MsgMakeSwapRequest, opts ...grpc.CallOption) (*MsgMakeSwapResponse, error)
	TakeSwap(ctx context.Context, in *MsgTakeSwapRequest, opts ...grpc.CallOption) (*MsgTakeSwapResponse, error)
	CancelSwap(ctx context.Context, in *MsgCancelSwapRequest, opts ...grpc.CallOption) (*MsgCancelSwapResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) MakeSwap(ctx context.Context, in *MsgMakeSwapRequest, opts ...grpc.CallOption) (*MsgMakeSwapResponse, error) {
	out := new(MsgMakeSwapResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.atomic_swap.v1.Msg/MakeSwap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) TakeSwap(ctx context.Context, in *MsgTakeSwapRequest, opts ...grpc.CallOption) (*MsgTakeSwapResponse, error) {
	out := new(MsgTakeSwapResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.atomic_swap.v1.Msg/TakeSwap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CancelSwap(ctx context.Context, in *MsgCancelSwapRequest, opts ...grpc.CallOption) (*MsgCancelSwapResponse, error) {
	out := new(MsgCancelSwapResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.atomic_swap.v1.Msg/CancelSwap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations should embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	MakeSwap(context.Context, *MsgMakeSwapRequest) (*MsgMakeSwapResponse, error)
	TakeSwap(context.Context, *MsgTakeSwapRequest) (*MsgTakeSwapResponse, error)
	CancelSwap(context.Context, *MsgCancelSwapRequest) (*MsgCancelSwapResponse, error)
}

// UnimplementedMsgServer should be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) MakeSwap(context.Context, *MsgMakeSwapRequest) (*MsgMakeSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeSwap not implemented")
}
func (UnimplementedMsgServer) TakeSwap(context.Context, *MsgTakeSwapRequest) (*MsgTakeSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeSwap not implemented")
}
func (UnimplementedMsgServer) CancelSwap(context.Context, *MsgCancelSwapRequest) (*MsgCancelSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelSwap not implemented")
}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_MakeSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMakeSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).MakeSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.atomic_swap.v1.Msg/MakeSwap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MakeSwap(ctx, req.(*MsgMakeSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_TakeSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgTakeSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).TakeSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.atomic_swap.v1.Msg/TakeSwap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).TakeSwap(ctx, req.(*MsgTakeSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.atomic_swap.v1.Msg/CancelSwap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelSwap(ctx, req.(*MsgCancelSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ibc.applications.atomic_swap.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakeSwap",
			Handler:    _Msg_MakeSwap_Handler,
		},
		{
			MethodName: "TakeSwap",
			Handler:    _Msg_TakeSwap_Handler,
		},
		{
			MethodName: "CancelSwap",
			Handler:    _Msg_CancelSwap_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ibc/applications/atomic_swap/v1/tx.proto",
}
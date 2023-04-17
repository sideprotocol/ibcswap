// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Params queries all parameters of the ibc-transfer module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel id.
	EscrowAddress(ctx context.Context, in *QueryEscrowAddressRequest, opts ...grpc.CallOption) (*QueryEscrowAddressResponse, error)
	// Queries a list of InterchainLiquidityPool items.
	InterchainLiquidityPool(ctx context.Context, in *QueryGetInterchainLiquidityPoolRequest, opts ...grpc.CallOption) (*QueryGetInterchainLiquidityPoolResponse, error)
	InterchainLiquidityPoolAll(ctx context.Context, in *QueryAllInterchainLiquidityPoolRequest, opts ...grpc.CallOption) (*QueryAllInterchainLiquidityPoolResponse, error)
	// Queries a list of InterchainMarketMaker items.
	InterchainMarketMaker(ctx context.Context, in *QueryGetInterchainMarketMakerRequest, opts ...grpc.CallOption) (*QueryGetInterchainMarketMakerResponse, error)
	InterchainMarketMakerAll(ctx context.Context, in *QueryAllInterchainMarketMakerRequest, opts ...grpc.CallOption) (*QueryAllInterchainMarketMakerResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EscrowAddress(ctx context.Context, in *QueryEscrowAddressRequest, opts ...grpc.CallOption) (*QueryEscrowAddressResponse, error) {
	out := new(QueryEscrowAddressResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/EscrowAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) InterchainLiquidityPool(ctx context.Context, in *QueryGetInterchainLiquidityPoolRequest, opts ...grpc.CallOption) (*QueryGetInterchainLiquidityPoolResponse, error) {
	out := new(QueryGetInterchainLiquidityPoolResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/InterchainLiquidityPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) InterchainLiquidityPoolAll(ctx context.Context, in *QueryAllInterchainLiquidityPoolRequest, opts ...grpc.CallOption) (*QueryAllInterchainLiquidityPoolResponse, error) {
	out := new(QueryAllInterchainLiquidityPoolResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/InterchainLiquidityPoolAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) InterchainMarketMaker(ctx context.Context, in *QueryGetInterchainMarketMakerRequest, opts ...grpc.CallOption) (*QueryGetInterchainMarketMakerResponse, error) {
	out := new(QueryGetInterchainMarketMakerResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/InterchainMarketMaker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) InterchainMarketMakerAll(ctx context.Context, in *QueryAllInterchainMarketMakerRequest, opts ...grpc.CallOption) (*QueryAllInterchainMarketMakerResponse, error) {
	out := new(QueryAllInterchainMarketMakerResponse)
	err := c.cc.Invoke(ctx, "/ibc.applications.interchain_swap.v1.Query/InterchainMarketMakerAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations should embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Params queries all parameters of the ibc-transfer module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel id.
	EscrowAddress(context.Context, *QueryEscrowAddressRequest) (*QueryEscrowAddressResponse, error)
	// Queries a list of InterchainLiquidityPool items.
	InterchainLiquidityPool(context.Context, *QueryGetInterchainLiquidityPoolRequest) (*QueryGetInterchainLiquidityPoolResponse, error)
	InterchainLiquidityPoolAll(context.Context, *QueryAllInterchainLiquidityPoolRequest) (*QueryAllInterchainLiquidityPoolResponse, error)
	// Queries a list of InterchainMarketMaker items.
	InterchainMarketMaker(context.Context, *QueryGetInterchainMarketMakerRequest) (*QueryGetInterchainMarketMakerResponse, error)
	InterchainMarketMakerAll(context.Context, *QueryAllInterchainMarketMakerRequest) (*QueryAllInterchainMarketMakerResponse, error)
}

// UnimplementedQueryServer should be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) EscrowAddress(context.Context, *QueryEscrowAddressRequest) (*QueryEscrowAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EscrowAddress not implemented")
}
func (UnimplementedQueryServer) InterchainLiquidityPool(context.Context, *QueryGetInterchainLiquidityPoolRequest) (*QueryGetInterchainLiquidityPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterchainLiquidityPool not implemented")
}
func (UnimplementedQueryServer) InterchainLiquidityPoolAll(context.Context, *QueryAllInterchainLiquidityPoolRequest) (*QueryAllInterchainLiquidityPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterchainLiquidityPoolAll not implemented")
}
func (UnimplementedQueryServer) InterchainMarketMaker(context.Context, *QueryGetInterchainMarketMakerRequest) (*QueryGetInterchainMarketMakerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterchainMarketMaker not implemented")
}
func (UnimplementedQueryServer) InterchainMarketMakerAll(context.Context, *QueryAllInterchainMarketMakerRequest) (*QueryAllInterchainMarketMakerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterchainMarketMakerAll not implemented")
}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EscrowAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEscrowAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EscrowAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/EscrowAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EscrowAddress(ctx, req.(*QueryEscrowAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_InterchainLiquidityPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetInterchainLiquidityPoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).InterchainLiquidityPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/InterchainLiquidityPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).InterchainLiquidityPool(ctx, req.(*QueryGetInterchainLiquidityPoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_InterchainLiquidityPoolAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllInterchainLiquidityPoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).InterchainLiquidityPoolAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/InterchainLiquidityPoolAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).InterchainLiquidityPoolAll(ctx, req.(*QueryAllInterchainLiquidityPoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_InterchainMarketMaker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetInterchainMarketMakerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).InterchainMarketMaker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/InterchainMarketMaker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).InterchainMarketMaker(ctx, req.(*QueryGetInterchainMarketMakerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_InterchainMarketMakerAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllInterchainMarketMakerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).InterchainMarketMakerAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.applications.interchain_swap.v1.Query/InterchainMarketMakerAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).InterchainMarketMakerAll(ctx, req.(*QueryAllInterchainMarketMakerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ibc.applications.interchain_swap.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "EscrowAddress",
			Handler:    _Query_EscrowAddress_Handler,
		},
		{
			MethodName: "InterchainLiquidityPool",
			Handler:    _Query_InterchainLiquidityPool_Handler,
		},
		{
			MethodName: "InterchainLiquidityPoolAll",
			Handler:    _Query_InterchainLiquidityPoolAll_Handler,
		},
		{
			MethodName: "InterchainMarketMaker",
			Handler:    _Query_InterchainMarketMaker_Handler,
		},
		{
			MethodName: "InterchainMarketMakerAll",
			Handler:    _Query_InterchainMarketMakerAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ibc/applications/interchain_swap/v1/query.proto",
}
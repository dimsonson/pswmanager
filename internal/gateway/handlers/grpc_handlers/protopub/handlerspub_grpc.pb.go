// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: handlerspub.proto

package grpc_handlers

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

// ClientRMQhandlersClient is the client API for ClientRMQhandlers service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientRMQhandlersClient interface {
	PublishText(ctx context.Context, in *PublishTextRequest, opts ...grpc.CallOption) (*PublishTextResponse, error)
	PublishLogins(ctx context.Context, in *PublishLoginsRequest, opts ...grpc.CallOption) (*PublishLoginsResponse, error)
	PublishBinary(ctx context.Context, in *PublishBinaryRequest, opts ...grpc.CallOption) (*PublishBinaryResponse, error)
	PublishCard(ctx context.Context, in *PublishCardRequest, opts ...grpc.CallOption) (*PublishCardResponse, error)
}

type clientRMQhandlersClient struct {
	cc grpc.ClientConnInterface
}

func NewClientRMQhandlersClient(cc grpc.ClientConnInterface) ClientRMQhandlersClient {
	return &clientRMQhandlersClient{cc}
}

func (c *clientRMQhandlersClient) PublishText(ctx context.Context, in *PublishTextRequest, opts ...grpc.CallOption) (*PublishTextResponse, error) {
	out := new(PublishTextResponse)
	err := c.cc.Invoke(ctx, "/proto.ClientRMQhandlers/PublishText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientRMQhandlersClient) PublishLogins(ctx context.Context, in *PublishLoginsRequest, opts ...grpc.CallOption) (*PublishLoginsResponse, error) {
	out := new(PublishLoginsResponse)
	err := c.cc.Invoke(ctx, "/proto.ClientRMQhandlers/PublishLogins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientRMQhandlersClient) PublishBinary(ctx context.Context, in *PublishBinaryRequest, opts ...grpc.CallOption) (*PublishBinaryResponse, error) {
	out := new(PublishBinaryResponse)
	err := c.cc.Invoke(ctx, "/proto.ClientRMQhandlers/PublishBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientRMQhandlersClient) PublishCard(ctx context.Context, in *PublishCardRequest, opts ...grpc.CallOption) (*PublishCardResponse, error) {
	out := new(PublishCardResponse)
	err := c.cc.Invoke(ctx, "/proto.ClientRMQhandlers/PublishCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientRMQhandlersServer is the server API for ClientRMQhandlers service.
// All implementations must embed UnimplementedClientRMQhandlersServer
// for forward compatibility
type ClientRMQhandlersServer interface {
	PublishText(context.Context, *PublishTextRequest) (*PublishTextResponse, error)
	PublishLogins(context.Context, *PublishLoginsRequest) (*PublishLoginsResponse, error)
	PublishBinary(context.Context, *PublishBinaryRequest) (*PublishBinaryResponse, error)
	PublishCard(context.Context, *PublishCardRequest) (*PublishCardResponse, error)
	mustEmbedUnimplementedClientRMQhandlersServer()
}

// UnimplementedClientRMQhandlersServer must be embedded to have forward compatible implementations.
type UnimplementedClientRMQhandlersServer struct {
}

func (UnimplementedClientRMQhandlersServer) PublishText(context.Context, *PublishTextRequest) (*PublishTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishText not implemented")
}
func (UnimplementedClientRMQhandlersServer) PublishLogins(context.Context, *PublishLoginsRequest) (*PublishLoginsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishLogins not implemented")
}
func (UnimplementedClientRMQhandlersServer) PublishBinary(context.Context, *PublishBinaryRequest) (*PublishBinaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishBinary not implemented")
}
func (UnimplementedClientRMQhandlersServer) PublishCard(context.Context, *PublishCardRequest) (*PublishCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishCard not implemented")
}
func (UnimplementedClientRMQhandlersServer) mustEmbedUnimplementedClientRMQhandlersServer() {}

// UnsafeClientRMQhandlersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientRMQhandlersServer will
// result in compilation errors.
type UnsafeClientRMQhandlersServer interface {
	mustEmbedUnimplementedClientRMQhandlersServer()
}

func RegisterClientRMQhandlersServer(s grpc.ServiceRegistrar, srv ClientRMQhandlersServer) {
	s.RegisterService(&ClientRMQhandlers_ServiceDesc, srv)
}

func _ClientRMQhandlers_PublishText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientRMQhandlersServer).PublishText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClientRMQhandlers/PublishText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientRMQhandlersServer).PublishText(ctx, req.(*PublishTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientRMQhandlers_PublishLogins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishLoginsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientRMQhandlersServer).PublishLogins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClientRMQhandlers/PublishLogins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientRMQhandlersServer).PublishLogins(ctx, req.(*PublishLoginsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientRMQhandlers_PublishBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientRMQhandlersServer).PublishBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClientRMQhandlers/PublishBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientRMQhandlersServer).PublishBinary(ctx, req.(*PublishBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientRMQhandlers_PublishCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientRMQhandlersServer).PublishCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClientRMQhandlers/PublishCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientRMQhandlersServer).PublishCard(ctx, req.(*PublishCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientRMQhandlers_ServiceDesc is the grpc.ServiceDesc for ClientRMQhandlers service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientRMQhandlers_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClientRMQhandlers",
	HandlerType: (*ClientRMQhandlersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishText",
			Handler:    _ClientRMQhandlers_PublishText_Handler,
		},
		{
			MethodName: "PublishLogins",
			Handler:    _ClientRMQhandlers_PublishLogins_Handler,
		},
		{
			MethodName: "PublishBinary",
			Handler:    _ClientRMQhandlers_PublishBinary_Handler,
		},
		{
			MethodName: "PublishCard",
			Handler:    _ClientRMQhandlers_PublishCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "handlerspub.proto",
}

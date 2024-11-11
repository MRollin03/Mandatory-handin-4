// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: proto/node.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TokenRing_PassTokenToNext_FullMethodName        = "/proto.TokenRing/PassTokenToNext"
	TokenRing_RequestCriticalSection_FullMethodName = "/proto.TokenRing/RequestCriticalSection"
	TokenRing_ReceiveToken_FullMethodName           = "/proto.TokenRing/ReceiveToken"
)

// TokenRingClient is the client API for TokenRing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenRingClient interface {
	PassTokenToNext(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error)
	RequestCriticalSection(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error)
	ReceiveToken(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error)
}

type tokenRingClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenRingClient(cc grpc.ClientConnInterface) TokenRingClient {
	return &tokenRingClient{cc}
}

func (c *tokenRingClient) PassTokenToNext(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AckMessage)
	err := c.cc.Invoke(ctx, TokenRing_PassTokenToNext_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenRingClient) RequestCriticalSection(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AckMessage)
	err := c.cc.Invoke(ctx, TokenRing_RequestCriticalSection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenRingClient) ReceiveToken(ctx context.Context, in *ReceiveMessageRequest, opts ...grpc.CallOption) (*AckMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AckMessage)
	err := c.cc.Invoke(ctx, TokenRing_ReceiveToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenRingServer is the server API for TokenRing service.
// All implementations must embed UnimplementedTokenRingServer
// for forward compatibility.
type TokenRingServer interface {
	PassTokenToNext(context.Context, *ReceiveMessageRequest) (*AckMessage, error)
	RequestCriticalSection(context.Context, *ReceiveMessageRequest) (*AckMessage, error)
	ReceiveToken(context.Context, *ReceiveMessageRequest) (*AckMessage, error)
	mustEmbedUnimplementedTokenRingServer()
}

// UnimplementedTokenRingServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTokenRingServer struct{}

func (UnimplementedTokenRingServer) PassTokenToNext(context.Context, *ReceiveMessageRequest) (*AckMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PassTokenToNext not implemented")
}
func (UnimplementedTokenRingServer) RequestCriticalSection(context.Context, *ReceiveMessageRequest) (*AckMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestCriticalSection not implemented")
}
func (UnimplementedTokenRingServer) ReceiveToken(context.Context, *ReceiveMessageRequest) (*AckMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveToken not implemented")
}
func (UnimplementedTokenRingServer) mustEmbedUnimplementedTokenRingServer() {}
func (UnimplementedTokenRingServer) testEmbeddedByValue()                   {}

// UnsafeTokenRingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenRingServer will
// result in compilation errors.
type UnsafeTokenRingServer interface {
	mustEmbedUnimplementedTokenRingServer()
}

func RegisterTokenRingServer(s grpc.ServiceRegistrar, srv TokenRingServer) {
	// If the following call pancis, it indicates UnimplementedTokenRingServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TokenRing_ServiceDesc, srv)
}

func _TokenRing_PassTokenToNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenRingServer).PassTokenToNext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenRing_PassTokenToNext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenRingServer).PassTokenToNext(ctx, req.(*ReceiveMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenRing_RequestCriticalSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenRingServer).RequestCriticalSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenRing_RequestCriticalSection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenRingServer).RequestCriticalSection(ctx, req.(*ReceiveMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenRing_ReceiveToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenRingServer).ReceiveToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenRing_ReceiveToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenRingServer).ReceiveToken(ctx, req.(*ReceiveMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenRing_ServiceDesc is the grpc.ServiceDesc for TokenRing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenRing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TokenRing",
	HandlerType: (*TokenRingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PassTokenToNext",
			Handler:    _TokenRing_PassTokenToNext_Handler,
		},
		{
			MethodName: "RequestCriticalSection",
			Handler:    _TokenRing_RequestCriticalSection_Handler,
		},
		{
			MethodName: "ReceiveToken",
			Handler:    _TokenRing_ReceiveToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/node.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// AuthorityServiceClient is the client API for AuthorityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorityServiceClient interface {
	Signup(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error)
	Signin(ctx context.Context, in *SigninRequest, opts ...grpc.CallOption) (*SigninResponse, error)
	ListPublicKeys(ctx context.Context, in *ListPublicKeysRequest, opts ...grpc.CallOption) (*ListPublicKeysResponse, error)
}

type authorityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorityServiceClient(cc grpc.ClientConnInterface) AuthorityServiceClient {
	return &authorityServiceClient{cc}
}

func (c *authorityServiceClient) Signup(ctx context.Context, in *SignupRequest, opts ...grpc.CallOption) (*SignupResponse, error) {
	out := new(SignupResponse)
	err := c.cc.Invoke(ctx, "/dnakano.microservice_handson.authority.AuthorityService/Signup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorityServiceClient) Signin(ctx context.Context, in *SigninRequest, opts ...grpc.CallOption) (*SigninResponse, error) {
	out := new(SigninResponse)
	err := c.cc.Invoke(ctx, "/dnakano.microservice_handson.authority.AuthorityService/Signin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorityServiceClient) ListPublicKeys(ctx context.Context, in *ListPublicKeysRequest, opts ...grpc.CallOption) (*ListPublicKeysResponse, error) {
	out := new(ListPublicKeysResponse)
	err := c.cc.Invoke(ctx, "/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorityServiceServer is the server API for AuthorityService service.
// All implementations must embed UnimplementedAuthorityServiceServer
// for forward compatibility
type AuthorityServiceServer interface {
	Signup(context.Context, *SignupRequest) (*SignupResponse, error)
	Signin(context.Context, *SigninRequest) (*SigninResponse, error)
	ListPublicKeys(context.Context, *ListPublicKeysRequest) (*ListPublicKeysResponse, error)
	mustEmbedUnimplementedAuthorityServiceServer()
}

// UnimplementedAuthorityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthorityServiceServer struct {
}

func (UnimplementedAuthorityServiceServer) Signup(context.Context, *SignupRequest) (*SignupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}
func (UnimplementedAuthorityServiceServer) Signin(context.Context, *SigninRequest) (*SigninResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signin not implemented")
}
func (UnimplementedAuthorityServiceServer) ListPublicKeys(context.Context, *ListPublicKeysRequest) (*ListPublicKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPublicKeys not implemented")
}
func (UnimplementedAuthorityServiceServer) mustEmbedUnimplementedAuthorityServiceServer() {}

// UnsafeAuthorityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorityServiceServer will
// result in compilation errors.
type UnsafeAuthorityServiceServer interface {
	mustEmbedUnimplementedAuthorityServiceServer()
}

func RegisterAuthorityServiceServer(s grpc.ServiceRegistrar, srv AuthorityServiceServer) {
	s.RegisterService(&AuthorityService_ServiceDesc, srv)
}

func _AuthorityService_Signup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorityServiceServer).Signup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dnakano.microservice_handson.authority.AuthorityService/Signup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorityServiceServer).Signup(ctx, req.(*SignupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorityService_Signin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SigninRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorityServiceServer).Signin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dnakano.microservice_handson.authority.AuthorityService/Signin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorityServiceServer).Signin(ctx, req.(*SigninRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorityService_ListPublicKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPublicKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorityServiceServer).ListPublicKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorityServiceServer).ListPublicKeys(ctx, req.(*ListPublicKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthorityService_ServiceDesc is the grpc.ServiceDesc for AuthorityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthorityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dnakano.microservice_handson.authority.AuthorityService",
	HandlerType: (*AuthorityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Signup",
			Handler:    _AuthorityService_Signup_Handler,
		},
		{
			MethodName: "Signin",
			Handler:    _AuthorityService_Signin_Handler,
		},
		{
			MethodName: "ListPublicKeys",
			Handler:    _AuthorityService_ListPublicKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/authority/proto/authority.proto",
}

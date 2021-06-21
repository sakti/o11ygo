// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// AddersvcClient is the client API for Addersvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AddersvcClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error)
}

type addersvcClient struct {
	cc grpc.ClientConnInterface
}

func NewAddersvcClient(cc grpc.ClientConnInterface) AddersvcClient {
	return &addersvcClient{cc}
}

func (c *addersvcClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error) {
	out := new(AddReply)
	err := c.cc.Invoke(ctx, "/api.Addersvc/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddersvcServer is the server API for Addersvc service.
// All implementations must embed UnimplementedAddersvcServer
// for forward compatibility
type AddersvcServer interface {
	Add(context.Context, *AddRequest) (*AddReply, error)
	mustEmbedUnimplementedAddersvcServer()
}

// UnimplementedAddersvcServer must be embedded to have forward compatible implementations.
type UnimplementedAddersvcServer struct {
}

func (UnimplementedAddersvcServer) Add(context.Context, *AddRequest) (*AddReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedAddersvcServer) mustEmbedUnimplementedAddersvcServer() {}

// UnsafeAddersvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AddersvcServer will
// result in compilation errors.
type UnsafeAddersvcServer interface {
	mustEmbedUnimplementedAddersvcServer()
}

func RegisterAddersvcServer(s grpc.ServiceRegistrar, srv AddersvcServer) {
	s.RegisterService(&Addersvc_ServiceDesc, srv)
}

func _Addersvc_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddersvcServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Addersvc/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddersvcServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Addersvc_ServiceDesc is the grpc.ServiceDesc for Addersvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Addersvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Addersvc",
	HandlerType: (*AddersvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Addersvc_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/addersvc.proto",
}

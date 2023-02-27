// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/api/proto/match.proto

package pb

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

// MatchClient is the client API for Match service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchClient interface {
	Next(ctx context.Context, in *NextRequest, opts ...grpc.CallOption) (*NextReply, error)
	UpdateRecommendations(ctx context.Context, in *UpdateRecommendationsRequest, opts ...grpc.CallOption) (*MatchEmpty, error)
}

type matchClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchClient(cc grpc.ClientConnInterface) MatchClient {
	return &matchClient{cc}
}

func (c *matchClient) Next(ctx context.Context, in *NextRequest, opts ...grpc.CallOption) (*NextReply, error) {
	out := new(NextReply)
	err := c.cc.Invoke(ctx, "/Match/Next", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchClient) UpdateRecommendations(ctx context.Context, in *UpdateRecommendationsRequest, opts ...grpc.CallOption) (*MatchEmpty, error) {
	out := new(MatchEmpty)
	err := c.cc.Invoke(ctx, "/Match/UpdateRecommendations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchServer is the server API for Match service.
// All implementations must embed UnimplementedMatchServer
// for forward compatibility
type MatchServer interface {
	Next(context.Context, *NextRequest) (*NextReply, error)
	UpdateRecommendations(context.Context, *UpdateRecommendationsRequest) (*MatchEmpty, error)
	mustEmbedUnimplementedMatchServer()
}

// UnimplementedMatchServer must be embedded to have forward compatible implementations.
type UnimplementedMatchServer struct {
}

func (UnimplementedMatchServer) Next(context.Context, *NextRequest) (*NextReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Next not implemented")
}
func (UnimplementedMatchServer) UpdateRecommendations(context.Context, *UpdateRecommendationsRequest) (*MatchEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRecommendations not implemented")
}
func (UnimplementedMatchServer) mustEmbedUnimplementedMatchServer() {}

// UnsafeMatchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchServer will
// result in compilation errors.
type UnsafeMatchServer interface {
	mustEmbedUnimplementedMatchServer()
}

func RegisterMatchServer(s grpc.ServiceRegistrar, srv MatchServer) {
	s.RegisterService(&Match_ServiceDesc, srv)
}

func _Match_Next_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServer).Next(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Match/Next",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServer).Next(ctx, req.(*NextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Match_UpdateRecommendations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRecommendationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServer).UpdateRecommendations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Match/UpdateRecommendations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServer).UpdateRecommendations(ctx, req.(*UpdateRecommendationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Match_ServiceDesc is the grpc.ServiceDesc for Match service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Match_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Match",
	HandlerType: (*MatchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Next",
			Handler:    _Match_Next_Handler,
		},
		{
			MethodName: "UpdateRecommendations",
			Handler:    _Match_UpdateRecommendations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/proto/match.proto",
}
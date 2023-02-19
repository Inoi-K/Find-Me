// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/api/proto/profile.proto

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

// ProfileClient is the client API for Profile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileClient interface {
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*Empty, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error)
	GetUserSphere(ctx context.Context, in *GetUserSphereRequest, opts ...grpc.CallOption) (*GetUserSphereReply, error)
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsReply, error)
	Edit(ctx context.Context, in *EditRequest, opts ...grpc.CallOption) (*Empty, error)
}

type profileClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileClient(cc grpc.ClientConnInterface) ProfileClient {
	return &profileClient{cc}
}

func (c *profileClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Profile/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error) {
	out := new(GetUserReply)
	err := c.cc.Invoke(ctx, "/Profile/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) GetUserSphere(ctx context.Context, in *GetUserSphereRequest, opts ...grpc.CallOption) (*GetUserSphereReply, error) {
	out := new(GetUserSphereReply)
	err := c.cc.Invoke(ctx, "/Profile/GetUserSphere", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsReply, error) {
	out := new(ExistsReply)
	err := c.cc.Invoke(ctx, "/Profile/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileClient) Edit(ctx context.Context, in *EditRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Profile/Edit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServer is the server API for Profile service.
// All implementations must embed UnimplementedProfileServer
// for forward compatibility
type ProfileServer interface {
	SignUp(context.Context, *SignUpRequest) (*Empty, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserReply, error)
	GetUserSphere(context.Context, *GetUserSphereRequest) (*GetUserSphereReply, error)
	Exists(context.Context, *ExistsRequest) (*ExistsReply, error)
	Edit(context.Context, *EditRequest) (*Empty, error)
	mustEmbedUnimplementedProfileServer()
}

// UnimplementedProfileServer must be embedded to have forward compatible implementations.
type UnimplementedProfileServer struct {
}

func (UnimplementedProfileServer) SignUp(context.Context, *SignUpRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedProfileServer) GetUser(context.Context, *GetUserRequest) (*GetUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedProfileServer) GetUserSphere(context.Context, *GetUserSphereRequest) (*GetUserSphereReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSphere not implemented")
}
func (UnimplementedProfileServer) Exists(context.Context, *ExistsRequest) (*ExistsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedProfileServer) Edit(context.Context, *EditRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedProfileServer) mustEmbedUnimplementedProfileServer() {}

// UnsafeProfileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServer will
// result in compilation errors.
type UnsafeProfileServer interface {
	mustEmbedUnimplementedProfileServer()
}

func RegisterProfileServer(s grpc.ServiceRegistrar, srv ProfileServer) {
	s.RegisterService(&Profile_ServiceDesc, srv)
}

func _Profile_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_GetUserSphere_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserSphereRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetUserSphere(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/GetUserSphere",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetUserSphere(ctx, req.(*GetUserSphereRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profile_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Profile/Edit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).Edit(ctx, req.(*EditRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Profile_ServiceDesc is the grpc.ServiceDesc for Profile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Profile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Profile",
	HandlerType: (*ProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _Profile_SignUp_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _Profile_GetUser_Handler,
		},
		{
			MethodName: "GetUserSphere",
			Handler:    _Profile_GetUserSphere_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _Profile_Exists_Handler,
		},
		{
			MethodName: "Edit",
			Handler:    _Profile_Edit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/proto/profile.proto",
}

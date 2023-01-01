// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: grpc/presentations/presentation.proto

package presentations

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

// PresentationClient is the client API for Presentation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PresentationClient interface {
	GetArtist(ctx context.Context, in *ArtistRequest, opts ...grpc.CallOption) (*Artist, error)
	GetRelease(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*Release, error)
}

type presentationClient struct {
	cc grpc.ClientConnInterface
}

func NewPresentationClient(cc grpc.ClientConnInterface) PresentationClient {
	return &presentationClient{cc}
}

func (c *presentationClient) GetArtist(ctx context.Context, in *ArtistRequest, opts ...grpc.CallOption) (*Artist, error) {
	out := new(Artist)
	err := c.cc.Invoke(ctx, "/Presentation/GetArtist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *presentationClient) GetRelease(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*Release, error) {
	out := new(Release)
	err := c.cc.Invoke(ctx, "/Presentation/GetRelease", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PresentationServer is the server API for Presentation service.
// All implementations must embed UnimplementedPresentationServer
// for forward compatibility
type PresentationServer interface {
	GetArtist(context.Context, *ArtistRequest) (*Artist, error)
	GetRelease(context.Context, *ReleaseRequest) (*Release, error)
	mustEmbedUnimplementedPresentationServer()
}

// UnimplementedPresentationServer must be embedded to have forward compatible implementations.
type UnimplementedPresentationServer struct {
}

func (UnimplementedPresentationServer) GetArtist(context.Context, *ArtistRequest) (*Artist, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArtist not implemented")
}
func (UnimplementedPresentationServer) GetRelease(context.Context, *ReleaseRequest) (*Release, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelease not implemented")
}
func (UnimplementedPresentationServer) mustEmbedUnimplementedPresentationServer() {}

// UnsafePresentationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PresentationServer will
// result in compilation errors.
type UnsafePresentationServer interface {
	mustEmbedUnimplementedPresentationServer()
}

func RegisterPresentationServer(s grpc.ServiceRegistrar, srv PresentationServer) {
	s.RegisterService(&Presentation_ServiceDesc, srv)
}

func _Presentation_GetArtist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArtistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresentationServer).GetArtist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Presentation/GetArtist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresentationServer).GetArtist(ctx, req.(*ArtistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Presentation_GetRelease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresentationServer).GetRelease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Presentation/GetRelease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresentationServer).GetRelease(ctx, req.(*ReleaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Presentation_ServiceDesc is the grpc.ServiceDesc for Presentation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Presentation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Presentation",
	HandlerType: (*PresentationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArtist",
			Handler:    _Presentation_GetArtist_Handler,
		},
		{
			MethodName: "GetRelease",
			Handler:    _Presentation_GetRelease_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/presentations/presentation.proto",
}

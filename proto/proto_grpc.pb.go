// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: proto/proto.proto

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

const (
	Replication_ConnectToServer_FullMethodName = "/Replication.Replication/ConnectToServer"
	Replication_Bid_FullMethodName             = "/Replication.Replication/Bid"
	Replication_Result_FullMethodName          = "/Replication.Replication/Result"
)

// ReplicationClient is the client API for Replication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationClient interface {
	ConnectToServer(ctx context.Context, in *User, opts ...grpc.CallOption) (Replication_ConnectToServerClient, error)
	Bid(ctx context.Context, in *PlaceBid, opts ...grpc.CallOption) (Replication_BidClient, error)
	Result(ctx context.Context, in *Close, opts ...grpc.CallOption) (Replication_ResultClient, error)
}

type replicationClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationClient(cc grpc.ClientConnInterface) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) ConnectToServer(ctx context.Context, in *User, opts ...grpc.CallOption) (Replication_ConnectToServerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Replication_ServiceDesc.Streams[0], Replication_ConnectToServer_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &replicationConnectToServerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Replication_ConnectToServerClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type replicationConnectToServerClient struct {
	grpc.ClientStream
}

func (x *replicationConnectToServerClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *replicationClient) Bid(ctx context.Context, in *PlaceBid, opts ...grpc.CallOption) (Replication_BidClient, error) {
	stream, err := c.cc.NewStream(ctx, &Replication_ServiceDesc.Streams[1], Replication_Bid_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &replicationBidClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Replication_BidClient interface {
	Recv() (*Acknowledgement, error)
	grpc.ClientStream
}

type replicationBidClient struct {
	grpc.ClientStream
}

func (x *replicationBidClient) Recv() (*Acknowledgement, error) {
	m := new(Acknowledgement)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *replicationClient) Result(ctx context.Context, in *Close, opts ...grpc.CallOption) (Replication_ResultClient, error) {
	stream, err := c.cc.NewStream(ctx, &Replication_ServiceDesc.Streams[2], Replication_Result_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &replicationResultClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Replication_ResultClient interface {
	Recv() (*Outcome, error)
	grpc.ClientStream
}

type replicationResultClient struct {
	grpc.ClientStream
}

func (x *replicationResultClient) Recv() (*Outcome, error) {
	m := new(Outcome)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReplicationServer is the server API for Replication service.
// All implementations must embed UnimplementedReplicationServer
// for forward compatibility
type ReplicationServer interface {
	ConnectToServer(*User, Replication_ConnectToServerServer) error
	Bid(*PlaceBid, Replication_BidServer) error
	Result(*Close, Replication_ResultServer) error
	mustEmbedUnimplementedReplicationServer()
}

// UnimplementedReplicationServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (UnimplementedReplicationServer) ConnectToServer(*User, Replication_ConnectToServerServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectToServer not implemented")
}
func (UnimplementedReplicationServer) Bid(*PlaceBid, Replication_BidServer) error {
	return status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedReplicationServer) Result(*Close, Replication_ResultServer) error {
	return status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedReplicationServer) mustEmbedUnimplementedReplicationServer() {}

// UnsafeReplicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServer will
// result in compilation errors.
type UnsafeReplicationServer interface {
	mustEmbedUnimplementedReplicationServer()
}

func RegisterReplicationServer(s grpc.ServiceRegistrar, srv ReplicationServer) {
	s.RegisterService(&Replication_ServiceDesc, srv)
}

func _Replication_ConnectToServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReplicationServer).ConnectToServer(m, &replicationConnectToServerServer{stream})
}

type Replication_ConnectToServerServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type replicationConnectToServerServer struct {
	grpc.ServerStream
}

func (x *replicationConnectToServerServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _Replication_Bid_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PlaceBid)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReplicationServer).Bid(m, &replicationBidServer{stream})
}

type Replication_BidServer interface {
	Send(*Acknowledgement) error
	grpc.ServerStream
}

type replicationBidServer struct {
	grpc.ServerStream
}

func (x *replicationBidServer) Send(m *Acknowledgement) error {
	return x.ServerStream.SendMsg(m)
}

func _Replication_Result_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Close)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReplicationServer).Result(m, &replicationResultServer{stream})
}

type Replication_ResultServer interface {
	Send(*Outcome) error
	grpc.ServerStream
}

type replicationResultServer struct {
	grpc.ServerStream
}

func (x *replicationResultServer) Send(m *Outcome) error {
	return x.ServerStream.SendMsg(m)
}

// Replication_ServiceDesc is the grpc.ServiceDesc for Replication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Replication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Replication.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectToServer",
			Handler:       _Replication_ConnectToServer_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Bid",
			Handler:       _Replication_Bid_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Result",
			Handler:       _Replication_Result_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/proto.proto",
}

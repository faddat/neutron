// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: interchaintxs/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgRegisterInterchainAccount is used to register an account on a remote zone.
type MsgRegisterInterchainAccount struct {
	FromAddress  string `protobuf:"bytes,1,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ConnectionId string `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	Owner        string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty" yaml:"owner"`
}

func (m *MsgRegisterInterchainAccount) Reset()         { *m = MsgRegisterInterchainAccount{} }
func (m *MsgRegisterInterchainAccount) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterInterchainAccount) ProtoMessage()    {}
func (*MsgRegisterInterchainAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_ecd987b66c8800e1, []int{0}
}
func (m *MsgRegisterInterchainAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterInterchainAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterInterchainAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterInterchainAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterInterchainAccount.Merge(m, src)
}
func (m *MsgRegisterInterchainAccount) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterInterchainAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterInterchainAccount.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterInterchainAccount proto.InternalMessageInfo

// MsgRegisterInterchainAccountResponse is the response type for
// MsgRegisterInterchainAccount.
type MsgRegisterInterchainAccountResponse struct {
}

func (m *MsgRegisterInterchainAccountResponse) Reset()         { *m = MsgRegisterInterchainAccountResponse{} }
func (m *MsgRegisterInterchainAccountResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterInterchainAccountResponse) ProtoMessage()    {}
func (*MsgRegisterInterchainAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ecd987b66c8800e1, []int{1}
}
func (m *MsgRegisterInterchainAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterInterchainAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterInterchainAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterInterchainAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterInterchainAccountResponse.Merge(m, src)
}
func (m *MsgRegisterInterchainAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterInterchainAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterInterchainAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterInterchainAccountResponse proto.InternalMessageInfo

// MsgSubmitTx defines the payload for Msg/SubmitTx
type MsgSubmitTx struct {
	FromAddress string `protobuf:"bytes,1,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	// Owner is supposed to be the unique identifier, e.g., neutron/kava.
	// This identifier will be a part of the portID that we'll claim our
	// capability for.
	Owner        string       `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	ConnectionId string       `protobuf:"bytes,3,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	Msgs         []*types.Any `protobuf:"bytes,4,rep,name=msgs,proto3" json:"msgs,omitempty"`
	Memo         string       `protobuf:"bytes,5,opt,name=memo,proto3" json:"memo,omitempty"`
}

func (m *MsgSubmitTx) Reset()         { *m = MsgSubmitTx{} }
func (m *MsgSubmitTx) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitTx) ProtoMessage()    {}
func (*MsgSubmitTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_ecd987b66c8800e1, []int{2}
}
func (m *MsgSubmitTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitTx.Merge(m, src)
}
func (m *MsgSubmitTx) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitTx) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitTx.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitTx proto.InternalMessageInfo

// MsgSubmitTxResponse defines the response for Msg/SubmitTx
type MsgSubmitTxResponse struct {
}

func (m *MsgSubmitTxResponse) Reset()         { *m = MsgSubmitTxResponse{} }
func (m *MsgSubmitTxResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitTxResponse) ProtoMessage()    {}
func (*MsgSubmitTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ecd987b66c8800e1, []int{3}
}
func (m *MsgSubmitTxResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitTxResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitTxResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitTxResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitTxResponse.Merge(m, src)
}
func (m *MsgSubmitTxResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitTxResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitTxResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitTxResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgRegisterInterchainAccount)(nil), "neutron.interchainadapter.interchaintxs.v1.MsgRegisterInterchainAccount")
	proto.RegisterType((*MsgRegisterInterchainAccountResponse)(nil), "neutron.interchainadapter.interchaintxs.v1.MsgRegisterInterchainAccountResponse")
	proto.RegisterType((*MsgSubmitTx)(nil), "neutron.interchainadapter.interchaintxs.v1.MsgSubmitTx")
	proto.RegisterType((*MsgSubmitTxResponse)(nil), "neutron.interchainadapter.interchaintxs.v1.MsgSubmitTxResponse")
}

func init() { proto.RegisterFile("interchaintxs/v1/tx.proto", fileDescriptor_ecd987b66c8800e1) }

var fileDescriptor_ecd987b66c8800e1 = []byte{
	// 537 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x41, 0x6b, 0x13, 0x41,
	0x18, 0xcd, 0x24, 0xa9, 0xd4, 0x49, 0x04, 0x59, 0x53, 0xd8, 0x84, 0xb2, 0xa9, 0xab, 0x94, 0x10,
	0x70, 0x86, 0x56, 0x41, 0x09, 0x14, 0x49, 0xbd, 0x58, 0x21, 0x2a, 0xab, 0x5e, 0xbc, 0x84, 0xc9,
	0x66, 0x3a, 0x59, 0xe8, 0xce, 0x2c, 0x33, 0x93, 0xba, 0xb9, 0x7a, 0xf2, 0x24, 0x82, 0x7f, 0xa0,
	0xbf, 0x41, 0xff, 0x84, 0xde, 0x0a, 0x5e, 0x3c, 0x15, 0x49, 0x3c, 0x78, 0xee, 0x2f, 0x90, 0xcc,
	0xee, 0xa6, 0x49, 0xc4, 0xd2, 0x8a, 0xb7, 0x99, 0xf7, 0xcd, 0x7b, 0xdf, 0x7b, 0xf3, 0xcd, 0xc0,
	0x6a, 0xc0, 0x35, 0x95, 0xfe, 0x80, 0x04, 0x5c, 0xc7, 0x0a, 0x1f, 0x6e, 0x61, 0x1d, 0xa3, 0x48,
	0x0a, 0x2d, 0xac, 0x26, 0xa7, 0x43, 0x2d, 0x05, 0x47, 0x67, 0x47, 0x48, 0x9f, 0x44, 0x9a, 0x4a,
	0xb4, 0x40, 0x42, 0x87, 0x5b, 0xb5, 0xaa, 0x2f, 0x54, 0x28, 0x54, 0xd7, 0x30, 0x71, 0xb2, 0x49,
	0x64, 0x6a, 0x15, 0x26, 0x98, 0x48, 0xf0, 0xe9, 0x2a, 0x45, 0xd7, 0x98, 0x10, 0xec, 0x80, 0x62,
	0x12, 0x05, 0x78, 0xa0, 0x75, 0x94, 0xc2, 0xeb, 0x73, 0x30, 0xe1, 0x5c, 0x68, 0xa2, 0x03, 0xc1,
	0x33, 0xa9, 0x6a, 0x5a, 0x35, 0xbb, 0xde, 0x70, 0x1f, 0x13, 0x3e, 0x4a, 0x4a, 0xee, 0x27, 0x00,
	0xd7, 0x3b, 0x8a, 0x79, 0x94, 0x05, 0x4a, 0x53, 0xb9, 0x37, 0x33, 0xd8, 0xf6, 0x7d, 0x31, 0xe4,
	0xda, 0xba, 0x09, 0xcb, 0xfb, 0x52, 0x84, 0x5d, 0xd2, 0xef, 0x4b, 0xaa, 0x94, 0x0d, 0x36, 0x40,
	0xe3, 0xaa, 0x57, 0x9a, 0x62, 0xed, 0x04, 0xb2, 0x76, 0xe0, 0x35, 0x5f, 0x70, 0x4e, 0xfd, 0x69,
	0xcf, 0x6e, 0xd0, 0xb7, 0xf3, 0xd3, 0x33, 0xbb, 0xf6, 0xe9, 0x49, 0xbd, 0x32, 0x22, 0xe1, 0x41,
	0xcb, 0x5d, 0x28, 0xbb, 0x5e, 0xf9, 0x6c, 0xbf, 0xd7, 0xb7, 0x36, 0xe1, 0x8a, 0x78, 0xc3, 0xa9,
	0xb4, 0x0b, 0x86, 0x76, 0xfd, 0xf4, 0xa4, 0x5e, 0x4e, 0x68, 0x06, 0x76, 0xbd, 0xa4, 0xdc, 0x5a,
	0x7d, 0x77, 0x54, 0xcf, 0xfd, 0x3a, 0xaa, 0xe7, 0xdc, 0x4d, 0x78, 0xfb, 0x3c, 0xcf, 0x1e, 0x55,
	0x91, 0xe0, 0x8a, 0xba, 0x9f, 0x01, 0x2c, 0x75, 0x14, 0x7b, 0x31, 0xec, 0x85, 0x81, 0x7e, 0x19,
	0x5f, 0x24, 0x4b, 0x25, 0x33, 0x63, 0x32, 0xa4, 0xad, 0xad, 0x5b, 0xcb, 0x09, 0x8d, 0xd5, 0xa5,
	0x1c, 0x0d, 0x58, 0x0c, 0x15, 0x53, 0x76, 0x71, 0xa3, 0xd0, 0x28, 0x6d, 0x57, 0x50, 0x72, 0xe9,
	0x28, 0xbb, 0x74, 0xd4, 0xe6, 0x23, 0xcf, 0x9c, 0xb0, 0x2c, 0x58, 0x0c, 0x69, 0x28, 0xec, 0x15,
	0xa3, 0x62, 0xd6, 0x73, 0xe9, 0xd6, 0xe0, 0x8d, 0x39, 0xd3, 0x59, 0x98, 0xed, 0x71, 0x01, 0x16,
	0x3a, 0x8a, 0x59, 0xef, 0xf3, 0xb0, 0xfa, 0xf7, 0x71, 0x3d, 0x46, 0x17, 0x7f, 0x7d, 0xe8, 0xbc,
	0x4b, 0xac, 0x3d, 0xff, 0x5f, 0x4a, 0xb3, 0x71, 0xbc, 0x7a, 0xfb, 0xed, 0xe7, 0xc7, 0xfc, 0x33,
	0xf7, 0x09, 0x4e, 0x95, 0xf1, 0x1f, 0xca, 0x78, 0xf1, 0x5b, 0xe9, 0x78, 0xfa, 0xb3, 0x64, 0xaa,
	0x3c, 0x57, 0xc4, 0x24, 0xd1, 0x6e, 0x81, 0xa6, 0xf5, 0x15, 0xc0, 0xd5, 0xd9, 0x88, 0xef, 0x5f,
	0xd2, 0x75, 0x46, 0xac, 0x3d, 0xfc, 0x47, 0xe2, 0x2c, 0xdd, 0x23, 0x93, 0x6e, 0xc7, 0x7d, 0x70,
	0xc9, 0x74, 0xca, 0x08, 0x61, 0x1d, 0xb7, 0x40, 0x73, 0xf7, 0xe9, 0x97, 0xb1, 0x03, 0x8e, 0xc7,
	0x0e, 0xf8, 0x31, 0x76, 0xc0, 0x87, 0x89, 0x93, 0x3b, 0x9e, 0x38, 0xb9, 0xef, 0x13, 0x27, 0xf7,
	0xfa, 0x1e, 0x0b, 0xf4, 0x60, 0xd8, 0x43, 0xbe, 0x08, 0xb3, 0x06, 0x77, 0x84, 0x64, 0xb3, 0x66,
	0xf1, 0xb2, 0xf8, 0x28, 0xa2, 0xaa, 0x77, 0xc5, 0xbc, 0xbe, 0xbb, 0xbf, 0x03, 0x00, 0x00, 0xff,
	0xff, 0x58, 0xe2, 0xf9, 0x17, 0xaf, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	RegisterInterchainAccount(ctx context.Context, in *MsgRegisterInterchainAccount, opts ...grpc.CallOption) (*MsgRegisterInterchainAccountResponse, error)
	SubmitTx(ctx context.Context, in *MsgSubmitTx, opts ...grpc.CallOption) (*MsgSubmitTxResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterInterchainAccount(ctx context.Context, in *MsgRegisterInterchainAccount, opts ...grpc.CallOption) (*MsgRegisterInterchainAccountResponse, error) {
	out := new(MsgRegisterInterchainAccountResponse)
	err := c.cc.Invoke(ctx, "/neutron.interchainadapter.interchaintxs.v1.Msg/RegisterInterchainAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SubmitTx(ctx context.Context, in *MsgSubmitTx, opts ...grpc.CallOption) (*MsgSubmitTxResponse, error) {
	out := new(MsgSubmitTxResponse)
	err := c.cc.Invoke(ctx, "/neutron.interchainadapter.interchaintxs.v1.Msg/SubmitTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	RegisterInterchainAccount(context.Context, *MsgRegisterInterchainAccount) (*MsgRegisterInterchainAccountResponse, error)
	SubmitTx(context.Context, *MsgSubmitTx) (*MsgSubmitTxResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) RegisterInterchainAccount(ctx context.Context, req *MsgRegisterInterchainAccount) (*MsgRegisterInterchainAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterInterchainAccount not implemented")
}
func (*UnimplementedMsgServer) SubmitTx(ctx context.Context, req *MsgSubmitTx) (*MsgSubmitTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitTx not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterInterchainAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterInterchainAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterInterchainAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutron.interchainadapter.interchaintxs.v1.Msg/RegisterInterchainAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterInterchainAccount(ctx, req.(*MsgRegisterInterchainAccount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neutron.interchainadapter.interchaintxs.v1.Msg/SubmitTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitTx(ctx, req.(*MsgSubmitTx))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "neutron.interchainadapter.interchaintxs.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterInterchainAccount",
			Handler:    _Msg_RegisterInterchainAccount_Handler,
		},
		{
			MethodName: "SubmitTx",
			Handler:    _Msg_SubmitTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interchaintxs/v1/tx.proto",
}

func (m *MsgRegisterInterchainAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterInterchainAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterInterchainAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.FromAddress) > 0 {
		i -= len(m.FromAddress)
		copy(dAtA[i:], m.FromAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.FromAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRegisterInterchainAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterInterchainAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterInterchainAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSubmitTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Memo) > 0 {
		i -= len(m.Memo)
		copy(dAtA[i:], m.Memo)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Memo)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Msgs) > 0 {
		for iNdEx := len(m.Msgs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Msgs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.FromAddress) > 0 {
		i -= len(m.FromAddress)
		copy(dAtA[i:], m.FromAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.FromAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSubmitTxResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitTxResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgRegisterInterchainAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FromAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRegisterInterchainAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSubmitTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FromAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.Msgs) > 0 {
		for _, e := range m.Msgs {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Memo)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSubmitTxResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgRegisterInterchainAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRegisterInterchainAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterInterchainAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgRegisterInterchainAccountResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRegisterInterchainAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterInterchainAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgSubmitTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSubmitTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msgs = append(m.Msgs, &types.Any{})
			if err := m.Msgs[len(m.Msgs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Memo", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Memo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgSubmitTxResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSubmitTxResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitTxResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)

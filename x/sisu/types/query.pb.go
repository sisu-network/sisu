// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/query.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type QueryAllPubKeysRequest struct {
}

func (m *QueryAllPubKeysRequest) Reset()         { *m = QueryAllPubKeysRequest{} }
func (m *QueryAllPubKeysRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllPubKeysRequest) ProtoMessage()    {}
func (*QueryAllPubKeysRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{0}
}
func (m *QueryAllPubKeysRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllPubKeysRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllPubKeysRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllPubKeysRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllPubKeysRequest.Merge(m, src)
}
func (m *QueryAllPubKeysRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllPubKeysRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllPubKeysRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllPubKeysRequest proto.InternalMessageInfo

type QueryAllPubKeysResponse struct {
	Pubkeys map[string][]byte `protobuf:"bytes,1,rep,name=pubkeys,proto3" json:"pubkeys,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *QueryAllPubKeysResponse) Reset()         { *m = QueryAllPubKeysResponse{} }
func (m *QueryAllPubKeysResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllPubKeysResponse) ProtoMessage()    {}
func (*QueryAllPubKeysResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{1}
}
func (m *QueryAllPubKeysResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllPubKeysResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllPubKeysResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllPubKeysResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllPubKeysResponse.Merge(m, src)
}
func (m *QueryAllPubKeysResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllPubKeysResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllPubKeysResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllPubKeysResponse proto.InternalMessageInfo

func (m *QueryAllPubKeysResponse) GetPubkeys() map[string][]byte {
	if m != nil {
		return m.Pubkeys
	}
	return nil
}

type QueryContractRequest struct {
	Chain string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Hash  string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *QueryContractRequest) Reset()         { *m = QueryContractRequest{} }
func (m *QueryContractRequest) String() string { return proto.CompactTextString(m) }
func (*QueryContractRequest) ProtoMessage()    {}
func (*QueryContractRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{2}
}
func (m *QueryContractRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryContractRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryContractRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryContractRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryContractRequest.Merge(m, src)
}
func (m *QueryContractRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryContractRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryContractRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryContractRequest proto.InternalMessageInfo

func (m *QueryContractRequest) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *QueryContractRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type QueryContractResponse struct {
	Contract *Contract `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`
}

func (m *QueryContractResponse) Reset()         { *m = QueryContractResponse{} }
func (m *QueryContractResponse) String() string { return proto.CompactTextString(m) }
func (*QueryContractResponse) ProtoMessage()    {}
func (*QueryContractResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{3}
}
func (m *QueryContractResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryContractResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryContractResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryContractResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryContractResponse.Merge(m, src)
}
func (m *QueryContractResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryContractResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryContractResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryContractResponse proto.InternalMessageInfo

func (m *QueryContractResponse) GetContract() *Contract {
	if m != nil {
		return m.Contract
	}
	return nil
}

type QueryTokenRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryTokenRequest) Reset()         { *m = QueryTokenRequest{} }
func (m *QueryTokenRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTokenRequest) ProtoMessage()    {}
func (*QueryTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{4}
}
func (m *QueryTokenRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTokenRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTokenRequest.Merge(m, src)
}
func (m *QueryTokenRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTokenRequest proto.InternalMessageInfo

func (m *QueryTokenRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type QueryTokenResponse struct {
	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *QueryTokenResponse) Reset()         { *m = QueryTokenResponse{} }
func (m *QueryTokenResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTokenResponse) ProtoMessage()    {}
func (*QueryTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a22fc0c3dad4d31a, []int{5}
}
func (m *QueryTokenResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTokenResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTokenResponse.Merge(m, src)
}
func (m *QueryTokenResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTokenResponse proto.InternalMessageInfo

func (m *QueryTokenResponse) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryAllPubKeysRequest)(nil), "types.QueryAllPubKeysRequest")
	proto.RegisterType((*QueryAllPubKeysResponse)(nil), "types.QueryAllPubKeysResponse")
	proto.RegisterMapType((map[string][]byte)(nil), "types.QueryAllPubKeysResponse.PubkeysEntry")
	proto.RegisterType((*QueryContractRequest)(nil), "types.QueryContractRequest")
	proto.RegisterType((*QueryContractResponse)(nil), "types.QueryContractResponse")
	proto.RegisterType((*QueryTokenRequest)(nil), "types.QueryTokenRequest")
	proto.RegisterType((*QueryTokenResponse)(nil), "types.QueryTokenResponse")
}

func init() { proto.RegisterFile("sisu/query.proto", fileDescriptor_a22fc0c3dad4d31a) }

var fileDescriptor_a22fc0c3dad4d31a = []byte{
	// 416 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xcd, 0xae, 0xd2, 0x40,
	0x14, 0xee, 0xf4, 0x8a, 0x5e, 0x0e, 0xa8, 0x38, 0xa2, 0xd6, 0xaa, 0x0d, 0x19, 0x37, 0x18, 0x62,
	0x49, 0x70, 0x43, 0x58, 0xa9, 0xc8, 0x4a, 0x4d, 0xb0, 0x61, 0xe5, 0xae, 0x2d, 0x13, 0xdb, 0x14,
	0xdb, 0xd2, 0x99, 0xaa, 0x7d, 0x0b, 0x1f, 0xc0, 0x07, 0x72, 0xc9, 0xd2, 0xa5, 0x81, 0xad, 0x0f,
	0x61, 0x3a, 0x33, 0x85, 0x8a, 0x70, 0x57, 0x9d, 0xd3, 0xef, 0xe7, 0x9c, 0xf3, 0xcd, 0x40, 0x87,
	0x85, 0x2c, 0x1f, 0xae, 0x73, 0x9a, 0x15, 0x76, 0x9a, 0x25, 0x3c, 0xc1, 0x0d, 0x5e, 0xa4, 0x94,
	0x99, 0x77, 0x05, 0xe0, 0x27, 0x31, 0xcf, 0x5c, 0x9f, 0x4b, 0xcc, 0x94, 0x6c, 0x9e, 0x44, 0x34,
	0x96, 0x7f, 0x88, 0x01, 0xf7, 0x3f, 0x94, 0xe2, 0x57, 0xab, 0xd5, 0x3c, 0xf7, 0xde, 0xd2, 0x82,
	0x39, 0x74, 0x9d, 0x53, 0xc6, 0xc9, 0x0f, 0x04, 0x0f, 0xfe, 0x83, 0x58, 0x9a, 0xc4, 0x8c, 0xe2,
	0x19, 0xdc, 0x48, 0x73, 0x2f, 0xa2, 0x05, 0x33, 0x50, 0xef, 0xa2, 0xdf, 0x1a, 0x0d, 0x6c, 0xd1,
	0xd5, 0x3e, 0x23, 0xb0, 0xe7, 0x92, 0x3d, 0x8b, 0x79, 0x56, 0x38, 0x95, 0xd6, 0x9c, 0x40, 0xbb,
	0x0e, 0xe0, 0x0e, 0x5c, 0x44, 0xb4, 0x30, 0x50, 0x0f, 0xf5, 0x9b, 0x4e, 0x79, 0xc4, 0x5d, 0x68,
	0x7c, 0x71, 0x57, 0x39, 0x35, 0xf4, 0x1e, 0xea, 0xb7, 0x1d, 0x59, 0x4c, 0xf4, 0x31, 0x22, 0x2f,
	0xa1, 0x2b, 0x9a, 0x4d, 0xd5, 0x86, 0x6a, 0xec, 0x52, 0xe1, 0x07, 0x6e, 0x18, 0x2b, 0x17, 0x59,
	0x60, 0x0c, 0xd7, 0x02, 0x97, 0x05, 0xc2, 0xa6, 0xe9, 0x88, 0x33, 0x79, 0x03, 0xf7, 0x8e, 0x1c,
	0xd4, 0x76, 0x03, 0xb8, 0xac, 0x72, 0x13, 0x2e, 0xad, 0xd1, 0x6d, 0xb5, 0xde, 0x9e, 0xba, 0x27,
	0x90, 0xa7, 0x70, 0x47, 0xb8, 0x2c, 0xca, 0x50, 0xab, 0x21, 0x6e, 0x81, 0x1e, 0x2e, 0xd5, 0x04,
	0x7a, 0xb8, 0x24, 0x63, 0xc0, 0x75, 0x92, 0xea, 0x43, 0xa0, 0x21, 0xae, 0x42, 0x35, 0x69, 0xab,
	0x26, 0x92, 0x24, 0xa1, 0xd1, 0x1f, 0x04, 0x97, 0x0b, 0xc6, 0x84, 0x1a, 0xbf, 0x07, 0x38, 0x64,
	0x8b, 0x9f, 0x9c, 0xcb, 0x5c, 0xcc, 0x60, 0x5a, 0x57, 0x5f, 0x09, 0xd1, 0xf0, 0x3b, 0xb8, 0xf9,
	0x4f, 0x00, 0xf8, 0x51, 0x5d, 0x72, 0x14, 0xac, 0xf9, 0xf8, 0x34, 0xb8, 0x77, 0x9b, 0x02, 0x1c,
	0x76, 0xc4, 0x46, 0x9d, 0x5d, 0xcf, 0xc6, 0x7c, 0x78, 0x02, 0xa9, 0x4c, 0x5e, 0x4f, 0x7f, 0x6e,
	0x2d, 0xb4, 0xd9, 0x5a, 0xe8, 0xf7, 0xd6, 0x42, 0xdf, 0x77, 0x96, 0xb6, 0xd9, 0x59, 0xda, 0xaf,
	0x9d, 0xa5, 0x7d, 0x7c, 0xf6, 0x29, 0xe4, 0x41, 0xee, 0xd9, 0x7e, 0xf2, 0x79, 0x58, 0xbe, 0xe2,
	0xe7, 0x31, 0xe5, 0x5f, 0x93, 0x2c, 0x12, 0xc5, 0xf0, 0x9b, 0xfc, 0x08, 0x67, 0xef, 0xba, 0x78,
	0xda, 0x2f, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb6, 0xa8, 0xa6, 0xa3, 0x1c, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TssQueryClient is the client API for TssQuery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TssQueryClient interface {
	AllPubKeys(ctx context.Context, in *QueryAllPubKeysRequest, opts ...grpc.CallOption) (*QueryAllPubKeysResponse, error)
	QueryContract(ctx context.Context, in *QueryContractRequest, opts ...grpc.CallOption) (*QueryContractResponse, error)
	QueryToken(ctx context.Context, in *QueryTokenRequest, opts ...grpc.CallOption) (*QueryTokenResponse, error)
}

type tssQueryClient struct {
	cc grpc1.ClientConn
}

func NewTssQueryClient(cc grpc1.ClientConn) TssQueryClient {
	return &tssQueryClient{cc}
}

func (c *tssQueryClient) AllPubKeys(ctx context.Context, in *QueryAllPubKeysRequest, opts ...grpc.CallOption) (*QueryAllPubKeysResponse, error) {
	out := new(QueryAllPubKeysResponse)
	err := c.cc.Invoke(ctx, "/types.TssQuery/AllPubKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tssQueryClient) QueryContract(ctx context.Context, in *QueryContractRequest, opts ...grpc.CallOption) (*QueryContractResponse, error) {
	out := new(QueryContractResponse)
	err := c.cc.Invoke(ctx, "/types.TssQuery/QueryContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tssQueryClient) QueryToken(ctx context.Context, in *QueryTokenRequest, opts ...grpc.CallOption) (*QueryTokenResponse, error) {
	out := new(QueryTokenResponse)
	err := c.cc.Invoke(ctx, "/types.TssQuery/QueryToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TssQueryServer is the server API for TssQuery service.
type TssQueryServer interface {
	AllPubKeys(context.Context, *QueryAllPubKeysRequest) (*QueryAllPubKeysResponse, error)
	QueryContract(context.Context, *QueryContractRequest) (*QueryContractResponse, error)
	QueryToken(context.Context, *QueryTokenRequest) (*QueryTokenResponse, error)
}

// UnimplementedTssQueryServer can be embedded to have forward compatible implementations.
type UnimplementedTssQueryServer struct {
}

func (*UnimplementedTssQueryServer) AllPubKeys(ctx context.Context, req *QueryAllPubKeysRequest) (*QueryAllPubKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllPubKeys not implemented")
}
func (*UnimplementedTssQueryServer) QueryContract(ctx context.Context, req *QueryContractRequest) (*QueryContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryContract not implemented")
}
func (*UnimplementedTssQueryServer) QueryToken(ctx context.Context, req *QueryTokenRequest) (*QueryTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryToken not implemented")
}

func RegisterTssQueryServer(s grpc1.Server, srv TssQueryServer) {
	s.RegisterService(&_TssQuery_serviceDesc, srv)
}

func _TssQuery_AllPubKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllPubKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TssQueryServer).AllPubKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.TssQuery/AllPubKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TssQueryServer).AllPubKeys(ctx, req.(*QueryAllPubKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TssQuery_QueryContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryContractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TssQueryServer).QueryContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.TssQuery/QueryContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TssQueryServer).QueryContract(ctx, req.(*QueryContractRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TssQuery_QueryToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TssQueryServer).QueryToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.TssQuery/QueryToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TssQueryServer).QueryToken(ctx, req.(*QueryTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TssQuery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "types.TssQuery",
	HandlerType: (*TssQueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AllPubKeys",
			Handler:    _TssQuery_AllPubKeys_Handler,
		},
		{
			MethodName: "QueryContract",
			Handler:    _TssQuery_QueryContract_Handler,
		},
		{
			MethodName: "QueryToken",
			Handler:    _TssQuery_QueryToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sisu/query.proto",
}

func (m *QueryAllPubKeysRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllPubKeysRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllPubKeysRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryAllPubKeysResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllPubKeysResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllPubKeysResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Pubkeys) > 0 {
		for k := range m.Pubkeys {
			v := m.Pubkeys[k]
			baseI := i
			if len(v) > 0 {
				i -= len(v)
				copy(dAtA[i:], v)
				i = encodeVarintQuery(dAtA, i, uint64(len(v)))
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintQuery(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintQuery(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryContractRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryContractRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryContractRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryContractResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryContractResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryContractResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Contract != nil {
		{
			size, err := m.Contract.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTokenRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTokenRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTokenRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTokenResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTokenResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTokenResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Token != nil {
		{
			size, err := m.Token.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryAllPubKeysRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryAllPubKeysResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pubkeys) > 0 {
		for k, v := range m.Pubkeys {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovQuery(uint64(len(v)))
			}
			mapEntrySize := 1 + len(k) + sovQuery(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovQuery(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *QueryContractRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryContractResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Contract != nil {
		l = m.Contract.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTokenRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTokenResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Token != nil {
		l = m.Token.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryAllPubKeysRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAllPubKeysRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllPubKeysRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryAllPubKeysResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAllPubKeysResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllPubKeysResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pubkeys == nil {
				m.Pubkeys = make(map[string][]byte)
			}
			var mapkey string
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowQuery
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowQuery
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthQuery
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthQuery
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowQuery
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthQuery
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex < 0 {
						return ErrInvalidLengthQuery
					}
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipQuery(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthQuery
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Pubkeys[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryContractRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryContractRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryContractRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryContractResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryContractResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryContractResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Contract == nil {
				m.Contract = &Contract{}
			}
			if err := m.Contract.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTokenRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTokenRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTokenRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTokenResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTokenResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTokenResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Token == nil {
				m.Token = &Token{}
			}
			if err := m.Token.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)

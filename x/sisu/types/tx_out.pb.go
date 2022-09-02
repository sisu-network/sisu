// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/tx_out.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type TxOutType int32

const (
	TxOutType_TRANSFER_OUT TxOutType = 0
)

var TxOutType_name = map[int32]string{
	0: "TRANSFER_OUT",
}

var TxOutType_value = map[string]int32{
	"TRANSFER_OUT": 0,
}

func (x TxOutType) String() string {
	return proto.EnumName(TxOutType_name, int32(x))
}

func (TxOutType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{0}
}

type TxOutMsg struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *TxOut `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *TxOutMsg) Reset()         { *m = TxOutMsg{} }
func (m *TxOutMsg) String() string { return proto.CompactTextString(m) }
func (*TxOutMsg) ProtoMessage()    {}
func (*TxOutMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{0}
}
func (m *TxOutMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxOutMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxOutMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxOutMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxOutMsg.Merge(m, src)
}
func (m *TxOutMsg) XXX_Size() int {
	return m.Size()
}
func (m *TxOutMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_TxOutMsg.DiscardUnknown(m)
}

var xxx_messageInfo_TxOutMsg proto.InternalMessageInfo

func (m *TxOutMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *TxOutMsg) GetData() *TxOut {
	if m != nil {
		return m.Data
	}
	return nil
}

type TxOut struct {
	TxType  TxOutType     `protobuf:"varint,1,opt,name=txType,proto3,enum=types.TxOutType" json:"txType,omitempty"`
	Content *TxOutContent `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Input   *TxOutInput   `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
}

func (m *TxOut) Reset()         { *m = TxOut{} }
func (m *TxOut) String() string { return proto.CompactTextString(m) }
func (*TxOut) ProtoMessage()    {}
func (*TxOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{1}
}
func (m *TxOut) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxOut.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxOut.Merge(m, src)
}
func (m *TxOut) XXX_Size() int {
	return m.Size()
}
func (m *TxOut) XXX_DiscardUnknown() {
	xxx_messageInfo_TxOut.DiscardUnknown(m)
}

var xxx_messageInfo_TxOut proto.InternalMessageInfo

func (m *TxOut) GetTxType() TxOutType {
	if m != nil {
		return m.TxType
	}
	return TxOutType_TRANSFER_OUT
}

func (m *TxOut) GetContent() *TxOutContent {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *TxOut) GetInput() *TxOutInput {
	if m != nil {
		return m.Input
	}
	return nil
}

type TxOutContent struct {
	OutChain string `protobuf:"bytes,1,opt,name=outChain,proto3" json:"outChain,omitempty"`
	OutHash  string `protobuf:"bytes,2,opt,name=outHash,proto3" json:"outHash,omitempty"`
	OutBytes []byte `protobuf:"bytes,3,opt,name=outBytes,proto3" json:"outBytes,omitempty"`
}

func (m *TxOutContent) Reset()         { *m = TxOutContent{} }
func (m *TxOutContent) String() string { return proto.CompactTextString(m) }
func (*TxOutContent) ProtoMessage()    {}
func (*TxOutContent) Descriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{2}
}
func (m *TxOutContent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxOutContent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxOutContent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxOutContent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxOutContent.Merge(m, src)
}
func (m *TxOutContent) XXX_Size() int {
	return m.Size()
}
func (m *TxOutContent) XXX_DiscardUnknown() {
	xxx_messageInfo_TxOutContent.DiscardUnknown(m)
}

var xxx_messageInfo_TxOutContent proto.InternalMessageInfo

func (m *TxOutContent) GetOutChain() string {
	if m != nil {
		return m.OutChain
	}
	return ""
}

func (m *TxOutContent) GetOutHash() string {
	if m != nil {
		return m.OutHash
	}
	return ""
}

func (m *TxOutContent) GetOutBytes() []byte {
	if m != nil {
		return m.OutBytes
	}
	return nil
}

type TxOutInput struct {
	// For transferOut
	TransferOutIds []string `protobuf:"bytes,1,rep,name=transferOutIds,proto3" json:"transferOutIds,omitempty"`
}

func (m *TxOutInput) Reset()         { *m = TxOutInput{} }
func (m *TxOutInput) String() string { return proto.CompactTextString(m) }
func (*TxOutInput) ProtoMessage()    {}
func (*TxOutInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{3}
}
func (m *TxOutInput) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxOutInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxOutInput.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxOutInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxOutInput.Merge(m, src)
}
func (m *TxOutInput) XXX_Size() int {
	return m.Size()
}
func (m *TxOutInput) XXX_DiscardUnknown() {
	xxx_messageInfo_TxOutInput.DiscardUnknown(m)
}

var xxx_messageInfo_TxOutInput proto.InternalMessageInfo

func (m *TxOutInput) GetTransferOutIds() []string {
	if m != nil {
		return m.TransferOutIds
	}
	return nil
}

// TxOut with and full transaction hash (including signature) to look up TxOut when a new tx comes in.
type TxOutSig struct {
	Chain       string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	HashWithSig string `protobuf:"bytes,2,opt,name=hashWithSig,proto3" json:"hashWithSig,omitempty"`
	HashNoSig   string `protobuf:"bytes,3,opt,name=hashNoSig,proto3" json:"hashNoSig,omitempty"`
}

func (m *TxOutSig) Reset()         { *m = TxOutSig{} }
func (m *TxOutSig) String() string { return proto.CompactTextString(m) }
func (*TxOutSig) ProtoMessage()    {}
func (*TxOutSig) Descriptor() ([]byte, []int) {
	return fileDescriptor_155b13ca5b94a7d7, []int{4}
}
func (m *TxOutSig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxOutSig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxOutSig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxOutSig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxOutSig.Merge(m, src)
}
func (m *TxOutSig) XXX_Size() int {
	return m.Size()
}
func (m *TxOutSig) XXX_DiscardUnknown() {
	xxx_messageInfo_TxOutSig.DiscardUnknown(m)
}

var xxx_messageInfo_TxOutSig proto.InternalMessageInfo

func (m *TxOutSig) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *TxOutSig) GetHashWithSig() string {
	if m != nil {
		return m.HashWithSig
	}
	return ""
}

func (m *TxOutSig) GetHashNoSig() string {
	if m != nil {
		return m.HashNoSig
	}
	return ""
}

func init() {
	proto.RegisterEnum("types.TxOutType", TxOutType_name, TxOutType_value)
	proto.RegisterType((*TxOutMsg)(nil), "types.TxOutMsg")
	proto.RegisterType((*TxOut)(nil), "types.TxOut")
	proto.RegisterType((*TxOutContent)(nil), "types.TxOutContent")
	proto.RegisterType((*TxOutInput)(nil), "types.TxOutInput")
	proto.RegisterType((*TxOutSig)(nil), "types.TxOutSig")
}

func init() { proto.RegisterFile("sisu/tx_out.proto", fileDescriptor_155b13ca5b94a7d7) }

var fileDescriptor_155b13ca5b94a7d7 = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xc1, 0xca, 0xd3, 0x40,
	0x10, 0xc7, 0xb3, 0xd6, 0xf4, 0xfb, 0x32, 0x0d, 0xa5, 0x5d, 0x45, 0x82, 0x68, 0x08, 0x39, 0x68,
	0x14, 0x9a, 0x42, 0xf5, 0x05, 0x6c, 0x55, 0xf4, 0x60, 0x0b, 0xdb, 0x88, 0xe0, 0xa5, 0x4d, 0xdb,
	0x98, 0x2c, 0x62, 0x36, 0x64, 0x27, 0x98, 0x3e, 0x82, 0x37, 0x1f, 0xcb, 0x63, 0x8f, 0x1e, 0xa5,
	0x7d, 0x11, 0xc9, 0x26, 0x6d, 0xc3, 0x77, 0x0a, 0xbf, 0xff, 0xff, 0x3f, 0xb3, 0x93, 0x61, 0x60,
	0x28, 0xb9, 0x2c, 0xc6, 0x58, 0xae, 0x44, 0x81, 0x7e, 0x96, 0x0b, 0x14, 0x54, 0xc7, 0x7d, 0x16,
	0x49, 0xf7, 0x2d, 0xdc, 0x06, 0xe5, 0xa2, 0xc0, 0x4f, 0x32, 0xa6, 0x8f, 0xa0, 0x2b, 0x79, 0x9c,
	0x46, 0xb9, 0x45, 0x1c, 0xe2, 0x19, 0xac, 0x21, 0xea, 0xc0, 0xfd, 0x5d, 0x88, 0xa1, 0x75, 0xcf,
	0x21, 0x5e, 0x6f, 0x62, 0xfa, 0xaa, 0xd2, 0x57, 0x65, 0x4c, 0x39, 0xee, 0x2f, 0x02, 0xba, 0x62,
	0xea, 0x41, 0x17, 0xcb, 0x60, 0x9f, 0x45, 0xaa, 0x47, 0x7f, 0x32, 0x68, 0xa7, 0x2b, 0x9d, 0x35,
	0x3e, 0x1d, 0xc1, 0xcd, 0x56, 0xa4, 0x18, 0xa5, 0xd8, 0x34, 0x7e, 0xd0, 0x8e, 0xce, 0x6a, 0x8b,
	0x9d, 0x33, 0xf4, 0x39, 0xe8, 0x3c, 0xcd, 0x0a, 0xb4, 0x3a, 0x2a, 0x3c, 0x6c, 0x87, 0x3f, 0x56,
	0x06, 0xab, 0x7d, 0x77, 0x0d, 0x66, 0xbb, 0x03, 0x7d, 0x0c, 0xb7, 0xa2, 0xc0, 0x59, 0x12, 0xf2,
	0xb4, 0xf9, 0xaf, 0x0b, 0x53, 0x0b, 0x6e, 0x44, 0x81, 0x1f, 0x42, 0x99, 0xa8, 0x19, 0x0c, 0x76,
	0xc6, 0xa6, 0x6a, 0xba, 0xc7, 0x48, 0xaa, 0x17, 0x4d, 0x76, 0x61, 0xf7, 0x35, 0xc0, 0xf5, 0x59,
	0xfa, 0x0c, 0xfa, 0x98, 0x87, 0xa9, 0xfc, 0x16, 0xe5, 0x95, 0xb6, 0x93, 0x16, 0x71, 0x3a, 0x9e,
	0xc1, 0xee, 0xa8, 0xee, 0xba, 0xd9, 0xf4, 0x92, 0xc7, 0xf4, 0x21, 0xe8, 0xdb, 0xd6, 0x40, 0x35,
	0x50, 0x07, 0x7a, 0x49, 0x28, 0x93, 0x2f, 0x1c, 0x93, 0x25, 0x8f, 0x9b, 0x89, 0xda, 0x12, 0x7d,
	0x02, 0x46, 0x85, 0x73, 0x51, 0xf9, 0x1d, 0xe5, 0x5f, 0x85, 0x97, 0x4f, 0xc1, 0xb8, 0xac, 0x99,
	0x0e, 0xc0, 0x0c, 0xd8, 0x9b, 0xf9, 0xf2, 0xfd, 0x3b, 0xb6, 0x5a, 0x7c, 0x0e, 0x06, 0xda, 0x74,
	0xf6, 0xe7, 0x68, 0x93, 0xc3, 0xd1, 0x26, 0xff, 0x8e, 0x36, 0xf9, 0x7d, 0xb2, 0xb5, 0xc3, 0xc9,
	0xd6, 0xfe, 0x9e, 0x6c, 0xed, 0xeb, 0x8b, 0x98, 0x63, 0x52, 0x6c, 0xfc, 0xad, 0xf8, 0x31, 0xae,
	0x2e, 0x65, 0x94, 0x46, 0xf8, 0x53, 0xe4, 0xdf, 0x15, 0x8c, 0xcb, 0xfa, 0xa3, 0xf6, 0xbd, 0xe9,
	0xaa, 0xeb, 0x79, 0xf5, 0x3f, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xf9, 0x64, 0xc4, 0x52, 0x02, 0x00,
	0x00,
}

func (m *TxOutMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxOutMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxOutMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Data != nil {
		{
			size, err := m.Data.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxOut(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TxOut) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxOut) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxOut) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Input != nil {
		{
			size, err := m.Input.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxOut(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Content != nil {
		{
			size, err := m.Content.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxOut(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.TxType != 0 {
		i = encodeVarintTxOut(dAtA, i, uint64(m.TxType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TxOutContent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxOutContent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxOutContent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OutBytes) > 0 {
		i -= len(m.OutBytes)
		copy(dAtA[i:], m.OutBytes)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.OutBytes)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OutHash) > 0 {
		i -= len(m.OutHash)
		copy(dAtA[i:], m.OutHash)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.OutHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.OutChain) > 0 {
		i -= len(m.OutChain)
		copy(dAtA[i:], m.OutChain)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.OutChain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TxOutInput) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxOutInput) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxOutInput) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TransferOutIds) > 0 {
		for iNdEx := len(m.TransferOutIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.TransferOutIds[iNdEx])
			copy(dAtA[i:], m.TransferOutIds[iNdEx])
			i = encodeVarintTxOut(dAtA, i, uint64(len(m.TransferOutIds[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *TxOutSig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxOutSig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxOutSig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.HashNoSig) > 0 {
		i -= len(m.HashNoSig)
		copy(dAtA[i:], m.HashNoSig)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.HashNoSig)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.HashWithSig) > 0 {
		i -= len(m.HashWithSig)
		copy(dAtA[i:], m.HashWithSig)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.HashWithSig)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxOut(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxOut(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxOutMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovTxOut(uint64(l))
	}
	return n
}

func (m *TxOut) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxType != 0 {
		n += 1 + sovTxOut(uint64(m.TxType))
	}
	if m.Content != nil {
		l = m.Content.Size()
		n += 1 + l + sovTxOut(uint64(l))
	}
	if m.Input != nil {
		l = m.Input.Size()
		n += 1 + l + sovTxOut(uint64(l))
	}
	return n
}

func (m *TxOutContent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OutChain)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	l = len(m.OutHash)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	l = len(m.OutBytes)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	return n
}

func (m *TxOutInput) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.TransferOutIds) > 0 {
		for _, s := range m.TransferOutIds {
			l = len(s)
			n += 1 + l + sovTxOut(uint64(l))
		}
	}
	return n
}

func (m *TxOutSig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	l = len(m.HashWithSig)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	l = len(m.HashNoSig)
	if l > 0 {
		n += 1 + l + sovTxOut(uint64(l))
	}
	return n
}

func sovTxOut(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxOut(x uint64) (n int) {
	return sovTxOut(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxOutMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxOut
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
			return fmt.Errorf("proto: TxOutMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxOutMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &TxOut{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxOut(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTxOut
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
func (m *TxOut) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxOut
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
			return fmt.Errorf("proto: TxOut: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxOut: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxType", wireType)
			}
			m.TxType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TxType |= TxOutType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Content == nil {
				m.Content = &TxOutContent{}
			}
			if err := m.Content.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Input", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Input == nil {
				m.Input = &TxOutInput{}
			}
			if err := m.Input.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxOut(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTxOut
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
func (m *TxOutContent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxOut
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
			return fmt.Errorf("proto: TxOutContent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxOutContent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutBytes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutBytes = append(m.OutBytes[:0], dAtA[iNdEx:postIndex]...)
			if m.OutBytes == nil {
				m.OutBytes = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxOut(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTxOut
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
func (m *TxOutInput) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxOut
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
			return fmt.Errorf("proto: TxOutInput: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxOutInput: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferOutIds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TransferOutIds = append(m.TransferOutIds, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxOut(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTxOut
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
func (m *TxOutSig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxOut
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
			return fmt.Errorf("proto: TxOutSig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxOutSig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HashWithSig", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HashWithSig = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HashNoSig", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxOut
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
				return ErrInvalidLengthTxOut
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxOut
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HashNoSig = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxOut(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTxOut
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
func skipTxOut(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxOut
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
					return 0, ErrIntOverflowTxOut
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
					return 0, ErrIntOverflowTxOut
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
				return 0, ErrInvalidLengthTxOut
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxOut
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxOut
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxOut        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxOut          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxOut = fmt.Errorf("proto: unexpected end of group")
)

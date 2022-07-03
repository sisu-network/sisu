// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/transfer_requests.proto

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

type TransferRequestsMsg struct {
	Signer string            `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *TransferRequests `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *TransferRequestsMsg) Reset()         { *m = TransferRequestsMsg{} }
func (m *TransferRequestsMsg) String() string { return proto.CompactTextString(m) }
func (*TransferRequestsMsg) ProtoMessage()    {}
func (*TransferRequestsMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_046e004cb5b2eb10, []int{0}
}
func (m *TransferRequestsMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferRequestsMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferRequestsMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferRequestsMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferRequestsMsg.Merge(m, src)
}
func (m *TransferRequestsMsg) XXX_Size() int {
	return m.Size()
}
func (m *TransferRequestsMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferRequestsMsg.DiscardUnknown(m)
}

var xxx_messageInfo_TransferRequestsMsg proto.InternalMessageInfo

func (m *TransferRequestsMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *TransferRequestsMsg) GetData() *TransferRequests {
	if m != nil {
		return m.Data
	}
	return nil
}

// BlockTransfers contains all transfer request in a block.
type TransferRequests struct {
	Chain    string             `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Height   int64              `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Hash     string             `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Requests []*TransferRequest `protobuf:"bytes,4,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (m *TransferRequests) Reset()         { *m = TransferRequests{} }
func (m *TransferRequests) String() string { return proto.CompactTextString(m) }
func (*TransferRequests) ProtoMessage()    {}
func (*TransferRequests) Descriptor() ([]byte, []int) {
	return fileDescriptor_046e004cb5b2eb10, []int{1}
}
func (m *TransferRequests) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferRequests) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferRequests.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferRequests) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferRequests.Merge(m, src)
}
func (m *TransferRequests) XXX_Size() int {
	return m.Size()
}
func (m *TransferRequests) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferRequests.DiscardUnknown(m)
}

var xxx_messageInfo_TransferRequests proto.InternalMessageInfo

func (m *TransferRequests) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *TransferRequests) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *TransferRequests) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *TransferRequests) GetRequests() []*TransferRequest {
	if m != nil {
		return m.Requests
	}
	return nil
}

type TransferRequest struct {
	Sender    string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	ToChain   string `protobuf:"bytes,2,opt,name=to_chain,json=toChain,proto3" json:"to_chain,omitempty"`
	Token     string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Recipient string `protobuf:"bytes,4,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount    string `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *TransferRequest) Reset()         { *m = TransferRequest{} }
func (m *TransferRequest) String() string { return proto.CompactTextString(m) }
func (*TransferRequest) ProtoMessage()    {}
func (*TransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_046e004cb5b2eb10, []int{2}
}
func (m *TransferRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferRequest.Merge(m, src)
}
func (m *TransferRequest) XXX_Size() int {
	return m.Size()
}
func (m *TransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TransferRequest proto.InternalMessageInfo

func (m *TransferRequest) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *TransferRequest) GetToChain() string {
	if m != nil {
		return m.ToChain
	}
	return ""
}

func (m *TransferRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *TransferRequest) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

func (m *TransferRequest) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func init() {
	proto.RegisterType((*TransferRequestsMsg)(nil), "types.TransferRequestsMsg")
	proto.RegisterType((*TransferRequests)(nil), "types.TransferRequests")
	proto.RegisterType((*TransferRequest)(nil), "types.TransferRequest")
}

func init() { proto.RegisterFile("sisu/transfer_requests.proto", fileDescriptor_046e004cb5b2eb10) }

var fileDescriptor_046e004cb5b2eb10 = []byte{
	// 318 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x4e, 0xc2, 0x40,
	0x10, 0xc6, 0x59, 0x68, 0x11, 0x86, 0x83, 0x66, 0x25, 0x58, 0x13, 0xd2, 0x10, 0x4e, 0x18, 0x63,
	0x49, 0xf0, 0x0d, 0xe4, 0xec, 0xa5, 0xf1, 0xc4, 0x85, 0x14, 0x58, 0xdb, 0x0d, 0x61, 0x17, 0x77,
	0xa7, 0x51, 0xdf, 0xc0, 0x23, 0x8f, 0xe5, 0x91, 0xa3, 0x47, 0x03, 0x2f, 0x62, 0x98, 0x5d, 0xff,
	0x84, 0x78, 0x6a, 0x7f, 0x3b, 0x5f, 0xbe, 0xf9, 0xbe, 0x0c, 0x74, 0xad, 0xb4, 0xe5, 0x10, 0x4d,
	0xa6, 0xec, 0xa3, 0x30, 0x53, 0x23, 0x9e, 0x4a, 0x61, 0xd1, 0x26, 0x6b, 0xa3, 0x51, 0xf3, 0x10,
	0x5f, 0xd7, 0xc2, 0xf6, 0x27, 0x70, 0xfe, 0xe0, 0x15, 0xa9, 0x17, 0xdc, 0xdb, 0x9c, 0x77, 0xa0,
	0x6e, 0x65, 0xae, 0x84, 0x89, 0x58, 0x8f, 0x0d, 0x9a, 0xa9, 0x27, 0x7e, 0x0d, 0xc1, 0x22, 0xc3,
	0x2c, 0xaa, 0xf6, 0xd8, 0xa0, 0x35, 0xba, 0x48, 0xc8, 0x24, 0x39, 0x76, 0x48, 0x49, 0xd4, 0x7f,
	0x63, 0x70, 0x76, 0x3c, 0xe2, 0x6d, 0x08, 0xe7, 0x45, 0x26, 0x95, 0x37, 0x76, 0x70, 0xd8, 0x57,
	0x08, 0x99, 0x17, 0x48, 0xce, 0xb5, 0xd4, 0x13, 0xe7, 0x10, 0x14, 0x99, 0x2d, 0xa2, 0x1a, 0x89,
	0xe9, 0x9f, 0x8f, 0xa0, 0xf1, 0xdd, 0x25, 0x0a, 0x7a, 0xb5, 0x41, 0x6b, 0xd4, 0xf9, 0x3f, 0x47,
	0xfa, 0xa3, 0xeb, 0x6f, 0x18, 0x9c, 0x1e, 0x4d, 0xa9, 0xa3, 0x50, 0x8b, 0x3f, 0x1d, 0x89, 0xf8,
	0x25, 0x34, 0x50, 0x4f, 0x5d, 0xc8, 0x2a, 0x4d, 0x4e, 0x50, 0x8f, 0x29, 0x66, 0x1b, 0x42, 0xd4,
	0x4b, 0xa1, 0x7c, 0x1e, 0x07, 0xbc, 0x0b, 0x4d, 0x23, 0xe6, 0x72, 0x2d, 0x85, 0xc2, 0x28, 0xa0,
	0xc9, 0xef, 0xc3, 0x61, 0x4d, 0xb6, 0xd2, 0xa5, 0xc2, 0x28, 0x74, 0x6b, 0x1c, 0xdd, 0x8d, 0xdf,
	0x77, 0x31, 0xdb, 0xee, 0x62, 0xf6, 0xb9, 0x8b, 0xd9, 0x66, 0x1f, 0x57, 0xb6, 0xfb, 0xb8, 0xf2,
	0xb1, 0x8f, 0x2b, 0x93, 0xab, 0x5c, 0x62, 0x51, 0xce, 0x92, 0xb9, 0x5e, 0x0d, 0x0f, 0x37, 0xbc,
	0x51, 0x02, 0x9f, 0xb5, 0x59, 0x12, 0x0c, 0x5f, 0xdc, 0x87, 0x1a, 0xcf, 0xea, 0x74, 0xcc, 0xdb,
	0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xc2, 0x3a, 0xe8, 0xec, 0x01, 0x00, 0x00,
}

func (m *TransferRequestsMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferRequestsMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferRequestsMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintTransferRequests(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TransferRequests) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferRequests) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferRequests) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTransferRequests(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Height != 0 {
		i = encodeVarintTransferRequests(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TransferRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ToChain) > 0 {
		i -= len(m.ToChain)
		copy(dAtA[i:], m.ToChain)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.ToChain)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTransferRequests(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTransferRequests(dAtA []byte, offset int, v uint64) int {
	offset -= sovTransferRequests(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TransferRequestsMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	return n
}

func (m *TransferRequests) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovTransferRequests(uint64(m.Height))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovTransferRequests(uint64(l))
		}
	}
	return n
}

func (m *TransferRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	l = len(m.ToChain)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovTransferRequests(uint64(l))
	}
	return n
}

func sovTransferRequests(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTransferRequests(x uint64) (n int) {
	return sovTransferRequests(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TransferRequestsMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferRequests
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
			return fmt.Errorf("proto: TransferRequestsMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferRequestsMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
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
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &TransferRequests{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransferRequests(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferRequests
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
func (m *TransferRequests) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferRequests
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
			return fmt.Errorf("proto: TransferRequests: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferRequests: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requests = append(m.Requests, &TransferRequest{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransferRequests(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferRequests
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
func (m *TransferRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferRequests
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
			return fmt.Errorf("proto: TransferRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ToChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferRequests
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
				return ErrInvalidLengthTransferRequests
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferRequests
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransferRequests(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferRequests
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
func skipTransferRequests(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransferRequests
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
					return 0, ErrIntOverflowTransferRequests
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
					return 0, ErrIntOverflowTransferRequests
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
				return 0, ErrInvalidLengthTransferRequests
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTransferRequests
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTransferRequests
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTransferRequests        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransferRequests          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTransferRequests = fmt.Errorf("proto: unexpected end of group")
)

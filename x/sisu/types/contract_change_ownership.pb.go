// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/contract_change_ownership.proto

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

type ChangeOwnershipContractMsg struct {
	Signer string           `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *ChangeOwnership `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ChangeOwnershipContractMsg) Reset()         { *m = ChangeOwnershipContractMsg{} }
func (m *ChangeOwnershipContractMsg) String() string { return proto.CompactTextString(m) }
func (*ChangeOwnershipContractMsg) ProtoMessage()    {}
func (*ChangeOwnershipContractMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_6116588f537b8fc7, []int{0}
}
func (m *ChangeOwnershipContractMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeOwnershipContractMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeOwnershipContractMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeOwnershipContractMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeOwnershipContractMsg.Merge(m, src)
}
func (m *ChangeOwnershipContractMsg) XXX_Size() int {
	return m.Size()
}
func (m *ChangeOwnershipContractMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeOwnershipContractMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeOwnershipContractMsg proto.InternalMessageInfo

func (m *ChangeOwnershipContractMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *ChangeOwnershipContractMsg) GetData() *ChangeOwnership {
	if m != nil {
		return m.Data
	}
	return nil
}

type ChangeOwnership struct {
	Chain    string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Hash     string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	NewOwner string `protobuf:"bytes,3,opt,name=newOwner,proto3" json:"newOwner,omitempty"`
	Index    int32  `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *ChangeOwnership) Reset()         { *m = ChangeOwnership{} }
func (m *ChangeOwnership) String() string { return proto.CompactTextString(m) }
func (*ChangeOwnership) ProtoMessage()    {}
func (*ChangeOwnership) Descriptor() ([]byte, []int) {
	return fileDescriptor_6116588f537b8fc7, []int{1}
}
func (m *ChangeOwnership) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeOwnership) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeOwnership.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeOwnership) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeOwnership.Merge(m, src)
}
func (m *ChangeOwnership) XXX_Size() int {
	return m.Size()
}
func (m *ChangeOwnership) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeOwnership.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeOwnership proto.InternalMessageInfo

func (m *ChangeOwnership) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *ChangeOwnership) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *ChangeOwnership) GetNewOwner() string {
	if m != nil {
		return m.NewOwner
	}
	return ""
}

func (m *ChangeOwnership) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*ChangeOwnershipContractMsg)(nil), "types.ChangeOwnershipContractMsg")
	proto.RegisterType((*ChangeOwnership)(nil), "types.ChangeOwnership")
}

func init() {
	proto.RegisterFile("sisu/contract_change_ownership.proto", fileDescriptor_6116588f537b8fc7)
}

var fileDescriptor_6116588f537b8fc7 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xc1, 0x4a, 0xc3, 0x30,
	0x18, 0xc7, 0x1b, 0x6d, 0x87, 0x8b, 0x07, 0x21, 0xc8, 0x28, 0x3b, 0x84, 0x32, 0x3c, 0x54, 0xc1,
	0x16, 0xf4, 0x0d, 0xec, 0x59, 0x84, 0x1e, 0xbd, 0xcc, 0x2c, 0x0b, 0x49, 0x90, 0x25, 0x25, 0xc9,
	0xe8, 0x7c, 0x0b, 0x1f, 0xcb, 0xe3, 0x8e, 0x1e, 0xa5, 0x7d, 0x11, 0xd9, 0xd7, 0xea, 0xa1, 0xa7,
	0xe4, 0x97, 0xef, 0xfb, 0xff, 0x20, 0x7f, 0x7c, 0xe3, 0xb5, 0xdf, 0x97, 0xdc, 0x9a, 0xe0, 0x18,
	0x0f, 0x6b, 0xae, 0x98, 0x91, 0x62, 0x6d, 0x5b, 0x23, 0x9c, 0x57, 0xba, 0x29, 0x1a, 0x67, 0x83,
	0x25, 0x49, 0xf8, 0x68, 0x84, 0x5f, 0xbd, 0xe1, 0x65, 0x05, 0x0b, 0x2f, 0x7f, 0xf3, 0x6a, 0x0c,
	0x3e, 0x7b, 0x49, 0x16, 0x78, 0xe6, 0xb5, 0x34, 0xc2, 0xa5, 0x28, 0x43, 0xf9, 0xbc, 0x1e, 0x89,
	0xdc, 0xe1, 0x78, 0xcb, 0x02, 0x4b, 0xcf, 0x32, 0x94, 0x5f, 0x3e, 0x2c, 0x0a, 0x70, 0x15, 0x13,
	0x51, 0x0d, 0x3b, 0xab, 0x1d, 0xbe, 0x9a, 0x0c, 0xc8, 0x35, 0x4e, 0xb8, 0x62, 0xda, 0x8c, 0xd6,
	0x01, 0x08, 0xc1, 0xb1, 0x62, 0x5e, 0x81, 0x74, 0x5e, 0xc3, 0x9d, 0x2c, 0xf1, 0x85, 0x11, 0x2d,
	0x24, 0xd3, 0x73, 0x78, 0xff, 0xe7, 0x93, 0x45, 0x9b, 0xad, 0x38, 0xa4, 0x71, 0x86, 0xf2, 0xa4,
	0x1e, 0xe0, 0xa9, 0xfa, 0xea, 0x28, 0x3a, 0x76, 0x14, 0xfd, 0x74, 0x14, 0x7d, 0xf6, 0x34, 0x3a,
	0xf6, 0x34, 0xfa, 0xee, 0x69, 0xf4, 0x7a, 0x2b, 0x75, 0x50, 0xfb, 0x4d, 0xc1, 0xed, 0xae, 0x3c,
	0x55, 0x74, 0x6f, 0x44, 0x68, 0xad, 0x7b, 0x07, 0x28, 0x0f, 0xc3, 0x01, 0x3f, 0xd9, 0xcc, 0xa0,
	0xa3, 0xc7, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x67, 0x06, 0xc7, 0x38, 0x4b, 0x01, 0x00, 0x00,
}

func (m *ChangeOwnershipContractMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeOwnershipContractMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeOwnershipContractMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintContractChangeOwnership(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintContractChangeOwnership(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ChangeOwnership) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeOwnership) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeOwnership) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintContractChangeOwnership(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x20
	}
	if len(m.NewOwner) > 0 {
		i -= len(m.NewOwner)
		copy(dAtA[i:], m.NewOwner)
		i = encodeVarintContractChangeOwnership(dAtA, i, uint64(len(m.NewOwner)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintContractChangeOwnership(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintContractChangeOwnership(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintContractChangeOwnership(dAtA []byte, offset int, v uint64) int {
	offset -= sovContractChangeOwnership(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChangeOwnershipContractMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovContractChangeOwnership(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovContractChangeOwnership(uint64(l))
	}
	return n
}

func (m *ChangeOwnership) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovContractChangeOwnership(uint64(l))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovContractChangeOwnership(uint64(l))
	}
	l = len(m.NewOwner)
	if l > 0 {
		n += 1 + l + sovContractChangeOwnership(uint64(l))
	}
	if m.Index != 0 {
		n += 1 + sovContractChangeOwnership(uint64(m.Index))
	}
	return n
}

func sovContractChangeOwnership(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozContractChangeOwnership(x uint64) (n int) {
	return sovContractChangeOwnership(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChangeOwnershipContractMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContractChangeOwnership
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
			return fmt.Errorf("proto: ChangeOwnershipContractMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeOwnershipContractMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeOwnership
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
				return ErrInvalidLengthContractChangeOwnership
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeOwnership
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
					return ErrIntOverflowContractChangeOwnership
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
				return ErrInvalidLengthContractChangeOwnership
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeOwnership
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &ChangeOwnership{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContractChangeOwnership(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContractChangeOwnership
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
func (m *ChangeOwnership) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContractChangeOwnership
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
			return fmt.Errorf("proto: ChangeOwnership: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeOwnership: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeOwnership
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
				return ErrInvalidLengthContractChangeOwnership
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeOwnership
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
					return ErrIntOverflowContractChangeOwnership
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
				return ErrInvalidLengthContractChangeOwnership
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeOwnership
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewOwner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeOwnership
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
				return ErrInvalidLengthContractChangeOwnership
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeOwnership
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewOwner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeOwnership
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipContractChangeOwnership(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContractChangeOwnership
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
func skipContractChangeOwnership(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContractChangeOwnership
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
					return 0, ErrIntOverflowContractChangeOwnership
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
					return 0, ErrIntOverflowContractChangeOwnership
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
				return 0, ErrInvalidLengthContractChangeOwnership
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupContractChangeOwnership
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthContractChangeOwnership
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthContractChangeOwnership        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContractChangeOwnership          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupContractChangeOwnership = fmt.Errorf("proto: unexpected end of group")
)

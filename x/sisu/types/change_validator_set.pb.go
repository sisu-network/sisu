// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/change_validator_set.proto

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

type ChangeValidatorSetMsg struct {
	Signer string                  `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *ChangeValidatorSetData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ChangeValidatorSetMsg) Reset()         { *m = ChangeValidatorSetMsg{} }
func (m *ChangeValidatorSetMsg) String() string { return proto.CompactTextString(m) }
func (*ChangeValidatorSetMsg) ProtoMessage()    {}
func (*ChangeValidatorSetMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6e7ac43babceb00, []int{0}
}
func (m *ChangeValidatorSetMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeValidatorSetMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeValidatorSetMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeValidatorSetMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeValidatorSetMsg.Merge(m, src)
}
func (m *ChangeValidatorSetMsg) XXX_Size() int {
	return m.Size()
}
func (m *ChangeValidatorSetMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeValidatorSetMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeValidatorSetMsg proto.InternalMessageInfo

func (m *ChangeValidatorSetMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *ChangeValidatorSetMsg) GetData() *ChangeValidatorSetData {
	if m != nil {
		return m.Data
	}
	return nil
}

type ChangeValidatorSetData struct {
	OldValidatorSet [][]byte `protobuf:"bytes,1,rep,name=oldValidatorSet,proto3" json:"oldValidatorSet,omitempty"`
	NewValidatorSet [][]byte `protobuf:"bytes,2,rep,name=newValidatorSet,proto3" json:"newValidatorSet,omitempty"`
	Index           int32    `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *ChangeValidatorSetData) Reset()         { *m = ChangeValidatorSetData{} }
func (m *ChangeValidatorSetData) String() string { return proto.CompactTextString(m) }
func (*ChangeValidatorSetData) ProtoMessage()    {}
func (*ChangeValidatorSetData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6e7ac43babceb00, []int{1}
}
func (m *ChangeValidatorSetData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeValidatorSetData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeValidatorSetData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeValidatorSetData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeValidatorSetData.Merge(m, src)
}
func (m *ChangeValidatorSetData) XXX_Size() int {
	return m.Size()
}
func (m *ChangeValidatorSetData) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeValidatorSetData.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeValidatorSetData proto.InternalMessageInfo

func (m *ChangeValidatorSetData) GetOldValidatorSet() [][]byte {
	if m != nil {
		return m.OldValidatorSet
	}
	return nil
}

func (m *ChangeValidatorSetData) GetNewValidatorSet() [][]byte {
	if m != nil {
		return m.NewValidatorSet
	}
	return nil
}

func (m *ChangeValidatorSetData) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*ChangeValidatorSetMsg)(nil), "types.ChangeValidatorSetMsg")
	proto.RegisterType((*ChangeValidatorSetData)(nil), "types.ChangeValidatorSetData")
}

func init() { proto.RegisterFile("sisu/change_validator_set.proto", fileDescriptor_c6e7ac43babceb00) }

var fileDescriptor_c6e7ac43babceb00 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0xce, 0x2c, 0x2e,
	0xd5, 0x4f, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0x8d, 0x2f, 0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0xc9,
	0x2f, 0x8a, 0x2f, 0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c,
	0x48, 0x2d, 0x56, 0x4a, 0xe2, 0x12, 0x75, 0x06, 0x2b, 0x0a, 0x83, 0xa9, 0x09, 0x4e, 0x2d, 0xf1,
	0x2d, 0x4e, 0x17, 0x12, 0xe3, 0x62, 0x2b, 0xce, 0x4c, 0xcf, 0x4b, 0x2d, 0x92, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x84, 0x0c, 0xb9, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x98, 0x14,
	0x18, 0x35, 0xb8, 0x8d, 0x64, 0xf5, 0xc0, 0xc6, 0xe8, 0x61, 0x9a, 0xe1, 0x92, 0x58, 0x92, 0x18,
	0x04, 0x56, 0xaa, 0xd4, 0xc4, 0xc8, 0x25, 0x86, 0x5d, 0x81, 0x90, 0x06, 0x17, 0x7f, 0x7e, 0x4e,
	0x0a, 0xb2, 0xb0, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x4f, 0x10, 0xba, 0x30, 0x48, 0x65, 0x5e, 0x6a,
	0x39, 0x8a, 0x4a, 0x26, 0x88, 0x4a, 0x34, 0x61, 0x21, 0x11, 0x2e, 0xd6, 0xcc, 0xbc, 0x94, 0xd4,
	0x0a, 0x09, 0x66, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x08, 0xc7, 0xc9, 0xf9, 0xc4, 0x23, 0x39, 0xc6,
	0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39,
	0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x34, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3,
	0x73, 0xf5, 0x41, 0xa1, 0xa6, 0x9b, 0x97, 0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x0d, 0xe6, 0xe8, 0x57,
	0x40, 0x28, 0xb0, 0x37, 0x93, 0xd8, 0xc0, 0x61, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xcf,
	0xae, 0x34, 0x4b, 0x5e, 0x01, 0x00, 0x00,
}

func (m *ChangeValidatorSetMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeValidatorSetMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeValidatorSetMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintChangeValidatorSet(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintChangeValidatorSet(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ChangeValidatorSetData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeValidatorSetData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeValidatorSetData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintChangeValidatorSet(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	if len(m.NewValidatorSet) > 0 {
		for iNdEx := len(m.NewValidatorSet) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NewValidatorSet[iNdEx])
			copy(dAtA[i:], m.NewValidatorSet[iNdEx])
			i = encodeVarintChangeValidatorSet(dAtA, i, uint64(len(m.NewValidatorSet[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.OldValidatorSet) > 0 {
		for iNdEx := len(m.OldValidatorSet) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.OldValidatorSet[iNdEx])
			copy(dAtA[i:], m.OldValidatorSet[iNdEx])
			i = encodeVarintChangeValidatorSet(dAtA, i, uint64(len(m.OldValidatorSet[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintChangeValidatorSet(dAtA []byte, offset int, v uint64) int {
	offset -= sovChangeValidatorSet(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChangeValidatorSetMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovChangeValidatorSet(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovChangeValidatorSet(uint64(l))
	}
	return n
}

func (m *ChangeValidatorSetData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.OldValidatorSet) > 0 {
		for _, b := range m.OldValidatorSet {
			l = len(b)
			n += 1 + l + sovChangeValidatorSet(uint64(l))
		}
	}
	if len(m.NewValidatorSet) > 0 {
		for _, b := range m.NewValidatorSet {
			l = len(b)
			n += 1 + l + sovChangeValidatorSet(uint64(l))
		}
	}
	if m.Index != 0 {
		n += 1 + sovChangeValidatorSet(uint64(m.Index))
	}
	return n
}

func sovChangeValidatorSet(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChangeValidatorSet(x uint64) (n int) {
	return sovChangeValidatorSet(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChangeValidatorSetMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChangeValidatorSet
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
			return fmt.Errorf("proto: ChangeValidatorSetMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeValidatorSetMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChangeValidatorSet
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
				return ErrInvalidLengthChangeValidatorSet
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChangeValidatorSet
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
					return ErrIntOverflowChangeValidatorSet
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
				return ErrInvalidLengthChangeValidatorSet
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthChangeValidatorSet
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &ChangeValidatorSetData{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipChangeValidatorSet(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChangeValidatorSet
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
func (m *ChangeValidatorSetData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChangeValidatorSet
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
			return fmt.Errorf("proto: ChangeValidatorSetData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeValidatorSetData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldValidatorSet", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChangeValidatorSet
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
				return ErrInvalidLengthChangeValidatorSet
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthChangeValidatorSet
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OldValidatorSet = append(m.OldValidatorSet, make([]byte, postIndex-iNdEx))
			copy(m.OldValidatorSet[len(m.OldValidatorSet)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewValidatorSet", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChangeValidatorSet
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
				return ErrInvalidLengthChangeValidatorSet
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthChangeValidatorSet
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewValidatorSet = append(m.NewValidatorSet, make([]byte, postIndex-iNdEx))
			copy(m.NewValidatorSet[len(m.NewValidatorSet)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChangeValidatorSet
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
			skippy, err := skipChangeValidatorSet(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChangeValidatorSet
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
func skipChangeValidatorSet(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChangeValidatorSet
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
					return 0, ErrIntOverflowChangeValidatorSet
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
					return 0, ErrIntOverflowChangeValidatorSet
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
				return 0, ErrInvalidLengthChangeValidatorSet
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChangeValidatorSet
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChangeValidatorSet
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChangeValidatorSet        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChangeValidatorSet          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChangeValidatorSet = fmt.Errorf("proto: unexpected end of group")
)
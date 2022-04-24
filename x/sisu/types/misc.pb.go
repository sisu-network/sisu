// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/misc.proto

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

type ProtoInt64 struct {
	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *ProtoInt64) Reset()         { *m = ProtoInt64{} }
func (m *ProtoInt64) String() string { return proto.CompactTextString(m) }
func (*ProtoInt64) ProtoMessage()    {}
func (*ProtoInt64) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d9fdbfa949063ec, []int{0}
}
func (m *ProtoInt64) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtoInt64) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtoInt64.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtoInt64) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtoInt64.Merge(m, src)
}
func (m *ProtoInt64) XXX_Size() int {
	return m.Size()
}
func (m *ProtoInt64) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtoInt64.DiscardUnknown(m)
}

var xxx_messageInfo_ProtoInt64 proto.InternalMessageInfo

func (m *ProtoInt64) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*ProtoInt64)(nil), "types.ProtoInt64")
}

func init() { proto.RegisterFile("sisu/misc.proto", fileDescriptor_8d9fdbfa949063ec) }

var fileDescriptor_8d9fdbfa949063ec = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xce, 0x2c, 0x2e,
	0xd5, 0xcf, 0xcd, 0x2c, 0x4e, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c,
	0x48, 0x2d, 0x56, 0x52, 0xe2, 0xe2, 0x0a, 0x00, 0xf1, 0x3d, 0xf3, 0x4a, 0xcc, 0x4c, 0x84, 0x44,
	0xb8, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x20, 0x1c,
	0x27, 0xe7, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2,
	0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0xd2, 0x4c, 0xcf, 0x2c,
	0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0x59, 0xa0, 0x9b, 0x97, 0x5a, 0x52, 0x9e,
	0x5f, 0x94, 0x0d, 0xe6, 0xe8, 0x57, 0x40, 0x28, 0xb0, 0x45, 0x49, 0x6c, 0x60, 0x6b, 0x8d, 0x01,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x73, 0x5c, 0x0f, 0x41, 0x89, 0x00, 0x00, 0x00,
}

func (m *ProtoInt64) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtoInt64) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtoInt64) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintMisc(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMisc(dAtA []byte, offset int, v uint64) int {
	offset -= sovMisc(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProtoInt64) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != 0 {
		n += 1 + sovMisc(uint64(m.Value))
	}
	return n
}

func sovMisc(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMisc(x uint64) (n int) {
	return sovMisc(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProtoInt64) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMisc
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
			return fmt.Errorf("proto: ProtoInt64: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtoInt64: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMisc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMisc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMisc
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
func skipMisc(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMisc
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
					return 0, ErrIntOverflowMisc
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
					return 0, ErrIntOverflowMisc
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
				return 0, ErrInvalidLengthMisc
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMisc
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMisc
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMisc        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMisc          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMisc = fmt.Errorf("proto: unexpected end of group")
)

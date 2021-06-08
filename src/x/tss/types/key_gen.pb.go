// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: key_gen.proto

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

type KeygenProposal struct {
	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	ChainSymbol string `protobuf:"bytes,2,opt,name=chainSymbol,proto3" json:"chainSymbol,omitempty"`
}

func (m *KeygenProposal) Reset()         { *m = KeygenProposal{} }
func (m *KeygenProposal) String() string { return proto.CompactTextString(m) }
func (*KeygenProposal) ProtoMessage()    {}
func (*KeygenProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d264d872070eabf, []int{0}
}
func (m *KeygenProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *KeygenProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_KeygenProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *KeygenProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeygenProposal.Merge(m, src)
}
func (m *KeygenProposal) XXX_Size() int {
	return m.Size()
}
func (m *KeygenProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_KeygenProposal.DiscardUnknown(m)
}

var xxx_messageInfo_KeygenProposal proto.InternalMessageInfo

func (m *KeygenProposal) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *KeygenProposal) GetChainSymbol() string {
	if m != nil {
		return m.ChainSymbol
	}
	return ""
}

func init() {
	proto.RegisterType((*KeygenProposal)(nil), "types.KeygenProposal")
}

func init() { proto.RegisterFile("key_gen.proto", fileDescriptor_3d264d872070eabf) }

var fileDescriptor_3d264d872070eabf = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4e, 0xad, 0x8c,
	0x4f, 0x4f, 0xcd, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d,
	0x56, 0xf2, 0xe2, 0xe2, 0xf3, 0x4e, 0xad, 0x4c, 0x4f, 0xcd, 0x0b, 0x28, 0xca, 0x2f, 0xc8, 0x2f,
	0x4e, 0xcc, 0x11, 0x12, 0xe3, 0x62, 0x2b, 0xce, 0x4c, 0xcf, 0x4b, 0x2d, 0x92, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x84, 0x14, 0xb8, 0xb8, 0x93, 0x33, 0x12, 0x33, 0xf3, 0x82, 0x2b,
	0x73, 0x93, 0xf2, 0x73, 0x24, 0x98, 0xc0, 0x92, 0xc8, 0x42, 0x4e, 0x12, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb,
	0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x90, 0xc4, 0x06, 0xb6, 0xd3, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff,
	0x99, 0x80, 0x8f, 0x76, 0x84, 0x00, 0x00, 0x00,
}

func (m *KeygenProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *KeygenProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *KeygenProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainSymbol) > 0 {
		i -= len(m.ChainSymbol)
		copy(dAtA[i:], m.ChainSymbol)
		i = encodeVarintKeyGen(dAtA, i, uint64(len(m.ChainSymbol)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintKeyGen(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintKeyGen(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeyGen(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *KeygenProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovKeyGen(uint64(l))
	}
	l = len(m.ChainSymbol)
	if l > 0 {
		n += 1 + l + sovKeyGen(uint64(l))
	}
	return n
}

func sovKeyGen(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeyGen(x uint64) (n int) {
	return sovKeyGen(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *KeygenProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeyGen
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
			return fmt.Errorf("proto: KeygenProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeygenProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeyGen
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
				return ErrInvalidLengthKeyGen
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthKeyGen
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainSymbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeyGen
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
				return ErrInvalidLengthKeyGen
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthKeyGen
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainSymbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeyGen(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthKeyGen
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
func skipKeyGen(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeyGen
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
					return 0, ErrIntOverflowKeyGen
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
					return 0, ErrIntOverflowKeyGen
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
				return 0, ErrInvalidLengthKeyGen
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeyGen
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeyGen
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthKeyGen        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowKeyGen          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupKeyGen = fmt.Errorf("proto: unexpected end of group")
)

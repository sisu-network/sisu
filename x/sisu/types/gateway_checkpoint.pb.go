// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/gateway_checkpoint.proto

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

// This meessage represents a data point where majority of nodes in Sisu network agrees to observe.
type GatewayCheckPoint struct {
	Chain       string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	BlockHeight int64  `protobuf:"varint,2,opt,name=blockHeight,proto3" json:"blockHeight,omitempty"`
	Nonce       int64  `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *GatewayCheckPoint) Reset()         { *m = GatewayCheckPoint{} }
func (m *GatewayCheckPoint) String() string { return proto.CompactTextString(m) }
func (*GatewayCheckPoint) ProtoMessage()    {}
func (*GatewayCheckPoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_07749ff65f38aa3b, []int{0}
}
func (m *GatewayCheckPoint) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GatewayCheckPoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GatewayCheckPoint.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GatewayCheckPoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GatewayCheckPoint.Merge(m, src)
}
func (m *GatewayCheckPoint) XXX_Size() int {
	return m.Size()
}
func (m *GatewayCheckPoint) XXX_DiscardUnknown() {
	xxx_messageInfo_GatewayCheckPoint.DiscardUnknown(m)
}

var xxx_messageInfo_GatewayCheckPoint proto.InternalMessageInfo

func (m *GatewayCheckPoint) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *GatewayCheckPoint) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *GatewayCheckPoint) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func init() {
	proto.RegisterType((*GatewayCheckPoint)(nil), "types.GatewayCheckPoint")
}

func init() { proto.RegisterFile("sisu/gateway_checkpoint.proto", fileDescriptor_07749ff65f38aa3b) }

var fileDescriptor_07749ff65f38aa3b = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0xce, 0x2c, 0x2e,
	0xd5, 0x4f, 0x4f, 0x2c, 0x49, 0x2d, 0x4f, 0xac, 0x8c, 0x4f, 0xce, 0x48, 0x4d, 0xce, 0x2e, 0xc8,
	0xcf, 0xcc, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d,
	0x56, 0x4a, 0xe4, 0x12, 0x74, 0x87, 0x28, 0x71, 0x06, 0xa9, 0x08, 0x00, 0xa9, 0x10, 0x12, 0xe1,
	0x62, 0x4d, 0xce, 0x48, 0xcc, 0xcc, 0x93, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x84,
	0x14, 0xb8, 0xb8, 0x93, 0x72, 0xf2, 0x93, 0xb3, 0x3d, 0x52, 0x33, 0xd3, 0x33, 0x4a, 0x24, 0x98,
	0x14, 0x18, 0x35, 0x98, 0x83, 0x90, 0x85, 0x40, 0xfa, 0xf2, 0xf2, 0xf3, 0x92, 0x53, 0x25, 0x98,
	0xc1, 0x72, 0x10, 0x8e, 0x93, 0xf3, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78,
	0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44,
	0x69, 0xa6, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x83, 0x5c, 0xab, 0x9b,
	0x97, 0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x0d, 0xe6, 0xe8, 0x57, 0x40, 0x28, 0xb0, 0x3b, 0x93, 0xd8,
	0xc0, 0xae, 0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa5, 0xdb, 0xa5, 0x85, 0xd6, 0x00, 0x00,
	0x00,
}

func (m *GatewayCheckPoint) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayCheckPoint) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GatewayCheckPoint) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Nonce != 0 {
		i = encodeVarintGatewayCheckpoint(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x18
	}
	if m.BlockHeight != 0 {
		i = encodeVarintGatewayCheckpoint(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintGatewayCheckpoint(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGatewayCheckpoint(dAtA []byte, offset int, v uint64) int {
	offset -= sovGatewayCheckpoint(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GatewayCheckPoint) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovGatewayCheckpoint(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovGatewayCheckpoint(uint64(m.BlockHeight))
	}
	if m.Nonce != 0 {
		n += 1 + sovGatewayCheckpoint(uint64(m.Nonce))
	}
	return n
}

func sovGatewayCheckpoint(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGatewayCheckpoint(x uint64) (n int) {
	return sovGatewayCheckpoint(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GatewayCheckPoint) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGatewayCheckpoint
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
			return fmt.Errorf("proto: GatewayCheckPoint: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayCheckPoint: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayCheckpoint
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
				return ErrInvalidLengthGatewayCheckpoint
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGatewayCheckpoint
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayCheckpoint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayCheckpoint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGatewayCheckpoint(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGatewayCheckpoint
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
func skipGatewayCheckpoint(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGatewayCheckpoint
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
					return 0, ErrIntOverflowGatewayCheckpoint
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
					return 0, ErrIntOverflowGatewayCheckpoint
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
				return 0, ErrInvalidLengthGatewayCheckpoint
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGatewayCheckpoint
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGatewayCheckpoint
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGatewayCheckpoint        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGatewayCheckpoint          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGatewayCheckpoint = fmt.Errorf("proto: unexpected end of group")
)
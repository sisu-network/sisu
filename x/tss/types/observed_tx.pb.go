// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: observed_tx.proto

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

// Note: all new fields in this proto MUST be included in SerializeWithoutSigner in the go file.
type ObservedTxWithSigner struct {
	Signer string      `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *ObservedTx `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ObservedTxWithSigner) Reset()         { *m = ObservedTxWithSigner{} }
func (m *ObservedTxWithSigner) String() string { return proto.CompactTextString(m) }
func (*ObservedTxWithSigner) ProtoMessage()    {}
func (*ObservedTxWithSigner) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b040b825bdb87c3, []int{0}
}
func (m *ObservedTxWithSigner) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ObservedTxWithSigner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ObservedTxWithSigner.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ObservedTxWithSigner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObservedTxWithSigner.Merge(m, src)
}
func (m *ObservedTxWithSigner) XXX_Size() int {
	return m.Size()
}
func (m *ObservedTxWithSigner) XXX_DiscardUnknown() {
	xxx_messageInfo_ObservedTxWithSigner.DiscardUnknown(m)
}

var xxx_messageInfo_ObservedTxWithSigner proto.InternalMessageInfo

func (m *ObservedTxWithSigner) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *ObservedTxWithSigner) GetData() *ObservedTx {
	if m != nil {
		return m.Data
	}
	return nil
}

type ObservedTx struct {
	Chain       string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	BlockHeight int64  `protobuf:"varint,2,opt,name=blockHeight,proto3" json:"blockHeight,omitempty"`
	TxHash      string `protobuf:"bytes,3,opt,name=txHash,proto3" json:"txHash,omitempty"`
	Serialized  []byte `protobuf:"bytes,4,opt,name=serialized,proto3" json:"serialized,omitempty"`
}

func (m *ObservedTx) Reset()         { *m = ObservedTx{} }
func (m *ObservedTx) String() string { return proto.CompactTextString(m) }
func (*ObservedTx) ProtoMessage()    {}
func (*ObservedTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b040b825bdb87c3, []int{1}
}
func (m *ObservedTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ObservedTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ObservedTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ObservedTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObservedTx.Merge(m, src)
}
func (m *ObservedTx) XXX_Size() int {
	return m.Size()
}
func (m *ObservedTx) XXX_DiscardUnknown() {
	xxx_messageInfo_ObservedTx.DiscardUnknown(m)
}

var xxx_messageInfo_ObservedTx proto.InternalMessageInfo

func (m *ObservedTx) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *ObservedTx) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *ObservedTx) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *ObservedTx) GetSerialized() []byte {
	if m != nil {
		return m.Serialized
	}
	return nil
}

func init() {
	proto.RegisterType((*ObservedTxWithSigner)(nil), "types.ObservedTxWithSigner")
	proto.RegisterType((*ObservedTx)(nil), "types.ObservedTx")
}

func init() { proto.RegisterFile("observed_tx.proto", fileDescriptor_2b040b825bdb87c3) }

var fileDescriptor_2b040b825bdb87c3 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcc, 0x4f, 0x2a, 0x4e,
	0x2d, 0x2a, 0x4b, 0x4d, 0x89, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d,
	0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x0a, 0xe5, 0x12, 0xf1, 0x87, 0xca, 0x85, 0x54, 0x84, 0x67, 0x96,
	0x64, 0x04, 0x67, 0xa6, 0xe7, 0xa5, 0x16, 0x09, 0x89, 0x71, 0xb1, 0x15, 0x83, 0x59, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x9e, 0x90, 0x2a, 0x17, 0x4b, 0x4a, 0x62, 0x49, 0xa2, 0x04,
	0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0xa0, 0x1e, 0xd8, 0x14, 0x3d, 0x84, 0x11, 0x41, 0x60, 0x69,
	0xa5, 0x1a, 0x2e, 0x2e, 0x84, 0x98, 0x90, 0x08, 0x17, 0x6b, 0x72, 0x46, 0x62, 0x66, 0x1e, 0xd4,
	0x2c, 0x08, 0x47, 0x48, 0x81, 0x8b, 0x3b, 0x29, 0x27, 0x3f, 0x39, 0xdb, 0x23, 0x35, 0x33, 0x3d,
	0xa3, 0x04, 0x6c, 0x22, 0x73, 0x10, 0xb2, 0x10, 0xc8, 0x11, 0x25, 0x15, 0x1e, 0x89, 0xc5, 0x19,
	0x12, 0xcc, 0x10, 0x47, 0x40, 0x78, 0x42, 0x72, 0x5c, 0x5c, 0xc5, 0xa9, 0x45, 0x99, 0x89, 0x39,
	0x99, 0x55, 0xa9, 0x29, 0x12, 0x2c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x48, 0x22, 0x4e, 0x12, 0x27,
	0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c,
	0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x90, 0xc4, 0x06, 0xf6, 0xbc, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0xf0, 0x0e, 0xee, 0x0d, 0x11, 0x01, 0x00, 0x00,
}

func (m *ObservedTxWithSigner) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ObservedTxWithSigner) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ObservedTxWithSigner) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintObservedTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ObservedTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ObservedTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ObservedTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Serialized) > 0 {
		i -= len(m.Serialized)
		copy(dAtA[i:], m.Serialized)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.Serialized)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x1a
	}
	if m.BlockHeight != 0 {
		i = encodeVarintObservedTx(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintObservedTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovObservedTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ObservedTxWithSigner) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovObservedTx(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovObservedTx(uint64(l))
	}
	return n
}

func (m *ObservedTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovObservedTx(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovObservedTx(uint64(m.BlockHeight))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovObservedTx(uint64(l))
	}
	l = len(m.Serialized)
	if l > 0 {
		n += 1 + l + sovObservedTx(uint64(l))
	}
	return n
}

func sovObservedTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozObservedTx(x uint64) (n int) {
	return sovObservedTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ObservedTxWithSigner) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowObservedTx
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
			return fmt.Errorf("proto: ObservedTxWithSigner: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ObservedTxWithSigner: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObservedTx
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
				return ErrInvalidLengthObservedTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthObservedTx
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
					return ErrIntOverflowObservedTx
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
				return ErrInvalidLengthObservedTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthObservedTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &ObservedTx{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipObservedTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthObservedTx
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
func (m *ObservedTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowObservedTx
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
			return fmt.Errorf("proto: ObservedTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ObservedTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObservedTx
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
				return ErrInvalidLengthObservedTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthObservedTx
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
					return ErrIntOverflowObservedTx
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObservedTx
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
				return ErrInvalidLengthObservedTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthObservedTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Serialized", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObservedTx
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
				return ErrInvalidLengthObservedTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthObservedTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Serialized = append(m.Serialized[:0], dAtA[iNdEx:postIndex]...)
			if m.Serialized == nil {
				m.Serialized = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipObservedTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthObservedTx
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
func skipObservedTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowObservedTx
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
					return 0, ErrIntOverflowObservedTx
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
					return 0, ErrIntOverflowObservedTx
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
				return 0, ErrInvalidLengthObservedTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupObservedTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthObservedTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthObservedTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowObservedTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupObservedTx = fmt.Errorf("proto: unexpected end of group")
)

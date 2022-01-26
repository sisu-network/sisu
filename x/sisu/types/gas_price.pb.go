// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/gas_price.proto

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

type GasPriceMsg struct {
	Chain       string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	BlockHeight int64  `protobuf:"varint,2,opt,name=blockHeight,proto3" json:"blockHeight,omitempty"`
	GasPrice    int64  `protobuf:"varint,3,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	Signer      string `protobuf:"bytes,4,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *GasPriceMsg) Reset()         { *m = GasPriceMsg{} }
func (m *GasPriceMsg) String() string { return proto.CompactTextString(m) }
func (*GasPriceMsg) ProtoMessage()    {}
func (*GasPriceMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_1a59567a734f71ce, []int{0}
}
func (m *GasPriceMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GasPriceMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GasPriceMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GasPriceMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPriceMsg.Merge(m, src)
}
func (m *GasPriceMsg) XXX_Size() int {
	return m.Size()
}
func (m *GasPriceMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPriceMsg.DiscardUnknown(m)
}

var xxx_messageInfo_GasPriceMsg proto.InternalMessageInfo

func (m *GasPriceMsg) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *GasPriceMsg) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *GasPriceMsg) GetGasPrice() int64 {
	if m != nil {
		return m.GasPrice
	}
	return 0
}

func (m *GasPriceMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

type GasPriceRecord struct {
	Messages []*GasPriceMsg `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (m *GasPriceRecord) Reset()         { *m = GasPriceRecord{} }
func (m *GasPriceRecord) String() string { return proto.CompactTextString(m) }
func (*GasPriceRecord) ProtoMessage()    {}
func (*GasPriceRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_1a59567a734f71ce, []int{1}
}
func (m *GasPriceRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GasPriceRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GasPriceRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GasPriceRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPriceRecord.Merge(m, src)
}
func (m *GasPriceRecord) XXX_Size() int {
	return m.Size()
}
func (m *GasPriceRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPriceRecord.DiscardUnknown(m)
}

var xxx_messageInfo_GasPriceRecord proto.InternalMessageInfo

func (m *GasPriceRecord) GetMessages() []*GasPriceMsg {
	if m != nil {
		return m.Messages
	}
	return nil
}

func init() {
	proto.RegisterType((*GasPriceMsg)(nil), "types.GasPriceMsg")
	proto.RegisterType((*GasPriceRecord)(nil), "types.GasPriceRecord")
}

func init() { proto.RegisterFile("sisu/gas_price.proto", fileDescriptor_1a59567a734f71ce) }

var fileDescriptor_1a59567a734f71ce = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xce, 0x2c, 0x2e,
	0xd5, 0x4f, 0x4f, 0x2c, 0x8e, 0x2f, 0x28, 0xca, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0xaa, 0xe4, 0xe2, 0x76, 0x4f, 0x2c, 0x0e, 0x00,
	0x49, 0xf8, 0x16, 0xa7, 0x0b, 0x89, 0x70, 0xb1, 0x26, 0x67, 0x24, 0x66, 0xe6, 0x49, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x42, 0x0a, 0x5c, 0xdc, 0x49, 0x39, 0xf9, 0xc9, 0xd9, 0x1e,
	0xa9, 0x99, 0xe9, 0x19, 0x25, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0xc8, 0x42, 0x42, 0x52,
	0x5c, 0x1c, 0xe9, 0x50, 0x63, 0x24, 0x98, 0xc1, 0xd2, 0x70, 0xbe, 0x90, 0x18, 0x17, 0x5b, 0x71,
	0x66, 0x7a, 0x5e, 0x6a, 0x91, 0x04, 0x0b, 0xd8, 0x50, 0x28, 0x4f, 0xc9, 0x81, 0x8b, 0x0f, 0x66,
	0x75, 0x50, 0x6a, 0x72, 0x7e, 0x51, 0x8a, 0x90, 0x1e, 0x17, 0x47, 0x6e, 0x6a, 0x71, 0x71, 0x62,
	0x7a, 0x6a, 0xb1, 0x04, 0xb3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0x90, 0x1e, 0xd8, 0x99, 0x7a, 0x48,
	0x6e, 0x0c, 0x82, 0xab, 0x71, 0x72, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07,
	0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86,
	0x28, 0xcd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x90, 0xf7, 0x75,
	0xf3, 0x52, 0x4b, 0xca, 0xf3, 0x8b, 0xb2, 0xc1, 0x1c, 0xfd, 0x0a, 0x08, 0x05, 0x36, 0x3a, 0x89,
	0x0d, 0x1c, 0x1e, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x92, 0x95, 0xdc, 0xa4, 0x27, 0x01,
	0x00, 0x00,
}

func (m *GasPriceMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GasPriceMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GasPriceMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x22
	}
	if m.GasPrice != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.GasPrice))
		i--
		dAtA[i] = 0x18
	}
	if m.BlockHeight != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GasPriceRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GasPriceRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GasPriceRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Messages) > 0 {
		for iNdEx := len(m.Messages) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Messages[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGasPrice(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGasPrice(dAtA []byte, offset int, v uint64) int {
	offset -= sovGasPrice(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GasPriceMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovGasPrice(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovGasPrice(uint64(m.BlockHeight))
	}
	if m.GasPrice != 0 {
		n += 1 + sovGasPrice(uint64(m.GasPrice))
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovGasPrice(uint64(l))
	}
	return n
}

func (m *GasPriceRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Messages) > 0 {
		for _, e := range m.Messages {
			l = e.Size()
			n += 1 + l + sovGasPrice(uint64(l))
		}
	}
	return n
}

func sovGasPrice(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGasPrice(x uint64) (n int) {
	return sovGasPrice(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GasPriceMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGasPrice
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
			return fmt.Errorf("proto: GasPriceMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GasPriceMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
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
					return ErrIntOverflowGasPrice
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
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			m.GasPrice = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasPrice |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGasPrice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGasPrice
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
func (m *GasPriceRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGasPrice
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
			return fmt.Errorf("proto: GasPriceRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GasPriceRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Messages = append(m.Messages, &GasPriceMsg{})
			if err := m.Messages[len(m.Messages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGasPrice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGasPrice
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
func skipGasPrice(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGasPrice
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
					return 0, ErrIntOverflowGasPrice
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
					return 0, ErrIntOverflowGasPrice
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
				return 0, ErrInvalidLengthGasPrice
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGasPrice
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGasPrice
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGasPrice        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGasPrice          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGasPrice = fmt.Errorf("proto: unexpected end of group")
)
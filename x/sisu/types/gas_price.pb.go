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

// TODO: Deprecated this and use ExternalData instead.
type GasPriceMsg struct {
	Signer   string   `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Chains   []string `protobuf:"bytes,2,rep,name=chains,proto3" json:"chains,omitempty"`
	GasPrices []int64  `protobuf:"varint,3,rep,packed,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	BaseFees  []int64  `protobuf:"varint,4,rep,packed,name=base_fee,json=baseFee,proto3" json:"base_fee,omitempty"`
	Tips      []int64  `protobuf:"varint,5,rep,packed,name=tip,proto3" json:"tip,omitempty"`
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

func (m *GasPriceMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *GasPriceMsg) GetChains() []string {
	if m != nil {
		return m.Chains
	}
	return nil
}

func (m *GasPriceMsg) GetGasPrice() []int64 {
	if m != nil {
		return m.GasPrices
	}
	return nil
}

func (m *GasPriceMsg) GetBaseFee() []int64 {
	if m != nil {
		return m.BaseFees
	}
	return nil
}

func (m *GasPriceMsg) GetTip() []int64 {
	if m != nil {
		return m.Tips
	}
	return nil
}

type GasPriceRecord struct {
	GasPrice int64 `protobuf:"varint,3,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	BaseFee  int64 `protobuf:"varint,4,opt,name=base_fee,json=baseFee,proto3" json:"base_fee,omitempty"`
	Tip      int64 `protobuf:"varint,5,opt,name=tip,proto3" json:"tip,omitempty"`
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

func (m *GasPriceRecord) GetGasPrice() int64 {
	if m != nil {
		return m.GasPrice
	}
	return 0
}

func (m *GasPriceRecord) GetBaseFee() int64 {
	if m != nil {
		return m.BaseFee
	}
	return 0
}

func (m *GasPriceRecord) GetTip() int64 {
	if m != nil {
		return m.Tip
	}
	return 0
}

func init() {
	proto.RegisterType((*GasPriceMsg)(nil), "types.GasPriceMsg")
	proto.RegisterType((*GasPriceRecord)(nil), "types.GasPriceRecord")
}

func init() { proto.RegisterFile("sisu/gas_price.proto", fileDescriptor_1a59567a734f71ce) }

var fileDescriptor_1a59567a734f71ce = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xce, 0x2c, 0x2e,
	0xd5, 0x4f, 0x4f, 0x2c, 0x8e, 0x2f, 0x28, 0xca, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x6a, 0x67, 0xe4, 0xe2, 0x76, 0x4f, 0x2c, 0x0e,
	0x00, 0xc9, 0xf8, 0x16, 0xa7, 0x0b, 0x89, 0x71, 0xb1, 0x15, 0x67, 0xa6, 0xe7, 0xa5, 0x16, 0x49,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x79, 0x20, 0xf1, 0xe4, 0x8c, 0xc4, 0xcc, 0xbc, 0x62,
	0x09, 0x26, 0x05, 0x66, 0x90, 0x38, 0x84, 0x27, 0x24, 0xcd, 0xc5, 0x09, 0x37, 0x59, 0x82, 0x59,
	0x81, 0x59, 0x83, 0x39, 0x88, 0x23, 0x1d, 0x6a, 0x9e, 0x90, 0x24, 0x17, 0x47, 0x52, 0x62, 0x71,
	0x6a, 0x7c, 0x5a, 0x6a, 0xaa, 0x04, 0x0b, 0x58, 0x8e, 0x1d, 0xc4, 0x77, 0x4b, 0x4d, 0x15, 0x12,
	0xe0, 0x62, 0x2e, 0xc9, 0x2c, 0x90, 0x60, 0x05, 0x8b, 0x82, 0x98, 0x4a, 0x51, 0x5c, 0x7c, 0x30,
	0x87, 0x04, 0xa5, 0x26, 0xe7, 0x17, 0xa5, 0xa0, 0x9b, 0xcd, 0x88, 0xc7, 0x6c, 0x46, 0xac, 0x66,
	0x33, 0x42, 0xcd, 0x76, 0x72, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f,
	0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28,
	0xcd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x50, 0x38, 0xe9, 0xe6,
	0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0x83, 0x39, 0xfa, 0x15, 0x10, 0x0a, 0x1c, 0x54, 0x49, 0x6c,
	0xe0, 0x80, 0x33, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x86, 0xb3, 0x8b, 0x6f, 0x50, 0x01, 0x00,
	0x00,
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
	if len(m.Tips) > 0 {
		dAtA2 := make([]byte, len(m.Tips)*10)
		var j1 int
		for _, num1 := range m.Tips {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintGasPrice(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.BaseFees) > 0 {
		dAtA4 := make([]byte, len(m.BaseFees)*10)
		var j3 int
		for _, num1 := range m.BaseFees {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		i -= j3
		copy(dAtA[i:], dAtA4[:j3])
		i = encodeVarintGasPrice(dAtA, i, uint64(j3))
		i--
		dAtA[i] = 0x22
	}
	if len(m.GasPrices) > 0 {
		dAtA6 := make([]byte, len(m.GasPrices)*10)
		var j5 int
		for _, num1 := range m.GasPrices {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA6[j5] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j5++
			}
			dAtA6[j5] = uint8(num)
			j5++
		}
		i -= j5
		copy(dAtA[i:], dAtA6[:j5])
		i = encodeVarintGasPrice(dAtA, i, uint64(j5))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Chains) > 0 {
		for iNdEx := len(m.Chains) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Chains[iNdEx])
			copy(dAtA[i:], m.Chains[iNdEx])
			i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Chains[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Signer)))
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
	if m.Tip != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.Tip))
		i--
		dAtA[i] = 0x28
	}
	if m.BaseFee != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.BaseFee))
		i--
		dAtA[i] = 0x20
	}
	if m.GasPrice != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.GasPrice))
		i--
		dAtA[i] = 0x18
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
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovGasPrice(uint64(l))
	}
	if len(m.Chains) > 0 {
		for _, s := range m.Chains {
			l = len(s)
			n += 1 + l + sovGasPrice(uint64(l))
		}
	}
	if len(m.GasPrices) > 0 {
		l = 0
		for _, e := range m.GasPrices {
			l += sovGasPrice(uint64(e))
		}
		n += 1 + sovGasPrice(uint64(l)) + l
	}
	if len(m.BaseFees) > 0 {
		l = 0
		for _, e := range m.BaseFees {
			l += sovGasPrice(uint64(e))
		}
		n += 1 + sovGasPrice(uint64(l)) + l
	}
	if len(m.Tips) > 0 {
		l = 0
		for _, e := range m.Tips {
			l += sovGasPrice(uint64(e))
		}
		n += 1 + sovGasPrice(uint64(l)) + l
	}
	return n
}

func (m *GasPriceRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GasPrice != 0 {
		n += 1 + sovGasPrice(uint64(m.GasPrice))
	}
	if m.BaseFee != 0 {
		n += 1 + sovGasPrice(uint64(m.BaseFee))
	}
	if m.Tip != 0 {
		n += 1 + sovGasPrice(uint64(m.Tip))
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chains", wireType)
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
			m.Chains = append(m.Chains, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.GasPrices = append(m.GasPrices, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGasPrice
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGasPrice
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.GasPrices) == 0 {
					m.GasPrices = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGasPrice
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.GasPrices = append(m.GasPrices, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
		case 4:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.BaseFees = append(m.BaseFees, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGasPrice
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGasPrice
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.BaseFees) == 0 {
					m.BaseFees = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGasPrice
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.BaseFees = append(m.BaseFees, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
		case 5:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Tips = append(m.Tips, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGasPrice
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGasPrice
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Tips) == 0 {
					m.Tips = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGasPrice
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Tips = append(m.Tips, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Tip", wireType)
			}
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			m.BaseFee = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaseFee |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tip", wireType)
			}
			m.Tip = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tip |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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

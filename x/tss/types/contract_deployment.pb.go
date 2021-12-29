// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: contract_deployment.proto

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

type ContractDeployment struct {
	Chain        string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	ContractHash string `protobuf:"bytes,2,opt,name=contractHash,proto3" json:"contractHash,omitempty"`
	TxOutHash    string `protobuf:"bytes,3,opt,name=txOutHash,proto3" json:"txOutHash,omitempty"`
	Address      string `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Status       string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

func (m *ContractDeployment) Reset()         { *m = ContractDeployment{} }
func (m *ContractDeployment) String() string { return proto.CompactTextString(m) }
func (*ContractDeployment) ProtoMessage()    {}
func (*ContractDeployment) Descriptor() ([]byte, []int) {
	return fileDescriptor_06004a942c1f0ae4, []int{0}
}
func (m *ContractDeployment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ContractDeployment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ContractDeployment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ContractDeployment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContractDeployment.Merge(m, src)
}
func (m *ContractDeployment) XXX_Size() int {
	return m.Size()
}
func (m *ContractDeployment) XXX_DiscardUnknown() {
	xxx_messageInfo_ContractDeployment.DiscardUnknown(m)
}

var xxx_messageInfo_ContractDeployment proto.InternalMessageInfo

func (m *ContractDeployment) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *ContractDeployment) GetContractHash() string {
	if m != nil {
		return m.ContractHash
	}
	return ""
}

func (m *ContractDeployment) GetTxOutHash() string {
	if m != nil {
		return m.TxOutHash
	}
	return ""
}

func (m *ContractDeployment) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *ContractDeployment) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*ContractDeployment)(nil), "types.ContractDeployment")
}

func init() { proto.RegisterFile("contract_deployment.proto", fileDescriptor_06004a942c1f0ae4) }

var fileDescriptor_06004a942c1f0ae4 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4c, 0xce, 0xcf, 0x2b,
	0x29, 0x4a, 0x4c, 0x2e, 0x89, 0x4f, 0x49, 0x2d, 0xc8, 0xc9, 0xaf, 0xcc, 0x4d, 0xcd, 0x2b, 0xd1,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x9a, 0xc7, 0xc8,
	0x25, 0xe4, 0x0c, 0x55, 0xe4, 0x02, 0x57, 0x23, 0x24, 0xc2, 0xc5, 0x9a, 0x9c, 0x91, 0x98, 0x99,
	0x27, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe1, 0x08, 0x29, 0x71, 0xf1, 0xc0, 0x0c, 0xf4,
	0x48, 0x2c, 0xce, 0x90, 0x60, 0x02, 0x4b, 0xa2, 0x88, 0x09, 0xc9, 0x70, 0x71, 0x96, 0x54, 0xf8,
	0x97, 0x42, 0x14, 0x30, 0x83, 0x15, 0x20, 0x04, 0x84, 0x24, 0xb8, 0xd8, 0x13, 0x53, 0x52, 0x8a,
	0x52, 0x8b, 0x8b, 0x25, 0x58, 0xc0, 0x72, 0x30, 0xae, 0x90, 0x18, 0x17, 0x5b, 0x71, 0x49, 0x62,
	0x49, 0x69, 0xb1, 0x04, 0x2b, 0x58, 0x02, 0xca, 0x73, 0x92, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2,
	0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1,
	0xc6, 0x63, 0x39, 0x86, 0x24, 0x36, 0xb0, 0x47, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3e,
	0xb9, 0xd3, 0x53, 0xe5, 0x00, 0x00, 0x00,
}

func (m *ContractDeployment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ContractDeployment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ContractDeployment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintContractDeployment(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintContractDeployment(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.TxOutHash) > 0 {
		i -= len(m.TxOutHash)
		copy(dAtA[i:], m.TxOutHash)
		i = encodeVarintContractDeployment(dAtA, i, uint64(len(m.TxOutHash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ContractHash) > 0 {
		i -= len(m.ContractHash)
		copy(dAtA[i:], m.ContractHash)
		i = encodeVarintContractDeployment(dAtA, i, uint64(len(m.ContractHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintContractDeployment(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintContractDeployment(dAtA []byte, offset int, v uint64) int {
	offset -= sovContractDeployment(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ContractDeployment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovContractDeployment(uint64(l))
	}
	l = len(m.ContractHash)
	if l > 0 {
		n += 1 + l + sovContractDeployment(uint64(l))
	}
	l = len(m.TxOutHash)
	if l > 0 {
		n += 1 + l + sovContractDeployment(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovContractDeployment(uint64(l))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovContractDeployment(uint64(l))
	}
	return n
}

func sovContractDeployment(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozContractDeployment(x uint64) (n int) {
	return sovContractDeployment(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ContractDeployment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContractDeployment
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
			return fmt.Errorf("proto: ContractDeployment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ContractDeployment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractDeployment
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
				return ErrInvalidLengthContractDeployment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractDeployment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractDeployment
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
				return ErrInvalidLengthContractDeployment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractDeployment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxOutHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractDeployment
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
				return ErrInvalidLengthContractDeployment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractDeployment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxOutHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractDeployment
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
				return ErrInvalidLengthContractDeployment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractDeployment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractDeployment
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
				return ErrInvalidLengthContractDeployment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractDeployment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContractDeployment(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContractDeployment
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
func skipContractDeployment(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContractDeployment
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
					return 0, ErrIntOverflowContractDeployment
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
					return 0, ErrIntOverflowContractDeployment
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
				return 0, ErrInvalidLengthContractDeployment
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupContractDeployment
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthContractDeployment
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthContractDeployment        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContractDeployment          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupContractDeployment = fmt.Errorf("proto: unexpected end of group")
)

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/contract_change_liquid_address.proto

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

type ChangeSetPoolAddressMsg struct {
	Signer string               `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *ChangeLiquidAddress `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ChangeSetPoolAddressMsg) Reset()         { *m = ChangeSetPoolAddressMsg{} }
func (m *ChangeSetPoolAddressMsg) String() string { return proto.CompactTextString(m) }
func (*ChangeSetPoolAddressMsg) ProtoMessage()    {}
func (*ChangeSetPoolAddressMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_23dc60e0fae21457, []int{0}
}
func (m *ChangeSetPoolAddressMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeSetPoolAddressMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeLiquidPoolAddressMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeSetPoolAddressMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeLiquidPoolAddressMsg.Merge(m, src)
}
func (m *ChangeSetPoolAddressMsg) XXX_Size() int {
	return m.Size()
}
func (m *ChangeSetPoolAddressMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeLiquidPoolAddressMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeLiquidPoolAddressMsg proto.InternalMessageInfo

func (m *ChangeSetPoolAddressMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *ChangeSetPoolAddressMsg) GetData() *ChangeLiquidAddress {
	if m != nil {
		return m.Data
	}
	return nil
}

type ChangeLiquidAddress struct {
	Chain            string `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Hash             string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	NewLiquidAddress string `protobuf:"bytes,3,opt,name=newLiquidAddress,proto3" json:"newLiquidAddress,omitempty"`
	Index            int32  `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *ChangeLiquidAddress) Reset()         { *m = ChangeLiquidAddress{} }
func (m *ChangeLiquidAddress) String() string { return proto.CompactTextString(m) }
func (*ChangeLiquidAddress) ProtoMessage()    {}
func (*ChangeLiquidAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_23dc60e0fae21457, []int{1}
}
func (m *ChangeLiquidAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChangeLiquidAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChangeLiquidAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChangeLiquidAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeLiquidAddress.Merge(m, src)
}
func (m *ChangeLiquidAddress) XXX_Size() int {
	return m.Size()
}
func (m *ChangeLiquidAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeLiquidAddress.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeLiquidAddress proto.InternalMessageInfo

func (m *ChangeLiquidAddress) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *ChangeLiquidAddress) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *ChangeLiquidAddress) GetNewLiquidAddress() string {
	if m != nil {
		return m.NewLiquidAddress
	}
	return ""
}

func (m *ChangeLiquidAddress) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*ChangeSetPoolAddressMsg)(nil), "types.ChangeSetPoolAddressMsg")
	proto.RegisterType((*ChangeLiquidAddress)(nil), "types.ChangeLiquidAddress")
}

func init() {
	proto.RegisterFile("sisu/contract_change_liquid_address.proto", fileDescriptor_23dc60e0fae21457)
}

var fileDescriptor_23dc60e0fae21457 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2c, 0xce, 0x2c, 0x2e,
	0xd5, 0x4f, 0xce, 0xcf, 0x2b, 0x29, 0x4a, 0x4c, 0x2e, 0x89, 0x4f, 0xce, 0x48, 0xcc, 0x4b, 0x4f,
	0x8d, 0xcf, 0xc9, 0x2c, 0x2c, 0xcd, 0x4c, 0x89, 0x4f, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x4a, 0xe1, 0x92,
	0x72, 0x06, 0xab, 0xf2, 0x01, 0x2b, 0x0a, 0xc8, 0xcf, 0xcf, 0x71, 0x84, 0xa8, 0xf3, 0x2d, 0x4e,
	0x17, 0x12, 0xe3, 0x62, 0x2b, 0xce, 0x4c, 0xcf, 0x4b, 0x2d, 0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x0c, 0x82, 0xf2, 0x84, 0xf4, 0xb8, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x98, 0x14, 0x18, 0x35,
	0xb8, 0x8d, 0xa4, 0xf4, 0xc0, 0x66, 0xe9, 0x21, 0x1b, 0x04, 0x35, 0x24, 0x08, 0xac, 0x4e, 0xa9,
	0x91, 0x91, 0x4b, 0x18, 0x8b, 0xac, 0x90, 0x08, 0x17, 0x6b, 0x72, 0x46, 0x62, 0x66, 0x1e, 0xd4,
	0x78, 0x08, 0x47, 0x48, 0x88, 0x8b, 0x25, 0x23, 0xb1, 0x38, 0x03, 0x6c, 0x3a, 0x67, 0x10, 0x98,
	0x2d, 0xa4, 0xc5, 0x25, 0x90, 0x97, 0x5a, 0x8e, 0xa2, 0x5b, 0x82, 0x19, 0x2c, 0x8f, 0x21, 0x0e,
	0x32, 0x35, 0x33, 0x2f, 0x25, 0xb5, 0x42, 0x82, 0x45, 0x81, 0x51, 0x83, 0x35, 0x08, 0xc2, 0x71,
	0x72, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c,
	0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xcd, 0xf4, 0xcc, 0x92,
	0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x50, 0x00, 0xea, 0xe6, 0xa5, 0x96, 0x94, 0xe7,
	0x17, 0x65, 0x83, 0x39, 0xfa, 0x15, 0x10, 0x0a, 0xec, 0xc5, 0x24, 0x36, 0x70, 0xe0, 0x19, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x1c, 0xbb, 0x9a, 0xeb, 0x69, 0x01, 0x00, 0x00,
}

func (m *ChangeSetPoolAddressMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeSetPoolAddressMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeSetPoolAddressMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ChangeLiquidAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChangeLiquidAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChangeLiquidAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x20
	}
	if len(m.NewLiquidAddress) > 0 {
		i -= len(m.NewLiquidAddress)
		copy(dAtA[i:], m.NewLiquidAddress)
		i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(len(m.NewLiquidAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintContractChangeLiquidAddress(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintContractChangeLiquidAddress(dAtA []byte, offset int, v uint64) int {
	offset -= sovContractChangeLiquidAddress(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChangeSetPoolAddressMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovContractChangeLiquidAddress(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovContractChangeLiquidAddress(uint64(l))
	}
	return n
}

func (m *ChangeLiquidAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovContractChangeLiquidAddress(uint64(l))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovContractChangeLiquidAddress(uint64(l))
	}
	l = len(m.NewLiquidAddress)
	if l > 0 {
		n += 1 + l + sovContractChangeLiquidAddress(uint64(l))
	}
	if m.Index != 0 {
		n += 1 + sovContractChangeLiquidAddress(uint64(m.Index))
	}
	return n
}

func sovContractChangeLiquidAddress(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozContractChangeLiquidAddress(x uint64) (n int) {
	return sovContractChangeLiquidAddress(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChangeSetPoolAddressMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContractChangeLiquidAddress
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
			return fmt.Errorf("proto: ChangeSetPoolAddressMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeSetPoolAddressMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeLiquidAddress
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
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
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
					return ErrIntOverflowContractChangeLiquidAddress
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
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &ChangeLiquidAddress{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContractChangeLiquidAddress(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
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
func (m *ChangeLiquidAddress) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContractChangeLiquidAddress
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
			return fmt.Errorf("proto: ChangeLiquidAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChangeLiquidAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeLiquidAddress
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
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
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
					return ErrIntOverflowContractChangeLiquidAddress
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
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewLiquidAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeLiquidAddress
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
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewLiquidAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContractChangeLiquidAddress
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
			skippy, err := skipContractChangeLiquidAddress(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContractChangeLiquidAddress
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
func skipContractChangeLiquidAddress(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContractChangeLiquidAddress
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
					return 0, ErrIntOverflowContractChangeLiquidAddress
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
					return 0, ErrIntOverflowContractChangeLiquidAddress
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
				return 0, ErrInvalidLengthContractChangeLiquidAddress
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupContractChangeLiquidAddress
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthContractChangeLiquidAddress
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthContractChangeLiquidAddress        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContractChangeLiquidAddress          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupContractChangeLiquidAddress = fmt.Errorf("proto: unexpected end of group")
)

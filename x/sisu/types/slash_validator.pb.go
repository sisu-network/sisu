// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/slash_validator.proto

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

type SlashValidatorMsg struct {
	Signer string              `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *SlashValidatorData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *SlashValidatorMsg) Reset()         { *m = SlashValidatorMsg{} }
func (m *SlashValidatorMsg) String() string { return proto.CompactTextString(m) }
func (*SlashValidatorMsg) ProtoMessage()    {}
func (*SlashValidatorMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a7dc40c13521e1, []int{0}
}
func (m *SlashValidatorMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SlashValidatorMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SlashValidatorMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SlashValidatorMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SlashValidatorMsg.Merge(m, src)
}
func (m *SlashValidatorMsg) XXX_Size() int {
	return m.Size()
}
func (m *SlashValidatorMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SlashValidatorMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SlashValidatorMsg proto.InternalMessageInfo

func (m *SlashValidatorMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *SlashValidatorMsg) GetData() *SlashValidatorData {
	if m != nil {
		return m.Data
	}
	return nil
}

type SlashValidatorData struct {
	NodeAddress string `protobuf:"bytes,1,opt,name=nodeAddress,proto3" json:"nodeAddress,omitempty"`
	SlashPoint  int64  `protobuf:"varint,2,opt,name=slashPoint,proto3" json:"slashPoint,omitempty"`
	Index       int32  `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *SlashValidatorData) Reset()         { *m = SlashValidatorData{} }
func (m *SlashValidatorData) String() string { return proto.CompactTextString(m) }
func (*SlashValidatorData) ProtoMessage()    {}
func (*SlashValidatorData) Descriptor() ([]byte, []int) {
	return fileDescriptor_47a7dc40c13521e1, []int{1}
}
func (m *SlashValidatorData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SlashValidatorData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SlashValidatorData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SlashValidatorData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SlashValidatorData.Merge(m, src)
}
func (m *SlashValidatorData) XXX_Size() int {
	return m.Size()
}
func (m *SlashValidatorData) XXX_DiscardUnknown() {
	xxx_messageInfo_SlashValidatorData.DiscardUnknown(m)
}

var xxx_messageInfo_SlashValidatorData proto.InternalMessageInfo

func (m *SlashValidatorData) GetNodeAddress() string {
	if m != nil {
		return m.NodeAddress
	}
	return ""
}

func (m *SlashValidatorData) GetSlashPoint() int64 {
	if m != nil {
		return m.SlashPoint
	}
	return 0
}

func (m *SlashValidatorData) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*SlashValidatorMsg)(nil), "types.SlashValidatorMsg")
	proto.RegisterType((*SlashValidatorData)(nil), "types.SlashValidatorData")
}

func init() { proto.RegisterFile("sisu/slash_validator.proto", fileDescriptor_47a7dc40c13521e1) }

var fileDescriptor_47a7dc40c13521e1 = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0xce, 0x2c, 0x2e,
	0xd5, 0x2f, 0xce, 0x49, 0x2c, 0xce, 0x88, 0x2f, 0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0xc9, 0x2f,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x8a, 0xe2,
	0x12, 0x0c, 0x06, 0xc9, 0x87, 0xc1, 0xa4, 0x7d, 0x8b, 0xd3, 0x85, 0xc4, 0xb8, 0xd8, 0x8a, 0x33,
	0xd3, 0xf3, 0x52, 0x8b, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xa0, 0x3c, 0x21, 0x5d, 0x2e,
	0x96, 0x94, 0xc4, 0x92, 0x44, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x6e, 0x23, 0x49, 0x3d, 0xb0, 0x11,
	0x7a, 0xa8, 0xfa, 0x5d, 0x12, 0x4b, 0x12, 0x83, 0xc0, 0xca, 0x94, 0x72, 0xb8, 0x84, 0x30, 0xe5,
	0x84, 0x14, 0xb8, 0xb8, 0xf3, 0xf2, 0x53, 0x52, 0x1d, 0x53, 0x52, 0x8a, 0x52, 0x8b, 0x8b, 0xa1,
	0x36, 0x20, 0x0b, 0x09, 0xc9, 0x71, 0x71, 0x81, 0xdd, 0x1c, 0x90, 0x9f, 0x99, 0x57, 0x02, 0xb6,
	0x8c, 0x39, 0x08, 0x49, 0x44, 0x48, 0x84, 0x8b, 0x35, 0x33, 0x2f, 0x25, 0xb5, 0x42, 0x82, 0x59,
	0x81, 0x51, 0x83, 0x35, 0x08, 0xc2, 0x71, 0x72, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39,
	0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63,
	0x39, 0x86, 0x28, 0xcd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x50,
	0x88, 0xe8, 0xe6, 0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0x83, 0x39, 0xfa, 0x15, 0x10, 0x0a, 0xec,
	0x97, 0x24, 0x36, 0x70, 0xe0, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x33, 0x1e, 0xfc, 0x3d,
	0x3a, 0x01, 0x00, 0x00,
}

func (m *SlashValidatorMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SlashValidatorMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SlashValidatorMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintSlashValidator(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintSlashValidator(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SlashValidatorData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SlashValidatorData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SlashValidatorData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintSlashValidator(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	if m.SlashPoint != 0 {
		i = encodeVarintSlashValidator(dAtA, i, uint64(m.SlashPoint))
		i--
		dAtA[i] = 0x10
	}
	if len(m.NodeAddress) > 0 {
		i -= len(m.NodeAddress)
		copy(dAtA[i:], m.NodeAddress)
		i = encodeVarintSlashValidator(dAtA, i, uint64(len(m.NodeAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSlashValidator(dAtA []byte, offset int, v uint64) int {
	offset -= sovSlashValidator(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SlashValidatorMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovSlashValidator(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovSlashValidator(uint64(l))
	}
	return n
}

func (m *SlashValidatorData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NodeAddress)
	if l > 0 {
		n += 1 + l + sovSlashValidator(uint64(l))
	}
	if m.SlashPoint != 0 {
		n += 1 + sovSlashValidator(uint64(m.SlashPoint))
	}
	if m.Index != 0 {
		n += 1 + sovSlashValidator(uint64(m.Index))
	}
	return n
}

func sovSlashValidator(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSlashValidator(x uint64) (n int) {
	return sovSlashValidator(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SlashValidatorMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSlashValidator
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
			return fmt.Errorf("proto: SlashValidatorMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SlashValidatorMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSlashValidator
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
				return ErrInvalidLengthSlashValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSlashValidator
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
					return ErrIntOverflowSlashValidator
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
				return ErrInvalidLengthSlashValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSlashValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &SlashValidatorData{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSlashValidator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSlashValidator
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
func (m *SlashValidatorData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSlashValidator
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
			return fmt.Errorf("proto: SlashValidatorData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SlashValidatorData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSlashValidator
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
				return ErrInvalidLengthSlashValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSlashValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashPoint", wireType)
			}
			m.SlashPoint = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSlashValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SlashPoint |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSlashValidator
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
			skippy, err := skipSlashValidator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSlashValidator
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
func skipSlashValidator(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSlashValidator
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
					return 0, ErrIntOverflowSlashValidator
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
					return 0, ErrIntOverflowSlashValidator
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
				return 0, ErrInvalidLengthSlashValidator
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSlashValidator
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSlashValidator
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSlashValidator        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSlashValidator          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSlashValidator = fmt.Errorf("proto: unexpected end of group")
)

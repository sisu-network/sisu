// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/reshare.proto

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

type ReshareData_Result int32

const (
	ReshareData_SUCCESS ReshareData_Result = 0
	ReshareData_FAILURE ReshareData_Result = 1
)

var ReshareData_Result_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILURE",
}

var ReshareData_Result_value = map[string]int32{
	"SUCCESS": 0,
	"FAILURE": 1,
}

func (x ReshareData_Result) String() string {
	return proto.EnumName(ReshareData_Result_name, int32(x))
}

func (ReshareData_Result) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4158c619501fbaca, []int{1, 0}
}

type ReshareResultWithSigner struct {
	Signer string       `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *ReshareData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ReshareResultWithSigner) Reset()         { *m = ReshareResultWithSigner{} }
func (m *ReshareResultWithSigner) String() string { return proto.CompactTextString(m) }
func (*ReshareResultWithSigner) ProtoMessage()    {}
func (*ReshareResultWithSigner) Descriptor() ([]byte, []int) {
	return fileDescriptor_4158c619501fbaca, []int{0}
}
func (m *ReshareResultWithSigner) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReshareResultWithSigner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReshareResultWithSigner.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReshareResultWithSigner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReshareResultWithSigner.Merge(m, src)
}
func (m *ReshareResultWithSigner) XXX_Size() int {
	return m.Size()
}
func (m *ReshareResultWithSigner) XXX_DiscardUnknown() {
	xxx_messageInfo_ReshareResultWithSigner.DiscardUnknown(m)
}

var xxx_messageInfo_ReshareResultWithSigner proto.InternalMessageInfo

func (m *ReshareResultWithSigner) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *ReshareResultWithSigner) GetData() *ReshareData {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReshareData struct {
	NewValidatorSetPubKeyBytes [][]byte           `protobuf:"bytes,1,rep,name=newValidatorSetPubKeyBytes,proto3" json:"newValidatorSetPubKeyBytes,omitempty"`
	Result                     ReshareData_Result `protobuf:"varint,2,opt,name=result,proto3,enum=types.ReshareData_Result" json:"result,omitempty"`
}

func (m *ReshareData) Reset()         { *m = ReshareData{} }
func (m *ReshareData) String() string { return proto.CompactTextString(m) }
func (*ReshareData) ProtoMessage()    {}
func (*ReshareData) Descriptor() ([]byte, []int) {
	return fileDescriptor_4158c619501fbaca, []int{1}
}
func (m *ReshareData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReshareData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReshareData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReshareData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReshareData.Merge(m, src)
}
func (m *ReshareData) XXX_Size() int {
	return m.Size()
}
func (m *ReshareData) XXX_DiscardUnknown() {
	xxx_messageInfo_ReshareData.DiscardUnknown(m)
}

var xxx_messageInfo_ReshareData proto.InternalMessageInfo

func (m *ReshareData) GetNewValidatorSetPubKeyBytes() [][]byte {
	if m != nil {
		return m.NewValidatorSetPubKeyBytes
	}
	return nil
}

func (m *ReshareData) GetResult() ReshareData_Result {
	if m != nil {
		return m.Result
	}
	return ReshareData_SUCCESS
}

func init() {
	proto.RegisterEnum("types.ReshareData_Result", ReshareData_Result_name, ReshareData_Result_value)
	proto.RegisterType((*ReshareResultWithSigner)(nil), "types.ReshareResultWithSigner")
	proto.RegisterType((*ReshareData)(nil), "types.ReshareData")
}

func init() { proto.RegisterFile("sisu/reshare.proto", fileDescriptor_4158c619501fbaca) }

var fileDescriptor_4158c619501fbaca = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0xce, 0x2c, 0x2e,
	0xd5, 0x2f, 0x4a, 0x2d, 0xce, 0x48, 0x2c, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x8a, 0xe4, 0x12, 0x0f, 0x82, 0x88, 0x07, 0xa5, 0x16, 0x97,
	0xe6, 0x94, 0x84, 0x67, 0x96, 0x64, 0x04, 0x67, 0xa6, 0xe7, 0xa5, 0x16, 0x09, 0x89, 0x71, 0xb1,
	0x15, 0x83, 0x59, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x9e, 0x90, 0x1a, 0x17, 0x4b,
	0x4a, 0x62, 0x49, 0xa2, 0x04, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x90, 0x1e, 0xd8, 0x20, 0x3d,
	0xa8, 0x29, 0x2e, 0x89, 0x25, 0x89, 0x41, 0x60, 0x79, 0xa5, 0x25, 0x8c, 0x5c, 0xdc, 0x48, 0xa2,
	0x42, 0x76, 0x5c, 0x52, 0x79, 0xa9, 0xe5, 0x61, 0x89, 0x39, 0x99, 0x29, 0x89, 0x25, 0xf9, 0x45,
	0xc1, 0xa9, 0x25, 0x01, 0xa5, 0x49, 0xde, 0xa9, 0x95, 0x4e, 0x95, 0x25, 0xa9, 0xc5, 0x12, 0x8c,
	0x0a, 0xcc, 0x1a, 0x3c, 0x41, 0x78, 0x54, 0x08, 0x19, 0x72, 0xb1, 0x15, 0x81, 0xdd, 0x08, 0xb6,
	0x99, 0xcf, 0x48, 0x12, 0xd3, 0x66, 0x3d, 0x88, 0x27, 0x82, 0xa0, 0x0a, 0x95, 0x94, 0xb8, 0xd8,
	0x20, 0x22, 0x42, 0xdc, 0x5c, 0xec, 0xc1, 0xa1, 0xce, 0xce, 0xae, 0xc1, 0xc1, 0x02, 0x0c, 0x20,
	0x8e, 0x9b, 0xa3, 0xa7, 0x4f, 0x68, 0x90, 0xab, 0x00, 0xa3, 0x93, 0xf3, 0x89, 0x47, 0x72, 0x8c,
	0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72,
	0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x69, 0xa6, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7,
	0xe7, 0xea, 0x83, 0x42, 0x50, 0x37, 0x2f, 0xb5, 0xa4, 0x3c, 0xbf, 0x28, 0x1b, 0xcc, 0xd1, 0xaf,
	0x80, 0x50, 0x60, 0x37, 0x24, 0xb1, 0x81, 0x03, 0xd5, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xa0,
	0x5f, 0x76, 0x5d, 0x6a, 0x01, 0x00, 0x00,
}

func (m *ReshareResultWithSigner) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReshareResultWithSigner) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReshareResultWithSigner) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintReshare(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintReshare(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ReshareData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReshareData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReshareData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Result != 0 {
		i = encodeVarintReshare(dAtA, i, uint64(m.Result))
		i--
		dAtA[i] = 0x10
	}
	if len(m.NewValidatorSetPubKeyBytes) > 0 {
		for iNdEx := len(m.NewValidatorSetPubKeyBytes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NewValidatorSetPubKeyBytes[iNdEx])
			copy(dAtA[i:], m.NewValidatorSetPubKeyBytes[iNdEx])
			i = encodeVarintReshare(dAtA, i, uint64(len(m.NewValidatorSetPubKeyBytes[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintReshare(dAtA []byte, offset int, v uint64) int {
	offset -= sovReshare(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ReshareResultWithSigner) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovReshare(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovReshare(uint64(l))
	}
	return n
}

func (m *ReshareData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.NewValidatorSetPubKeyBytes) > 0 {
		for _, b := range m.NewValidatorSetPubKeyBytes {
			l = len(b)
			n += 1 + l + sovReshare(uint64(l))
		}
	}
	if m.Result != 0 {
		n += 1 + sovReshare(uint64(m.Result))
	}
	return n
}

func sovReshare(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozReshare(x uint64) (n int) {
	return sovReshare(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ReshareResultWithSigner) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReshare
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
			return fmt.Errorf("proto: ReshareResultWithSigner: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReshareResultWithSigner: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReshare
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
				return ErrInvalidLengthReshare
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReshare
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
					return ErrIntOverflowReshare
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
				return ErrInvalidLengthReshare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthReshare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &ReshareData{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReshare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReshare
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
func (m *ReshareData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReshare
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
			return fmt.Errorf("proto: ReshareData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReshareData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewValidatorSetPubKeyBytes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReshare
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
				return ErrInvalidLengthReshare
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReshare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewValidatorSetPubKeyBytes = append(m.NewValidatorSetPubKeyBytes, make([]byte, postIndex-iNdEx))
			copy(m.NewValidatorSetPubKeyBytes[len(m.NewValidatorSetPubKeyBytes)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			m.Result = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReshare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Result |= ReshareData_Result(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipReshare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReshare
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
func skipReshare(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReshare
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
					return 0, ErrIntOverflowReshare
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
					return 0, ErrIntOverflowReshare
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
				return 0, ErrInvalidLengthReshare
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupReshare
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthReshare
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthReshare        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReshare          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupReshare = fmt.Errorf("proto: unexpected end of group")
)

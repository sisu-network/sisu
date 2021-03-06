// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/keysign.proto

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

// Cosmos message to broadcast KeysignResult
type KeysignResult struct {
	Signer    string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Success   bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	OutChain  string `protobuf:"bytes,3,opt,name=outChain,proto3" json:"outChain,omitempty"`
	OutHash   string `protobuf:"bytes,4,opt,name=outHash,proto3" json:"outHash,omitempty"`
	Tx        []byte `protobuf:"bytes,5,opt,name=tx,proto3" json:"tx,omitempty"`
	Signature []byte `protobuf:"bytes,6,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *KeysignResult) Reset()         { *m = KeysignResult{} }
func (m *KeysignResult) String() string { return proto.CompactTextString(m) }
func (*KeysignResult) ProtoMessage()    {}
func (*KeysignResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c231a9642d6e9c8, []int{0}
}
func (m *KeysignResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *KeysignResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_KeysignResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *KeysignResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeysignResult.Merge(m, src)
}
func (m *KeysignResult) XXX_Size() int {
	return m.Size()
}
func (m *KeysignResult) XXX_DiscardUnknown() {
	xxx_messageInfo_KeysignResult.DiscardUnknown(m)
}

var xxx_messageInfo_KeysignResult proto.InternalMessageInfo

func (m *KeysignResult) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *KeysignResult) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *KeysignResult) GetOutChain() string {
	if m != nil {
		return m.OutChain
	}
	return ""
}

func (m *KeysignResult) GetOutHash() string {
	if m != nil {
		return m.OutHash
	}
	return ""
}

func (m *KeysignResult) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *KeysignResult) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*KeysignResult)(nil), "types.KeysignResult")
}

func init() { proto.RegisterFile("sisu/keysign.proto", fileDescriptor_1c231a9642d6e9c8) }

var fileDescriptor_1c231a9642d6e9c8 = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0xce, 0x2c, 0x2e,
	0xd5, 0xcf, 0x4e, 0xad, 0x2c, 0xce, 0x4c, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x5a, 0xca, 0xc8, 0xc5, 0xeb, 0x0d, 0x91, 0x08, 0x4a, 0x2d,
	0x2e, 0xcd, 0x29, 0x11, 0x12, 0xe3, 0x62, 0x03, 0xf1, 0x52, 0x8b, 0x24, 0x18, 0x15, 0x18, 0x35,
	0x38, 0x83, 0xa0, 0x3c, 0x21, 0x09, 0x2e, 0xf6, 0xe2, 0xd2, 0xe4, 0xe4, 0xd4, 0xe2, 0x62, 0x09,
	0x26, 0x05, 0x46, 0x0d, 0x8e, 0x20, 0x18, 0x57, 0x48, 0x8a, 0x8b, 0x23, 0xbf, 0xb4, 0xc4, 0x39,
	0x23, 0x31, 0x33, 0x4f, 0x82, 0x19, 0xac, 0x07, 0xce, 0x07, 0xe9, 0xca, 0x2f, 0x2d, 0xf1, 0x48,
	0x2c, 0xce, 0x90, 0x60, 0x01, 0x4b, 0xc1, 0xb8, 0x42, 0x7c, 0x5c, 0x4c, 0x25, 0x15, 0x12, 0xac,
	0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x4c, 0x25, 0x15, 0x42, 0x32, 0x5c, 0x9c, 0x20, 0x9b, 0x12, 0x4b,
	0x4a, 0x8b, 0x52, 0x25, 0xd8, 0xc0, 0xc2, 0x08, 0x01, 0x27, 0xe7, 0x13, 0x8f, 0xe4, 0x18, 0x2f,
	0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18,
	0x6e, 0x3c, 0x96, 0x63, 0x88, 0xd2, 0x4c, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf,
	0xd5, 0x07, 0xf9, 0x53, 0x37, 0x2f, 0xb5, 0xa4, 0x3c, 0xbf, 0x28, 0x1b, 0xcc, 0xd1, 0xaf, 0x80,
	0x50, 0x60, 0xcf, 0x26, 0xb1, 0x81, 0xbd, 0x6e, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x97, 0x19,
	0xff, 0x0b, 0x10, 0x01, 0x00, 0x00,
}

func (m *KeysignResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *KeysignResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *KeysignResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintKeysign(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Tx) > 0 {
		i -= len(m.Tx)
		copy(dAtA[i:], m.Tx)
		i = encodeVarintKeysign(dAtA, i, uint64(len(m.Tx)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.OutHash) > 0 {
		i -= len(m.OutHash)
		copy(dAtA[i:], m.OutHash)
		i = encodeVarintKeysign(dAtA, i, uint64(len(m.OutHash)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.OutChain) > 0 {
		i -= len(m.OutChain)
		copy(dAtA[i:], m.OutChain)
		i = encodeVarintKeysign(dAtA, i, uint64(len(m.OutChain)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintKeysign(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintKeysign(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeysign(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *KeysignResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovKeysign(uint64(l))
	}
	if m.Success {
		n += 2
	}
	l = len(m.OutChain)
	if l > 0 {
		n += 1 + l + sovKeysign(uint64(l))
	}
	l = len(m.OutHash)
	if l > 0 {
		n += 1 + l + sovKeysign(uint64(l))
	}
	l = len(m.Tx)
	if l > 0 {
		n += 1 + l + sovKeysign(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovKeysign(uint64(l))
	}
	return n
}

func sovKeysign(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeysign(x uint64) (n int) {
	return sovKeysign(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *KeysignResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeysign
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
			return fmt.Errorf("proto: KeysignResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeysignResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
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
				return ErrInvalidLengthKeysign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthKeysign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Success = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
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
				return ErrInvalidLengthKeysign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthKeysign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
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
				return ErrInvalidLengthKeysign
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthKeysign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tx", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
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
				return ErrInvalidLengthKeysign
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeysign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tx = append(m.Tx[:0], dAtA[iNdEx:postIndex]...)
			if m.Tx == nil {
				m.Tx = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeysign
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
				return ErrInvalidLengthKeysign
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeysign
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeysign(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthKeysign
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
func skipKeysign(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeysign
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
					return 0, ErrIntOverflowKeysign
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
					return 0, ErrIntOverflowKeysign
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
				return 0, ErrInvalidLengthKeysign
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeysign
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeysign
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthKeysign        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowKeysign          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupKeysign = fmt.Errorf("proto: unexpected end of group")
)

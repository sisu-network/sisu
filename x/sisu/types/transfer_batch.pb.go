// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sisu/transfer_batch.proto

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

type TransferBatchMsg struct {
	Signer string         `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Data   *TransferBatch `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *TransferBatchMsg) Reset()         { *m = TransferBatchMsg{} }
func (m *TransferBatchMsg) String() string { return proto.CompactTextString(m) }
func (*TransferBatchMsg) ProtoMessage()    {}
func (*TransferBatchMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_f426c263829dc8e8, []int{0}
}
func (m *TransferBatchMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferBatchMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferBatchMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferBatchMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferBatchMsg.Merge(m, src)
}
func (m *TransferBatchMsg) XXX_Size() int {
	return m.Size()
}
func (m *TransferBatchMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferBatchMsg.DiscardUnknown(m)
}

var xxx_messageInfo_TransferBatchMsg proto.InternalMessageInfo

func (m *TransferBatchMsg) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *TransferBatchMsg) GetData() *TransferBatch {
	if m != nil {
		return m.Data
	}
	return nil
}

type TransferBatch struct {
	Chain      string      `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Transfers  []*Transfer `protobuf:"bytes,2,rep,name=transfers,proto3" json:"transfers,omitempty"`
	StartBlock int64       `protobuf:"varint,3,opt,name=startBlock,proto3" json:"startBlock,omitempty"`
	Attemp     int32       `protobuf:"varint,4,opt,name=attemp,proto3" json:"attemp,omitempty"`
}

func (m *TransferBatch) Reset()         { *m = TransferBatch{} }
func (m *TransferBatch) String() string { return proto.CompactTextString(m) }
func (*TransferBatch) ProtoMessage()    {}
func (*TransferBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_f426c263829dc8e8, []int{1}
}
func (m *TransferBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferBatch.Merge(m, src)
}
func (m *TransferBatch) XXX_Size() int {
	return m.Size()
}
func (m *TransferBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferBatch.DiscardUnknown(m)
}

var xxx_messageInfo_TransferBatch proto.InternalMessageInfo

func (m *TransferBatch) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *TransferBatch) GetTransfers() []*Transfer {
	if m != nil {
		return m.Transfers
	}
	return nil
}

func (m *TransferBatch) GetStartBlock() int64 {
	if m != nil {
		return m.StartBlock
	}
	return 0
}

func (m *TransferBatch) GetAttemp() int32 {
	if m != nil {
		return m.Attemp
	}
	return 0
}

type Transfer struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Recipient string `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Token     string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Amount    string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *Transfer) Reset()         { *m = Transfer{} }
func (m *Transfer) String() string { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()    {}
func (*Transfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_f426c263829dc8e8, []int{2}
}
func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(m, src)
}
func (m *Transfer) XXX_Size() int {
	return m.Size()
}
func (m *Transfer) XXX_DiscardUnknown() {
	xxx_messageInfo_Transfer.DiscardUnknown(m)
}

var xxx_messageInfo_Transfer proto.InternalMessageInfo

func (m *Transfer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Transfer) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

func (m *Transfer) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Transfer) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func init() {
	proto.RegisterType((*TransferBatchMsg)(nil), "types.TransferBatchMsg")
	proto.RegisterType((*TransferBatch)(nil), "types.TransferBatch")
	proto.RegisterType((*Transfer)(nil), "types.Transfer")
}

func init() { proto.RegisterFile("sisu/transfer_batch.proto", fileDescriptor_f426c263829dc8e8) }

var fileDescriptor_f426c263829dc8e8 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0x4e, 0x02, 0x31,
	0x10, 0xc6, 0xe9, 0xf2, 0x27, 0xee, 0x10, 0xff, 0xa4, 0x21, 0x66, 0x4d, 0x4c, 0xb3, 0xe1, 0xb4,
	0x1e, 0x58, 0x12, 0x7c, 0x03, 0x3c, 0x7b, 0x69, 0x38, 0x79, 0x31, 0x65, 0x29, 0xd0, 0x20, 0xed,
	0xa6, 0x1d, 0xa2, 0xbe, 0x83, 0x07, 0x1f, 0xcb, 0x23, 0x47, 0x8f, 0x06, 0x5e, 0xc4, 0x6c, 0xb7,
	0x04, 0xf1, 0xd4, 0xfc, 0x66, 0xa6, 0xdf, 0xf7, 0x4d, 0x06, 0x6e, 0x9c, 0x72, 0x9b, 0x21, 0x5a,
	0xa1, 0xdd, 0x5c, 0xda, 0xe7, 0xa9, 0xc0, 0x62, 0x99, 0x97, 0xd6, 0xa0, 0xa1, 0x6d, 0x7c, 0x2f,
	0xa5, 0xeb, 0x4f, 0xe0, 0x6a, 0x12, 0xda, 0xe3, 0xaa, 0xfb, 0xe8, 0x16, 0xf4, 0x1a, 0x3a, 0x4e,
	0x2d, 0xb4, 0xb4, 0x09, 0x49, 0x49, 0x16, 0xf3, 0x40, 0x34, 0x83, 0xd6, 0x4c, 0xa0, 0x48, 0xa2,
	0x94, 0x64, 0xdd, 0x51, 0x2f, 0xf7, 0x0a, 0xf9, 0xc9, 0x77, 0xee, 0x27, 0xfa, 0x1f, 0x04, 0xce,
	0x4f, 0xea, 0xb4, 0x07, 0xed, 0x62, 0x29, 0x94, 0x0e, 0x92, 0x35, 0xd0, 0x01, 0xc4, 0x87, 0x70,
	0x2e, 0x89, 0xd2, 0x66, 0xd6, 0x1d, 0x5d, 0xfe, 0x93, 0xe5, 0xc7, 0x09, 0xca, 0x00, 0x1c, 0x0a,
	0x8b, 0xe3, 0x17, 0x53, 0xac, 0x92, 0x66, 0x4a, 0xb2, 0x26, 0xff, 0x53, 0xa9, 0x82, 0x0b, 0x44,
	0xb9, 0x2e, 0x93, 0x56, 0x4a, 0xb2, 0x36, 0x0f, 0xd4, 0x9f, 0xc3, 0xd9, 0x41, 0x8e, 0x5e, 0x40,
	0xa4, 0x66, 0x21, 0x45, 0xa4, 0x66, 0xf4, 0x16, 0x62, 0x2b, 0x0b, 0x55, 0x2a, 0xa9, 0xd1, 0x6f,
	0x16, 0xf3, 0x63, 0xa1, 0x8a, 0x8d, 0x66, 0x25, 0xb5, 0x37, 0x8b, 0x79, 0x0d, 0xde, 0x67, 0x6d,
	0x36, 0x1a, 0xbd, 0x4f, 0xcc, 0x03, 0x8d, 0x1f, 0xbe, 0x76, 0x8c, 0x6c, 0x77, 0x8c, 0xfc, 0xec,
	0x18, 0xf9, 0xdc, 0xb3, 0xc6, 0x76, 0xcf, 0x1a, 0xdf, 0x7b, 0xd6, 0x78, 0xba, 0x5b, 0x28, 0x5c,
	0x6e, 0xa6, 0x79, 0x61, 0xd6, 0xc3, 0xea, 0x26, 0x03, 0x2d, 0xf1, 0xd5, 0xd8, 0x95, 0x87, 0xe1,
	0x5b, 0xfd, 0xf8, 0xc5, 0xa7, 0x1d, 0x7f, 0x9f, 0xfb, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3e,
	0x05, 0xb8, 0xe5, 0xbc, 0x01, 0x00, 0x00,
}

func (m *TransferBatchMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferBatchMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferBatchMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintTransferBatch(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TransferBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Attemp != 0 {
		i = encodeVarintTransferBatch(dAtA, i, uint64(m.Attemp))
		i--
		dAtA[i] = 0x20
	}
	if m.StartBlock != 0 {
		i = encodeVarintTransferBatch(dAtA, i, uint64(m.StartBlock))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Transfers) > 0 {
		for iNdEx := len(m.Transfers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Transfers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTransferBatch(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Transfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Transfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Transfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTransferBatch(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTransferBatch(dAtA []byte, offset int, v uint64) int {
	offset -= sovTransferBatch(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TransferBatchMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	return n
}

func (m *TransferBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Chain)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	if len(m.Transfers) > 0 {
		for _, e := range m.Transfers {
			l = e.Size()
			n += 1 + l + sovTransferBatch(uint64(l))
		}
	}
	if m.StartBlock != 0 {
		n += 1 + sovTransferBatch(uint64(m.StartBlock))
	}
	if m.Attemp != 0 {
		n += 1 + sovTransferBatch(uint64(m.Attemp))
	}
	return n
}

func (m *Transfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovTransferBatch(uint64(l))
	}
	return n
}

func sovTransferBatch(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTransferBatch(x uint64) (n int) {
	return sovTransferBatch(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TransferBatchMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferBatch
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
			return fmt.Errorf("proto: TransferBatchMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferBatchMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
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
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &TransferBatch{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransferBatch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferBatch
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
func (m *TransferBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferBatch
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
			return fmt.Errorf("proto: TransferBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transfers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Transfers = append(m.Transfers, &Transfer{})
			if err := m.Transfers[len(m.Transfers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartBlock", wireType)
			}
			m.StartBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attemp", wireType)
			}
			m.Attemp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Attemp |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTransferBatch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferBatch
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
func (m *Transfer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransferBatch
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
			return fmt.Errorf("proto: Transfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Transfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransferBatch
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
				return ErrInvalidLengthTransferBatch
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransferBatch
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransferBatch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransferBatch
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
func skipTransferBatch(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransferBatch
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
					return 0, ErrIntOverflowTransferBatch
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
					return 0, ErrIntOverflowTransferBatch
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
				return 0, ErrInvalidLengthTransferBatch
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTransferBatch
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTransferBatch
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTransferBatch        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransferBatch          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTransferBatch = fmt.Errorf("proto: unexpected end of group")
)

package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenProposal{}

func NewMsgKeygenProposal(signer string, chain string, id string, expireBlock int64) *KeygenProposal {
	return &KeygenProposal{
		Signer:          signer,
		Chain:           chain,
		Id:              id,
		ExpirationBlock: expireBlock,
	}
}

// Route ...
func (msg *KeygenProposal) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenProposal) Type() string {
	return MSG_TYPE_KEYGEN_PROPOSAL
}

// GetSigners ...
func (msg *KeygenProposal) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenProposal) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// Serialize this message without the signer. This is similar to MarshalToSizedBuffer with the
// signer encoding removed. Any change in the proto file should also change this function.
func (m *KeygenProposal) SerializeWithoutSigner() []byte {
	size := m.Size()
	dAtA := make([]byte, size)

	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ExpirationBlock != 0 {
		i = encodeVarintKeygen(dAtA, i, uint64(m.ExpirationBlock))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintKeygen(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintKeygen(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0x12
	}

	return dAtA[:len(dAtA)-i]
}

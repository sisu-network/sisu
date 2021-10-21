package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ObservedTx{}

func NewObservedTxs(signer string, chain string, txHash string, blockHeight int64, serialized []byte) *ObservedTx {
	return &ObservedTx{
		Signer:      signer,
		Chain:       chain,
		TxHash:      txHash,
		BlockHeight: blockHeight,
		Serialized:  serialized,
	}
}

// Route ...
func (msg *ObservedTx) Route() string {
	return RouterKey
}

// Type ...
func (msg *ObservedTx) Type() string {
	return MSG_TYPE_OBSERVED_TX
}

// GetSigners ...
func (msg *ObservedTx) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ObservedTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ObservedTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ObservedTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// Serialize this message without the signer. This is similar to MarshalToSizedBuffer with the
// signer encoding removed. Any change in the proto file should also change this function.
func (m *ObservedTx) SerializeWithoutSigner() []byte {
	size := m.Size()
	dAtA := make([]byte, size)

	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Serialized) > 0 {
		i -= len(m.Serialized)
		copy(dAtA[i:], m.Serialized)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.Serialized)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x22
	}
	if m.BlockHeight != 0 {
		i = encodeVarintObservedTx(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Chain) > 0 {
		i -= len(m.Chain)
		copy(dAtA[i:], m.Chain)
		i = encodeVarintObservedTx(dAtA, i, uint64(len(m.Chain)))
		i--
		dAtA[i] = 0x12
	}

	// No signer.

	return dAtA[:len(dAtA)-i]
}

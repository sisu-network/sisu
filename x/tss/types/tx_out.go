package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/utils"
)

// TxOutStatus alias for tx out status
type TxOutStatus string

const (
	TxOutStatusConfirmed TxOutStatus = "confirmed"
	// TxOutStatusPreBroadcast after produce from observedTx, the txOut is ready to broadcast to Cosmos chain
	TxOutStatusPreBroadcast TxOutStatus = "pre_broadcast"
	// TxOutStatusBroadcasted TxOut has broadcasted to Cosmos chain successfully
	TxOutStatusBroadcasted TxOutStatus = "broadcasted"

	TxOutStatusDeployingToBlockchain TxOutStatus = "deploying_to_chain"
	TxOutStatusDeployedToBlockchain  TxOutStatus = "deployed_to_chain"

	// TxOutStatusPreSigning txOut is ready to be signed
	TxOutStatusPreSigning TxOutStatus = "pre_signing"
	// TxOutStatusSigning txOut is in singing progress
	TxOutStatusSigning TxOutStatus = "signing"
	// TxOutStatusSigned txOut is signed successfully
	TxOutStatusSigned TxOutStatus = "signed"
	// TxOutStatusSignFailed signing progress is failed
	TxOutStatusSignFailed TxOutStatus = "sign_failed"
)

var _ sdk.Msg = &TxOut{}

func NewMsgTxOut(signer string, inBlockHeight int64, inChain string, inHash string, outChain string, outBytes []byte) *TxOut {
	return &TxOut{
		Signer:        signer,
		InBlockHeight: inBlockHeight,
		InChain:       inChain,
		OutChain:      outChain,
		InHash:        inHash,
		OutBytes:      outBytes,
	}
}

// Route ...
func (msg *TxOut) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOut) Type() string {
	return MsgTypeTxOut
}

// GetSigners ...
func (msg *TxOut) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOut) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

func (msg *TxOut) GetHash() string {
	return utils.KeccakHash32(msg.OutChain + string(msg.OutBytes))
}

// Serialize this message without the signer. This is similar to MarshalToSizedBuffer with the
// signer encoding removed. Any change in the proto file should also change this function.
func (m *TxOut) SerializeWithoutSigner() []byte {
	size := m.Size()
	dAtA := make([]byte, size)

	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OutBytes) > 0 {
		i -= len(m.OutBytes)
		copy(dAtA[i:], m.OutBytes)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.OutBytes)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.InHash) > 0 {
		i -= len(m.InHash)
		copy(dAtA[i:], m.InHash)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.InHash)))
		i--
		dAtA[i] = 0x32
	}
	if m.InBlockHeight != 0 {
		i = encodeVarintTxOut(dAtA, i, uint64(m.InBlockHeight))
		i--
		dAtA[i] = 0x28
	}
	if len(m.OutChain) > 0 {
		i -= len(m.OutChain)
		copy(dAtA[i:], m.OutChain)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.OutChain)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.InChain) > 0 {
		i -= len(m.InChain)
		copy(dAtA[i:], m.InChain)
		i = encodeVarintTxOut(dAtA, i, uint64(len(m.InChain)))
		i--
		dAtA[i] = 0x1a
	}

	// No signer.

	if m.TxType != 0 {
		i = encodeVarintTxOut(dAtA, i, uint64(m.TxType))
		i--
		dAtA[i] = 0x8
	}

	return dAtA[:len(dAtA)-i]
}

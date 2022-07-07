package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TxOutStatus alias for tx out status
type TxOutStatus string

const (
	// TxOutStatusSigning txOut is in singing progress
	TxOutStatusSigning TxOutStatus = "signing"
)

var _ sdk.Msg = &TxOutWithSigner{}

func NewMsgTxOutWithSigner(signer string, txType TxOutType, inChains []string, inHashes []string,
	outChain string, outHash string, outBytes []byte, contractHash string) *TxOutWithSigner {
	return &TxOutWithSigner{
		Signer: signer,
		Data: &TxOut{
			TxType:       txType,
			OutChain:     outChain,
			OutHash:      outHash,
			InChains:     inChains,
			InHashes:     inHashes,
			OutBytes:     outBytes,
			ContractHash: contractHash,
		},
	}
}

// Route ...
func (msg *TxOutWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutWithSigner) Type() string {
	return MsgTypeTxOutWithSigner
}

// GetSigners ...
func (msg *TxOutWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

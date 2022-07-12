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

var _ sdk.Msg = &TxOutMsg{}

func NewTxOutMsg(signer string, txType TxOutType, inHashes []string,
	outChain string, outHash string, outBytes []byte, contractHash string) *TxOutMsg {
	return &TxOutMsg{
		Signer: signer,
		Data: &TxOut{
			TxType:       txType,
			OutChain:     outChain,
			OutHash:      outHash,
			InHashes:     inHashes,
			OutBytes:     outBytes,
			ContractHash: contractHash,
		},
	}
}

// Route ...
func (msg *TxOutMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutMsg) Type() string {
	return MsgTypeTxOutMsg
}

// GetSigners ...
func (msg *TxOutMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

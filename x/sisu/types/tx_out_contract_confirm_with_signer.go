package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutContractConfirmWithSigner{}

func NewTxOutContractConfirmWithSigner(signer string, data *TxOutContractConfirm) *TxOutContractConfirmWithSigner {
	return &TxOutContractConfirmWithSigner{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxOutContractConfirmWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutContractConfirmWithSigner) Type() string {
	return MsgTypeContractConfirmWithSigner
}

// GetSigners ...
func (msg *TxOutContractConfirmWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutContractConfirmWithSigner) GetSender() sdk.AccAddress {
	return msg.GetSigners()[0]
}

func (msg *TxOutContractConfirmWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutContractConfirmWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutContractConfirmWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

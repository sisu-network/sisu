package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutResultMsg{}

func NewTxOutResultMsg(signer string, data *TxOutResult) *TxOutResultMsg {
	return &TxOutResultMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxOutResultMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutResultMsg) Type() string {
	return MsgTypeTxOutResult
}

// GetSigners ...
func (msg *TxOutResultMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutResultMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutResultMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutResultMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

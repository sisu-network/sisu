package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutConfirmMsg{}

func NewTxOutConfirmMsg(signer string, data *TxOutConfirm) *TxOutConfirmMsg {
	return &TxOutConfirmMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxOutConfirmMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutConfirmMsg) Type() string {
	return MsgTypeTxOutConfirmMsg
}

// GetSigners ...
func (msg *TxOutConfirmMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutConfirmMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutConfirmMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutConfirmMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

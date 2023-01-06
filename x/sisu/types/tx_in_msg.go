package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewTxInMsg(signer string, data *TxInOld) *TxInMsgOld {
	return &TxInMsgOld{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxInMsgOld) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxInMsgOld) Type() string {
	return MsgTxIn
}

// GetSigners ...
func (msg *TxInMsgOld) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxInMsgOld) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxInMsgOld) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxInMsgOld) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

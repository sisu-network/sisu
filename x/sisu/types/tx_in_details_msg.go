package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewTxInDetailsMsg(signer string, data *TxInDetails) *TxInDetailsMsg {
	return &TxInDetailsMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxInDetailsMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxInDetailsMsg) Type() string {
	return MsgTxInDetails
}

// GetSigners ...
func (msg *TxInDetailsMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxInDetailsMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxInDetailsMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxInDetailsMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

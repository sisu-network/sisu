package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TransferOutsMsg{}

func NewTransferOutsMsg(signer string, data *TransferOuts) *TransferOutsMsg {
	return &TransferOutsMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TransferOutsMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TransferOutsMsg) Type() string {
	return MsgTransferOutsMsg
}

// GetSigners ...
func (msg *TransferOutsMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TransferOutsMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TransferOutsMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TransferOutsMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

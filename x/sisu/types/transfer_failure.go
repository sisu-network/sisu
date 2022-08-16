package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TransferFailureMsg{}

func NewTransferFailureMsg(signer string, data *TransferFailure) *TransferFailureMsg {
	return &TransferFailureMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TransferFailureMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TransferFailureMsg) Type() string {
	return MsgTypeTransferBatchMsg
}

// GetSigners ...
func (msg *TransferFailureMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TransferFailureMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TransferFailureMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TransferFailureMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

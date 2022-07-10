package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TransferBatchMsg{}

func NewTransferBatchMsg(signer string, data *TransferBatch) *TransferBatchMsg {
	return &TransferBatchMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TransferBatchMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TransferBatchMsg) Type() string {
	return MsgTypeTransferBatchMsg
}

// GetSigners ...
func (msg *TransferBatchMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TransferBatchMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TransferBatchMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TransferBatchMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

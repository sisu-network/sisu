package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TransferRequestsMsg{}

func NewBlockTransfersMsg(signer string, data *TransferRequests) *TransferRequestsMsg {
	return &TransferRequestsMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TransferRequestsMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TransferRequestsMsg) Type() string {
	return MsgTypeTransferRequestsMsg
}

// GetSigners ...
func (msg *TransferRequestsMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TransferRequestsMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TransferRequestsMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TransferRequestsMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

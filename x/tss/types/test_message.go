package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TestMessage{}

func NewTestMessage(signer string, index int, number int) *TestMessage {
	return &TestMessage{
		Signer: signer,
		Index:  int32(index),
		Number: int32(number),
	}
}

// Route ...
func (msg *TestMessage) Route() string {
	return RouterKey
}

// Type ...
func (msg *TestMessage) Type() string {
	return MsgTypeTestMessage
}

// GetSigners ...
func (msg *TestMessage) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TestMessage) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TestMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TestMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

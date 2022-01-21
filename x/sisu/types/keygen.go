package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenWithSigner{}

func NewMsgKeygenWithSigner(signer string, keyType string, index int) *KeygenWithSigner {
	data := &Keygen{
		KeyType: keyType,
		Index:   int32(index),
	}

	return &KeygenWithSigner{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *KeygenWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenWithSigner) Type() string {
	return MsgTypeKeygenWithSigner
}

// GetSigners ...
func (msg *KeygenWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

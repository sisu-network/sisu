package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenResultWithSigner{}

func NewKeygenResultWithSigner(signer string, keyType string, index int, result KeygenResult_Result, pubKeyBytes []byte, address string) *KeygenResultWithSigner {
	return &KeygenResultWithSigner{
		Signer: signer,
		Keygen: &Keygen{
			KeyType:     keyType,
			Index:       int32(index),
			PubKeyBytes: pubKeyBytes,
			Address:     address,
		},
		Data: &KeygenResult{
			From:   signer,
			Result: result,
		},
	}
}

// Route ...
func (msg *KeygenResultWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenResultWithSigner) Type() string {
	return MsgTypeKeygenResultWithSigner
}

// GetSigners ...
func (msg *KeygenResultWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

// GetSender ...
func (msg *KeygenResultWithSigner) GetSender() sdk.AccAddress {
	return msg.GetSigners()[0]
}

func (msg *KeygenResultWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenResultWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenResultWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeysignResultMsg{}

func NewKeysignResult(signer, id string, success bool, signature []byte) *KeysignResultMsg {
	return &KeysignResultMsg{
		Signer: signer,
		Data: &KeysignResult{
			TxOutId:   id,
			Success:   success,
			Signature: signature,
		},
	}
}

// Route ...
func (msg *KeysignResultMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeysignResultMsg) Type() string {
	return MsgTypeKeysignResult
}

// GetSigners ...
func (msg *KeysignResultMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeysignResultMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeysignResultMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeysignResultMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

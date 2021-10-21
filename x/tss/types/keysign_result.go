package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeysignResult{}

func NewKeysignResult(signer, outChain, outHash string, success bool, signature []byte) *KeysignResult {
	return &KeysignResult{
		Signer:    signer,
		OutChain:  outChain,
		OutHash:   outHash,
		Success:   success,
		Signature: signature,
	}
}

// Route ...
func (msg *KeysignResult) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeysignResult) Type() string {
	return MSG_TYPE_KEYGEN_RESULT
}

// GetSigners ...
func (msg *KeysignResult) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeysignResult) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeysignResult) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeysignResult) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

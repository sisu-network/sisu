package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenResultWithSigner{}

func NewKeygenResultWithSigner(signer string, keyType string, result KeygenResult_Result, pubKeyBytes []byte, address string) *KeygenResultWithSigner {
	keygenResult := &KeygenResult{
		Keygen: &Keygen{
			KeyType:     keyType,
			PubKeyBytes: pubKeyBytes,
			Address:     address,
		},

		Result: result,
	}

	return &KeygenResultWithSigner{
		Signer: signer,
		Data:   keygenResult,
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

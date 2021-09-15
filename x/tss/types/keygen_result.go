package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenResult{}

func NewKeygenResult(signer string, chainSymbol string, result KeygenResult_Result, pubKeyBytes []byte, address string) *KeygenResult {
	return &KeygenResult{
		Signer:      signer,
		ChainSymbol: chainSymbol,
		Result:      result,
		PubKeyBytes: pubKeyBytes,
		Address:     address,
	}
}

// Route ...
func (msg *KeygenResult) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenResult) Type() string {
	return MSG_TYPE_KEYGEN_RESULT
}

// GetSigners ...
func (msg *KeygenResult) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenResult) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenResult) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenResult) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

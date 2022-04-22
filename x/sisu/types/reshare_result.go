package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ReshareResultWithSigner{}

func NewReshareResultWithSigner(signer string, newValidatorSet [][]byte, result ReshareData_Result) *ReshareResultWithSigner {
	return &ReshareResultWithSigner{
		Signer: signer,
		Data: &ReshareData{
			NewValidatorSetPubKeyBytes: newValidatorSet,
			Result:                     result,
		},
	}
}

// Route ...
func (msg *ReshareResultWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *ReshareResultWithSigner) Type() string {
	return MsgTypeKeygenResultWithSigner
}

// GetSigners ...
func (msg *ReshareResultWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ReshareResultWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ReshareResultWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ReshareResultWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

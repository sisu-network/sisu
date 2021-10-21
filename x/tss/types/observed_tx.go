package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ObservedTx{}

func NewObservedTxs(signer string, tx *ObservedTx) *ObservedTx {
	return &ObservedTx{
		Signer: signer,
	}
}

// Route ...
func (msg *ObservedTx) Route() string {
	return RouterKey
}

// Type ...
func (msg *ObservedTx) Type() string {
	return MSG_TYPE_OBSERVED_TXS
}

// GetSigners ...
func (msg *ObservedTx) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ObservedTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ObservedTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ObservedTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

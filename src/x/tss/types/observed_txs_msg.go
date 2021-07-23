package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ObservedTxs{}

func NewObservedTxs(signer string, Txs []*ObservedTx) *ObservedTxs {
	return &ObservedTxs{
		Signer: signer,
		Txs:    Txs,
	}
}

// Route ...
func (msg *ObservedTxs) Route() string {
	return RouterKey
}

// Type ...
func (msg *ObservedTxs) Type() string {
	return MSG_TYPE_OBSERVED_TXS
}

// GetSigners ...
func (msg *ObservedTxs) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ObservedTxs) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ObservedTxs) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ObservedTxs) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

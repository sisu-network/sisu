package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &UpdateTokenPrice{}

func NewUpdateTokenPrice(signer string, prices []*TokenPrice) *UpdateTokenPrice {
	return &UpdateTokenPrice{
		Signer:      signer,
		TokenPrices: prices,
	}
}

// Route ...
func (msg *UpdateTokenPrice) Route() string {
	return RouterKey
}

// Type ...
func (msg *UpdateTokenPrice) Type() string {
	return MsgTypeUpdateTokenPrice
}

// GetSigners ...
func (msg *UpdateTokenPrice) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *UpdateTokenPrice) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *UpdateTokenPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *UpdateTokenPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

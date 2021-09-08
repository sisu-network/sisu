package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &EthTx{}

func NewMsgEthTx(author string, data []byte) *EthTx {
	return &EthTx{
		Author: author,
		Data:   data,
	}
}

// Route ...
func (msg *EthTx) Route() string {
	return RouterKey
}

// Type ...
func (msg *EthTx) Type() string {
	return "EthTx"
}

// GetSigners ...
func (msg *EthTx) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *EthTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *EthTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *EthTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

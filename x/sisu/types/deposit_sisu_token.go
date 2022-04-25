package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &DepositSisuTokenMsg{}

func NewDepositSisuTokenMsg(signer string, amount int64) *DepositSisuTokenMsg {
	return &DepositSisuTokenMsg{
		Signer: signer,
		Data:   &DepositSisuData{Amount: amount},
	}
}

// Route ...
func (msg *DepositSisuTokenMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *DepositSisuTokenMsg) Type() string {
	return MsgTypeDepositSisuTokenWithSigner
}

// GetSigners ...
func (msg *DepositSisuTokenMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *DepositSisuTokenMsg) GetSender() sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return author
}

func (msg *DepositSisuTokenMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *DepositSisuTokenMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *DepositSisuTokenMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

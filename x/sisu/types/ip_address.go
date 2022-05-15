package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &SetDheartIpAddressMsg{}

func NewSetDheartIPAddressMsg(signer string, ip string) *SetDheartIpAddressMsg {
	return &SetDheartIpAddressMsg{
		Signer: signer,
		Data:   &SetDheartIPAddressData{IPAddress: ip},
	}
}

// Route ...
func (msg *SetDheartIpAddressMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *SetDheartIpAddressMsg) Type() string {
	return MsgTypeSetDheartIPAddress
}

// GetSigners ...
func (msg *SetDheartIpAddressMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *SetDheartIpAddressMsg) GetSender() sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return author
}

func (msg *SetDheartIpAddressMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *SetDheartIpAddressMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *SetDheartIpAddressMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

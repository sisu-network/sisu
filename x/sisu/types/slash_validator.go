package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &SlashValidatorMsg{}

func NewSlashValidatorMsg(signer string, nodeAddress string, slashPoint int64, index int32) *SlashValidatorMsg {
	return &SlashValidatorMsg{
		Signer: signer,
		Data:   &SlashValidatorData{NodeAddress: nodeAddress, SlashPoint: slashPoint, Index: index},
	}
}

// Route ...
func (msg *SlashValidatorMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *SlashValidatorMsg) Type() string {
	return MsgTypeSlashValidator
}

// GetSigners ...
func (msg *SlashValidatorMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *SlashValidatorMsg) GetSender() sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return author
}

func (msg *SlashValidatorMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *SlashValidatorMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *SlashValidatorMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

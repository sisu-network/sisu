package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &PauseContractMsg{}

func NewPauseContractMsg(signer string, chain string, hash string, index int) *PauseContractMsg {
	return &PauseContractMsg{
		Signer: signer,
		Data: &PauseContract{
			Chain: chain,
			Hash:  hash,
			Index: int32(index),
		},
	}
}

// Route ...
func (msg *PauseContractMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *PauseContractMsg) Type() string {
	return MsgTypePauseContract
}

// GetSigners ...
func (msg *PauseContractMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *PauseContractMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *PauseContractMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *PauseContractMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

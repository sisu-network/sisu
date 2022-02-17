package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ResumeContractMsg{}

func NewResumeContractMsg(signer string, chain string, hash string, index int) *ResumeContractMsg {
	return &ResumeContractMsg{
		Signer: signer,
		Data: &ResumeContract{
			Chain: chain,
			Hash:  hash,
			Index: int32(index),
		},
	}
}

// Route ...
func (msg *ResumeContractMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *ResumeContractMsg) Type() string {
	return MsgTypeResumeContract
}

// GetSigners ...
func (msg *ResumeContractMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ResumeContractMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ResumeContractMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ResumeContractMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

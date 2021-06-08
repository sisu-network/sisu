package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenProposal{}

func NewMsgKeygenProposal(signer string, chainSymbol string) *KeygenProposal {
	return &KeygenProposal{
		Signer:      signer,
		ChainSymbol: chainSymbol,
	}
}

// Route ...
func (msg *KeygenProposal) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenProposal) Type() string {
	return "KeygenProposal"
}

// GetSigners ...
func (msg *KeygenProposal) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenProposal) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ChangeOwnershipContractMsg{}

func NewChangeOwnershipMsg(signer, chain, hash, newOwner string) *ChangeOwnershipContractMsg {
	return &ChangeOwnershipContractMsg{
		Signer: signer,
		Data: &ChangeOwnership{
			Chain:    chain,
			Hash:     hash,
			NewOwner: newOwner,
		},
	}
}

// Route ...
func (msg *ChangeOwnershipContractMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *ChangeOwnershipContractMsg) Type() string {
	return MsgTypeContractChangeOwnership
}

// GetSigners ...
func (msg *ChangeOwnershipContractMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ChangeOwnershipContractMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ChangeOwnershipContractMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ChangeOwnershipContractMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

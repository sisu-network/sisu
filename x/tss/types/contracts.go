package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ContractsWithSigner{}

func NewContractsWithSigner(signer string, contracts []*Contract) *ContractsWithSigner {
	data := &Contracts{
		Contracts: contracts,
	}

	return &ContractsWithSigner{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *ContractsWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *ContractsWithSigner) Type() string {
	return MsgTypeContractsWithSigner
}

// GetSigners ...
func (msg *ContractsWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ContractsWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ContractsWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ContractsWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

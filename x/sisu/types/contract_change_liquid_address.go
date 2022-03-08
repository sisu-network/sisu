package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ChangeSetPoolAddressMsg{}

func NewChangePoolAddressMsg(signer, chain, hash, newLiquidPoolAddress string, index int32) *ChangeSetPoolAddressMsg {
	return &ChangeSetPoolAddressMsg{
		Signer: signer,
		Data: &ChangeLiquidAddress{
			Chain:            chain,
			Hash:             hash,
			NewLiquidAddress: newLiquidPoolAddress,
			Index:            index,
		},
	}
}

// Route ...
func (msg *ChangeSetPoolAddressMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *ChangeSetPoolAddressMsg) Type() string {
	return MsgTypeContractChangeLiquidityAddress
}

// GetSigners ...
func (msg *ChangeSetPoolAddressMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ChangeSetPoolAddressMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ChangeSetPoolAddressMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ChangeSetPoolAddressMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

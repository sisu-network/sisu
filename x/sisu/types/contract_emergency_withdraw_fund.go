package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &EmergencyWithdrawFundMsg{}

func NewEmergencyWithdrawFundMsg(signer, chain, hash string, tokens []string, newOwner string, index int32) *EmergencyWithdrawFundMsg {
	return &EmergencyWithdrawFundMsg{
		Signer: signer,
		Data: &WithdrawFund{
			Chain:          chain,
			Hash:           hash,
			TokenAddresses: tokens,
			NewOwner:       newOwner,
			Index:          index,
		},
	}
}

// Route ...
func (msg *EmergencyWithdrawFundMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *EmergencyWithdrawFundMsg) Type() string {
	return MsgTypeContractEmergencyWithdrawFund
}

// GetSigners ...
func (msg *EmergencyWithdrawFundMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *EmergencyWithdrawFundMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *EmergencyWithdrawFundMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *EmergencyWithdrawFundMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

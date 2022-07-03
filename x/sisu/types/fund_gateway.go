package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &FundGatewayMsg{}

func NewFundGatewayMsg(signer string, data *FundGateway) *FundGatewayMsg {
	return &FundGatewayMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *FundGatewayMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *FundGatewayMsg) Type() string {
	return MsgTypeFundGatewayMsg
}

// GetSigners ...
func (msg *FundGatewayMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *FundGatewayMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *FundGatewayMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *FundGatewayMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

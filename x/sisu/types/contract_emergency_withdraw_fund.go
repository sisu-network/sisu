package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &LiquidityWithdrawFundMsg{}

func NewLiquidityWithdrawFundMsg(signer, chain, hash string, tokens []string, newOwner string, index int32) *LiquidityWithdrawFundMsg {
	return &LiquidityWithdrawFundMsg{
		Signer: signer,
		Data: &LiquidityWithdrawFund{
			Chain:          chain,
			Hash:           hash,
			TokenAddresses: tokens,
			NewOwner:       newOwner,
			Index:          index,
		},
	}
}

// Route ...
func (msg *LiquidityWithdrawFundMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *LiquidityWithdrawFundMsg) Type() string {
	return MsgTypeContractLiquidityWithdrawFund
}

// GetSigners ...
func (msg *LiquidityWithdrawFundMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *LiquidityWithdrawFundMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *LiquidityWithdrawFundMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *LiquidityWithdrawFundMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

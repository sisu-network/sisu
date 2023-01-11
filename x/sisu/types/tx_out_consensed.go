package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutConsensedMsg{}

func NewTxOutConsensedMsg(signer string, data *TxOutConsensed) *TxOutConsensedMsg {
	return &TxOutConsensedMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxOutConsensedMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutConsensedMsg) Type() string {
	return MsgTxOutConsensed
}

// GetSigners ...
func (msg *TxOutConsensedMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutConsensedMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutConsensedMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutConsensedMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

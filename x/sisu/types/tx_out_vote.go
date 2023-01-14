package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutVoteMsg{}

func NewTxOutVoteMsg(signer string, data *TxOutVote) *TxOutVoteMsg {
	return &TxOutVoteMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *TxOutVoteMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutVoteMsg) Type() string {
	return MsgTxOutVote
}

// GetSigners ...
func (msg *TxOutVoteMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutVoteMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutVoteMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutVoteMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenProposalVote{}

func NewMsgKeygenProposalVote(signer string, chainSymbol string, vote KeygenProposalVote_Vote) *KeygenProposalVote {
	return &KeygenProposalVote{
		Signer:      signer,
		ChainSymbol: chainSymbol,
		Vote:        vote,
	}
}

// Route ...
func (msg *KeygenProposalVote) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenProposalVote) Type() string {
	return MSG_TYPE_KEYGEN_PROPOSAL_VOTE
}

// GetSigners ...
func (msg *KeygenProposalVote) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenProposalVote) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenProposalVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenProposalVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

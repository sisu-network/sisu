package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &KeygenProposalWithSigner{}

func NewMsgKeygenProposalWithSigner(signer string, keyType string, id string, craetedBlock int64) *KeygenProposalWithSigner {
	data := &KeygenProposal{
		KeyType:      keyType,
		Id:           id,
		CreatedBlock: craetedBlock,
	}

	return &KeygenProposalWithSigner{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *KeygenProposalWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *KeygenProposalWithSigner) Type() string {
	return MsgTypeKeygenProposalWithSigner
}

// GetSigners ...
func (msg *KeygenProposalWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *KeygenProposalWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *KeygenProposalWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *KeygenProposalWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

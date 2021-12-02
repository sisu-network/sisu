package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &PresignResult{}

func NewPresignResult(signer, chain string, success bool, pubkeyBytes []byte, address string, culprits []*PartyID) *PresignResult {
	return &PresignResult{
		Signer:      signer,
		Chain:       chain,
		Success:     success,
		PubkeyBytes: pubkeyBytes,
		Address:     address,
		Culprits:    culprits,
	}
}

// Route ...
func (msg *PresignResult) Route() string {
	return RouterKey
}

// Type ...
func (msg *PresignResult) Type() string {
	return MsgTypePresignResult
}

// GetSigners ...
func (msg *PresignResult) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *PresignResult) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *PresignResult) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *PresignResult) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

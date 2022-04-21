package types

import (
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/x/sisu/helper"
)

var _ sdk.Msg = &ChangeValidatorSetMsg{}

func NewChangeValidatorSetMsg(signer string, oldValidatorSet, newValidatorSet [][]byte, index int32) *ChangeValidatorSetMsg {
	return &ChangeValidatorSetMsg{
		Signer: signer,
		Data: &ChangeValidatorSetData{
			OldValidatorSet: oldValidatorSet,
			NewValidatorSet: newValidatorSet,
			Index:           index,
		},
	}
}

// Route ...
func (msg *ChangeValidatorSetMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *ChangeValidatorSetMsg) Type() string {
	return MsgTypeChangeValidatorSet
}

// GetSigners ...
func (msg *ChangeValidatorSetMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ChangeValidatorSetMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ChangeValidatorSetMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ChangeValidatorSetMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *ChangeValidatorSetMsg) GetOldAndNewValidatorSet() ([]types.PubKey, []types.PubKey) {
	oldValSet := make([]types.PubKey, 0, len(msg.Data.OldValidatorSet))
	newValSet := make([]types.PubKey, 0, len(msg.Data.NewValidatorSet))

	for _, val := range msg.Data.OldValidatorSet {
		oldValSet = append(oldValSet, helper.BytesToValPubKey(val))
	}

	for _, val := range msg.Data.NewValidatorSet {
		newValSet = append(newValSet, helper.BytesToValPubKey(val))
	}

	return oldValSet, newValSet
}

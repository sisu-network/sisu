package types

import (
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
)

var _ sdk.Msg = &ChangeValidatorSetMsg{}

func NewChangeValidatorSetMsg(signer string, oldValidatorSet, newValidatorSet [][]byte, index int) *ChangeValidatorSetMsg {
	return &ChangeValidatorSetMsg{
		Signer: signer,
		Data: &ChangeValidatorSetData{
			OldValidatorSet: oldValidatorSet,
			NewValidatorSet: newValidatorSet,
			Index:           int32(index),
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

func (msg *ChangeValidatorSetMsg) GetOldAndNewValidatorSet() ([]types.PubKey, []types.PubKey, error) {
	oldValSet := make([]types.PubKey, 0, len(msg.Data.OldValidatorSet))
	newValSet := make([]types.PubKey, 0, len(msg.Data.NewValidatorSet))

	for _, val := range msg.Data.OldValidatorSet {
		pk, err := utils.GetCosmosPubKey("ed25519", val)
		if err != nil {
			log.Error("error when get cosmos pubkey: ", err)
			return nil, nil, err
		}

		oldValSet = append(oldValSet, pk)
	}

	for _, val := range msg.Data.NewValidatorSet {
		pk, err := utils.GetCosmosPubKey("ed25519", val)
		if err != nil {
			log.Error("error when get cosmos pubkey: ", err)
			return nil, nil, err
		}

		newValSet = append(newValSet, pk)
	}

	return oldValSet, newValSet, nil
}

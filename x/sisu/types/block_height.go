package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &BlockHeightMsg{}

func ConvertBlockHeightsMapToArray(m map[string]*BlockHeight) ([]string, []*BlockHeight) {
	signers := make([]string, 0, len(m))
	blockHeights := make([]*BlockHeight, 0, len(m))

	for signer, blockHeight := range m {
		signers = append(signers, signer)
		blockHeights = append(blockHeights, blockHeight)
	}

	return signers, blockHeights
}

func NewBlockHeightMsg(signer string, data *BlockHeight) *BlockHeightMsg {
	return &BlockHeightMsg{
		Signer: signer,
		Data:   data,
	}
}

// Route ...
func (msg *BlockHeightMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *BlockHeightMsg) Type() string {
	return MsgTypeBlockHeight
}

// GetSigners ...
func (msg *BlockHeightMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *BlockHeightMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *BlockHeightMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *BlockHeightMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

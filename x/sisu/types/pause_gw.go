package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

var _ sdk.Msg = &MsgPauseGw{}

func NewMsgPauseGw(signer sdk.AccAddress) *MsgPauseGw {
	return &MsgPauseGw{
		Signer: signer.String(),
	}
}

// Route ...
func (msg *MsgPauseGw) Route() string {
	return RouterKey
}

// Type ...
func (msg *MsgPauseGw) Type() string {
	return MsgTypePauseGwWithSigner
}

// GetSigners ...
func (msg *MsgPauseGw) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *MsgPauseGw) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *MsgPauseGw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *MsgPauseGw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (rc *PauseGwRecord) HasSigned(signer string) bool {
	for _, m := range rc.Messages {
		if strings.EqualFold(strings.ToLower(m.Signer), strings.ToLower(signer)) {
			return true
		}
	}

	return false
}

func (rc *PauseGwRecord) ReachConsensus(nbActiveValidator int) bool {
	return len(rc.Messages) >= nbActiveValidator*2/3
}

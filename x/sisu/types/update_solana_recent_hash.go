package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	libchain "github.com/sisu-network/lib/chain"
)

func NewUpdateSolanaRecentHashMsg(signer, chain, solanaHash string, solanaHeight int64) *UpdateSolanaRecentHashMsg {
	return &UpdateSolanaRecentHashMsg{
		Signer: signer,
		Data: &UpdateSolanaRecentHash{
			Chain:  chain,
			Hash:   solanaHash,
			Height: solanaHeight,
		},
	}
}

// Route ...
func (msg *UpdateSolanaRecentHashMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *UpdateSolanaRecentHashMsg) Type() string {
	return MsgUpdateSolanaRecentHash
}

// GetSigners ...
func (msg *UpdateSolanaRecentHashMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *UpdateSolanaRecentHashMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *UpdateSolanaRecentHashMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *UpdateSolanaRecentHashMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if libchain.IsSolanaChain(msg.Data.Chain) {
		// TODO: Make sure that the block hash is valid
	}

	return nil
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewAdjustEthNonceMsg(signer string, chain string, nonce int64, index int64) *AdjustEthNonceMsg {
	return &AdjustEthNonceMsg{
		Signer: signer,
		Data: &AdjustEthNonce{
			Chain:    chain,
			Nonce:    nonce,
			MsgIndex: index,
		},
	}
}

// Route ...
func (msg *AdjustEthNonceMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *AdjustEthNonceMsg) Type() string {
	return MsgAdjustEthNonce
}

// GetSigners ...
func (msg *AdjustEthNonceMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *AdjustEthNonceMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *AdjustEthNonceMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *AdjustEthNonceMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

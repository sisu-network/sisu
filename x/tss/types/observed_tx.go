package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &ObservedTxWithSigner{}

func NewObservedTxs(signer string, chain string, txHash string, blockHeight int64, serialized []byte) *ObservedTxWithSigner {
	return &ObservedTxWithSigner{
		Signer: signer,
		Data: &ObservedTx{
			Chain:       chain,
			TxHash:      txHash,
			BlockHeight: blockHeight,
			Serialized:  serialized,
		},
	}
}

// Route ...
func (msg *ObservedTxWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *ObservedTxWithSigner) Type() string {
	return MsgTypeObservedTxWithSigner
}

// GetSigners ...
func (msg *ObservedTxWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *ObservedTxWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *ObservedTxWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *ObservedTxWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

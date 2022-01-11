package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &TxOutWithSigner{}

func NewTxOutConfirmWithSigner(signer string, txType TxOutType, outChain string, hash string, blockHeight int64, contractAddress string) *TxOutConfirmWithSigner {
	return &TxOutConfirmWithSigner{
		Signer: signer,
		Data: &TxOutConfirm{
			TxType:          txType,
			OutChain:        outChain,
			OutHash:         hash,
			ContractAddress: contractAddress,
		},
	}
}

// Route ...
func (msg *TxOutConfirmWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutConfirmWithSigner) Type() string {
	return MsgTypeTxOutConfirmationWithSigner
}

// GetSigners ...
func (msg *TxOutConfirmWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutConfirmWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutConfirmWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutConfirmWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

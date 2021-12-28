package types

import (
	fmt "fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/utils"
)

var _ sdk.Msg = &TxInWithSigner{}

func NewTxInWithSigner(signer string, chain string, txHash string, blockHeight int64, serialized []byte) *TxInWithSigner {
	return &TxInWithSigner{
		Signer: signer,
		Data: &TxIn{
			Chain:       chain,
			TxHash:      txHash,
			BlockHeight: blockHeight,
			Serialized:  serialized,
		},
	}
}

// Route ...
func (msg *TxInWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxInWithSigner) Type() string {
	return MsgTypeTxInWithSigner
}

// GetSigners ...
func (msg *TxInWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxInWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxInWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxInWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Make sure the hash matches.
	txIn := msg.Data
	if txIn.TxHash != utils.GetTxInHash(txIn.BlockHeight, txIn.Chain, txIn.Serialized) {
		return fmt.Errorf("Hash does not match")
	}

	return nil
}

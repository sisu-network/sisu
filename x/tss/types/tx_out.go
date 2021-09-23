package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/utils"
)

var _ sdk.Msg = &TxOut{}

func NewMsgTxOut(txType TxOut_Type, signer string, inBlockHeight int64, inChain string, inHash string, outChain string, outBytes []byte) *TxOut {
	return &TxOut{
		TxType:        txType,
		Signer:        signer,
		InBlockHeight: inBlockHeight,
		InChain:       inChain,
		OutChain:      outChain,
		InHash:        inHash,
		OutBytes:      outBytes,
	}
}

// Route ...
func (msg *TxOut) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOut) Type() string {
	return MSG_TYPE_TX_OUT
}

// GetSigners ...
func (msg *TxOut) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOut) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

func (msg *TxOut) GetHash() string {
	return utils.KeccakHash32(msg.OutChain + string(msg.OutBytes))
}

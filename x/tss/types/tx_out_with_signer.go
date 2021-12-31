package types

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
)

// TxOutStatus alias for tx out status
type TxOutStatus string

const (
	TxOutStatusConfirmed TxOutStatus = "confirmed"
	// TxOutStatusPreBroadcast after produce from txIn, the txOut is ready to broadcast to Cosmos chain
	TxOutStatusPreBroadcast TxOutStatus = "pre_broadcast"
	// TxOutStatusBroadcasted TxOut has broadcasted to Cosmos chain successfully
	TxOutStatusBroadcasted TxOutStatus = "broadcasted"

	TxOutStatusDeployingToBlockchain TxOutStatus = "deploying_to_chain"
	TxOutStatusDeployedToBlockchain  TxOutStatus = "deployed_to_chain"

	// TxOutStatusPreSigning txOut is ready to be signed
	TxOutStatusPreSigning TxOutStatus = "pre_signing"
	// TxOutStatusSigning txOut is in singing progress
	TxOutStatusSigning TxOutStatus = "signing"
	// TxOutStatusSigned txOut is signed successfully
	TxOutStatusSigned TxOutStatus = "signed"
	// TxOutStatusSignFailed signing progress is failed
	TxOutStatusSignFailed TxOutStatus = "sign_failed"
)

var _ sdk.Msg = &TxOutWithSigner{}

func NewMsgTxOutWithSigner(signer string, txType TxOutType, inBlockHeight int64,
	inChain string, inHash string, outChain string, outHash string, outBytes []byte, contractHash string) *TxOutWithSigner {
	return &TxOutWithSigner{
		Signer: signer,
		Data: &TxOut{
			TxType:        txType,
			OutChain:      outChain,
			OutHash:       outHash,
			InBlockHeight: inBlockHeight,
			InChain:       inChain,
			InHash:        inHash,
			OutBytes:      outBytes,
			ContractHash:  contractHash,
		},
	}
}

// Route ...
func (msg *TxOutWithSigner) Route() string {
	return RouterKey
}

// Type ...
func (msg *TxOutWithSigner) Type() string {
	return MsgTypeTxOutWithSigner
}

// GetSigners ...
func (msg *TxOutWithSigner) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TxOutWithSigner) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TxOutWithSigner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TxOutWithSigner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

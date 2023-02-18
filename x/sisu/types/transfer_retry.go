package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewTransferRetryMsg(signer string, transferId string, nonce int64) *TransferRetryMsg {
	return &TransferRetryMsg{
		Signer: signer,
		Data: &TransferRetry{
			TransferId: transferId,
			Nonce:      nonce,
		},
	}
}

// Route ...
func (msg *TransferRetryMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *TransferRetryMsg) Type() string {
	return MsgTxIn
}

// GetSigners ...
func (msg *TransferRetryMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *TransferRetryMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *TransferRetryMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *TransferRetryMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

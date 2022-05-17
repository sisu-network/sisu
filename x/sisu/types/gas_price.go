package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &GasPriceMsg{}

func NewGasPriceMsg(signer, chain string, blockHeight, gasPrice int64) *GasPriceMsg {
	return &GasPriceMsg{
		Chain:       chain,
		BlockHeight: blockHeight,
		GasPrice:    gasPrice,
		Signer:      signer,
	}
}

// Route ...
func (msg *GasPriceMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *GasPriceMsg) Type() string {
	return MsgTypeGasPriceWithSigner
}

// GetSigners ...
func (msg *GasPriceMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

// GetSender ...
func (msg *GasPriceMsg) GetSender() sdk.AccAddress {
	return msg.GetSigners()[0]
}

func (msg *GasPriceMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *GasPriceMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *GasPriceMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (record *GasPriceRecord) ReachConsensus(nbActiveNodes int) bool {
	// TODO: Check signers is actually validator
	return len(record.Messages) >= nbActiveNodes*2/3
}

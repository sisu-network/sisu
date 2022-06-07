package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &CardanoMintMultiAssetMsg{}

func NewCardanoMintMultiAssetMsg(signer string,
	lovelace int64, assetName string, assetAmount int64, tssPubkey string, index int32) *CardanoMintMultiAssetMsg {
	return &CardanoMintMultiAssetMsg{
		Signer: signer,
		Data: &CardanoMintMultiAssetData{
			Lovelace:    lovelace,
			AssetName:   assetName,
			AssetAmount: assetAmount,
			TssPubkey:   tssPubkey,
			Index:       index,
		},
	}
}

// Route ...
func (msg *CardanoMintMultiAssetMsg) Route() string {
	return RouterKey
}

// Type ...
func (msg *CardanoMintMultiAssetMsg) Type() string {
	return MsgTypeCardanoMintMultiAsset
}

// GetSigners ...
func (msg *CardanoMintMultiAssetMsg) GetSigners() []sdk.AccAddress {
	author, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{author}
}

func (msg *CardanoMintMultiAssetMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSignBytes ...
func (msg *CardanoMintMultiAssetMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *CardanoMintMultiAssetMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

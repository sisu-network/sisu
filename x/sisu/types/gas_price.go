package types

import (
	"fmt"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &GasPriceMsg{}

func NewGasPriceMsg(signer string, chains []string, prices, baseFees, tip []int64) *GasPriceMsg {
	return &GasPriceMsg{
		Signer:    signer,
		Chains:    chains,
		GasPrices: prices,
		BaseFees:  baseFees,
		Tips:      tip,
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

	if msg.Chains == nil {
		return fmt.Errorf("chains array is nil")
	}

	for _, chain := range msg.Chains {
		if !libchain.IsETHBasedChain(chain) {
			return fmt.Errorf("Chain %s is not an ETH based chain", chain)
		}
	}

	if msg.GasPrices == nil && (msg.BaseFees == nil || msg.Tips == nil) {
		return fmt.Errorf("Either gas prices or base fee or tip array is nil")
	}

	l := len(msg.Chains)
	if l != len(msg.GasPrices) && (l != len(msg.BaseFees) || l != len(msg.Tips)) {
		return fmt.Errorf("Chains array does not have the same length with gas pricess, base fee or tip")
	}

	return nil
}

package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TssValidator interface {
	CheckTx(msgs []sdk.Msg) error
}

package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TssValidator interface {
	CheckTx(ctx sdk.Context, msgs []sdk.Msg) error
}

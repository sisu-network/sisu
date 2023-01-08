package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type Bridge interface {
	ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOutMsg, error)
	ParseIncomginTx(ctx sdk.Context, chain string, tx *eyesTypes.Tx) ([]*types.TransferDetails, error)
}

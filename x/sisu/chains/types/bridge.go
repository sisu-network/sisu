package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/sisu/x/sisu/types"
)

type Bridge interface {
	ProcessCommand(ctx sdk.Context, cmd *types.Command) (*types.TxOutMsg, error)
	ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOutMsg, error)
	ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) ([]*types.TransferDetails, error)
}

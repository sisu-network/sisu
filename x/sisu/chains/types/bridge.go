package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sisu-network/sisu/x/sisu/types"
)

type Bridge interface {
	ProcessCommand(ctx sdk.Context, cmd *types.Command) (*types.TxOut, error)
	ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) (*types.TxOut, error)
	ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) ([]*types.TransferDetails, error)
	ValidateTxOut(ctx sdk.Context, txOut *types.TxOut, transfers []*types.TransferDetails) error
}

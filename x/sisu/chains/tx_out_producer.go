package chains

import (
	"math/big"

	"github.com/sisu-network/sisu/x/sisu/components"

	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	// GetTxOuts returns a list of TxOut message and a list of un-processed transfer out request that
	// needs to be processed next time.
	GetTxOuts(ctx sdk.Context, chain string, transfers []*types.TransferDetails) ([]*types.TxOut, error)
}

type DefaultTxOutputProducer struct {
	signer    string
	keeper    keeper.Keeper
	txTracker components.TxTracker

	// Only use for cardano chain
	bridgeManager BridgeManager
}

type transferInData struct {
	token     ethcommon.Address
	recipient string
	amount    *big.Int
}

func NewTxOutputProducer(appKeys components.AppKeys, keeper keeper.Keeper,
	bridgeManager BridgeManager,
	txTracker components.TxTracker) TxOutputProducer {
	return &DefaultTxOutputProducer{
		signer:        appKeys.GetSignerAddress().String(),
		keeper:        keeper,
		txTracker:     txTracker,
		bridgeManager: bridgeManager,
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string,
	transfers []*types.TransferDetails) ([]*types.TxOut, error) {
	bridge := p.bridgeManager.GetBridge(ctx, chain)
	msgs, err := bridge.ProcessTransfers(ctx, transfers)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

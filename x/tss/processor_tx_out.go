package tss

import (
	"strconv"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"

	"github.com/sisu-network/lib/log"
)

// checkTxOut checks if a TxOut message is valid before it is added into Sisu block.
func (p *Processor) checkTxOut(ctx sdk.Context, msg *types.TxOutWithSigner) error {
	if !p.privateDb.IsTxOutExisted(msg.Data) {
		return ErrCannotFindMessage
	}

	if p.keeper.IsTxOutExisted(ctx, msg.Data) {
		return ErrMessageHasBeenProcessed
	}

	return nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (p *Processor) deliverTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	if p.keeper.IsTxOutExisted(ctx, txOut) {
		// The message has been processed
		return nil, nil
	}

	log.Info("Delivering TxOut")

	// Save this to KVStore
	p.keeper.SaveTxOut(ctx, txOut)
	p.privateDb.SaveTxOut(txOut)

	// If this is a txout deployment,

	return nil, nil
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}

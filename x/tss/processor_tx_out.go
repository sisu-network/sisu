package tss

import (
	"fmt"
	"strconv"

	hTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) createTxOuts(ctx sdk.Context, tx *types.ObservedTx) []*tssTypes.TxOut {
	outMsgs := p.txOutputProducer.GetTxOuts(ctx, p.currentHeight, tx)

	for _, msg := range outMsgs {
		p.storage.AddTxOut(msg)

		go func(m *tssTypes.TxOut) {
			p.txSubmit.SubmitMessage(m)
		}(msg)
	}

	return outMsgs
}

func (p *Processor) CheckTxOut(ctx sdk.Context, msg *types.TxOut) error {
	txOut := p.storage.GetTxOut(msg.GetHash())
	if txOut == nil {
		return fmt.Errorf("txout not found, hash = %s", txOut)
	}

	return nil
}

func (p *Processor) DeliverTxOut(ctx sdk.Context, msg *types.TxOut) ([]byte, error) {
	outHash := msg.GetHash()

	utils.LogVerbose("Delivering TXOUT")

	// 4. Broadcast it to Dheart for processing.
	keysignReq := &hTypes.KeysignRequest{
		Id:             p.getKeysignRequestId(msg.OutChain, ctx.BlockHeight(), outHash),
		OutChain:       msg.OutChain,
		OutBlockHeight: p.currentHeight,
		OutHash:        outHash,
		OutBytes:       msg.OutBytes,
	}

	// TODO: check if this tx has been requested to be signed.
	err := p.dheartClient.KeySign(keysignReq)
	if err != nil {
		utils.LogError("Keysign: err =", err)
		return nil, err
	}

	return nil, nil
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}

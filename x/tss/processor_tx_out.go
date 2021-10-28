package tss

import (
	"strconv"

	hTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) createAndBroadcastTxOuts(ctx sdk.Context, tx *types.ObservedTx) []*tssTypes.TxOut {
	outMsgs, outEntities := p.txOutputProducer.GetTxOuts(ctx, p.currentHeight, tx)

	// Save this to database
	utils.LogVerbose("len(outEntities) = ", len(outEntities))
	if len(outEntities) > 0 {
		p.db.InsertTxOuts(outEntities)
	}

	for _, msg := range outMsgs {
		go func(m *tssTypes.TxOut) {
			p.txSubmit.SubmitMessage(m)
		}(msg)
	}

	return outMsgs
}

func (p *Processor) CheckTxOut(ctx sdk.Context, msg *types.TxOut) error {

	return nil
}

func (p *Processor) DeliverTxOut(ctx sdk.Context, tx *types.TxOut) ([]byte, error) {
	// TODO: check if this tx has been requested to be signed.
	outHash := tx.GetHash()

	utils.LogVerbose("Delivering TXOUT")

	// 4. Broadcast it to Dheart for processing.
	keysignReq := &hTypes.KeysignRequest{
		Id:             p.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), outHash),
		OutChain:       tx.OutChain,
		OutBlockHeight: p.currentHeight,
		OutHash:        outHash,
		OutBytes:       tx.OutBytes,
	}

	pubKeys := p.partyManager.GetActivePartyPubkeys()
	err := p.dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		utils.LogError("Keysign: err =", err)
		return nil, err
	}

	return nil, nil
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}

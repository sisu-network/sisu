package tss

import (
	"bytes"
	"encoding/hex"
	"fmt"

	tTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) CreateTxOuts(ctx sdk.Context, tx *types.ObservedTx) {
	outMsgs := p.txOutputProducer.GetOutputs(ctx, p.currentHeight, tx)

	for _, msg := range outMsgs {
		p.storage.AddPendingTxOut(
			p.currentHeight,
			msg.InChain,
			msg.InHash,
			msg.OutChain,
			msg.OutBytes,
		)

		go func(inHeight int64, m *tssTypes.TxOut) {
			p.txSubmit.SubmitMessage(
				tssTypes.NewMsgTxOut(
					p.appKeys.GetSignerAddress().String(),
					inHeight,
					m.InChain,
					m.InHash,
					m.OutChain,
					m.OutBytes,
				),
			)
		}(p.currentHeight, msg)
	}
}

func (p *Processor) CheckTxOut(ctx sdk.Context, msg *types.TxOut) error {
	txWrapper := p.storage.GetPendingTxOUt(msg.InBlockHeight, msg.InHash)
	if txWrapper == nil {
		utils.LogError("Cannot find txWrapper", msg.InBlockHeight, msg.InHash)
		return fmt.Errorf("Transaction not found")
	}

	if bytes.Compare(txWrapper.OutBytes, msg.OutBytes) != 0 {
		utils.LogError("Txouts do not match.")
		return fmt.Errorf("OutBytes do not match")
	}

	return nil
}

func (p *Processor) DeliverTxOut(ctx sdk.Context, msg *types.TxOut) ([]byte, error) {
	utils.LogVerbose("Delivering TXOUT")

	outHash, err := utils.GetTxHash(msg.OutChain, msg.OutBytes)
	if err != nil {
		utils.LogCritical("Cannot get tx hash for tx with serialized data: ", hex.EncodeToString(msg.OutBytes), "err = ", err)
		return nil, err
	}

	// 4. Broadcast it to Dheart for processing.
	err = p.dheartClient.KeySign(&tTypes.KeysignRequest{
		OutChain:       msg.OutChain,
		OutBlockHeight: p.currentHeight,
		OutHash:        outHash,
		OutBytes:       msg.OutBytes,
	})
	if err != nil {
		utils.LogError("Keysign: err =", err)
		return nil, err
	}

	return nil, nil
}

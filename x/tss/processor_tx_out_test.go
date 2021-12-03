package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"
	"testing"
)

func TestDeliverTxOutEth(t *testing.T) {
	p := &Processor{}	
	ctx := sdk.Context{}
	txOut := types.TxOut{
		TxType:        0,
		Signer:        "",
		InChain:       "",
		OutChain:      "",
		InBlockHeight: 0,
		InHash:        "",
		OutBytes:      nil,
	}

	p.deliverTxOutEth(ctx, &txOut)
}

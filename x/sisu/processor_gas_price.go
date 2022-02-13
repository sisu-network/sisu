package sisu

import (
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *Processor) OnUpdateGasPriceRequest(request *etypes.GasPriceRequest) {
	gasPriceMsg := types.NewGasPriceMsg(p.appKeys.GetSignerAddress().String(), request.Chain, request.Height, request.GasPrice)
	p.txSubmit.SubmitMessageAsync(gasPriceMsg)
}

func (p *Processor) deliverGasPriceMsg(ctx sdk.Context, msg *types.GasPriceMsg) ([]byte, error) {
	log.Debug("Setting gas price ...")
	currentPriceRecord := p.publicDb.GetGasPriceRecord(msg.Chain, msg.BlockHeight)
	if currentPriceRecord != nil {
		for _, m := range currentPriceRecord.Messages {
			if strings.EqualFold(strings.ToLower(m.Signer), strings.ToLower(msg.Signer)) {
				return nil, nil
			}
		}
	}

	p.publicDb.SetGasPrice(msg)
	savedRecord := p.publicDb.GetGasPriceRecord(msg.Chain, msg.BlockHeight)
	totalValidator := len(p.globalData.GetValidatorSet())
	if savedRecord == nil || !savedRecord.ReachConsensus(totalValidator) {
		return nil, nil
	}

	// Only save network gas price if reached consensus
	listGasPrices := make([]int64, 0)
	for _, m := range savedRecord.Messages {
		listGasPrices = append(listGasPrices, m.GasPrice)
	}

	sort.SliceStable(listGasPrices, func(i, j int) bool {
		return listGasPrices[i] < listGasPrices[j]
	})

	median := listGasPrices[len(listGasPrices)/2]

	// Save to db
	chain := p.publicDb.GetChain(msg.Chain)
	if chain == nil {
		chain = new(types.Chain)
	}
	chain.GasPrice = median
	p.publicDb.SaveChain(chain)

	// Save to the world state
	p.worldState.SetChain(chain)

	return nil, nil
}

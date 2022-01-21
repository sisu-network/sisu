package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	"sort"
)

func (p *Processor) OnUpdateGasPriceRequest(request *etypes.GasPriceRequest) {
	gasPriceMsg := types.NewGasPriceMsg(p.appKeys.GetSignerAddress().String(), request.Chain, request.Height, request.GasPrice)
	go func() {
		if err := p.txSubmit.SubmitMessage(gasPriceMsg); err != nil {
			log.Error(err)
			return
		}

		log.Debug("Submit gas price msg successfully")
	}()
}

func (p *Processor) deliverGasPriceMsg(ctx sdk.Context, msg *types.GasPriceMsg) ([]byte, error) {
	log.Debug("Setting gas price ...")
	savedRecord := p.privateDb.SetGasPrice(msg)
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
	p.privateDb.SaveNetworkGasPrice(savedRecord.Chain, median)
	return nil, nil
}

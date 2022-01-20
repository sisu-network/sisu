package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
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
	p.privateDb.SetGasPrice(msg, len(p.globalData.GetValidatorSet()))
	return nil, nil
}

package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	// The number of block interval that we should update all token prices.
	UPDATE_BLOCK_INTERVAL = 600 // About 30 mins for 3s block.
)

// OnUpdateTokenPrice is called when there is a token price update from deyes. Post to the network
// until we reach a consensus about token price. The token price is only used to calculate gas price
// fee and not used for actual swapping calculation.
func (p *Processor) OnUpdateTokenPrice(tokenPrices etypes.TokenPrices) {
	prices := make([]*types.TokenPrice, 0, len(tokenPrices))

	// Convert from deyes type to msg type
	for _, token := range tokenPrices {
		prices = append(prices, &types.TokenPrice{
			Id:    token.Id,
			Price: token.Price,
		})
	}

	msg := types.NewUpdateTokenPrice(p.appKeys.GetSignerAddress().String(), prices)
	p.txSubmit.SubmitMessageAsync(msg)
}

func (p *Processor) deliverUpdateTokenPrice(ctx sdk.Context, msg *types.UpdateTokenPrice) ([]byte, error) {
	p.publicDb.SetTokenPrices(uint64(ctx.BlockHeight()), msg)

	return nil, nil
}

// calculateTokenPrices gets all token prices posted from all validators and calculate the median.
func (p *Processor) calculateTokenPrices(ctx sdk.Context) {
	valArr := p.globalData.GetValidatorSet()
	valSet := make(map[string]bool)

	for _, value := range valArr {
		valSet[value.Address.String()] = true
	}

	// records := p.publicDb.GetAllTokenPricesRecord()
}

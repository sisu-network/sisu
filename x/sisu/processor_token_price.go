package sisu

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const (
	// The number of block interval that we should update all token prices.
	UPDATE_BLOCK_INTERVAL = 600 // About 30 mins for 3s block.
)

// OnUpdateTokenPrice is called when there is a token price update from deyes. Post to the network
// until we reach a consensus about token price. The token price is only used to calculate gas price
// fee and not used for actual swapping calculation.
func (p *Processor) OnUpdateTokenPrice(tokenPrices []*etypes.TokenPrice) {
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
	log.Info("Delivering token price update ....")

	p.publicDb.SetTokenPrices(uint64(ctx.BlockHeight()), msg)

	return nil, nil
}

// calculateTokenPrices gets all token prices posted from all validators and calculate the median.
func (p *Processor) calculateTokenPrices(ctx sdk.Context) {
	curBlock := ctx.BlockHeight()

	// We wait for 5 more blocks after we get prices from deyes so that any record can be posted
	// onto the blockchain.
	if curBlock%UPDATE_BLOCK_INTERVAL != 5 {
		return
	}

	log.Info("Calcuating token prices....")

	// TODO: Fix the signer set.
	records := p.publicDb.GetAllTokenPricesRecord()

	tokenPrices := make(map[string][]float32)
	for _, record := range records {
		for token, pair := range record.Prices {
			// Only calculate token prices that has been updated recently.
			if curBlock-int64(pair.BlockHeight) > UPDATE_BLOCK_INTERVAL {
				continue
			}

			m := tokenPrices[token]
			if m == nil {
				m = make([]float32, 0)
			}

			m = append(m, pair.Price)

			tokenPrices[token] = m
		}
	}

	// Now sort all the array and get the median
	medians := make(map[string]float32)
	for token, list := range tokenPrices {
		if len(list) == 0 {
			log.Error("cannot find price list for token ", token)
			continue
		}

		sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
		median := list[len(list)/2]
		medians[token] = median
	}

	log.Verbose("Calculated prices = ", medians)

	p.publicDb.SetCalculatedTokenPrice(medians)

	// Update the world state
	p.worldState.SetTokenPrices(medians)
}

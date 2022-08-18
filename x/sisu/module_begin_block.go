package sisu

import (
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (a AppModule) beginBlock(ctx sdk.Context, blockHeight int64) {
	// Check keygen proposal
	if blockHeight > 1 {
		// We need to wait till block 2 for multistore of the app to be updated with latest account info
		// for signing.
		a.checkTssKeygen(ctx, blockHeight)
	}

	oldValue := a.globalData.IsCatchingUp()
	a.globalData.UpdateCatchingUp()
	newValue := a.globalData.IsCatchingUp()

	if oldValue && !newValue {
		log.Info("Setting Sisu readiness for dheart.")
		// This node has fully catched up with the blockchain, we need to inform dheart about this.
		a.mc.DheartClient().SetSisuReady(true)
		a.mc.DeyesClient().SetSisuReady(true)
	}

	// Calculate token prices
	a.calculateTokenPrices(ctx)
}

/**
Process for generating a new key:
- Wait for the app to catch up
- If there is no support for a particular chain, creates a proposal to include a chain
- When other nodes receive the proposal, top N validator nodes vote to see if it should accept that.
- After M blocks (M is a constant) since a proposal is sent, count the number of yes vote. If there
are enough validator supporting the new chain, send a message to TSS engine to do keygen.
*/
func (a AppModule) checkTssKeygen(ctx sdk.Context, blockHeight int64) {
	// TODO: We can replace this by sending command from client instead of running at the beginning
	// of each block.
	if a.globalData.IsCatchingUp() || ctx.BlockHeight()%50 != 2 {
		return
	}

	keyTypes := []string{libchain.KEY_TYPE_ECDSA, libchain.KEY_TYPE_EDDSA}
	for _, keyType := range keyTypes {
		if a.keeper.IsKeygenExisted(ctx, keyType, 0) {
			continue
		}

		// Broadcast a message.
		signer := a.appKeys.GetSignerAddress()
		proposal := types.NewMsgKeygenWithSigner(
			signer.String(),
			keyType,
			0,
		)

		log.Info("Submitting proposal message for ", keyType)
		a.txSubmit.SubmitMessageAsync(proposal)
	}
}

// calculateTokenPrices gets all token prices posted from all validators and calculate the median.
func (a AppModule) calculateTokenPrices(ctx sdk.Context) {
	curBlock := ctx.BlockHeight()

	// We wait for 5 more blocks after we get prices from deyes so that any record can be posted
	// onto the blockchain.
	if curBlock%TokenPriceUpdateInterval != 5 {
		return
	}

	log.Info("Calcuating token prices....")

	// TODO: Fix the signer set.
	records := a.keeper.GetAllTokenPricesRecord(ctx)

	tokenPrices := make(map[string][]*big.Int)
	for _, data := range records {
		for _, record := range data.Records {
			// Only calculate token prices that has been updated recently.
			if curBlock-int64(record.BlockHeight) > TokenPriceUpdateInterval {
				continue
			}

			m := tokenPrices[record.Token]
			if m == nil {
				m = make([]*big.Int, 0)
			}

			value, _ := new(big.Int).SetString(record.Price, 10)
			m = append(m, value)

			tokenPrices[record.Token] = m
		}
	}

	// Now sort all the array and get the median
	medians := make(map[string]*big.Int)
	for token, list := range tokenPrices {
		if len(list) == 0 {
			log.Error("cannot find price list for token ", token)
			continue
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].Cmp(list[j]) < 0
		})
		median := list[len(list)/2]
		medians[token] = median
	}

	log.Verbose("Calculated prices = ", medians)

	// Update all the token data.
	arr := make([]string, 0, len(medians))
	for token, _ := range medians {
		arr = append(arr, token)
	}

	savedTokens := a.keeper.GetTokens(ctx, arr)

	for tokenId, price := range medians {
		savedTokens[tokenId].Price = price.String()
	}

	a.keeper.SetTokens(ctx, savedTokens)
}

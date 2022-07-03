package sisu

import (
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"

	hutils "github.com/sisu-network/dheart/utils"
)

func (p *DefaultTxOutputProducer) processCardanoBatches(ctx sdk.Context, chain string, batches [][]*transferOutData) []*types.TxOutWithSigner {
	return nil

	// // Find the highest block where majority of the validator nodes has reach to.
	// currentVals := p.valsManager.GetValAccAddrs()

	// blockHeightsMap := p.keeper.GetBlockHeightsForChain(ctx, chain, currentVals)
	// blockHeights := make([]*types.BlockHeight, 0, len(blockHeightsMap))
	// for _, blockHeight := range blockHeightsMap {
	// 	blockHeights = append(blockHeights, blockHeight)
	// }

	// // Sort by block heights.
	// sort.Slice(blockHeights, func(i, j int) bool {
	// 	return blockHeights[i].Height > blockHeights[j].Height
	// })

	// tssParams := p.keeper.GetParams(ctx)
	// majority := int(tssParams.MajorityThreshold)

	// var maxBlock int64
	// if majority >= len(blockHeights) {
	// 	maxBlock = scardano.MaxBlockHeight
	// } else {
	// 	maxBlock = blockHeights[majority].Height
	// }

	// // MaxBlock is the max block height of cardano chain.
	// for _, batch := range batches {
	// 	for _, transferOut := range batch {
	// 		tx, err := p.getCardanoTx(ctx, data, uint64(maxBlock))
	// 		if err != nil {
	// 			return nil
	// 		}
	// 	}
	// }

	// return nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (p *DefaultTxOutputProducer) getCardanoTx(ctx sdk.Context, data *transferOutData, maxBlock uint64,
	assetAmount uint64) (*cardano.Tx, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	receiverAddr, err := cardano.NewAddress(data.recipient)
	if err != nil {
		log.Error("error when parsing receiver addr: ", err)
		return nil, err
	}

	multiAsset, err := scardano.GetCardanoMultiAsset(data.destChain, data.token, assetAmount)
	if err != nil {
		return nil, err
	}

	// We need at least 1 ada to send multi assets.
	tx, err := scardano.BuildTx(p.cardanoClient, senderAddr, receiverAddr,
		cardano.NewValueWithAssets(cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset), nil,
		maxBlock)

	if err != nil {
		log.Error("error when building tx: ", err)
		return nil, err
	}

	for _, i := range tx.Body.Inputs {
		log.Debugf("tx input = %v\n", i)
	}

	for _, o := range tx.Body.Outputs {
		log.Debugf("tx output = %v\n", o)
	}

	log.Debug("tx fee = ", tx.Body.Fee)

	return tx, nil
}

package sisu

import (
	"fmt"
	"sort"

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
	// Find the highest block where majority of the validator nodes has reach to.
	currentVals := p.valsManager.GetValAccAddrs()

	blockHeightsMap := p.keeper.GetBlockHeightsForChain(ctx, chain, currentVals)
	blockHeights := make([]*types.BlockHeight, 0, len(blockHeightsMap))
	for _, blockHeight := range blockHeightsMap {
		blockHeights = append(blockHeights, blockHeight)
	}

	// Sort by block heights.
	sort.Slice(blockHeights, func(i, j int) bool {
		return blockHeights[i].Height > blockHeights[j].Height
	})

	tssParams := p.keeper.GetParams(ctx)
	majority := int(tssParams.MajorityThreshold)

	var maxBlock int64
	if majority >= len(blockHeights) {
		maxBlock = scardano.MaxBlockHeight
	} else {
		maxBlock = blockHeights[majority].Height
	}

	outMgs := make([]*types.TxOutWithSigner, 0)
	utxos, err := p.getUtxos(ctx, chain, maxBlock)
	if err != nil {
		return nil
	}
	usedUtxos := make([]*cardano.UTxO, 0)

	fmt.Println("utxos = ", len(utxos), utxos)

	// MaxBlock is the max block height of cardano chain.
	for _, batch := range batches {
		for _, transferOut := range batch {
			lovelaceAmount := utils.WeiToLovelace(transferOut.amount)
			fmt.Println("transferOut.amount = ", transferOut.amount)
			fmt.Println("lovelaceAmount = ", lovelaceAmount, lovelaceAmount.Int64(), lovelaceAmount.Uint64())

			tx, err := p.getCardanoTx(ctx, transferOut, lovelaceAmount.Uint64(), utxos, uint64(maxBlock))
			if err != nil {
				return nil
			}

			bz, err := tx.MarshalCBOR()
			if err != nil {
				return nil
			}

			for _, input := range tx.Body.Inputs {
				usedUtxos = append(usedUtxos, &cardano.UTxO{
					TxHash: input.TxHash,
					Index:  input.Index,
				})
			}

			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				types.TxOutType_TRANSFER_OUT,
				transferOut.blockHeight,
				transferOut.destChain,
				"",
				transferOut.destChain,
				utils.KeccakHash32Bytes(bz),
				bz,
				"",
			)
			outMgs = append(outMgs, outMsg)

			// TODO: Make this track multiple transactions.
			// // Track the txout
			// p.txTracker.AddTransaction(
			// 	outMsg.Data,
			// 	transfer.txIn,
			// )
		}
	}

	for _, utxo := range usedUtxos {
		p.privateDb.AddUtxo(chain, utxo.TxHash.String(), int(utxo.Index))
	}

	return outMgs
}

func (p *DefaultTxOutputProducer) getUtxos(ctx sdk.Context, chain string, height int64) ([]cardano.UTxO, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)

	utxos, err := p.cardanoClient.UTxOs(senderAddr, uint64(height))
	if err != nil {
		return nil, err
	}

	// Filter used utxo
	filtered := make([]cardano.UTxO, 0)
	for _, utxo := range utxos {
		if !p.privateDb.IsUtxoExisted(chain, utxo.TxHash.String(), int(utxo.Index)) {
			fmt.Println("Utxo amount = ", utxo.Amount)
			filtered = append(filtered, utxo)
		}
	}

	return filtered, nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (p *DefaultTxOutputProducer) getCardanoTx(ctx sdk.Context, data *transferOutData,
	assetAmount uint64, utxos []cardano.UTxO, maxBlock uint64) (*cardano.Tx, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	receiverAddr, err := cardano.NewAddress(data.recipient)
	if err != nil {
		log.Error("error when parsing receiver addr: ", err)
		return nil, err
	}

	fmt.Println("getCardanoTx: assetAmount = ", assetAmount)

	multiAsset, err := scardano.GetCardanoMultiAsset(data.destChain, data.token, assetAmount)
	if err != nil {
		return nil, err
	}

	fmt.Println("utils.ONE_ADA_IN_LOVELACE.Uint64() = ", utils.ONE_ADA_IN_LOVELACE.Uint64())

	// We need at least 1 ada to send multi assets.
	tx, err := scardano.BuildTx(p.cardanoClient, senderAddr, receiverAddr,
		cardano.NewValueWithAssets(cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset), nil,
		utxos, maxBlock)

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

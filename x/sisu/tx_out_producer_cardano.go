package sisu

import (
	"math"
	"math/big"

	scardano "github.com/sisu-network/sisu/x/sisu/cardano"
	"github.com/sisu-network/sisu/x/sisu/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"

	hutils "github.com/sisu-network/dheart/utils"
)

func (p *DefaultTxOutputProducer) processCardanoBatches(ctx sdk.Context, k keeper.Keeper, destChain string,
	transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	// Find the highest block where majority of the validator nodes has reach to.
	outMgs := make([]*types.TxOutMsg, 0)
	inHashes := make([]string, len(transfers))

	for _, transfer := range transfers {
		inHashes = append(inHashes, transfer.Id)
	}

	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)

	var maxBlockHeight uint64
	checkPoint := k.GetGatewayCheckPoint(ctx, destChain)
	if checkPoint == nil {
		maxBlockHeight = math.MaxUint64
	} else {
		maxBlockHeight = uint64(checkPoint.BlockHeight)
	}
	utxos, err := p.cardanoClient.UTxOs(senderAddr, maxBlockHeight)
	if err != nil {
		return nil, err
	}

	tx, err := p.getCardanoTx(ctx, destChain, transfers, utxos, maxBlockHeight)
	if err != nil {
		return nil, err
	}

	bz, err := tx.MarshalCBOR()
	if err != nil {
		return nil, err
	}

	hash, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		p.signer,
		types.TxOutType_TRANSFER_OUT,
		inHashes,
		destChain,
		hash.String(),
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

	return outMgs, nil
}

func (p *DefaultTxOutputProducer) getUtxos(ctx sdk.Context, chain string, height int64) ([]cardano.UTxO, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)

	return p.cardanoClient.UTxOs(senderAddr, uint64(height))
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (p *DefaultTxOutputProducer) getCardanoTx(ctx sdk.Context, chain string, transfers []*types.Transfer,
	utxos []cardano.UTxO, maxBlock uint64) (*cardano.Tx, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	allTokens := p.keeper.GetAllTokens(ctx)
	receiverAddrs := make([]cardano.Address, 0)
	amounts := make([]*cardano.Value, 0, len(transfers))
	for _, transfer := range transfers {
		// Receivers
		receiverAddr, err := cardano.NewAddress(transfer.Recipient)
		if err != nil {
			log.Error("error when parsing receiver addr: ", err)
			continue
		}
		receiverAddrs = append(receiverAddrs, receiverAddr)

		token := allTokens[transfer.Token]
		if token == nil {
			continue
		}

		amountInt, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warnf("Cannot create big.Int value from amount %s on chain %s", transfer.Amount, chain)
			continue
		}

		// amounts
		lovelaceAmount := utils.WeiToLovelace(amountInt)
		multiAsset, err := scardano.GetCardanoMultiAsset(chain, token, lovelaceAmount.Uint64())
		if err != nil {
			return nil, err
		}
		amount := cardano.NewValueWithAssets(1_600_000, multiAsset)
		amounts = append(amounts, amount)
	}

	// We need at least 1 ada to send multi assets.
	tx, err := scardano.BuildTx(p.cardanoClient, senderAddr, receiverAddrs, amounts, nil, utxos, maxBlock)

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

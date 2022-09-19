package cardano

import (
	"fmt"
	"math"
	"math/big"

	"github.com/sisu-network/lib/log"

	"github.com/echovl/cardano-go"
	hutils "github.com/sisu-network/dheart/utils"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	chaintypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type bridge struct {
	chain  string
	signer string
	keeper keeper.Keeper
	client CardanoClient
}

func NewBridge(chain string, signer string, keeper keeper.Keeper, client CardanoClient) chaintypes.Bridge {
	return &bridge{
		keeper: keeper,
		chain:  chain,
		signer: signer,
		client: client,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	// Find the highest block where majority of the validator nodes has reach to.
	outMgs := make([]*types.TxOutMsg, 0)
	inHashes := make([]string, len(transfers))

	for _, transfer := range transfers {
		inHashes = append(inHashes, transfer.Id)
	}

	pubkey := b.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)

	var maxBlockHeight uint64
	checkPoint := b.keeper.GetGatewayCheckPoint(ctx, b.chain)
	if checkPoint == nil {
		maxBlockHeight = math.MaxUint64
	} else {
		maxBlockHeight = uint64(checkPoint.BlockHeight)
	}
	utxos, err := b.client.UTxOs(senderAddr, maxBlockHeight)
	if err != nil {
		return nil, err
	}

	tx, err := b.getCardanoTx(ctx, b.chain, transfers, utxos, maxBlockHeight)
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
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  hash.String(),
			OutBytes: bz,
		},
		&types.TxOutInput{
			TransferIds: inHashes,
		},
	)
	outMgs = append(outMgs, outMsg)

	return outMgs, nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (b *bridge) getCardanoTx(ctx sdk.Context, chain string, transfers []*types.Transfer,
	utxos []cardano.UTxO, maxBlock uint64) (*cardano.Tx, error) {
	pubkey := b.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	allTokens := b.keeper.GetAllTokens(ctx)
	receiverAddrs := make([]cardano.Address, 0)
	amounts := make([]*cardano.Value, 0, len(transfers))
	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("Commission rate is invalid, rate = %d", commissionRate)
	}
	for _, transfer := range transfers {
		// Receivers
		receiverAddr, err := cardano.NewAddress(transfer.ToRecipient)
		if err != nil {
			log.Error("error when parsing receiver addr: ", err)
			continue
		}
		receiverAddrs = append(receiverAddrs, receiverAddr)

		token := allTokens[transfer.Token]
		if token == nil {
			continue
		}

		amountOut, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warnf("Cannot create big.Int value from amount %s on chain %s", transfer.Amount, chain)
			continue
		}

		// Subtract commission rate
		amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

		// Subtract the 1.6 ADA for multi asset transaction
		if token.Id == "ADA" {
			amountOut = amountOut.Sub(amountOut, utils.OnePointSixEthToWei)
		} else {
			// Convert the price of 1.6 ADA to token unit
			adaToken := allTokens["ADA"]
			adaInUsd, ok := new(big.Int).SetString(adaToken.Price, 10)
			if !ok {
				return nil, fmt.Errorf("Invalid ada price %s", adaToken.Price)
			}
			// Times 1.6
			tmp := new(big.Int).Mul(adaInUsd, big.NewInt(16))
			adaInUsd = tmp.Div(tmp, big.NewInt(10))
			// Get the token amount from ada price
			tokenPrice, ok := new(big.Int).SetString(token.Price, 10)
			if !ok {
				return nil, fmt.Errorf("Invalid token price %s", adaToken.Price)
			}

			// Amount of ADA fee in Token price
			amountInToken := adaInUsd.Mul(adaInUsd, utils.EthToWei)
			amountInToken = amountInToken.Div(amountInToken, tokenPrice)

			amountOut = amountOut.Sub(amountOut, amountInToken)
		}

		// If amountOut is smaller or equal 0, quit
		if amountOut.Cmp(utils.ZeroBigInt) < 0 {
			return nil, common.InsufficientFundErr
		}

		// Convert from Wei unit to lovelace unit
		lovelaceAmount := utils.WeiToLovelace(amountOut)

		var amount *cardano.Value
		if token.Id == "ADA" {
			// Transfer native ADA instead of wrapped token
			amount = cardano.NewValue(cardano.Coin(amountOut.Uint64()))
		} else {
			multiAsset, err := GetCardanoMultiAsset(chain, token, lovelaceAmount.Uint64())
			if err != nil {
				return nil, err
			}
			amount = cardano.NewValueWithAssets(1_600_000, multiAsset)
		}
		amounts = append(amounts, amount)
	}

	// We need at least 1 ada to send multi assets.
	tx, err := BuildTx(b.client, senderAddr, receiverAddrs, amounts, nil, utxos, maxBlock)

	if err != nil {
		log.Error("error when building tx: ", err)
		return nil, err
	}

	for _, i := range tx.Body.Inputs {
		log.Verbosef("tx input = %v\n", i)
	}

	for _, o := range tx.Body.Outputs {
		log.Verbosef("tx output = %v\n", o)
	}

	log.Verbose("tx fee = ", tx.Body.Fee)

	return tx, nil
}
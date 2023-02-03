package cardano

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	eyesTypes "github.com/sisu-network/deyes/types"

	"github.com/sisu-network/lib/log"

	"github.com/echovl/cardano-go"
	hutils "github.com/sisu-network/dheart/utils"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	chaintypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type bridge struct {
	chain       string
	signer      string
	keeper      keeper.Keeper
	deyesClient external.DeyesClient
}

func NewBridge(chain string, signer string, keeper keeper.Keeper, deyesClient external.DeyesClient) chaintypes.Bridge {
	return &bridge{
		keeper:      keeper,
		chain:       chain,
		signer:      signer,
		deyesClient: deyesClient,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOut, error) {
	// Find the highest block where majority of the validator nodes has reach to.
	outMgs := make([]*types.TxOut, 0)
	inHashes := make([]string, len(transfers))

	for _, transfer := range transfers {
		inHashes = append(inHashes, transfer.Id)
	}

	pubkey := b.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	sisuAddr := hutils.GetAddressFromCardanoPubkey(pubkey)

	var maxBlockHeight uint64
	block := b.keeper.GetBlockHeight(ctx, b.chain)
	if block == nil {
		maxBlockHeight = math.MaxUint64
	} else {
		maxBlockHeight = uint64(block.Height)
	}

	utxos, err := b.deyesClient.CardanoUtxos(b.chain, sisuAddr.String(), maxBlockHeight)
	for _, utxo := range utxos {
		fmt.Printf("utxo = %v+\n", utxo.Amount.MultiAsset)
	}

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

	outMsg := &types.TxOut{
		TxType: types.TxOutType_TRANSFER_OUT,
		Content: &types.TxOutContent{
			OutChain: b.chain,
			OutHash:  hash.String(),
			OutBytes: bz,
		},
		Input: &types.TxOutInput{
			TransferIds: inHashes,
		},
	}

	outMgs = append(outMgs, outMsg)

	return outMgs, nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (b *bridge) getCardanoTx(ctx sdk.Context, chain string, transfers []*types.TransferDetails,
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

		// Convert from Wei unit to lovelace unit
		amountOut = utils.WeiToLovelace(amountOut)

		// Subtract commission rate
		amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

		if token.Id == "ADA" {
			// Subtract 0.2 ADA for transaction fee.
			amountOut = amountOut.Sub(amountOut, new(big.Int).Div(big.NewInt(utils.OneAdaInLoveLace),
				big.NewInt(5)))
		} else {
			// Subtract the 1.6 ADA for multi asset transaction

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

			if tokenPrice.Cmp(big.NewInt(0)) == 0 {
				return nil, fmt.Errorf("Token %s has price 0", token.Id)
			}

			// Amount of ADA fee in Token price
			amountInToken := adaInUsd.Mul(adaInUsd, big.NewInt(utils.OneAdaInLoveLace))
			amountInToken = amountInToken.Div(amountInToken, tokenPrice)

			amountOut = amountOut.Sub(amountOut, amountInToken)
		}

		// If amountOut is smaller or equal 0, quit
		if amountOut.Cmp(utils.ZeroBigInt) <= 0 {
			return nil, components.InsufficientFundErr
		}

		var amount *cardano.Value
		if token.Id == "ADA" {
			// Minimum ADA per UTXO is 1,000,000 lovelace.
			if amountOut.Cmp(big.NewInt(utils.OneAdaInLoveLace)) < 0 {
				return nil, fmt.Errorf("Lovelace output is %s, min requirement is 1_000_000 lovelace",
					amountOut.String())
			}

			// Transfer native ADA instead of wrapped token
			amount = cardano.NewValue(cardano.Coin(amountOut.Uint64()))
		} else {
			multiAsset, err := GetCardanoMultiAsset(chain, token, amountOut.Uint64())
			if err != nil {
				return nil, err
			}
			amount = cardano.NewValueWithAssets(1_600_000, multiAsset)
		}
		amounts = append(amounts, amount)
	}

	// We need at least 1 ada to send multi assets.
	tx, err := BuildTx(b.deyesClient, chain, senderAddr, receiverAddrs, amounts, nil, utxos, maxBlock)

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

func (b *bridge) ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) ([]*types.TransferDetails, error) {
	ret := make([]*types.TransferDetails, 0)
	cardanoTx := &eyesTypes.CardanoTransactionUtxo{}
	err := json.Unmarshal(serialized, cardanoTx)
	if err != nil {
		return nil, err
	}

	if cardanoTx.Metadata != nil {
		nativeTransfer := cardanoTx.Metadata.NativeAda != 0
		log.Verbose("cardanoTx.Amount = ", cardanoTx.Amount)

		// Convert from ADA unit (10^6) to our standard unit (10^18)
		for _, amount := range cardanoTx.Amount {
			tokenUnit := amount.Unit
			if tokenUnit != "lovelace" {
				token := GetTokenFromCardanoAsset(ctx, b.keeper, tokenUnit, chain)
				if token == nil {
					log.Error("Failed to find token with id: ", tokenUnit)
					continue
				}
				tokenUnit = token.Id
			} else {
				if !nativeTransfer {
					// This ADA is for transaction transfer fee. It is not meant to be transfered.
					continue
				}
				tokenUnit = "ADA"
			}

			// Converte token quantity to Sisu decimals
			quantity, ok := new(big.Int).SetString(amount.Quantity, 10)
			if !ok {
				return nil, fmt.Errorf("Failed to get amount quantity in cardano tx, amount = %s", amount.Quantity)
			}
			tokens := b.keeper.GetTokens(ctx, []string{tokenUnit})
			if len(tokens) == 0 {
				return nil, fmt.Errorf("Cannot find token %s", tokenUnit)
			}
			amountBig, err := tokens[tokenUnit].ChainAmountToSisuAmount(chain, quantity)
			if err != nil {
				log.Warnf("Cannot convert amount %d on chain %s", amount.Quantity, chain)
				continue
			}

			log.Verbose("tokenUnit = ", tokenUnit, " quantity = ", quantity)
			log.Verbose("cardanoTx.Metadata = ", cardanoTx.Metadata)

			ret = append(ret, &types.TransferDetails{
				FromHash:    cardanoTx.Hash,
				Token:       tokenUnit,
				Amount:      amountBig.String(),
				ToChain:     cardanoTx.Metadata.Chain,
				ToRecipient: cardanoTx.Metadata.Recipient,
			})
		}
	}

	return ret, nil
}

func (b *bridge) ProcessCommand(ctx sdk.Context, cmd *types.Command) (*types.TxOut, error) {
	// TODO: Implement this
	return nil, types.NewErrNotImplemented(
		fmt.Sprintf("ProcessCommand not implemented for chain %s", b.chain))
}

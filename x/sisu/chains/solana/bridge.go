package solana

import (
	"encoding/json"
	"fmt"
	"math/big"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"

	"github.com/mr-tron/base58"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/sisu-network/sisu/x/sisu/external"

	sdk "github.com/cosmos/cosmos-sdk/types"
	chaintypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	// The amount needs to fix into an uin64 integer. This is the max amount user can transfer.
	MaxTransferAmountInteger = new(big.Int).Exp(big.NewInt(2), big.NewInt(63), nil)
)

type defaultBridge struct {
	chain    string
	keeper   keeper.Keeper
	signer   string
	config   config.Config
	deyesCli external.DeyesClient
}

func NewBridge(chain, signer string, keeper keeper.Keeper, cfg config.Config,
	deyesCli external.DeyesClient) chaintypes.Bridge {
	return &defaultBridge{
		chain:    chain,
		keeper:   keeper,
		signer:   signer,
		config:   cfg,
		deyesCli: deyesCli,
	}
}

func (b *defaultBridge) ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) ([]*types.TxOutMsg, error) {
	allTokens := b.keeper.GetAllTokens(ctx)
	tokens := make([]*types.Token, 0, len(transfers))
	recipients := make([]string, 0, len(transfers))
	amounts := make([]*big.Int, 0, len(transfers))
	inHashes := make([]string, 0, len(transfers))

	for _, transfer := range transfers {
		token := allTokens[transfer.Token]
		if token == nil {
			log.Warn("cannot find token ", transfer.Token)
			continue
		}

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amount ", transfer.Amount)
			continue
		}

		tokens = append(tokens, token)
		recipients = append(recipients, transfer.ToRecipient)
		amounts = append(amounts, amount)
		inHashes = append(inHashes, transfer.Id)
	}

	responseTx, err := b.buildTransferInResponse(ctx, tokens, recipients, amounts)
	if err != nil {
		log.Error("Failed to build solana transfer in, err = ", err)
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  utils.KeccakHash32(string(responseTx.RawBytes)),
			OutBytes: responseTx.RawBytes,
		},
		&types.TxOutInput{
			TransferIds: inHashes,
		},
	)

	return []*types.TxOutMsg{outMsg}, nil
}

func (b *defaultBridge) buildTransferInResponse(
	ctx sdk.Context,
	tokens []*types.Token,
	recipients []string,
	amounts []*big.Int,
) (*types.TxResponse, error) {
	tx, err := b.getTransaction(ctx, tokens, recipients, amounts)
	if err != nil {
		return nil, err
	}

	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{
		OutChain: b.chain,
		RawBytes: messageContent,
	}, nil
}

func (b *defaultBridge) getTransaction(
	ctx sdk.Context,
	tokens []*types.Token,
	recipients []string,
	amounts []*big.Int,
) (*solanago.Transaction, error) {
	chain := b.chain
	// Get mpc address
	mpcAddr := b.keeper.GetMpcAddress(ctx, chain)
	if mpcAddr == "" {
		return nil, fmt.Errorf("Cannot find mpc address for chain %s", chain)
	}
	mpcPubkey, err := solanago.PublicKeyFromBase58(mpcAddr)
	if err != nil {
		return nil, err
	}

	tokenAddrs := make([]string, 0)
	for _, token := range tokens {
		addr := token.GetAddressForChain(chain)
		if addr == "" {
			return nil, fmt.Errorf("Cannot find token %s address for chain %s", token.Id, chain)
		}
		tokenAddrs = append(tokenAddrs, addr)
	}

	nonce := b.keeper.GetMpcNonce(ctx, chain)
	if nonce == nil {
		return nil, fmt.Errorf("Nonce is nil for chain %s", chain)
	}

	// Convert amount into token with correct decimal
	solAmounts := make([]uint64, 0)
	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("Commission rate is invalid, rate = %d", commissionRate)
	}

	for i, amount := range amounts {
		// Get token decimals
		decimals := tokens[i].GetDecimalsForChain(chain)
		base := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
		amountOut := new(big.Int).Mul(amount, base)
		amountOut = amountOut.Div(amountOut, utils.EthToWei)
		amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

		if amountOut.Cmp(MaxTransferAmountInteger) > 0 {
			return nil, fmt.Errorf(
				"TransferExceedMax amount, original amount = %s, token decimals decimals = %d",
				amount.String(),
				decimals,
			)
		}

		solAmounts = append(solAmounts, amountOut.Uint64())
	}

	// TODO: Don't hardcode token program id here. Make each token has different token program ID
	transferInIx, err := solanatypes.NewTransferInIx(
		b.config.Solana.BridgeProgramId,
		mpcAddr,
		uint64(nonce.Nonce),
		"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		b.config.Solana.BridgePda,
		tokenAddrs,
		recipients,
		solAmounts,
	)
	if err != nil {
		return nil, err
	}

	// Get recennt block hash
	result, err := b.deyesCli.SolanaQueryRecentBlock(b.chain)
	if err != nil {
		return nil, err
	}

	hash, err := solanago.HashFromBase58(result.Hash)
	if err != nil {
		return nil, err
	}

	tx, err := solanago.NewTransaction(
		[]solanago.Instruction{transferInIx},
		hash,
		solanago.TransactionPayer(mpcPubkey),
	)

	return tx, err
}

func (b *defaultBridge) ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) ([]*types.TransferDetails, error) {
	log.Verbose("Parsing solana incomgin tx...")
	ret := make([]*types.TransferDetails, 0)

	outerTx := new(eyessolanatypes.Transaction)
	err := json.Unmarshal(serialized, outerTx)
	if err != nil {
		return nil, err
	}

	if outerTx.TransactionInner == nil || outerTx.TransactionInner.Message == nil || outerTx.TransactionInner.Message.AccountKeys == nil {
		return nil, fmt.Errorf("Invalid outerTx")
	}

	accounts := outerTx.TransactionInner.Message.AccountKeys

	allTokens := b.keeper.GetAllTokens(ctx)
	// Check that there is at least one instruction sent to the program id
	for _, ix := range outerTx.TransactionInner.Message.Instructions {
		if accounts[ix.ProgramIdIndex] != b.config.Solana.BridgeProgramId {
			continue
		}

		if len(outerTx.TransactionInner.Signatures) == 0 {
			continue
		}

		// Decode the instruction
		bytesArr, err := base58.Decode(ix.Data)
		if err != nil {
			return nil, err
		}

		if len(bytesArr) == 0 {
			return nil, fmt.Errorf("Data is empty")
		}

		instruction := bytesArr[0]

		switch instruction {
		case solanatypes.TransferOut:
			transferOut := new(solanatypes.TransferOutData)
			err = transferOut.Deserialize(bytesArr)
			if err != nil {
				return nil, err
			}

			// look up the token in the keeper
			log.Verbose("Transfer data on solana = ", *transferOut)
			token := utils.GetTokenOnChain(allTokens, transferOut.TokenAddress, chain)
			if token == nil {
				continue
			}

			amount, err := token.ChainAmountToSisuAmount(chain, big.NewInt(int64(transferOut.Amount)))
			if err != nil {
				log.Warnf("Cannot convert amount %d on chain %s", transferOut.Amount, chain)
				continue
			}

			ret = append(ret, &types.TransferDetails{
				FromChain:   chain,
				FromHash:    outerTx.TransactionInner.Signatures[0],
				Token:       token.Id,
				Amount:      amount.String(),
				ToChain:     libchain.GetChainNameFromInt(big.NewInt(int64(transferOut.ChainId))),
				ToRecipient: transferOut.Recipient,
			})

		case solanatypes.TransferIn:
			transferIn := new(solanatypes.TransferInData)
			err = borsh.Deserialize(transferIn, bytesArr)
			if err != nil {
				return nil, err
			}

			log.Warn("This is a transfer in. Do nothing. It should be confirmed by Sisu")
		}
	}

	return ret, nil
}

func (b *defaultBridge) ProcessCommand(ctx sdk.Context, cmd *types.Command) (*types.TxOutMsg, error) {
	// TODO: Implement this
	return nil, types.NewErrNotImplemented(
		fmt.Sprintf("ProcessCommand not implemented for chain %s", b.chain))
}

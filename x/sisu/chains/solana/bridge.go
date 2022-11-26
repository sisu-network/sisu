package solana

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	solanago "github.com/gagliardetto/solana-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"

	"github.com/mr-tron/base58"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	eyestypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	chaintypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type bridge struct {
	chain  string
	keeper keeper.Keeper
	signer string
	config config.Config
}

func NewBridge(chain, signer string, keeper keeper.Keeper, cfg config.Config) chaintypes.Bridge {
	return &bridge{
		chain:  chain,
		keeper: keeper,
		signer: signer,
		config: cfg,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	allTokens := b.keeper.GetAllTokens(ctx)
	tokens := make([]*types.Token, 0, len(transfers))
	recipients := make([]string, 0, len(transfers))
	amounts := make([]*big.Int, 0, len(transfers))
	inHashes := make([]string, 0, len(transfers))

	for _, transfer := range transfers {
		token := allTokens[transfer.Token]
		if token == nil {
			log.Warn("cannot find token", transfer.Token)
			continue
		}

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
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

func (b *bridge) buildTransferInResponse(
	ctx sdk.Context,
	tokens []*types.Token,
	recipients []string,
	amounts []*big.Int,
) (*types.TxResponse, error) {
	chain := b.chain
	// Get mpc address
	mpcAddr := b.keeper.GetMpcAddress(ctx, chain)
	if mpcAddr == "" {
		return nil, fmt.Errorf("Cannot find mpc address for chain %s", chain)
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
		return nil, fmt.Errorf("Nonce is nil for chainn %s", chain)
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
		amounts,
	)

	if err != nil {
		return nil, err
	}

	// Get recennt block hash
	hashStr, err := b.getRecentBlockHash(ctx, chain)
	if err != nil {
		return nil, err
	}
	hash, err := solanago.HashFromBase58(hashStr)
	if err != nil {
		return nil, err
	}

	tx, err := solanago.NewTransaction(
		[]solanago.Instruction{transferInIx},
		hash,
	)

	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{
		OutChain: b.chain,
		RawBytes: messageContent,
	}, nil
}

func (b *bridge) getRecentBlockHash(ctx sdk.Context, chain string) (string, error) {
	metas := b.keeper.GetAllSolanaConfirmedBlock(ctx, chain)
	if len(metas) == 0 {
		return "", fmt.Errorf("Empty metas array")
	}

	arr := make([]*types.ChainMetadata, 0)
	for _, value := range metas {
		arr = append(arr, value)
	}

	// Sort arr
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].SolanaRecentBlockHeight < arr[j].SolanaRecentBlockHeight
	})

	return arr[len(arr)/2].SolanaRecentBlockHash, nil
}

func (b *bridge) ParseIncomginTx(ctx sdk.Context, chain string, tx *eyestypes.Tx) ([]*types.Transfer, error) {
	ret := make([]*types.Transfer, 0)

	outerTx := new(eyessolanatypes.Transaction)
	err := json.Unmarshal(tx.Serialized, outerTx)
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

		transferOut := new(solanatypes.TransferOutData)
		err = transferOut.Deserialize(bytesArr)
		if err != nil {
			return nil, err
		}

		switch transferOut.Instruction {
		case solanatypes.TransferOut:
			// look up the token in the keeper
			log.Verbose("Transfer data on solana = ", *transferOut)
			token := utils.GetTokenOnChain(allTokens, transferOut.TokenAddress, chain)
			if token == nil {
				continue
			}

			amount, err := token.ConvertAmountToSisuAmount(chain, big.NewInt(int64(transferOut.Amount)))
			if err != nil {
				log.Warnf("Cannot convert amount %d on chain %s", transferOut.Amount, chain)
				continue
			}

			ret = append(ret, &types.Transfer{
				FromChain:   chain,
				FromHash:    outerTx.TransactionInner.Signatures[0],
				Token:       token.Id,
				Amount:      amount.String(),
				ToChain:     libchain.GetChainNameFromInt(big.NewInt(int64(transferOut.ChainId))),
				ToRecipient: transferOut.Recipient,
			})
		}
	}

	return ret, nil
}

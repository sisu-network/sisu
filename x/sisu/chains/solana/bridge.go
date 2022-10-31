package solana

import (
	"encoding/json"
	"fmt"
	"math/big"

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
	config config.Config
}

func NewBridge(chain string, keeper keeper.Keeper, cfg config.Config) chaintypes.Bridge {
	return &bridge{
		chain:  chain,
		keeper: keeper,
		config: cfg,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	return nil, nil
}

func (b *bridge) ParseIncomginTx(ctx sdk.Context, chain string, tx *eyestypes.Tx) ([]*types.Transfer, error) {
	ret := make([]*types.Transfer, 0)

	fmt.Println("AAAAAA 000000")
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

		fmt.Println("AAAAAA 111111")

		ix := new(solanatypes.TransferOutInstruction)
		err = ix.Deserialize(bytesArr)
		if err != nil {
			return nil, err
		}

		transferData := ix.Data

		fmt.Println("AAAAAA 222222")

		switch ix.Instruction {
		case solanatypes.TranserOut:
			// look up the token in the keeper
			log.Verbose("Transfer data on solana = ", transferData)
			token := utils.GetTokenOnChain(allTokens, transferData.TokenAddress, chain)
			if token == nil {
				fmt.Println("Cannot find token with address ", transferData.TokenAddress)
				continue
			}

			fmt.Println("Adding to ret")

			ret = append(ret, &types.Transfer{
				FromChain:   chain,
				FromHash:    outerTx.TransactionInner.Signatures[0],
				Token:       token.Id,
				Amount:      transferData.Amount.String(), // TODO Convert token amount
				ToChain:     libchain.GetChainNameFromInt(big.NewInt(int64(transferData.ChainId))),
				ToRecipient: transferData.Recipient,
			})
		}
	}

	fmt.Println("len ret = ", len(ret))

	return ret, nil
}

package solana

import (
	"encoding/json"
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/mr-tron/base58"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	eyestypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
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

	outerTx := new(eyessolanatypes.Transaction)
	err := json.Unmarshal(tx.Serialized, outerTx)
	if err != nil {
		return nil, err
	}

	accounts := outerTx.TransactionInner.Message.AccountKeys

	// Check that there is at least one instruction sent to the program id
	for _, ix := range outerTx.TransactionInner.Message.Instructions {
		if accounts[ix.ProgramIdIndex] != b.config.Solana.BridgeProgramId {
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

		ix := new(solanatypes.TransferOutInstruction)
		err = ix.Deserialize(bytesArr)
		if err != nil {
			return nil, err
		}

		transferData := ix.Data

		switch ix.Instruction {
		case solanatypes.TranserOut:
			ret = append(ret, &types.Transfer{
				FromHash:    outerTx.TransactionInner.Signatures[0],
				Token:       transferData.TokenAddress,
				Amount:      transferData.Amount.String(),
				ToChain:     libchain.GetChainNameFromInt(big.NewInt(int64(transferData.ChainId))),
				ToRecipient: transferData.Recipient,
			})
		}
	}

	return ret, nil
}

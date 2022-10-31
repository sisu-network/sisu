package solana

import (
	"encoding/json"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func mockForBridgeTest() (sdk.Context, keeper.Keeper) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)

	return ctx, k
}

func TestParseSolana(t *testing.T) {
	bridgeProgramId := "HguMTvmDfspHuEWycDSP1XtVQJi47hVNAyLbFEf2EJEQ"

	ctx, k := mockForBridgeTest()

	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = bridgeProgramId

	bridge := NewBridge("solana-devnet", k, cfg)

	transferOut := solanatypes.TransferOutInstruction{
		Instruction: solanatypes.TranserOut,
		Data: solanatypes.TransferOutData{
			Amount: *big.NewInt(100),
		},
	}

	bz, err := transferOut.Serialize()
	require.Nil(t, err)

	outerTx := &eyessolanatypes.Transaction{
		Meta: &eyessolanatypes.TransactionMeta{},
		TransactionInner: &eyessolanatypes.TransactionInner{
			Signatures: []string{"Signature"},
			Message: &eyessolanatypes.TransactionMessage{
				AccountKeys: []string{bridgeProgramId},
				Instructions: []eyessolanatypes.Instruction{
					{
						ProgramIdIndex: 0,
						Data:           base58.Encode(bz),
					},
				},
			},
		},
	}

	bz, err = json.Marshal(outerTx)
	require.Nil(t, err)

	eyesTx := &eyesTypes.Tx{
		Hash:       outerTx.TransactionInner.Signatures[0],
		Serialized: bz,
		To:         outerTx.TransactionInner.Message.AccountKeys[0],
		Success:    true,
	}

	transfers, err := bridge.ParseIncomginTx(ctx, "solana-devnet", eyesTx)
	require.Nil(t, err)

	require.Equal(t, 1, len(transfers))
}

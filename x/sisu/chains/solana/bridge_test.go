package solana

import (
	"encoding/json"
	"testing"

	bin "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func mockForBridgeTest() (sdk.Context, keeper.Keeper) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)

	return ctx, k
}

func TestParseIncoming(t *testing.T) {
	bridgeProgramId := "HguMTvmDfspHuEWycDSP1XtVQJi47hVNAyLbFEf2EJEQ"

	ctx, k := mockForBridgeTest()

	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = bridgeProgramId

	bridge := NewBridge("solana-devnet", "signer", k, cfg)

	transferOut := solanatypes.TransferOutData{
		Instruction:  solanatypes.TransferOut,
		Amount:       100,
		TokenAddress: "8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX",
		ChainId:      1,
		Recipient:    "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988",
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

func TestProcessTransfer(t *testing.T) {
	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = "HguMTvmDfspHuEWycDSP1XtVQJi47hVNAyLbFEf2EJEQ"
	cfg.Solana.BridgePda = "6GQbD1BxAiog9Ym1YLnaci4BpcR1HLNNdc7wRrBqPA2D"

	chain := "solana-devnent"
	mpcAddress := "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
	tokenProgramId := "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"

	ctx, k := mockForBridgeTest()
	k.SetMpcAddress(ctx, chain, mpcAddress)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer", "Q6XprfkF8RQQKoQVG33xT88H7wi8Uk1B1CC7YAs69Gi", 100)
	k.SetMpcNonce(ctx, &types.MpcNonce{Chain: chain, Nonce: 1})

	transfer := &types.Transfer{
		Id:          "transfer_123",
		Token:       "8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX",
		ToRecipient: "CLKfJz4bTt9a2ZEJdv2FLwBziUBAtDH1bDfrUKENkf1H",
	}

	bridge := NewBridge(chain, "signer", k, cfg)
	msgs, err := bridge.ProcessTransfers(ctx, []*types.Transfer{transfer})
	require.Nil(t, err)

	require.Equal(t, 1, len(msgs))

	// Let's deserialize the tx
	message := solanago.Message{}
	err = message.UnmarshalLegacy(bin.NewCompactU16Decoder(msgs[0].Data.Content.OutBytes))
	require.Nil(t, err)
	require.Equal(t, 1, len(message.Instructions))
	require.Equal(t, []solanago.PublicKey{
		solanago.MustPublicKeyFromBase58(mpcAddress),
		solanago.MustPublicKeyFromBase58(tokenProgramId),
		solanago.MustPublicKeyFromBase58(cfg.Solana.BridgePda),
		solanago.MustPublicKeyFromBase58(cfg.Solana.BridgeProgramId),
	}, message.AccountKeys)
}

package solana

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/decred/dcrd/dcrec/edwards/v2"
	bin "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eyessolanatypes "github.com/sisu-network/deyes/chains/solana/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
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

	transfers, err := bridge.ParseIncomingTx(ctx, "solana-devnet", eyesTx)
	require.Nil(t, err)

	require.Equal(t, 1, len(transfers))
}

func TestProcessTransfer(t *testing.T) {
	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = "HguMTvmDfspHuEWycDSP1XtVQJi47hVNAyLbFEf2EJEQ"
	cfg.Solana.BridgePda = "6GQbD1BxAiog9Ym1YLnaci4BpcR1HLNNdc7wRrBqPA2D"

	chain := "solana-devnet"
	mpcAddress := "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
	tokenProgramId := "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	nonce := 1
	amountInt := 5

	ctx, k := mockForBridgeTest()
	k.SetMpcAddress(ctx, chain, mpcAddress)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer", "Q6XprfkF8RQQKoQVG33xT88H7wi8Uk1B1CC7YAs69Gi", 100)
	k.SetMpcNonce(ctx, &types.MpcNonce{Chain: chain, Nonce: int64(nonce)})

	tokens := k.GetTokens(ctx, []string{"SISU"})
	tokenAddr := tokens["SISU"].GetAddressForChain(chain)
	receiverAta, err := solanatypes.GetAtaPubkey("5s3YB3BzLKNxT4bKjxfXTeQnNuokkH5J68tHMN7uqV8q", tokenAddr)
	require.Nil(t, err)
	transfer := &types.Transfer{
		Id:          "transfer_123",
		Token:       "SISU",
		ToRecipient: receiverAta.String(),
		Amount:      new(big.Int).Mul(big.NewInt(int64(amountInt)), utils.EthToWei).String(),
	}

	bridge := NewBridge(chain, "signer", k, cfg)
	msgs, err := bridge.ProcessTransfers(ctx, []*types.Transfer{transfer})
	require.Nil(t, err)

	require.Equal(t, 1, len(msgs))

	// Find the bridge ata
	bridgeAta, err := solanatypes.GetAtaPubkey(cfg.Solana.BridgePda, tokenAddr)
	require.Nil(t, err)

	// Let's deserialize the tx
	message := solanago.Message{}
	err = message.UnmarshalLegacy(bin.NewCompactU16Decoder(msgs[0].Data.Content.OutBytes))
	require.Nil(t, err)

	// Verify message accounts
	require.Equal(t, 6, len(message.AccountKeys))
	require.Contains(t, message.AccountKeys, solanago.MustPublicKeyFromBase58(mpcAddress))
	require.Contains(t, message.AccountKeys, solanago.MustPublicKeyFromBase58(tokenProgramId))
	require.Contains(t, message.AccountKeys, solanago.MustPublicKeyFromBase58(cfg.Solana.BridgePda))
	require.Contains(t, message.AccountKeys, solanago.MustPublicKeyFromBase58(cfg.Solana.BridgeProgramId))
	require.Contains(t, message.AccountKeys, bridgeAta)
	require.Contains(t, message.AccountKeys, receiverAta)

	// Verify instructions accounts
	require.Equal(t, 1, len(message.Instructions))
	instruction := message.Instructions[0]
	require.Equal(t, 5, len(message.Instructions[0].Accounts))
	require.Equal(t,
		[]solanago.PublicKey{
			solanago.MustPublicKeyFromBase58(mpcAddress),
			solanago.MustPublicKeyFromBase58(tokenProgramId),
			solanago.MustPublicKeyFromBase58(cfg.Solana.BridgePda),
			bridgeAta,
			receiverAta,
		},
		[]solanago.PublicKey{
			message.AccountKeys[instruction.Accounts[0]],
			message.AccountKeys[instruction.Accounts[1]],
			message.AccountKeys[instruction.Accounts[2]],
			message.AccountKeys[instruction.Accounts[3]],
			message.AccountKeys[instruction.Accounts[4]],
		},
	)

	// Verify instruction data
	transferIn := solanatypes.TransferInData{}
	err = borsh.Deserialize(&transferIn, instruction.Data)
	require.Nil(t, err)
	require.Equal(t, uint64(nonce), transferIn.Nonce)
	require.Equal(t, []uint64{uint64(amountInt * 100_000_000)}, transferIn.Amounts)
}

func TestGetRecentBlockHash(t *testing.T) {
	ctx, k := mockForBridgeTest()
	chain := "solana-devnet"
	// Medium of [0, 3, 5, 7, 8, 8, 9, 10, 100] is 8
	k.SetSolanaConfirmedBlock(ctx, chain, "signer1", "Hash1", 8)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer2", "Hash2", 3)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer3", "Hash3", 5)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer4", "Hash4", 100)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer5", "Hash5", 0)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer6", "Hash6", 10)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer7", "Hash7", 9)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer8", "Hash8", 8)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer9", "Hash9", 7)

	b := NewBridge(chain, "signer", k, config.Config{}).(*defaultBridge)
	hash, err := b.getRecentBlockHash(ctx, chain)
	require.Nil(t, err)
	require.Equal(t, "Hash1", hash)
}

// This test verifies that the ED25519 of solana is compatible with the eddsa curve that Sisu uses.
func TestSigning(t *testing.T) {
	msg := []byte("This is a test message")

	edwardPrivkey, _ := edwards.GeneratePrivateKey()
	require.NotNil(t, edwardPrivkey)
	pubkey := solanago.PublicKeyFromBytes(edwardPrivkey.PubKey().Serialize())

	sig2, err := edwardPrivkey.Sign(msg)
	require.Nil(t, err)
	require.True(t, pubkey.Verify(msg, solanago.SignatureFromBytes(sig2.Serialize())))
}

func TestDeserializeTx(t *testing.T) {
	mnemonic := utils.LOCALHOST_MNEMONIC
	admin := GetSolanaPrivateKey(mnemonic)
	mpcKey := admin.PublicKey()

	chain := "solana-devnet"

	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.SetMpcAddress(ctx, chain, mpcKey.String())
	k.SetMpcNonce(ctx, &types.MpcNonce{Chain: chain, Nonce: 1})
	k.SetSolanaConfirmedBlock(ctx, chain, "signer", "EnzRJ6ojbs5GDEtv4vDRuNnMSJyYKfeD7ATNLwZQLWHe", int64(163_936_646))
	m := k.GetTokens(ctx, []string{"SISU"})
	token := m["SISU"]
	for i, c := range token.Chains {
		if c == chain {
			token.Addresses[i] = "DEfbTuKfeXxXYkXFU6eGgyzDNbTbsAD9U6z7xexK4nUd" // Put your token address here.
		}
	}
	m["SISU"] = token
	k.SetTokens(ctx, m)

	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = "3pqWds7QP82yfxykgrvLszdmkQv6Vb5bukPZzzhYAYec"
	cfg.Solana.BridgePda = "CvocQ9ivbdz5rUnTh6zBgxaiR4asMNbXRrG2VPUYpoau"
	receiverAta := "LkxVSjLH4mjxndDQKrG1a7FYTU7zGFYE5tDzr3PLd3i"

	bridge := NewBridge(chain, "signer", k, cfg).(*defaultBridge)
	tx, err := bridge.getTransaction(
		ctx,
		[]*types.Token{token},
		[]string{receiverAta}, // Receiver account ata here.
		[]*big.Int{new(big.Int).Mul(big.NewInt(3), utils.SisuDecimalBase)},
	)
	if err != nil {
		panic(err)
	}

	tx.Sign(
		func(key solanago.PublicKey) *solanago.PrivateKey {
			if admin.PublicKey().Equals(key) {
				return &admin
			}

			return nil
		},
	)

	bz, err := tx.MarshalBinary()
	require.Nil(t, err)

	decoder := bin.NewBinDecoder(bz)
	decodedTx := solanago.Transaction{}

	err = decodedTx.UnmarshalWithDecoder(decoder)
	require.Nil(t, err)
}

package solana

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	deyestypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

// This files contains test that connect to a real network. They are skipped in CI mode.

func doTestTransferIn(t *testing.T, privateKey solanago.PrivateKey, mpcKey solanago.PublicKey,
	cfg config.Config, tokenAddr, receiverAta string) {
	client, wsClient := GetBasicData("localhost")
	chain := "solana-devnet"

	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.SetMpcAddress(ctx, chain, mpcKey.String())
	k.SetMpcNonce(ctx, &types.MpcNonce{Chain: chain, Nonce: 1})
	m := k.GetTokens(ctx, []string{"SISU"})
	token := m["SISU"]
	for i, c := range token.Chains {
		if c == chain {
			token.Addresses[i] = tokenAddr // Put your token address here.
		}
	}
	m["SISU"] = token
	k.SetTokens(ctx, m)

	deyesCli := &external.MockDeyesClient{
		SolanaQueryRecentBlockFunc: func(chain string) (*deyestypes.SolanaQueryRecentBlockResult, error) {
			result := getRecentBlockHash(client)

			return &deyestypes.SolanaQueryRecentBlockResult{
				Hash: result.Value.Blockhash.String(),
			}, nil
		},
	}

	bridge := NewBridge(chain, "signer", k, cfg, deyesCli).(*defaultBridge)
	tx, err := bridge.getTransaction(
		ctx,
		[]*types.Token{token},
		[]string{receiverAta}, // Receiver account ata here.
		[]*big.Int{new(big.Int).Mul(big.NewInt(3), utils.SisuDecimalBase)},
	)
	if err != nil {
		panic(err)
	}

	require.Equal(t, 1, int(tx.Message.Header.NumRequiredSignatures))

	signTx(tx, privateKey)
	sig, err := confirm.SendAndConfirmTransaction(context.Background(), client, wsClient, tx)
	if err != nil {
		panic(err)
	}
	log.Verbose("Final sig = ", sig)
}

func getRecentBlockHash(client *rpc.Client) *rpc.GetRecentBlockhashResult {
	// Get blockhash
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	return result
}

func signTx(tx *solanago.Transaction, privateKey solanago.PrivateKey) {
	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("unable to encode message for signing: %w", err))
	}
	signerKeys := tx.Message.AccountKeys[0:tx.Message.Header.NumRequiredSignatures]

	signedSignatures := []solanago.Signature{}
	for _, key := range signerKeys {
		if privateKey != nil {
			s, err := privateKey.Sign(messageContent)
			if err != nil {
				panic(fmt.Errorf("failed to signed with key %q: %w", key.String(), err))
			}
			signedSignatures = append(signedSignatures, s)
		}
	}
	tx.Signatures = append(tx.Signatures, signedSignatures...)
}

func doSetSpender(t *testing.T, mnemonic string, cfg config.Config, spenderKey solanago.PublicKey) {
	ownerKey := GetSolanaPrivateKey(mnemonic)
	ix, err := solanatypes.NewAddSpenderIx(cfg.Solana.BridgeProgramId, ownerKey.PublicKey().String(),
		cfg.Solana.BridgePda, spenderKey.String())
	if err != nil {
		panic(err)
	}

	client, wsClient := GetBasicData("localhost")
	err = SignAndSubmit([]*rpc.Client{client}, []*ws.Client{wsClient}, []solanago.Instruction{ix}, ownerKey)
	if err != nil {
		panic(err)
	}
}

// Set the mnemonic to run this test.
// MNEMONIC=YOUR_MNEMONIC go test -v -run TestTransferIn
func TestTransferIn(t *testing.T) {
	t.Skip()
	mnemonic := os.Getenv("MNEMONIC") // use your mnemonic here or pass it from the environment
	admin := GetSolanaPrivateKey(mnemonic)
	mpcKey := admin.PublicKey()

	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = "3pqWds7QP82yfxykgrvLszdmkQv6Vb5bukPZzzhYAYec" // Use your program id here
	cfg.Solana.BridgePda = "CvocQ9ivbdz5rUnTh6zBgxaiR4asMNbXRrG2VPUYpoau"       // Bridge pda
	token := "DEfbTuKfeXxXYkXFU6eGgyzDNbTbsAD9U6z7xexK4nUd"
	receiverAta := "LkxVSjLH4mjxndDQKrG1a7FYTU7zGFYE5tDzr3PLd3i"
	doTestTransferIn(t, admin, mpcKey, cfg, token, receiverAta)
}

func TestSetSpenderAndTransferIn(t *testing.T) {
	t.Skip()
	mnemonic := os.Getenv("MNEMONIC") // use your mnemonic here or pass it from the environment

	// private key of pubkey H4MctVS4MUteTAmLLUZfheCy9voAUngCY7zVqYKJQStG. Rember to fund this
	// address before running this test.
	bz, err := hex.DecodeString("3c3a2b82283b7691f3a3b9c11507559c592209a6cbac8416f9cbf2f8d1ed202eee970c04d8411b751449bc231eaf2748756c9c9331273590df2768897088f0c9")
	require.Nil(t, err)
	privKey := solanago.PrivateKey(bz)

	// privKey := GetSolanaPrivateKey(mnemonic)

	cfg := config.Config{}
	cfg.Solana.BridgeProgramId = "3pqWds7QP82yfxykgrvLszdmkQv6Vb5bukPZzzhYAYec" // Use your program id here
	cfg.Solana.BridgePda = "CvocQ9ivbdz5rUnTh6zBgxaiR4asMNbXRrG2VPUYpoau"       // Bridge pda
	token := "DEfbTuKfeXxXYkXFU6eGgyzDNbTbsAD9U6z7xexK4nUd"
	receiverAta := "LkxVSjLH4mjxndDQKrG1a7FYTU7zGFYE5tDzr3PLd3i"

	spenderPubkey := privKey.PublicKey()
	doSetSpender(t, mnemonic, cfg, spenderPubkey)
	doTestTransferIn(t, privKey, spenderPubkey, cfg, token, receiverAta)
}

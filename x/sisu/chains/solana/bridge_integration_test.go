package solana

import (
	"context"
	"math/big"
	"os"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// This files contains test that connect to a real network. They are skipped in CI mode.

// Set the mnemonic to run this test.
// MNEMONIC=YOUR_MNEMONIC go test -v -run TestTransferIn
func TestTransferIn(t *testing.T) {
	t.Skip()
	mnemonic := os.Getenv("MNEMONIC") // use your mnemonic here or pass it from the environment
	feePayer := GetSolanaPrivateKey(mnemonic)
	mpcKey := feePayer.PublicKey()
	client, wsClient := GetBasicData("localhost")
	chain := "solana-devnet"

	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.SetMpcAddress(ctx, chain, mpcKey.String())
	k.SetMpcNonce(ctx, &types.MpcNonce{Chain: chain, Nonce: 1})
	recentHash := getRecentBlockHash(client)
	k.SetSolanaConfirmedBlock(ctx, chain, "signer", recentHash.Value.Blockhash.String(), int64(recentHash.Context.Slot))
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
	cfg.Solana.BridgeProgramId = "3pqWds7QP82yfxykgrvLszdmkQv6Vb5bukPZzzhYAYec" // Use your program id here
	cfg.Solana.BridgePda = "CvocQ9ivbdz5rUnTh6zBgxaiR4asMNbXRrG2VPUYpoau"

	bridge := NewBridge(chain, "signer", k, cfg).(*defaultBridge)
	tx, err := bridge.getTransaction(
		ctx,
		[]*types.Token{token},
		[]string{"LkxVSjLH4mjxndDQKrG1a7FYTU7zGFYE5tDzr3PLd3i"}, // Receiver account ata here.
		[]*big.Int{new(big.Int).Mul(big.NewInt(3), utils.SisuDecimalBase)},
	)
	if err != nil {
		panic(err)
	}

	signTx(tx, mnemonic)
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

func signTx(tx *solanago.Transaction, mnemonic string) {
	feePayer := GetSolanaPrivateKey(mnemonic)

	tx.Sign(
		func(key solanago.PublicKey) *solanago.PrivateKey {
			if feePayer.PublicKey().Equals(key) {
				return &feePayer
			}

			return nil
		},
	)
}

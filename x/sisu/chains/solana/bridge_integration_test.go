package solana

import (
	"context"
	"math/big"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// This files contains test that connect to a real network. They are skipped in CI mode.
func TestTransferIn(t *testing.T) {
	mnemonic := "" // use your mnemonic here or pass it from the environment
	client, wsClient := GetBasicData("devnet")
	chain := "solana-devenet"

	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	m := k.GetTokens(ctx, []string{"SISU"})
	token := m["SISU"]
	for i, c := range token.Chains {
		if c == chain {
			token.Addresses[i] = "DEfbTuKfeXxXYkXFU6eGgyzDNbTbsAD9U6z7xexK4nUd" // Put your address here.
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
		[]string{"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"}, // Use account ata here.
		[]*big.Int{big.NewInt(3 * 100_000_000)},
	)
	if err != nil {
		panic(err)
	}

	signTx(tx, mnemonic)
	sig, err := confirm.SendAndConfirmTransaction(context.Background(), client, wsClient, tx)
	log.Verbose("Final sig = ", sig)
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

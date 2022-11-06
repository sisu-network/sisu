package dev

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/cosmos/go-bip39"
	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
)

func GetSolanaPrivateKey(mnemonic string) solanago.PrivateKey {
	seed := bip39.NewSeed(mnemonic, "")[:32]
	key := ed25519.NewKeyFromSeed(seed)
	privKey := solanago.PrivateKey(key)

	return privKey
}

func (c *fundAccountCmd) fundSolana(genesisFolder, mnemonic string, allPubKeys map[string][]byte) {
	privateKey := GetSolanaPrivateKey(mnemonic)
	faucet := privateKey.PublicKey()

	// Get Mpc address
	bz := allPubKeys[libchain.KEY_TYPE_EDDSA]
	mpcAddress := utils.GetSolanaAddressFromPubkey(bz)
	mpcAccount := solanago.MustPublicKeyFromBase58(mpcAddress)

	// Get all tokens
	tokens := helper.GetTokens(filepath.Join(genesisFolder, "tokens.json"))
	solanaConfig, err := config.ReadSolanaConfig(filepath.Join(genesisFolder, "solana_config.json"))
	if err != nil {
		panic(err)
	}

	// Create client & ws connector
	client := rpc.New(solanaConfig.Rpc)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), solanaConfig.Ws)
	if err != nil {
		panic(err)
	}

	// Get all ATA address created from mpc address and token address
	for _, token := range tokens {
		for i := range token.Chains {
			if len(token.Addresses[i]) == 0 {
				continue
			}

			if libchain.IsSolanaChain(token.Chains[i]) {
				tokentMintPubKey := solanago.MustPublicKeyFromBase58(token.Addresses[i])

				// Create source ata
				sourceAta, err := c.createAssociatedAccount(client, wsClient, mnemonic, faucet, tokentMintPubKey)
				if err != nil {
					panic(err)
				}

				// Create mpc ata
				mpcAta, err := c.createAssociatedAccount(client, wsClient, mnemonic, mpcAccount, tokentMintPubKey)
				if err != nil {
					panic(err)
				}

				// Fund the address
				c.transferSolanaToken(client, wsClient, mnemonic, token.Addresses[i], sourceAta.String(), mpcAta.String())
			}
		}
	}
}

func (c *fundAccountCmd) transferSolanaToken(client *rpc.Client, wsClient *ws.Client, mnemonic, token,
	sourceAta, dstAta string) {
	feePayer := GetSolanaPrivateKey(mnemonic)
	feePayerPubkey := feePayer.PublicKey()

	log.Verbosef("Funding token = %s, source = %s, destination = %s\n", token, sourceAta, dstAta)

	// Get blockhash
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	ix := solanatypes.NewTransferTokenIx(
		solanago.MustPublicKeyFromBase58(sourceAta),
		solanago.MustPublicKeyFromBase58(token),
		solanago.MustPublicKeyFromBase58(dstAta),
		feePayerPubkey,
		big.NewInt(100),
		8,
	)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		result.Value.Blockhash,
		solana.TransactionPayer(feePayerPubkey),
	)
	if err != nil {
		panic(err)
	}

	err = c.signAndSubmit(client, wsClient, tx, feePayer)
	if err != nil {
		panic(err)
	}
}

// signAndSubmit signs a transaction and submits it to the network.
func (c *fundAccountCmd) signAndSubmit(client *rpc.Client, wsClient *ws.Client,
	tx *solanago.Transaction, feePayer solanago.PrivateKey) error {
	tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if feePayer.PublicKey().Equals(key) {
				return &feePayer
			}

			return nil
		},
	)

	// Send transaction, and wait for confirmation
	sig, err := confirm.SendAndConfirmTransaction(
		context.Background(),
		client,
		wsClient,
		tx,
	)
	log.Verbose("sig = ", sig)

	return err
}

func (c *fundAccountCmd) createAssociatedAccount(client *rpc.Client, wsClient *ws.Client, mnemonic string,
	owner, tokenMint solanago.PublicKey) (solanago.PublicKey, error) {
	privateKey := GetSolanaPrivateKey(mnemonic)
	feePayer := privateKey.PublicKey()

	// Check if the ata account existed. If not create a new one.
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(
		owner,
		tokenMint,
	)

	if err != nil {
		panic(err)
	}

	_, err = querySolanaAccountBalance(client, ownerAta.String())
	if err == nil {
		// Account already existed, do nothing
		fmt.Printf("Accounts %s has been created\n", ownerAta.String())
		return ownerAta, nil
	}

	log.Verbosef("Creating new ATA account, owner = %s, ownerAta = %s, tokenMint = %s", owner.String(),
		ownerAta.String(), tokenMint.String())

	// Create a new account
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	ix := solanatypes.NewCreateAssociatedAccountIx(feePayer, owner, ownerAta, tokenMint)
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		result.Value.Blockhash,
		solana.TransactionPayer(feePayer),
	)
	if err != nil {
		panic(err)
	}

	return ownerAta, c.signAndSubmit(client, wsClient, tx, privateKey)
}

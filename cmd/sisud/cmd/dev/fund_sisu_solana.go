package dev

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
)

func (c *fundAccountCmd) fundSolana(genesisFolder, mnemonic string, mpcPubKey []byte) {
	privateKey := solana.GetSolanaPrivateKey(mnemonic)
	faucet := privateKey.PublicKey()

	// Get all tokens
	tokens := helper.GetTokens(filepath.Join(genesisFolder, "tokens.json"))
	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
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
		if len(token.Addresses) == 0 {
			continue
		}

		for i := range token.Chains {
			if len(token.Addresses[i]) == 0 {
				continue
			}

			if libchain.IsSolanaChain(token.Chains[i]) {
				decimals := token.GetDecimalsForChain(solanaConfig.Chain)
				tokentMintPubKey := solanago.MustPublicKeyFromBase58(token.Addresses[i])

				// Create source ata
				sourceAta, created, err := c.createAssociatedAccount(client, wsClient, mnemonic, faucet, tokentMintPubKey)
				if err != nil {
					panic(err)
				}

				// Mint token for the source if needed.
				fmt.Println("Minting token ", tokentMintPubKey, " to ", sourceAta)
				err = c.mintToken(client, wsClient, mnemonic, created, tokentMintPubKey, byte(token.Decimals[i]),
					sourceAta, 1_000_000*100_000_000)
				if err != nil {
					panic(err)
				}

				// Create bridge ata
				bridgePda := solanago.MustPublicKeyFromBase58(solanaConfig.BridgePda)
				bridgeAta, _, err := c.createAssociatedAccount(client, wsClient, mnemonic, bridgePda, tokentMintPubKey)
				if err != nil {
					panic(err)
				}

				// Fund the address
				c.transferSolanaToken(client, wsClient, mnemonic, token.Addresses[i],
					byte(decimals), sourceAta.String(), bridgeAta.String(), 10_000*100_000_000)

				// Set the spender for the vault.
				c.setSpender(client, wsClient, genesisFolder, mnemonic,
					utils.GetSolanaAddressFromPubkey(mpcPubKey))
			}
		}
	}
}

func (c *fundAccountCmd) transferSolanaToken(client *rpc.Client, wsClient *ws.Client, mnemonic,
	token string, tokenDecimals byte, sourceAta, receiverAta string, amount uint64) {
	balance, err := solana.QuerySolanaAccountBalance(client, receiverAta)
	if balance.Cmp(big.NewInt(int64(amount))) >= 0 {
		log.Verbose("Account already has enough balance. Skip transferring")
		return
	}

	feePayer := solana.GetSolanaPrivateKey(mnemonic)
	feePayerPubkey := feePayer.PublicKey()

	log.Verbosef("Funding token = %s, source = %s, destination = %s\n", token, sourceAta, receiverAta)

	ix := solanatypes.NewTransferTokenIx(
		solanago.MustPublicKeyFromBase58(sourceAta),
		solanago.MustPublicKeyFromBase58(token),
		solanago.MustPublicKeyFromBase58(receiverAta),
		feePayerPubkey,
		amount,
		tokenDecimals,
	)

	err = solana.SignAndSubmit(client, wsClient, []solanago.Instruction{ix}, feePayer)
	if err != nil {
		panic(err)
	}
}

func (c *fundAccountCmd) createAssociatedAccount(client *rpc.Client, wsClient *ws.Client, mnemonic string,
	owner, tokenMint solanago.PublicKey) (solanago.PublicKey, bool, error) {
	privateKey := solana.GetSolanaPrivateKey(mnemonic)
	feePayer := privateKey.PublicKey()

	// Check if the ata account existed. If not create a new one.
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(
		owner,
		tokenMint,
	)

	if err != nil {
		panic(err)
	}

	_, err = solana.QuerySolanaAccountBalance(client, ownerAta.String())
	if err == nil {
		// Account already existed, do nothing
		log.Verbosef("Accounts %s has been created", ownerAta.String())
		return ownerAta, false, nil
	}

	log.Verbosef("Creating new ATA account, owner = %s, ownerAta = %s, tokenMint = %s", owner.String(),
		ownerAta.String(), tokenMint.String())

	// Create a new account
	ix := solanatypes.NewCreateAssociatedAccountIx(feePayer, owner, ownerAta, tokenMint)

	return ownerAta, true, solana.SignAndSubmit(client, wsClient, []solanago.Instruction{ix}, privateKey)
}

func (c *fundAccountCmd) mintToken(client *rpc.Client, wsClient *ws.Client, mnemonic string,
	newAccount bool, tokenMint solanago.PublicKey, tokenDecimals byte, receiverAta solanago.PublicKey,
	amount uint64) error {
	// Check if we need to mint token for this account.
	shouldMint := false
	if newAccount {
		shouldMint = true
	} else {
		balance, err := solana.QuerySolanaAccountBalance(client, receiverAta.String())
		if err != nil {
			panic(err)
		}
		if balance.Cmp(big.NewInt(int64(amount/2))) < 0 {
			shouldMint = true
		}
	}

	if !shouldMint {
		return nil
	}

	privateKey := solana.GetSolanaPrivateKey(mnemonic)
	owner := privateKey.PublicKey()

	mintTokenIx := solanatypes.NewMintTokenIx(
		tokenMint,
		receiverAta,
		owner,
		tokenDecimals,
		amount,
	)

	return solana.SignAndSubmitWithOptions(client, wsClient, []solanago.Instruction{mintTokenIx},
		privateKey, rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		})
}

func (c *fundAccountCmd) setSpender(client *rpc.Client, wsClient *ws.Client, genesisFolder,
	mnemonic string, mpc string) {
	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}

	ownerKey := solana.GetSolanaPrivateKey(mnemonic)
	ix, err := solanatypes.NewAddSpenderIx(solanaConfig.BridgeProgramId, ownerKey.PublicKey().String(),
		solanaConfig.BridgePda, mpc)
	if err != nil {
		panic(err)
	}

	err = solana.SignAndSubmit(client, wsClient, []solanago.Instruction{ix}, ownerKey)
	if err != nil {
		panic(err)
	}
}

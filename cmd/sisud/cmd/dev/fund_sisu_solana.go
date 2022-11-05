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
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
)

func GetSolanaPrivateKey(mnemonic string) solanago.PrivateKey {
	seed := bip39.NewSeed(mnemonic, "")[:32]
	key := ed25519.NewKeyFromSeed(seed)
	privKey := solanago.PrivateKey(key)

	return privKey
}

func (c *fundAccountCmd) fundOnSolana(genesisFolder, mnemonic, sisuRpc string) {
	// privateKey := GetSolanaPrivateKey(mnemonic)
	// funderPubkey := privateKey.PublicKey()
	// funderPubkey.Address

	allPubKeys := queryPubKeys(context.Background(), sisuRpc)
	bz := allPubKeys[libchain.KEY_TYPE_EDDSA]

	// Get Mpc address
	mpcAddress := utils.GetSolanaAddressFromPubkey(bz)

	// Get all tokens
	tokens := helper.GetTokens(filepath.Join(genesisFolder, "tokens.json"))

	// Get all ATA address created from mpc address and token address
	for _, token := range tokens {
		for i := range token.Chains {
			if libchain.IsSolanaChain(token.Chains[i]) {
				// Get ATA from mpc and token address.
				solanago.FindAssociatedTokenAddress(
					solanago.MustPublicKeyFromBase58(mpcAddress),
					solanago.MustPublicKeyFromBase58(token.Addresses[i]),
				)

				// Fund the address

			}
		}
	}
}

func (c *fundAccountCmd) transferSolanaToken(client *rpc.Client, wsClient *ws.Client, mnemonic, token,
	sourceAta, dstAta string, recentBlockHash solanago.Hash) {
	feePayer := GetSolanaPrivateKey(mnemonic)
	feePayerPubkey := feePayer.PublicKey()

	// This is the key source code in JS.
	// 	const keys = addSigners(
	// 		[
	// 				{ pubkey: source, isSigner: false, isWritable: true },
	// 				{ pubkey: mint, isSigner: false, isWritable: false },
	// 				{ pubkey: destination, isSigner: false, isWritable: true },
	// 		],
	// 		owner,
	// 		multiSigners
	// );

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(solanago.MustPublicKeyFromBase58(sourceAta), true, false),
		solana.NewAccountMeta(solanago.MustPublicKeyFromBase58(token), false, false),
		solana.NewAccountMeta(solanago.MustPublicKeyFromBase58(dstAta), true, false),
		solana.NewAccountMeta(feePayerPubkey, false, true),
	}

	ix := solanatypes.NewTransferTokenIx(accounts, big.NewInt(100), 8)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recentBlockHash,
		solana.TransactionPayer(feePayerPubkey),
	)
	if err != nil {
		panic(err)
	}

	tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if feePayer.PublicKey().Equals(key) {
				return &feePayer
			}

			fmt.Println("Private key is nil for ", key)
			return nil
		},
	)

	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.Background(),
		client,
		wsClient,
		tx,
	)

	fmt.Println("sig = ", sig)
	if err != nil {
		panic(err)
	}
}

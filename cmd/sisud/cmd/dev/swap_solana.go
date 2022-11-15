package dev

import (
	"context"
	"path/filepath"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
)

func (c *swapCommand) swapFromSolana(genesisFolder, chain, mnemonic, tokenAddr, recipient string,
	dstChain uint64, amount uint64) {
	feePayer := GetSolanaPrivateKey(mnemonic)

	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}

	client := rpc.New(solanaConfig.Rpc)
	wsClient, err := ws.Connect(context.Background(), solanaConfig.Ws)
	if err != nil {
		panic(err)
	}

	approveIx := c.approveSolanaIx(genesisFolder, chain, mnemonic, tokenAddr, amount)
	transferIx := c.transferTokenIx(genesisFolder, mnemonic, tokenAddr, recipient, dstChain, amount)

	err = solana.SignAndSubmit(client, wsClient, []solanago.Instruction{approveIx, transferIx}, feePayer)
	if err != nil {
		panic(err)
	}
}

func (c *swapCommand) approveSolanaIx(genesisFolder, chain, mnemonic, tokenAddr string, amount uint64) solanago.Instruction {
	tokenMintPubkey := solanago.MustPublicKeyFromBase58(tokenAddr)

	ownerPrivKey := GetSolanaPrivateKey(mnemonic)
	ownerPubkey := ownerPrivKey.PublicKey()
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(ownerPubkey, tokenMintPubkey)
	if err != nil {
		panic(err)
	}

	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}
	bridgePda := solanago.MustPublicKeyFromBase58(solanaConfig.BridgePda)

	// Get token config
	var decimal byte
	tokens := helper.GetTokens(filepath.Join(genesisFolder, "tokens.json"))
	for _, token := range tokens {
		for j, c := range token.Chains {
			if c == chain && token.Addresses[j] == tokenAddr {
				decimal = token.Decimals[j]
			}
		}
	}

	if decimal == 0 {
		panic("Invalid decimals")
	}

	ix := solanatypes.NewApproveCheckedIx(ownerPubkey, ownerAta, tokenMintPubkey, bridgePda, amount,
		decimal)

	return ix
}

func (c *swapCommand) transferTokenIx(genesisFolder, mnemonic, tokenAddr, recipient string, dstChainId, amount uint64) solanago.Instruction {
	tokenMintPubkey := solanago.MustPublicKeyFromBase58(tokenAddr)

	ownerPrivKey := GetSolanaPrivateKey(mnemonic)
	ownerPubkey := ownerPrivKey.PublicKey()
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(ownerPubkey, tokenMintPubkey)
	if err != nil {
		panic(err)
	}

	solanaConfig, err := helper.ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}
	bridgeProgramId := solanago.MustPublicKeyFromBase58(solanaConfig.BridgeProgramId)

	bridgePda := solanago.MustPublicKeyFromBase58(solanaConfig.BridgePda)
	bridgeAta, _, err := solanago.FindAssociatedTokenAddress(bridgePda, tokenMintPubkey)
	if err != nil {
		panic(err)
	}

	data := solanatypes.NewTransferOutData(amount, tokenAddr, dstChainId, recipient)

	return solanatypes.NewTransferOutInstruction(bridgeProgramId, ownerPubkey, ownerAta, bridgeAta, bridgePda, data)
}

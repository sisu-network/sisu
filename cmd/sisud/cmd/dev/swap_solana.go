package dev

import (
	"math/big"
	"path/filepath"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/sisu-network/sisu/config"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
)

func (c *swapCommand) swapFromSolana(genesisFolder, chain, mnemonic, tokenAddr string, amount *big.Int, allPubKeys map[string][]byte) {
	// approveIx := c.approveSolanaIx(genesisFolder, chain, mnemonic, tokenAddr, amount.Uint64())
}

func (c *swapCommand) approveSolanaIx(genesisFolder, chain, mnemonic, tokenAddr string, amount uint64) solanago.Instruction {
	tokenMintPubkey := solanago.MustPublicKeyFromBase58(tokenAddr)

	ownerPrivKey := GetSolanaPrivateKey(mnemonic)
	ownerPubkey := ownerPrivKey.PublicKey()
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(ownerPubkey, tokenMintPubkey)
	if err != nil {
		panic(err)
	}

	solanaConfig, err := config.ReadSolanaConfig(filepath.Join(genesisFolder, "solana_config.json"))
	if err != nil {
		panic(err)
	}
	bridgePda := solanago.MustPublicKeyFromBase58(solanaConfig.BridgePda)
	// ix := solanatypes.NewApproveCheckedIx(ownerAta, tokenMintPubkey, bridgePda, amount,
	// 	byte(token.Decimals))

	// TODO: Use custom decimals for the token in solana chain.
	ix := solanatypes.NewApproveCheckedIx(ownerPubkey, ownerAta, tokenMintPubkey, bridgePda, amount,
		byte(8))

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

	solanaConfig, err := config.ReadSolanaConfig(filepath.Join(genesisFolder, "solana_config.json"))
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

package dev

import (
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
)

func TestApproveToken(t *testing.T) {
	t.Skip()

	mnemonic, client, wsClient := getBasicData("localhost")
	feePayer := GetSolanaPrivateKey(mnemonic)

	cmd := &swapCommand{}
	ix := cmd.approveSolanaIx("../../../../misc/test", "solana-devnet", mnemonic, "AJdUMt177iQ19J63ybkXtUVD6sK8dxD5ibietQANuv9S", 1000)

	err := solana.SignAndSubmit(client, wsClient, []solanago.Instruction{ix}, feePayer)
	if err != nil {
		panic(err)
	}
}

func TestTransferOut(t *testing.T) {
	mnemonic, client, wsClient := getBasicData("devnet")
	feePayer := GetSolanaPrivateKey(mnemonic)

	cmd := &swapCommand{}

	ix := cmd.transferTokenIx("../../../../misc/test", mnemonic, "AJdUMt177iQ19J63ybkXtUVD6sK8dxD5ibietQANuv9S", "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988", 189985, 300)
	err := solana.SignAndSubmit(client, wsClient, []solanago.Instruction{ix}, feePayer)
	if err != nil {
		panic(err)
	}
}

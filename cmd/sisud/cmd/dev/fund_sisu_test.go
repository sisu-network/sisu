package dev

import (
	"context"
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/near/borsh-go"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/stretchr/testify/require"
)

func TestQueryPubKeys(t *testing.T) {
	m := queryPubKeys(context.Background(), "0.0.0.0:9090")

	for chain, key := range m {
		fmt.Println("chain, ", chain, len(key))
	}
}

func TestSerializeTransferIxData(t *testing.T) {
	data := solanatypes.TransferSplTokenData{
		Instruction: 12,
		Amount:      100,
		Decimals:    8,
	}

	bz, err := borsh.Serialize(data)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 10, len(bz))
}

func TestGetSolanaPrivateKey(t *testing.T) {
	mnemonic := utils.LOCALHOST_MNEMONIC
	privateKey := GetSolanaPrivateKey(mnemonic)

	require.Equal(t, "Cy4RyK92aQHuaPgw6PdSYJ5GbcAw9uL8fTPawEtZwiWw", privateKey.PublicKey().String())
}

func TestTransferToken(t *testing.T) {
	t.Skip()

	mnemonic := utils.LOCALHOST_MNEMONIC

	// Transfer token
	cmd := &fundAccountCmd{}

	endpoint := rpc.LocalNet_RPC
	client := rpc.New(endpoint)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), rpc.LocalNet_WS)
	if err != nil {
		panic(err)
	}

	// Block hash
	result, err := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	tokenMintPubkey := "8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX"
	srcAta := "BPRyt1DwNCzMpbnMkzxbkj1A6sNRN5KP8Ej4iGeudtLm"
	dstAta := "BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a"

	cmd.transferSolanaToken(client, wsClient, mnemonic, tokenMintPubkey, srcAta, dstAta, result.Value.Blockhash)
}

package dev

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/near/borsh-go"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/chains/solana"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/stretchr/testify/require"
)

var TokenMintPubkey = solanago.MustPublicKeyFromBase58("AJdUMt177iQ19J63ybkXtUVD6sK8dxD5ibietQANuv9S")

func getBasicData(network string) (string, *rpc.Client, *ws.Client) {
	var rpcEndpoint string
	var wssEndpoint string

	if network == "devnet" {
		rpcEndpoint = rpc.DevNet_RPC
		wssEndpoint = rpc.DevNet.WS
	} else {
		rpcEndpoint = rpc.LocalNet_RPC
		wssEndpoint = rpc.LocalNet_WS
	}

	client := rpc.New(rpcEndpoint)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), wssEndpoint)
	if err != nil {
		panic(err)
	}

	return utils.LOCALHOST_MNEMONIC, client, wsClient
}

func TestQueryPubKeys(t *testing.T) {
	t.Skip()

	queryPubKeys(context.Background(), "0.0.0.0:9090")
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
	privateKey := solana.GetSolanaPrivateKey(mnemonic)

	require.Equal(t, "Cy4RyK92aQHuaPgw6PdSYJ5GbcAw9uL8fTPawEtZwiWw", privateKey.PublicKey().String())
}

// Sanity check on localhost. Disabled by default. Enable if you want to debug the fund command.
func TestTransferToken(t *testing.T) {
	t.Skip()

	// Transfer token
	mnemonic, client, wsClient := getBasicData("localhost")

	srcAta := "BPRyt1DwNCzMpbnMkzxbkj1A6sNRN5KP8Ej4iGeudtLm"
	dstAta := "BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a"

	transferSolanaToken(client, wsClient, mnemonic, TokenMintPubkey.String(), 8, srcAta, dstAta, 1000)
}

// Sanity check on localhost. Disabled by default. Enable if you want to debug the fund command.
func TestFundOnSolana(t *testing.T) {
	t.Skip()

	mnemonic := os.Getenv("MNEMONIC")
	cmd := &fundAccountCmd{}
	cmd.fundSolana("../../../../misc/test", mnemonic, utils.RandomBytes(32))
}

func TestCreateAssociatedProgram(t *testing.T) {
	t.Skip()

	mnemonic, client, wsClient := getBasicData("localhost")

	// Generate a random private key
	privKey, err := solanago.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	ownerPubkey := privKey.PublicKey()
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(ownerPubkey, TokenMintPubkey)

	// Query owner ata. This should return error
	_, err = solana.QuerySolanaAccountBalance(client, ownerAta.String())
	require.True(t, strings.Contains(err.Error(), "could not find account"))

	createSolanaAta(client, wsClient, mnemonic, ownerPubkey, TokenMintPubkey)

	// Query account ata
	balance, err := solana.QuerySolanaAccountBalance(client, ownerAta.String())
	require.Nil(t, err)
	require.Equal(t, big.NewInt(0), balance)
}

func TestMintSolanaToken(t *testing.T) {
	t.Skip()

	cmd := &fundAccountCmd{}
	mnemonic, client, wsClient := getBasicData("devnet")

	err := cmd.mintToken(client, wsClient, mnemonic, true, TokenMintPubkey, 8,
		solanago.MustPublicKeyFromBase58("B6hYN4gKqN5iRThXuP7rAX4aH7vY2HnJLFLvsGqiuiKd"),
		100*100_000_000,
	)

	if err != nil {
		panic(err)
	}
}

func TestSolanaSetSpender(t *testing.T) {
	t.Skip()
	mnemonic := os.Getenv("MNEMONIC")
	cmd := &fundAccountCmd{}

	_, client, wsClient := getBasicData("localhost")

	pubkey := []byte{78, 114, 255, 58, 70, 231, 143, 6, 154, 69, 54, 90, 87, 89, 180, 208, 71, 88,
		209, 74, 207, 217, 103, 218, 227, 238, 151, 136, 200, 253, 217, 17}
	mockMpcAddr := utils.GetSolanaAddressFromPubkey(pubkey)

	cmd.setSpender(client, wsClient, "../../../../misc/test", mnemonic, mockMpcAddr)
}

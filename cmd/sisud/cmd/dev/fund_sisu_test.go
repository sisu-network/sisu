package dev

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/decred/dcrd/dcrec/edwards/v2"
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
	cmd := &fundAccountCmd{}
	mnemonic, client, wsClient := getBasicData("localhost")

	srcAta := "BPRyt1DwNCzMpbnMkzxbkj1A6sNRN5KP8Ej4iGeudtLm"
	dstAta := "BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a"

	cmd.transferSolanaToken(client, wsClient, mnemonic, TokenMintPubkey.String(), 8, srcAta, dstAta, 1000)
}

// Sanity check on localhost. Disabled by default. Enable if you want to debug the fund command.
func TestFundOnSolana(t *testing.T) {
	t.Skip()
	// This is the code to generate a new private key
	// privateKey, err := edwards.GeneratePrivateKey()
	// require.Nil(t, err)

	bz, err := hex.DecodeString("00c5fb9d911b4cb3adf209bfa532e3004692c888f71bc6857095ba6674dc2d7b")
	require.Nil(t, err)
	privateKey, _, err := edwards.PrivKeyFromScalar(bz)
	require.NotNil(t, privateKey)
	require.Nil(t, err)

	cmd := &fundAccountCmd{}
	cmd.fundSolana("../../../../misc/test", utils.LOCALHOST_MNEMONIC)
}

func TestCreateAssociatedProgram(t *testing.T) {
	t.Skip()

	cmd := &fundAccountCmd{}
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

	cmd.createAssociatedAccount(client, wsClient, mnemonic, ownerPubkey, TokenMintPubkey)

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

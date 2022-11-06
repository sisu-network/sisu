package dev

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/decred/dcrd/dcrec/edwards/v2"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/near/borsh-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/utils"
	solanatypes "github.com/sisu-network/sisu/x/sisu/chains/solana/types"
	"github.com/stretchr/testify/require"
)

func getBasicData() (string, *rpc.Client, *ws.Client) {
	mnemonic := utils.LOCALHOST_MNEMONIC

	endpoint := rpc.LocalNet_RPC
	client := rpc.New(endpoint)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), rpc.LocalNet_WS)
	if err != nil {
		panic(err)
	}

	return mnemonic, client, wsClient
}

func TestQueryPubKeys(t *testing.T) {
	t.Skip()

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

// Sanity check on localhost. Disabled by default. Enable if you want to debug the fund command.
func TestTransferToken(t *testing.T) {
	t.Skip()

	// Transfer token
	cmd := &fundAccountCmd{}
	mnemonic, client, wsClient := getBasicData()

	tokenMintPubkey := "8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX"
	srcAta := "BPRyt1DwNCzMpbnMkzxbkj1A6sNRN5KP8Ej4iGeudtLm"
	dstAta := "BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a"

	cmd.transferSolanaToken(client, wsClient, mnemonic, tokenMintPubkey, srcAta, dstAta)
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

	pubKeyBytes := privateKey.PubKey().Serialize()

	allPubkey := map[string][]byte{
		libchain.KEY_TYPE_EDDSA: pubKeyBytes,
	}

	cmd := &fundAccountCmd{}
	cmd.fundSolana("../../../../misc/test", utils.LOCALHOST_MNEMONIC, allPubkey)
}

func TestCreateAssociatedProgram(t *testing.T) {
	t.Skip()

	cmd := &fundAccountCmd{}
	mnemonic, client, wsClient := getBasicData()

	tokenMintPubkey := solanago.MustPublicKeyFromBase58("8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX")

	// Generate a random private key
	privKey, err := solanago.NewRandomPrivateKey()
	if err != nil {
		panic(err)
	}

	ownerPubkey := privKey.PublicKey()
	ownerAta, _, err := solanago.FindAssociatedTokenAddress(ownerPubkey, tokenMintPubkey)

	// Query owner ata. This should return error
	_, err = querySolanaAccountBalance(client, ownerAta.String())
	require.True(t, strings.Contains(err.Error(), "could not find account"))

	cmd.createAssociatedAccount(client, wsClient, mnemonic, ownerPubkey, tokenMintPubkey)

	// Query account ata
	balance, err := querySolanaAccountBalance(client, ownerAta.String())
	require.Nil(t, err)
	require.Equal(t, big.NewInt(0), balance)
}

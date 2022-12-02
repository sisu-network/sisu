package solana

import (
	"context"
	"crypto/ed25519"

	"github.com/cosmos/go-bip39"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func GetBasicData(network string) (*rpc.Client, *ws.Client) {
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

	return client, wsClient
}

func GetSolanaPrivateKey(mnemonic string) solanago.PrivateKey {
	seed := bip39.NewSeed(mnemonic, "")[:32]
	key := ed25519.NewKeyFromSeed(seed)
	privKey := solanago.PrivateKey(key)

	return privKey
}

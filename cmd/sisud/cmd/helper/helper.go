package helper

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func GetDevPrivateKey() *ecdsa.PrivateKey {
	// This is the private key for account 0xbeF23B2AC7857748fEA1f499BE8227c5fD07E70c
	bz, err := hex.DecodeString("9f575b88940d452da46a6ceec06a108fcd5863885524aec7fb0bc4906eb63ab1")
	if err != nil {
		panic(err)
	}

	privateKey, err := ethcrypto.ToECDSA(bz)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func GetTokens(file string) []*types.Token {
	tokens := []*types.Token{}

	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &tokens); err != nil {
		panic(err)
	}

	return tokens
}

func GetSolanaClientAndWss(genesisFolder string) ([]*rpc.Client, []*ws.Client) {
	// Create client & ws connector
	clients := make([]*rpc.Client, 0)
	wsClients := make([]*ws.Client, 0)

	// Read RPCs from deyes_chains
	deyesChains := ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))
	for _, chain := range deyesChains {
		if libchain.IsSolanaChain(chain.Chain) {
			for i, rpcUrl := range chain.Rpcs {
				wsClient, err := ws.Connect(context.Background(), chain.Wss[i])
				if err != nil {
					continue
				}

				clients = append(clients, rpc.New(rpcUrl))
				wsClients = append(wsClients, wsClient)
			}
		}
	}

	if len(clients) == 0 {
		panic(fmt.Errorf("Cannot find config for solaan chain, genesis folder = %s", genesisFolder))
	}

	return clients, wsClients
}

func GetEthClient(genesisFolder string, chain string) *ethclient.Client {
	deyesChains := ReadDeyesChainConfigs(filepath.Join(genesisFolder, "deyes_chains.json"))
	for _, cfg := range deyesChains {
		if cfg.Chain == chain {
			client, err := ethclient.Dial(cfg.Rpcs[0])
			if err != nil {
				panic(err)
			}
			return client
		}
	}

	panic(fmt.Errorf("Cannot find chain %s in the config", chain))
}

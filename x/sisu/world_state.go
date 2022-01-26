package sisu

import (
	"github.com/ethereum/go-ethereum/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

// This is an interface of a struct that stores all data of the world viewed by this node (but
// not by the network). Examples of world state data are token price, nonce of each interested
// address, etc.
type WorldState interface {
	UseAndIncreaseNonce(chain string) int64
}

type DefaultWorldState struct {
	privateDb    keeper.Storage
	tssConfig    config.TssConfig
	nonces       map[string]int64
	deyesClients map[string]*tssclients.DeyesClient
}

func NewWorldState(tssConfig config.TssConfig, privateDb keeper.Storage, deyesClients map[string]*tssclients.DeyesClient) WorldState {
	return &DefaultWorldState{
		tssConfig:    tssConfig,
		privateDb:    privateDb,
		nonces:       make(map[string]int64, 0),
		deyesClients: deyesClients,
	}
}

func (ws *DefaultWorldState) UseAndIncreaseNonce(chain string) int64 {
	keyType := libchain.GetKeyTypeForChain(chain)

	pubKeyBytes := ws.privateDb.GetKeygenPubkey(keyType)
	if pubKeyBytes == nil {
		log.Error("cannot find pub key for keyType", chain)
		return -1
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		log.Error("Cannot unmarshal pubkey, err = ", err)
		return -1
	}

	pubKeyAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	if ws.nonces[chain] == 0 {
		eyeClient := ws.deyesClients[chain]
		if eyeClient == nil {
			return -1
		}

		nonce := eyeClient.GetNonce(chain, pubKeyAddress)

		if nonce == -1 {
			return -1
		}
		ws.nonces[chain] = nonce
	}

	nonce := ws.nonces[chain]
	ws.nonces[chain] = ws.nonces[chain] + 1

	log.Verbose("World state, nonce for chain ", chain, " is ", nonce)

	return nonce
}
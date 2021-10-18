package tss

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/tssclients"
)

// This is an interface of a struct that stores all data of the world viewed by this node (but
// not by the network). Examples of world state data are token price, nonce of each interested
// address, etc.
type WorldState interface {
	UseAndIncreaseNonce(chain string) int64
}

type DefaultWorldState struct {
	db           db.Database
	storage      *TssStorage
	tssConfig    config.TssConfig
	nonces       map[string]int64
	deyesClients map[string]*tssclients.DeyesClient
}

func NewWorldState(tssConfig config.TssConfig, db db.Database, storage *TssStorage, deyesClients map[string]*tssclients.DeyesClient) WorldState {
	return &DefaultWorldState{
		tssConfig:    tssConfig,
		db:           db,
		storage:      storage,
		nonces:       make(map[string]int64, 0),
		deyesClients: deyesClients,
	}
}

func (ws *DefaultWorldState) UseAndIncreaseNonce(chain string) int64 {
	pubKeyBytes := ws.storage.GetPubKey(chain)
	if pubKeyBytes == nil {
		utils.LogError("cannot find pub key for chain", chain)
		return -1
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
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

	return nonce
}

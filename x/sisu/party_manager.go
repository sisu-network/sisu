package sisu

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/sisu/common"
)

//go:generate mockgen -source=party_manager.go -destination=../../tests/mock/party_manager.go -package=mock

type Party struct {
	pubKey ctypes.PubKey
}

type PartyManager interface {
	GetActivePartyPubkeys() []ctypes.PubKey
}

type DefaultPartyManager struct {
	globalData    common.GlobalData
	activePubkeys ctypes.PubKey
}

func NewPartyManager(globalData common.GlobalData) PartyManager {
	return &DefaultPartyManager{
		globalData: globalData,
	}
}

func (pm *DefaultPartyManager) GetActivePartyPubkeys() []ctypes.PubKey {
	// TODO: Load this from database or update the list of active parties on memory. For now, use
	// all validator set.
	validators := pm.globalData.GetValidatorSet()
	pubkeys := make([]ctypes.PubKey, len(validators))
	for i, v := range validators {
		pubkeys[i] = v.PubKey
	}

	return pubkeys
}

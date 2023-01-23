package components

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type Party struct {
	pubKey ctypes.PubKey
}

// TODO: Merge this with validator manager.
type PartyManager interface {
	GetActivePartyPubkeys() []ctypes.PubKey
}

type DefaultPartyManager struct {
	globalData    GlobalData
	activePubkeys ctypes.PubKey
}

func NewPartyManager(globalData GlobalData) PartyManager {
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

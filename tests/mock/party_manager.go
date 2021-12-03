package mock

import (
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
)

type PartyManager struct{}

func (pm *PartyManager) GetActivePartyPubkeys() []ctypes.PubKey {
	return []ctypes.PubKey{}
}

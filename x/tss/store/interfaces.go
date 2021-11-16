package store

import sdk "github.com/sisu-network/cosmos-sdk/types"

type ChainStore interface {
	GetChainID() ChainId
	SaveKeyAddrs(ctx sdk.Context, keyAddrs map[string]bool) error
	GetAllKeyAddrs(ctx sdk.Context) (map[string]map[string]bool, error)
}

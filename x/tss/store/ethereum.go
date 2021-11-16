package store

import (
	"encoding/json"
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
)

var (
	PREFIX_ETH_KEY_ADDRESS = []byte{0x08}

	KEY_ETH_KEY_ADDRESS = "eth_key_address_%s" // chain
)

type EthereumStore struct {
	storeKey sdk.StoreKey
}

func NewEthereumStore(storeKey sdk.StoreKey) *EthereumStore {
	return &EthereumStore{storeKey: storeKey}
}

func (s *EthereumStore) GetChainID() ChainId {
	return Ethereum
}

func (s *EthereumStore) SaveKeyAddrs(ctx sdk.Context, keyAddrs map[string]bool) error {
	key := fmt.Sprintf(KEY_ETH_KEY_ADDRESS, s.GetChainID())
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PREFIX_ETH_KEY_ADDRESS)

	bz, err := json.Marshal(keyAddrs)
	if err != nil {
		log.Error("cannot marshal key addrs, err =", err)
		return err
	}

	store.Set([]byte(key), bz)
	return nil
}

func (s *EthereumStore) GetAllKeyAddrs(ctx sdk.Context) (map[string]map[string]bool, error) {
	m := make(map[string]map[string]bool)
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PREFIX_ETH_KEY_ADDRESS)

	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		m2 := make(map[string]bool)
		err := json.Unmarshal(iter.Value(), &m2)
		if err != nil {
			log.Error("cannot unmarshal value with key", iter.Key())
			return nil, err
		}
		m[string(iter.Key())] = m2
	}

	return m, nil
}

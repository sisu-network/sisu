package keeper

import (
	"fmt"
	"strings"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	chainstore "github.com/sisu-network/sisu/x/tss/store"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

const (
	OBSERVED_TX_CACHE_SIZE = 2500
)

// TODO: clean up this list
var (
	PREFIX_RECORDED_CHAIN            = []byte{0x01}
	PREFIX_OBSERVED_TX               = []byte{0x02}
	PREFIX_OBSERVED_TX_VALIDATOR_SET = []byte{0x03}
	PREFIX_PENDING_OBSERVED_TX       = []byte{0x04}
	PREFIX_PROCESSED_OBSERVED_TX     = []byte{0x05}
	PREFIX_PUBLIC_KEY_BYTES          = []byte{0x06}
	PREFIX_PENDING_KEYGEN_TX         = []byte{0x07}
	PREFIX_TX_OUT                    = []byte{0x09}

	// List of on memory keys. These data are not persisted into kvstore.
	// List of contracts that need to be deployed to a chain.
	KEY_CONTRACT_QUEUE     = "contract_queue_%s_%s"     // chain
	KEY_DEPLOYING_CONTRACT = "deploying_contract_%s_%s" // chain - contract hash
)

type ISmartContract interface {
	GetData() []byte
	GetCreatedBlock() int64
	GetDesignatedValidator() string
}

type deployContractWrapper struct {
	data         []byte
	createdBlock int64 // Sisu block when the contract is created.
	// id of the designated validator that is supposed to post the tx out to the Sisu chain.
	designatedValidator string
}

func (d *deployContractWrapper) GetData() []byte {
	return d.data
}

func (d *deployContractWrapper) GetCreatedBlock() int64 {
	return d.createdBlock
}

func (d *deployContractWrapper) GetDesignatedValidator() string {
	return d.designatedValidator
}

type Keeper interface {
	GetKeyStore() sdk.StoreKey
	GetContractQueue(address string) (map[string]string, error)
	GetDeployingContract(address string) (ISmartContract, error)
	GetDeployedContract(address string) (ISmartContract, error)

	// TODO: move those methods to a new struct?
	SaveObservedTx(ctx sdk.Context, tx *tssTypes.ObservedTx)
	GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte
	SavePubKey(ctx sdk.Context, chain string, keyBytes []byte)
	GetAllPubKeys(ctx sdk.Context) map[string][]byte

	GetChainStore(chainId chainstore.ChainId) (chainstore.ChainStore, error)
}

type DefaultKeeper struct {
	storeKey sdk.StoreKey

	// List of contracts that waits to be deployed.
	contractQueue map[string]string

	// List of contracts that are being deployed.
	deployingContracts map[string]*deployContractWrapper

	// List of contracts that have been deployed. (This should be saved in a KV store?)
	deployedContracts map[string]*deployContractWrapper

	chainStores map[chainstore.ChainId]chainstore.ChainStore
}

func NewDefaultKeeper(storeKey sdk.StoreKey) *DefaultKeeper {
	keeper := &DefaultKeeper{
		storeKey:           storeKey,
		contractQueue:      make(map[string]string),
		deployingContracts: make(map[string]*deployContractWrapper),
		deployedContracts:  make(map[string]*deployContractWrapper),
	}
	keeper.initChainStores()

	return keeper
}

func (k *DefaultKeeper) initChainStores() {
	chainStores := make(map[chainstore.ChainId]chainstore.ChainStore)

	ethStore := chainstore.NewEthereumStore(k.GetKeyStore())
	chainStores[ethStore.GetChainID()] = ethStore

	k.chainStores = chainStores
}

func (k *DefaultKeeper) getKey(chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	chain = strings.Replace(chain, "_", "*", -1)
	return []byte(fmt.Sprintf("%s_%d_%s", chain, height, hash))
}

func (k *DefaultKeeper) SaveObservedTx(ctx sdk.Context, tx *tssTypes.ObservedTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(tx.Chain, tx.BlockHeight, tx.TxHash)

	store.Set(key, tx.Serialized)
}

func (k *DefaultKeeper) GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(chain, blockHeight, hash)

	return store.Get(key)
}

func (k *DefaultKeeper) SavePubKey(ctx sdk.Context, chain string, keyBytes []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	store.Set([]byte(chain), keyBytes)
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	iter := store.Iterator(nil, nil)
	ret := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		ret[string(iter.Key())] = iter.Value()
	}

	return ret
}

func (k *DefaultKeeper) GetKeyStore() sdk.StoreKey {
	return k.storeKey
}

func (k *DefaultKeeper) GetContractQueue(address string) (map[string]string, error) {
	panic("implement me")
}

func (k *DefaultKeeper) GetDeployingContract(address string) (ISmartContract, error) {
	panic("implement me")
}

func (k *DefaultKeeper) GetDeployedContract(address string) (ISmartContract, error) {
	panic("implement me")
}

func (k *DefaultKeeper) GetChainStore(chainId chainstore.ChainId) (chainstore.ChainStore, error) {
	chainStore, ok := k.chainStores[chainId]
	if !ok {
		err := fmt.Errorf("chain store not found: %s", chainId)
		log.Error(err)
		return nil, err
	}

	return chainStore, nil
}

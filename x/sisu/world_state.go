package sisu

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

var (
	chainToTokens = map[string]string{
		"bsc":     "BNB",
		"eth":     "ETH",
		"ropsten": "ETH",
	}

	defaultTokenPrices = map[string]float32{
		"BNB": 100.0,
		"ETH": 2000.0,
	}

	defaultGasPrice = map[string]*big.Int{
		"ganache1":            big.NewInt(2_000_000_000),
		"ganache2":            big.NewInt(2_000_000_000),
		"eth-ropsten":         big.NewInt(4_000_000_000),
		"eth-binance-testnet": big.NewInt(10_000_000_000),
	}

	ErrChainNotFound      = fmt.Errorf("chain not found")
	ErrTokenPriceNotFound = fmt.Errorf("token price not found")
)

// This is an interface of a struct that stores all data of the world data. Examples of world state
// data are token price, nonce of addresses, etc.
type WorldState interface {
	UseAndIncreaseNonce(chain string) int64

	SetGasPrice(chain string, price *big.Int)
	GetGasPrice(chain string) (*big.Int, error)

	SetTokenPrices(tokenPrices map[string]float32)
	GetNativeTokenPriceForChain(chain string) (float32, error)
}

type DefaultWorldState struct {
	publicDb    keeper.Storage
	tssConfig   config.TssConfig
	nonces      map[string]int64
	deyesClient tssclients.DeyesClient

	gasPrices   *sync.Map
	tokenPrices *sync.Map
}

func NewWorldState(tssConfig config.TssConfig, publicDb keeper.Storage, deyesClients tssclients.DeyesClient) WorldState {
	return &DefaultWorldState{
		tssConfig:   tssConfig,
		publicDb:    publicDb,
		nonces:      make(map[string]int64, 0),
		deyesClient: deyesClients,
		tokenPrices: &sync.Map{},
		gasPrices:   &sync.Map{},
	}
}

func (ws *DefaultWorldState) UseAndIncreaseNonce(chain string) int64 {
	keyType := libchain.GetKeyTypeForChain(chain)

	pubKeyBytes := ws.publicDb.GetKeygenPubkey(keyType)
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
		nonce := ws.deyesClient.GetNonce(chain, pubKeyAddress)

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

func (ws *DefaultWorldState) SetGasPrice(chain string, price *big.Int) {
	ws.gasPrices.Store(chain, price)
}

func (ws *DefaultWorldState) GetGasPrice(chain string) (*big.Int, error) {
	val, ok := ws.gasPrices.Load(chain)
	if ok {
		return val.(*big.Int), nil
	}

	return nil, ErrChainNotFound
}

// getCommissionFee returns the amount of fee user needs to pay the Sisu network (often a percentage
// of the transaction amount).
func (ws *DefaultWorldState) getCommissionFee(amount *big.Int) *big.Int {
	return big.NewInt(0)
}

func (ws *DefaultWorldState) SetTokenPrices(tokenPrices map[string]float32) {
	for token, price := range tokenPrices {
		ws.tokenPrices.Store(token, price)
	}
}

func (ws *DefaultWorldState) GetNativeTokenPriceForChain(chain string) (float32, error) {
	token := chainToTokens[chain]
	if len(token) == 0 {
		return 0, ErrChainNotFound
	}

	val, ok := ws.tokenPrices.Load(token)
	if ok {
		price := val.(float32)
		return price, nil
	}

	return 0, ErrTokenPriceNotFound
}

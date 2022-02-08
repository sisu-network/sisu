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
	"github.com/sisu-network/sisu/x/sisu/types"
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
// go:generate mockgen -source x/sisu/world_state.go -destination=tests/mock/x/sisu/world_state.go -package=mock
type WorldState interface {
	UseAndIncreaseNonce(chain string) int64

	SetGasPrice(chain string, price *big.Int)
	GetGasPrice(chain string) (*big.Int, error)

	SetTokens(tokenPrices map[string]*types.Token)
	GetTokenPrice(token string) (int64, error)
	GetNativeTokenPriceForChain(chain string) (int64, error)

	GetTokenFromAddress(chain string, tokenAddr string) *types.Token
}

type DefaultWorldState struct {
	publicDb    keeper.Storage
	tssConfig   config.TssConfig
	nonces      map[string]int64
	deyesClient tssclients.DeyesClient

	gasPrices   *sync.Map
	tokens      *sync.Map
	addrToToken *sync.Map
}

func NewWorldState(tssConfig config.TssConfig, publicDb keeper.Storage, deyesClients tssclients.DeyesClient) WorldState {
	return &DefaultWorldState{
		tssConfig:   tssConfig,
		publicDb:    publicDb,
		nonces:      make(map[string]int64, 0),
		deyesClient: deyesClients,
		tokens:      &sync.Map{},
		gasPrices:   &sync.Map{},
		addrToToken: &sync.Map{},
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

func (ws *DefaultWorldState) SetTokens(tokens map[string]*types.Token) {
	for tokenId, token := range tokens {
		ws.tokens.Store(tokenId, token)

		// Save the mapping of token address on each chain to token for later retrieval.
		if token.Addresses != nil {
			for chain, addr := range token.Addresses {
				key := fmt.Sprintf("%s__%s", chain, addr)
				ws.addrToToken.Store(key, token)
			}
		}
	}
}

func (ws *DefaultWorldState) GetNativeTokenPriceForChain(chain string) (int64, error) {
	tokenId := chainToTokens[chain]
	if len(tokenId) == 0 {
		return 0, ErrChainNotFound
	}

	return ws.GetTokenPrice(tokenId)
}

func (ws *DefaultWorldState) GetTokenPrice(tokenId string) (int64, error) {
	val, ok := ws.tokens.Load(tokenId)
	if ok {
		token := val.(*types.Token)
		return token.Price, nil
	}

	return 0, ErrTokenPriceNotFound
}

func (ws *DefaultWorldState) GetTokenFromAddress(chain string, tokenAddr string) *types.Token {
	key := fmt.Sprintf("%s__%s", chain, tokenAddr)
	val, ok := ws.addrToToken.Load(key)
	if !ok {
		return nil
	}

	return val.(*types.Token)
}

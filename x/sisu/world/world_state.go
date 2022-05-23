package world

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"go.uber.org/atomic"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// This is a mapping between chain name and the native token that it uses.
	chainToTokens = map[string]string{
		"bsc":              "BNB",
		"eth":              "ETH",
		"ganache1":         "NATIVE_GANACHE1",
		"ganache2":         "NATIVE_GANACHE2",
		"ropsten-testnet":  "ETH",
		"binance-testnet":  "BNB",
		"polygon-testnet":  "MATIC",
		"xdai":             "xDai",
		"goerli-testnet":   "ETH",
		"arbitrum-testnet": "ETH",
		"fantom-testnet":   "FTM",
	}
)

// This is an interface of a struct that stores all data of the world data. Examples of world state
// data are token price, nonce of addresses, network gas fee, etc.
// go:generate mockgen -source x/sisu/world/world_state.go -destination=tests/mock/x/sisu/world_state.go -package=mock
type WorldState interface {
	IsDataInitialized() bool
	InitData(ctx sdk.Context)

	UseAndIncreaseNonce(ctx sdk.Context, chain string) int64

	SetChain(chain *types.Chain)
	GetGasPrice(chain string) (*big.Int, error)

	SetTokens(tokenPrices map[string]*types.Token)
	GetTokenPrice(token string) (int64, error)
	GetNativeTokenPriceForChain(chain string) (int64, error)
	GetGasCostInToken(tokenId, chainId string) (int64, error)

	GetTokenFromAddress(chain string, tokenAddr string) *types.Token
}

type DefaultWorldState struct {
	keeper      keeper.Keeper
	nonces      map[string]int64
	deyesClient tssclients.DeyesClient

	isDataInit  *atomic.Bool
	chains      *sync.Map
	tokens      *sync.Map
	addrToToken *sync.Map // chain__addr => *types.Token
}

func NewWorldState(keeper keeper.Keeper, deyesClients tssclients.DeyesClient) WorldState {
	return &DefaultWorldState{
		keeper:      keeper,
		nonces:      make(map[string]int64, 0),
		deyesClient: deyesClients,
		tokens:      &sync.Map{},
		chains:      &sync.Map{},
		addrToToken: &sync.Map{},
		isDataInit:  atomic.NewBool(false),
	}
}

func (ws *DefaultWorldState) IsDataInitialized() bool {
	return ws.isDataInit.Load()
}

func (ws *DefaultWorldState) InitData(ctx sdk.Context) {
	log.Info("Initializing world state data")

	// Get saved tokens
	tokens := ws.keeper.GetAllTokens(ctx)
	ws.SetTokens(tokens)

	// Map between address and tokens
	for _, token := range tokens {
		for i, chain := range token.Chains {
			addr := token.Addresses[i]
			key := ws.getChainAddrKey(chain, addr)
			ws.addrToToken.Store(key, token)
		}
	}

	// Get saved network gas
	chains := ws.keeper.GetAllChains(ctx)
	for _, chain := range chains {
		ws.SetChain(chain)
	}

	ws.isDataInit.Store(true)
}

func (ws *DefaultWorldState) UseAndIncreaseNonce(ctx sdk.Context, chain string) int64 {
	keyType := libchain.GetKeyTypeForChain(chain)

	pubKeyBytes := ws.keeper.GetKeygenPubkey(ctx, keyType)
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

	log.Verbose("World state, nonce for address ", pubKeyAddress, " on chain ", chain, " is ", nonce)

	return nonce
}

func (ws *DefaultWorldState) SetChain(chain *types.Chain) {
	ws.chains.Store(chain.Id, chain)
}

func (ws *DefaultWorldState) GetGasPrice(chainId string) (*big.Int, error) {
	val, ok := ws.chains.Load(chainId)
	if !ok {
		return nil, NewErrChainNotFound(chainId)
	}

	chain := val.(*types.Chain)

	return big.NewInt(chain.GasPrice), nil
}

func (ws *DefaultWorldState) SetTokens(tokens map[string]*types.Token) {
	for tokenId, token := range tokens {
		ws.tokens.Store(tokenId, token)

		// Save the mapping of token address on each chain to token for later retrieval.
		if token.Addresses != nil {
			for chain, addr := range token.Addresses {
				key := fmt.Sprintf("%d__%s", chain, addr)
				ws.addrToToken.Store(key, token)
			}
		}
	}
}

func (ws *DefaultWorldState) GetNativeTokenPriceForChain(chain string) (int64, error) {
	tokenId := chainToTokens[chain]
	if len(tokenId) == 0 {
		return 0, NewErrTokenNotFound(tokenId)
	}

	return ws.GetTokenPrice(tokenId)
}

func (ws *DefaultWorldState) GetTokenPrice(tokenId string) (int64, error) {
	val, ok := ws.tokens.Load(tokenId)
	if ok {
		token := val.(*types.Token)
		return token.Price, nil
	}

	return 0, NewErrTokenNotFound(tokenId)
}

func (ws *DefaultWorldState) GetTokenFromAddress(chain string, tokenAddr string) *types.Token {
	key := ws.getChainAddrKey(chain, tokenAddr)
	val, ok := ws.addrToToken.Load(key)
	if !ok {
		return nil
	}

	return val.(*types.Token)
}

func (ws *DefaultWorldState) GetGasCostInToken(tokenId, chainId string) (int64, error) {
	gasPrice, err := ws.GetGasPrice(chainId)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	// TODO: correct gasUnit here
	gasUnit := big.NewInt(60_000) // Estimated cost for swapping.
	tokenPrice, err := ws.GetTokenPrice(tokenId)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	if tokenPrice < 0 {
		return 0, fmt.Errorf("Token price is negative, token id = %s, token price = %d", tokenId, tokenPrice)
	}

	nativeTokenPrice, err := ws.GetNativeTokenPriceForChain(chainId)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	gasCost, err := helper.GetGasCostInToken(gasUnit, gasPrice, big.NewInt(tokenPrice), big.NewInt(nativeTokenPrice))
	if err != nil {
		log.Error(err)
		return -1, err
	}

	return gasCost.Int64(), nil
}

func (ws *DefaultWorldState) getChainAddrKey(chain, addr string) string {
	return fmt.Sprintf("%s__%s", chain, addr)
}

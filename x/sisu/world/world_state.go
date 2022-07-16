package world

import (
	"fmt"
	"math/big"
	"sync"

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
type WorldState interface {
	IsDataInitialized() bool
	InitData(ctx sdk.Context)

	SetChain(chain *types.Chain)
	GetGasPrice(chain string) (*big.Int, error)

	SetTokens(tokenPrices map[string]*types.Token)
	GetTokenPrice(token string) (*big.Int, error)
	GetNativeTokenPriceForChain(chain string) (*big.Int, error)
	GetGasCostInToken(tokenId, chainId string) (*big.Int, error)

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

func (ws *DefaultWorldState) GetNativeTokenPriceForChain(chain string) (*big.Int, error) {
	tokenId := chainToTokens[chain]
	if len(tokenId) == 0 {
		return big.NewInt(0), NewErrTokenNotFound(tokenId)
	}

	return ws.GetTokenPrice(tokenId)
}

func (ws *DefaultWorldState) GetTokenPrice(tokenId string) (*big.Int, error) {
	val, ok := ws.tokens.Load(tokenId)
	if ok {
		token := val.(*types.Token)
		price, ok := new(big.Int).SetString(token.Price, 10)
		if !ok {
			return nil, fmt.Errorf("invalid token price %s", token.Price)
		}
		return price, nil
	}

	return big.NewInt(0), NewErrTokenNotFound(tokenId)
}

func (ws *DefaultWorldState) GetTokenFromAddress(chain string, tokenAddr string) *types.Token {
	key := ws.getChainAddrKey(chain, tokenAddr)
	val, ok := ws.addrToToken.Load(key)
	if !ok {
		return nil
	}

	return val.(*types.Token)
}

func (ws *DefaultWorldState) GetGasCostInToken(tokenId, chainId string) (*big.Int, error) {
	gasPrice, err := ws.GetGasPrice(chainId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	gasUnit := big.NewInt(80_000) // Estimated cost for swapping is 60k. We add some redundancy here.
	tokenPrice, err := ws.GetTokenPrice(tokenId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if big.NewInt(0).Cmp(tokenPrice) == 0 {
		return nil, fmt.Errorf("Token %s has price 0", tokenId)
	}

	if tokenPrice.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("Token price is negative, token id = %s, token price = %d", tokenId, tokenPrice)
	}

	nativeTokenPrice, err := ws.GetNativeTokenPriceForChain(chainId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	gasCost, err := helper.GetGasCostInToken(gasUnit, gasPrice, tokenPrice, nativeTokenPrice)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return gasCost, nil
}

func (ws *DefaultWorldState) getChainAddrKey(chain, addr string) string {
	return fmt.Sprintf("%s__%s", chain, addr)
}

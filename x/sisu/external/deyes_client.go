package external

import (
	"context"

	"github.com/echovl/cardano-go"
	"github.com/ethereum/go-ethereum/rpc"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
)

type DeyesClient interface {
	Ping(source string) error
	Dispatch(request *etypes.DispatchedTxRequest) (*etypes.DispatchedTxResult, error)
	SetVaultAddress(chain string, addr string, token string) error
	GetNonce(chain string, address string) (int64, error)
	SetSisuReady(isReady bool) error
	GetGasPrices(chains []string) ([]int64, error)

	// Cardano
	CardanoProtocolParams(chain string) (*cardano.ProtocolParams, error)
	CardanoUtxos(chain string, addr string, maxBlock uint64) ([]cardano.UTxO, error)
	CardanoBalance(chain string, address string, maxBlock int64) (*cardano.Value, error)
	CardanoSubmitTx(chain string, tx *cardano.Tx) (*cardano.Hash32, error)
	CardanoTip(chain string, blockHeight uint64) (*cardano.NodeTip, error)

	// Solana
	SolanaQueryRecentBlock(chain string) (*etypes.SolanaQueryRecentBlockResult, error)
}

type defaultDeyesClient struct {
	client *rpc.Client
}

func DialDeyes(rawurl string) (DeyesClient, error) {
	return dialDeyesContext(context.Background(), rawurl)
}

func dialDeyesContext(ctx context.Context, rawurl string) (DeyesClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDeyesClient(c), nil
}

func newDeyesClient(c *rpc.Client) DeyesClient {
	return &defaultDeyesClient{c}
}

func (c *defaultDeyesClient) Ping(source string) error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "deyes_ping", source)
	if err != nil {
		log.Error("Cannot ping deyes, err = ", err)
		return err
	}

	return nil
}

// Informs the deyes that Sisu server is ready to accept transaction.
func (c *defaultDeyesClient) SetSisuReady(isReady bool) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setSisuReady", isReady)
	if err != nil {
		log.Error("Cannot Set readiness for deyes, err = ", err)
		return err
	}

	return nil
}

func (c *defaultDeyesClient) SetVaultAddress(chain string, addr string, token string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setVaultAddress", chain, addr, token)
	if err != nil {
		log.Error("Cannot set gateway address for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

func (c *defaultDeyesClient) Dispatch(request *etypes.DispatchedTxRequest) (*etypes.DispatchedTxResult, error) {
	var result = &etypes.DispatchedTxResult{}
	err := c.client.CallContext(context.Background(), &result, "deyes_dispatchTx", request)
	if err != nil {
		log.Error("Cannot Dispatch tx to the chain", request.Chain, "err =", err)
		return result, err
	}

	log.Verbose("Tx has been dispatched")

	return result, nil
}

func (c *defaultDeyesClient) GetNonce(chain string, address string) (int64, error) {
	var result int64
	err := c.client.CallContext(context.Background(), &result, "deyes_getNonce", chain, address)
	if err != nil {
		log.Error("Cannot get nonce for chain and address", chain, address, "err =", err)
	}

	return result, err
}

func (c *defaultDeyesClient) GetGasPrices(chains []string) ([]int64, error) {
	result := make([]int64, 0)
	err := c.client.CallContext(context.Background(), &result, "deyes_getGasPrices", chains)
	if err != nil {
		log.Error("Cannot get gas price for chains = ", chains, "err = ", err)
		return nil, err
	}

	return result, nil
}

///// Carnado

func (c *defaultDeyesClient) CardanoProtocolParams(chain string) (*cardano.ProtocolParams, error) {
	result := &cardano.ProtocolParams{}

	err := c.client.CallContext(context.Background(), &result, "deyes_cardanoProtocolParams", chain)
	if err != nil {
		log.Error("Cannot get cardano protocol params = ", chain, "err = ", err)
		return nil, err
	}

	return result, nil
}

func (c *defaultDeyesClient) CardanoUtxos(chain string, addr string, maxBlock uint64) ([]cardano.UTxO, error) {
	result := struct {
		Utxos []cardano.UTxO
		Bytes [][]byte // Amounts
	}{}

	err := c.client.CallContext(context.Background(), &result, "deyes_cardanoUtxos", chain, addr, maxBlock)
	if err != nil {
		log.Errorf("Cannot get cardano utxos, chain = %s, addr = %s, err = %s", chain, addr, err.Error())
		return nil, err
	}

	// Unmarshal amount since it's not serializable through network.
	utxos := result.Utxos
	for i, bz := range result.Bytes {
		utxos[i].Amount = cardano.NewValue(0)
		err := utxos[i].Amount.UnmarshalCBOR(bz)
		if err != nil {
			return nil, err
		}
	}

	return utxos, nil
}

// Balance returns the current balance of an account.
func (c *defaultDeyesClient) CardanoBalance(chain string, address string, maxBlock int64) (*cardano.Value, error) {
	result := new(cardano.Value)

	err := c.client.CallContext(context.Background(), &result, "deyes_cardanoBalance", chain, address, maxBlock)
	if err != nil {
		log.Error("Cannot get cardano balance, chain = ", chain, "err = ", err)
		return nil, err
	}

	return result, nil
}

// Tip returns the node's current tip
func (c *defaultDeyesClient) CardanoTip(chain string, blockHeight uint64) (*cardano.NodeTip, error) {
	result := new(cardano.NodeTip)

	err := c.client.CallContext(context.Background(), &result, "deyes_cardanoTip", chain, blockHeight)
	if err != nil {
		log.Errorf("Cannot get cardano block tip, chain = %s, height = %s, err = %s", chain, blockHeight, err.Error())
		return nil, err
	}

	return result, nil
}

func (c *defaultDeyesClient) CardanoSubmitTx(chain string, tx *cardano.Tx) (*cardano.Hash32, error) {
	result := new(cardano.Hash32)

	err := c.client.CallContext(context.Background(), &result, "deyes_cardanoSubmitTx", chain, tx)
	if err != nil {
		log.Error("Cannot submit cardano transaction, chain = ", chain, "err = ", err)
		return nil, err
	}

	return result, nil
}

/////
func (c *defaultDeyesClient) SolanaQueryRecentBlock(chain string) (*etypes.SolanaQueryRecentBlockResult, error) {
	result := &etypes.SolanaQueryRecentBlockResult{}

	err := c.client.CallContext(context.Background(), &result, "deyes_solanaQueryRecentBlock", chain)
	if err != nil {
		log.Error("Cannot query recent solana block, chain = ", chain, "err = ", err)
		return nil, err
	}

	return result, nil
}

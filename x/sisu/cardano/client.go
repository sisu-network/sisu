package cardano

import (
	"context"
	"fmt"

	"github.com/blockfrost/blockfrost-go"
	"github.com/echovl/cardano-go"
	cblockfrost "github.com/echovl/cardano-go/blockfrost"
)

const (
	MaxBlockHeight = 9223372036854775807
)

var (
	InsufficientFundErr = fmt.Errorf("Insufficient Fund")
)

// CardanoClient is a interface that interface with Cardano blockchain.
type CardanoClient interface {
	// Balance returns the current balance of an account.
	Balance(address cardano.Address) (*cardano.Value, error)

	// UTxOs returns a list of unspent transaction outputs for a given address
	UTxOs(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error)

	// Tip returns the node's current tip
	Tip() (*cardano.NodeTip, error)

	// ProtocolParams returns the Node's Protocol Parameters
	ProtocolParams() (*cardano.ProtocolParams, error)

	SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error)
}

// blockFrostClient is a struct that implements CardanoClient interface
type blockFrostClient struct {
	cardanoNode cardano.Node
	bfClient    blockfrost.APIClient
}

func NewBlockfrostClient(network cardano.Network, secret string) CardanoClient {
	server := blockfrost.CardanoMainNet
	if network == cardano.Testnet {
		server = blockfrost.CardanoTestNet
	}

	return &blockFrostClient{
		cardanoNode: cblockfrost.NewNode(network, secret),
		bfClient: blockfrost.NewAPIClient(blockfrost.APIClientOptions{
			ProjectID: secret,
			Server:    server,
		}),
	}
}

func (c *blockFrostClient) Balance(address cardano.Address) (*cardano.Value, error) {
	balance := cardano.NewValue(0)
	utxos, err := c.UTxOs(address, MaxBlockHeight)
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		balance = balance.Add(utxo.Amount)
	}

	return balance, nil
}

func (c *blockFrostClient) UTxOs(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error) {
	utxos, err := c.cardanoNode.UTxOs(addr)
	if err != nil {
		return utxos, err
	}

	if maxBlock == MaxBlockHeight {
		return utxos, nil
	}

	cached := make(map[string]blockfrost.TransactionContent)
	filteredUtxos := make([]cardano.UTxO, 0)

	// Filter utxos
	for _, utxo := range utxos {
		tx, ok := cached[utxo.TxHash.String()]
		if !ok {
			tx, err = c.bfClient.Transaction(context.Background(), string(utxo.TxHash))
			if err != nil {
				continue
			}
			cached[utxo.TxHash.String()] = tx
		}

		if tx.BlockHeight > int(maxBlock) {
			continue
		}

		filteredUtxos = append(filteredUtxos, utxo)
	}

	return filteredUtxos, nil
}

// Tip returns the node's current tip
func (c *blockFrostClient) Tip() (*cardano.NodeTip, error) {
	return c.cardanoNode.Tip()
}

// ProtocolParams returns the Node's Protocol Parameters
func (c *blockFrostClient) ProtocolParams() (*cardano.ProtocolParams, error) {
	return c.cardanoNode.ProtocolParams()
}

func (c *blockFrostClient) SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error) {
	return c.cardanoNode.SubmitTx(tx)
}

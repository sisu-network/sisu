package cardano

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/blockfrost/blockfrost-go"
	"github.com/echovl/cardano-go"
	ecore "github.com/sisu-network/deyes/chains/cardano/core"
	econfig "github.com/sisu-network/deyes/config"

	_ "github.com/lib/pq"
)

var _ CardanoClient = (*SyncDbClient)(nil)

type SyncDbClient struct {
	cfg         econfig.SyncDbConfig
	submitTxURL string
	syncDB      *ecore.SyncDB
}

func NewSyncDBClient(cfg econfig.SyncDbConfig, submitTxURL string) *SyncDbClient {
	db, err := ecore.ConnectDB(cfg)
	if err != nil {
		panic(err)
	}

	syncDB := ecore.NewSyncDBConnector(db)
	return &SyncDbClient{
		cfg:         cfg,
		syncDB:      syncDB,
		submitTxURL: submitTxURL,
	}
}

func (s *SyncDbClient) Balance(address cardano.Address, maxBlock uint64) (*cardano.Value, error) {
	balance := cardano.NewValue(0)
	utxos, err := s.UTxOs(address, maxBlock)
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		balance = balance.Add(utxo.Amount)
	}

	return balance, nil
}

func (s *SyncDbClient) UTxOs(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error) {
	butxos, err := s.syncDB.AddressUTXOs(context.Background(), addr.Bech32(), blockfrost.APIQueryParams{To: strconv.Itoa(int(maxBlock))})
	if err != nil {
		return nil, err
	}

	utxos := make([]cardano.UTxO, len(butxos))

	for i, butxo := range butxos {
		txHash, err := cardano.NewHash32(butxo.TxHash)
		if err != nil {
			return nil, err
		}

		amount := cardano.NewValue(0)
		for _, a := range butxo.Amount {
			if a.Unit == "lovelace" {
				lovelace, err := strconv.ParseUint(a.Quantity, 10, 64)
				if err != nil {
					return nil, err
				}
				amount.Coin += cardano.Coin(lovelace)
			} else {
				unitBytes, err := hex.DecodeString(a.Unit)
				if err != nil {
					return nil, err
				}
				policyID := cardano.NewPolicyIDFromHash(unitBytes[:28])
				assetName := string(unitBytes[28:])
				assetValue, err := strconv.ParseUint(a.Quantity, 10, 64)
				if err != nil {
					return nil, err
				}
				currentAssets := amount.MultiAsset.Get(policyID)
				if currentAssets != nil {
					currentAssets.Set(
						cardano.NewAssetName(assetName),
						cardano.BigNum(assetValue),
					)
				} else {
					amount.MultiAsset.Set(
						policyID,
						cardano.NewAssets().
							Set(
								cardano.NewAssetName(string(assetName)),
								cardano.BigNum(assetValue),
							),
					)
				}
			}
		}

		utxos[i] = cardano.UTxO{
			Spender: addr,
			TxHash:  txHash,
			Amount:  amount,
			Index:   uint64(butxo.OutputIndex),
		}
	}

	return utxos, nil
}

func (s *SyncDbClient) Tip() (*cardano.NodeTip, error) {
	block, err := s.syncDB.BlockLatest(context.Background())
	if err != nil {
		return nil, err
	}

	return &cardano.NodeTip{
		Block: uint64(block.Height),
		Epoch: uint64(block.Epoch),
		Slot:  uint64(block.Slot),
	}, nil
}

func (s *SyncDbClient) ProtocolParams() (*cardano.ProtocolParams, error) {
	eparams, err := s.syncDB.LatestEpochParameters(context.Background())
	if err != nil {
		return nil, err
	}

	minUTXO, err := strconv.ParseUint(eparams.MinUtxo, 10, 64)
	if err != nil {
		return nil, err
	}

	poolDeposit, err := strconv.ParseUint(eparams.PoolDeposit, 10, 64)
	if err != nil {
		return nil, err
	}
	keyDeposit, err := strconv.ParseUint(eparams.KeyDeposit, 10, 64)
	if err != nil {
		return nil, err
	}

	pparams := &cardano.ProtocolParams{
		MinFeeA:            cardano.Coin(eparams.MinFeeA),
		MinFeeB:            cardano.Coin(eparams.MinFeeB),
		MaxBlockBodySize:   uint(eparams.MaxBlockSize),
		MaxTxSize:          uint(eparams.MaxTxSize),
		MaxBlockHeaderSize: uint(eparams.MaxBlockHeaderSize),
		KeyDeposit:         cardano.Coin(keyDeposit),
		PoolDeposit:        cardano.Coin(poolDeposit),
		MaxEpoch:           uint(eparams.Epoch),
		NOpt:               uint(eparams.NOpt),
		CoinsPerUTXOWord:   cardano.Coin(minUTXO),
	}

	return pparams, nil
}

func (s *SyncDbClient) SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error) {
	txBytes := tx.Bytes()

	url := s.submitTxURL
	req, err := http.NewRequest("POST", url, bytes.NewReader(txBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/cbor")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, errors.New(string(respBody))
	}

	txHash, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	return &txHash, nil
}

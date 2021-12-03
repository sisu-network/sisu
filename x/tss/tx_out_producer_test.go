package tss

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/tests/mock"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	t.Parallel()

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
	contractEntity := &types.ContractEntity{
		Chain: "eth",
		Hash:  hash,
	}

	worldState := NewWorldState(config.TssConfig{}, nil, nil)
	txOutProducer := DefaultTxOutputProducer{
		worldState: worldState,
		tssConfig: config.TssConfig{
			Enable: true,
			SupportedChains: map[string]config.TssChainConfig{
				"ganache": {
					Symbol:   "ganache",
					DeyesUrl: "http://0.0.0.0:1234",
				},
			},
		},
	}

	tx := txOutProducer.getContractTx(contractEntity, 100)
	require.NotNil(t, tx)
	require.EqualValues(t, 100, tx.Nonce())
	require.EqualValues(t, *big.NewInt(1700000000), *tx.GasPrice())
	require.EqualValues(t, *big.NewInt(1700000000), *tx.GasFeeCap())
}

func TestTxOutProducer_getEthResponse(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	contractEntities := []*types.ContractEntity{
		{
			Chain: "eth",
			Hash:  SupportedContracts[ContractErc20].AbiHash,
		},
	}

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	pubkeyBytes := crypto.FromECDSAPub(&privKey.PublicKey)
	mockDb := mock.NewMockDatabase(ctrl)
	mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(true).Times(1)
	mockDb.EXPECT().GetPendingDeployContracts(gomock.Any()).Return(contractEntities).Times(1)
	mockDb.EXPECT().GetPubKey("eth").Return(pubkeyBytes).Times(1)
	mockDb.EXPECT().UpdateContractsStatus(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockAppKeys := mock.NewMockAppKeys(ctrl)
	accAddress := []byte{1, 2, 3}
	mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).Times(1)

	amount := big.NewInt(rand.Int63())
	gasLimit := uint64(rand.Int63())
	gasPrice := big.NewInt(rand.Int63())
	ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &common.Address{},
		Value:    amount,
	})
	binary, err := ethTransaction.MarshalBinary()
	require.NoError(t, err)

	observedTx := types.ObservedTx{
		BlockHeight: 1,
		Serialized:  binary,
	}

	worldState := DefaultWorldState{
		db:        mockDb,
		tssConfig: config.TssConfig{},
		nonces: map[string]int64{
			"eth": rand.Int63(),
		},
		deyesClients: nil,
	}
	txOutProducer := DefaultTxOutputProducer{
		worldState: &worldState,
		tssConfig: config.TssConfig{
			Enable: true,
			SupportedChains: map[string]config.TssChainConfig{
				"ganache": {
					Symbol:   "ganache",
					DeyesUrl: "http://0.0.0.0:1234",
				},
			},
		},
		db:      mockDb,
		appKeys: mockAppKeys,
	}

	ctx := sdk.Context{}
	txOuts, txOutEntities, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
	require.NoError(t, err)
	require.Len(t, txOuts, 1)
	require.Len(t, txOutEntities, 1)
}

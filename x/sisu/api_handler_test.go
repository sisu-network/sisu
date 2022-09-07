package sisu

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	ecommon "github.com/ethereum/go-ethereum/common"

	eyesTypes "github.com/sisu-network/deyes/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/stretchr/testify/require"
)

func mockForApiHandlerTest() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestAfterContractDeployed(ctx)

	globalData := &common.MockGlobalData{
		GetReadOnlyContextFunc: func() sdk.Context {
			return ctx
		},
	}
	pmm := NewPostedMessageManager(k)
	txSubmit := &common.MockTxSubmit{}
	txTracker := &MockTxTracker{}

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}
	appKeys := common.NewMockAppKeys()

	mc := MockManagerContainer(k, pmm, globalData, partyManager, dheartClient, txSubmit, appKeys, ctx, txTracker)
	return ctx, mc
}

func signEthTx(rawTx *ethtypes.Transaction) *ethtypes.Transaction {
	signer := libchain.GetEthChainSigner("ganache1")
	if signer == nil {
		err := fmt.Errorf("cannot find signer for chain %s", "ganache1")
		log.Error(err)
	}

	hash := signer.Hash(rawTx)
	privateKey, _ := utils.GetEcdsaPrivateKey(utils.LOCALHOST_MNEMONIC)
	sig, err := ethcrypto.Sign(hash[:], privateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := rawTx.WithSignature(ethtypes.NewLondonSigner(big.NewInt(189985)), sig)
	if err != nil {
		panic(err)
	}

	return signedTx
}

func createTransferOutEthTx(gateway string, dstChain string, srcToken string,
	recipient string, amount *big.Int) *ethtypes.Transaction {
	srcTokenAddr := ecommon.HexToAddress(srcToken)

	vaultAbi, err := abi.JSON(strings.NewReader(vault.VaultABI))
	input, err := vaultAbi.Pack(
		MethodTransferOut,
		srcTokenAddr,
		dstChain,
		ecommon.HexToAddress(recipient),
		amount,
	)
	if err != nil {
		panic(err)
	}

	rawTx := ethTypes.NewTransaction(
		uint64(1),
		ecommon.HexToAddress(gateway),
		big.NewInt(0),
		100_000, // 100k for swapping operation.
		big.NewInt(1_000),
		input,
	)
	signedTx := signEthTx(rawTx)

	return signedTx
}

func TestApiHandler_OnTxIns(t *testing.T) {
	t.Run("empty_tx", func(t *testing.T) {
		_, mc := mockForApiHandlerTest()
		processor := NewApiHandler(nil, mc)

		require.NoError(t, processor.OnTxIns(&eyesTypes.Txs{}))
	})

	t.Run("eth_transfer", func(t *testing.T) {
		_, mc := mockForApiHandlerTest()

		srcChain := "ganache1"
		toAddress := "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80"
		ethTx := createTransferOutEthTx(toAddress, "ganache2",
			"0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c", "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988",
			new(big.Int).Mul(big.NewInt(1), utils.EthToWei))

		bz, err := ethTx.MarshalBinary()
		if err != nil {
			panic(err)
		}

		txs := &eyesTypes.Txs{
			Chain: srcChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: bz,
				To:         toAddress,
				Success:    true,
			}},
		}

		submitCount := 0
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			require.Equal(t, "Transfers", msg.Type())
			submitCount = 1
			return nil
		}

		processor := NewApiHandler(nil, mc)
		err = processor.OnTxIns(txs)

		require.NoError(t, err)
		require.Equal(t, 1, submitCount)
	})

	t.Run("tx_sent_from_our_sisu_account", func(t *testing.T) {
		// There should be no tx out created
		ctx, mc := mockForApiHandlerTest()

		privateKey, _ := utils.GetEcdsaPrivateKey(utils.LOCALHOST_MNEMONIC)
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		mc.Keeper().SetMpcAddress(ctx, "ganache1", address)
		mc.Keeper().SetMpcAddress(ctx, "ganache2", address)

		srcChain := "ganache1"
		toAddress := "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80"
		ethTx := createTransferOutEthTx(toAddress, "ganache2",
			"0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c", "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988",
			new(big.Int).Mul(big.NewInt(1), utils.EthToWei))

		bz, err := ethTx.MarshalBinary()
		if err != nil {
			panic(err)
		}

		txs := &eyesTypes.Txs{
			Chain: srcChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				From:       address,
				Hash:       utils.RandomHeximalString(64),
				Serialized: bz,
				To:         toAddress,
			}},
		}

		submitCount := 0
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			submitCount = 1
			return nil
		}

		processor := NewApiHandler(nil, mc)
		err = processor.OnTxIns(txs)

		require.NoError(t, err)
		require.Equal(t, 0, submitCount)
	})
}

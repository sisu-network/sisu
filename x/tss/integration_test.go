package tss

import (
	"context"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_createERC20TransferIn(t *testing.T) {
	// Comment t.Skip() to run this test
	t.Skip()

	// Please run ganache and deploy ERC20 gateway + ERC20 token before running this test.
	// See repo sisu-network/smart-contracts to get instructions
	txOutProducer := DefaultTxOutputProducer{}
	gatewayAddr := ethcommon.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	tokenAddr := ethcommon.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
	recipient := ethcommon.HexToAddress("0xbcd4042de499d14e55001ccbb24a551f3b954096")
	amount := big.NewInt(9999)
	txResponse, err := txOutProducer.callERC20TransferIn(gatewayAddr, tokenAddr, recipient, amount, "eth")
	require.NoError(t, err)

	signer := ethTypes.NewEIP2930Signer(big.NewInt(31337))
	txHash := signer.Hash(txResponse.EthTx)

	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	require.NoError(t, err)

	sig, err := crypto.Sign(txHash.Bytes(), privKey)
	require.NoError(t, err)

	signedTx, err := txResponse.EthTx.WithSignature(signer, sig)
	require.NoError(t, err)

	ethClient := initETHClient(t, "http://localhost:8545")
	err = ethClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
}

func initETHClient(t *testing.T, rawURL string) *ethclient.Client {
	ethClient, err := ethclient.Dial(rawURL)
	require.NoError(t, err)
	return ethClient
}

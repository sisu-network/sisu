package types

import (
	"math/big"
	"testing"

	"github.com/sisu-network/deyes/utils"
	"github.com/stretchr/testify/require"
)

func TestConvertAmountToSisuAmount(t *testing.T) {
	token := &Token{
		Chains: []string{
			"ganache1",
			"cardano",
			"solana-devnet",
		},
		Decimals: []uint32{
			18, 6, 8,
		},
	}

	expected := big.NewInt(1_500_000_000_000_000_000)

	amount, err := token.ChainAmountToSisuAmount("solana-devnet", big.NewInt(150_000_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ChainAmountToSisuAmount("ganache1", big.NewInt(1_500_000_000_000_000_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ChainAmountToSisuAmount("cardano", big.NewInt(1_500_000))
	require.Nil(t, err)
	require.Equal(t, expected, amount)

	amount, err = token.ChainAmountToSisuAmount("unknow", big.NewInt(1_500_000))
	require.NotNil(t, err)
}

func TestSisuAmountToChainAmount(t *testing.T) {
	token := &Token{
		Chains: []string{
			"ganache1",
			"cardano",
			"solana-devnet",
		},
		Decimals: []uint32{
			18, 6, 8,
		},
	}

	sisuAmount := big.NewInt(utils.OneEtherInWei)
	amount, err := token.SisuAmountToChainAmount("solana-devnet", sisuAmount)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(100_000_000), amount)
}

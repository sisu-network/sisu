package sisu

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mocksisu "github.com/sisu-network/sisu/tests/mock/x/sisu"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func TestTxOutProducerErc20_getGasCostInToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWorldState := mocksisu.NewMockWorldState(ctrl)

	p := &DefaultTxOutputProducer{
		worldState: mockWorldState,
	}

	chain := "ganache1"
	token := &types.Token{
		Id:    "SISU",
		Price: int64(4 * utils.DecinmalUnit),
	}
	mockWorldState.EXPECT().GetNativeTokenPriceForChain(chain).Return(int64(2*utils.DecinmalUnit), nil).Times(1)

	gas := big.NewInt(8_000_000)
	gasPrice := big.NewInt(10 * 1_000_000_000) // 10 gwei
	amount, err := p.getGasCostInToken(gas, gasPrice, chain, token)

	require.Equal(t, nil, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}

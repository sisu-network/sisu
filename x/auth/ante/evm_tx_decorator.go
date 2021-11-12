package ante

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	dtypes "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	evmTypes "github.com/sisu-network/sisu/x/evm/types"
)

var (
	ERR_NO_ETH_TX = fmt.Errorf("No ETH message is found in the cosmos transaction")
)

type EvmTxDecorator struct {
	validator ethchain.EthValidator
}

func NewEvmTxDecorator(validator ethchain.EthValidator) EvmTxDecorator {
	return EvmTxDecorator{
		validator: validator,
	}
}

func (decorator EvmTxDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	log.Debug("Checking ETH transaction")
	ethTxs := decorator.getEthTxs(tx.GetMsgs())
	if len(ethTxs) == 0 {
		return ctx, ERR_NO_ETH_TX
	}

	if err := decorator.validator.CheckTx(ethTxs); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (decorator EvmTxDecorator) getEthTxs(msgs []sdk.Msg) []*dtypes.Transaction {
	txs := make([]*dtypes.Transaction, 0)
	for _, msg := range msgs {
		etxMsg, ok := msg.(*evmTypes.EthTx)
		if !ok {
			continue
		}

		etx := new(dtypes.Transaction)
		err := etx.UnmarshalJSON(etxMsg.Data)
		if err != nil {
			continue
		}

		txs = append(txs, etx)
	}

	return txs
}

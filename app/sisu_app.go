package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *App) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	return abci.ResponseDeliverTx{}
}

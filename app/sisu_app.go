package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *App) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	app.BaseApp.DeliverTx(req)

	return abci.ResponseDeliverTx{Code: 0}
}

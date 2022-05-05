package sisu

import (
	"sync"

	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
)

type NetworkHealthListener interface {
	OnPing(source string)
}

type AppLogicListener interface {
	OnKeygenResult(result dhtypes.KeygenResult)
	OnTxIns(txs *eyesTypes.Txs) error
	OnKeysignResult(result *htypes.KeysignResult)
	OnReshareResult(result *htypes.ReshareResult)
	OnTxDeploymentResult(result *etypes.DispatchedTxResult)
	OnUpdateGasPriceRequest(request *etypes.GasPriceRequest)
	OnUpdateTokenPrice(tokenPrices []*etypes.TokenPrice)
}

// TODO: Rename this to API endPoint.
type ApiHandler struct {
	lock                  *sync.RWMutex
	networkHealthListener NetworkHealthListener
	appLogicListener      AppLogicListener
}

func NewApi(appLogicHandler AppLogicListener) *ApiHandler {
	return &ApiHandler{
		appLogicListener: appLogicHandler,
		lock:             &sync.RWMutex{},
	}
}

func (a *ApiHandler) getAppLogicListener() AppLogicListener {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.appLogicListener
}

func (a *ApiHandler) SetAppLogicListener(handler AppLogicListener) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.appLogicListener = handler
}

func (a *ApiHandler) getNetworkHealthListener() NetworkHealthListener {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.networkHealthListener
}

func (a *ApiHandler) SetNetworkHealthListener(listener NetworkHealthListener) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.networkHealthListener = listener
}

func (a *ApiHandler) Version() string {
	return "1.0"
}

///// Network health

// Empty function for checking health only.
func (a *ApiHandler) Ping(source string) error {
	listener := a.getNetworkHealthListener()
	if listener != nil {
		listener.OnPing(source)
	}

	return nil
}

///// Application logic

func (a *ApiHandler) KeygenResult(result htypes.KeygenResult) bool {
	log.Info("There is a Keygen Result")

	listener := a.getAppLogicListener()
	if listener != nil {
		listener.OnKeygenResult(result)
	}

	return true
}

// This is a API endpoint to receive transactions with To address we are interested in.
func (a *ApiHandler) PostObservedTxs(txs *etypes.Txs) {
	log.Debug("There is new list of transactions from deyes from chain ", txs.Chain)

	// There is a new transaction that we are interested in.
	listener := a.getAppLogicListener()
	if listener != nil {
		listener.OnTxIns(txs)
	}
}

func (a *ApiHandler) KeysignResult(result *htypes.KeysignResult) {
	log.Info("There is keysign result")
	handler := a.getAppLogicListener()
	if handler != nil {
		go handler.OnKeysignResult(result)
	}
}

func (a *ApiHandler) ReshareResult(result *htypes.ReshareResult) {
	log.Info("There is reshare result")
	handler := a.getAppLogicListener()
	if handler != nil {
		log.Debug("On going reshare result...")
		go handler.OnReshareResult(result)
	}
}

func (a *ApiHandler) PostDeploymentResult(result *etypes.DispatchedTxResult) {
	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnTxDeploymentResult(result)
	}
}

func (a *ApiHandler) UpdateGasPrice(request *etypes.GasPriceRequest) {
	log.Info("Received update gas price request")

	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnUpdateGasPriceRequest(request)
	}
}

func (a *ApiHandler) UpdateTokenPrices(prices []*etypes.TokenPrice) {
	log.Info("Received token prices update")

	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnUpdateTokenPrice(prices)
	}
}

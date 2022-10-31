package sisu

import (
	"sync"

	chainstypes "github.com/sisu-network/deyes/chains/types"

	etypes "github.com/sisu-network/deyes/types"
	eyestypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
)

type NetworkHealthListener interface {
	OnPing(source string)
}

type AppLogicListener interface {
	OnKeygenResult(result htypes.KeygenResult)
	OnTxIns(txs *eyestypes.Txs) error
	OnKeysignResult(result *htypes.KeysignResult)
	OnTxDeploymentResult(result *etypes.DispatchedTxResult)
	OnUpdateTokenPrice(tokenPrices []*etypes.TokenPrice)
	OnTxIncludedInBlock(txTrack *chainstypes.TrackUpdate)
}

type ApiEndPoint struct {
	lock                  *sync.RWMutex
	networkHealthListener NetworkHealthListener
	appLogicListener      AppLogicListener
}

func NewApi(appLogicHandler AppLogicListener) *ApiEndPoint {
	return &ApiEndPoint{
		appLogicListener: appLogicHandler,
		lock:             &sync.RWMutex{},
	}
}

func (a *ApiEndPoint) getAppLogicListener() AppLogicListener {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.appLogicListener
}

func (a *ApiEndPoint) SetAppLogicListener(handler AppLogicListener) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.appLogicListener = handler
}

func (a *ApiEndPoint) getNetworkHealthListener() NetworkHealthListener {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.networkHealthListener
}

func (a *ApiEndPoint) SetNetworkHealthListener(listener NetworkHealthListener) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.networkHealthListener = listener
}

func (a *ApiEndPoint) Version() string {
	return "1.0"
}

///// Network health

// Empty function for checking health only.
func (a *ApiEndPoint) Ping(source string) error {
	listener := a.getNetworkHealthListener()
	if listener != nil {
		listener.OnPing(source)
	}

	return nil
}

///// Application logic

func (a *ApiEndPoint) KeygenResult(result htypes.KeygenResult) bool {
	log.Info("There is a Keygen Result")

	listener := a.getAppLogicListener()
	if listener != nil {
		listener.OnKeygenResult(result)
	}

	return true
}

// This is a API endpoint to receive transactions with To address we are interested in.
func (a *ApiEndPoint) PostObservedTxs(txs *etypes.Txs) {
	log.Infof("There are %d transactions from deyes from chain %s", len(txs.Arr), txs.Chain)
	// There is a new transaction that we are interested in.
	listener := a.getAppLogicListener()
	if listener != nil {
		listener.OnTxIns(txs)
	}
}

func (a *ApiEndPoint) KeysignResult(result *htypes.KeysignResult) {
	log.Info("There is keysign result")
	handler := a.getAppLogicListener()
	if handler != nil {
		go handler.OnKeysignResult(result)
	}
}

func (a *ApiEndPoint) PostDeploymentResult(result *etypes.DispatchedTxResult) {
	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnTxDeploymentResult(result)
	}
}

func (a *ApiEndPoint) UpdateTokenPrices(prices []*etypes.TokenPrice) {
	log.Info("Received token prices update")

	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnUpdateTokenPrice(prices)
	}
}

func (a *ApiEndPoint) OnTxIncludedInBlock(txTrack *chainstypes.TrackUpdate) {
	listener := a.getAppLogicListener()
	if listener != nil {
		go listener.OnTxIncludedInBlock(txTrack)
	}
}

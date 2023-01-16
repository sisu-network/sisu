package common

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/atomic"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptoCdc "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/sisu-network/lib/log"
	tcrypto "github.com/tendermint/tendermint/crypto"
	pvm "github.com/tendermint/tendermint/privval"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// TODO: Refactor this interface.
type GlobalData interface {
	Init()
	UpdateCatchingUp() bool
	UpdateValidatorSets()
	IsCatchingUp() bool
	GetMyPubkey() tcrypto.PubKey
	GetValidatorSet() []rpc.ValidatorOutput
	SetReadOnlyContext(ctx sdk.Context)
	GetReadOnlyContext() sdk.Context
	AppInitialized() bool
	SetAppInitialized()

	// TxIn
	ConfirmTxIn(txIn *types.TxInDetails)
	GetTxInQueue() []*types.TxInDetails
	ResetTxInQueue()
}

// This is a struct to store global in-memory data. For persistent private data, use the private db.
type GlobalDataDefault struct {
	isCatchingUp    bool
	lock            *sync.RWMutex
	httpClient      *retryablehttp.Client
	cfg             config.Config
	myPubkey        tcrypto.PubKey
	cdc             *codec.LegacyAmino
	readOnlyContext atomic.Value
	isDataInit      *atomic.Bool
	txInQueue       []*types.TxInDetails

	validatorSets *rpc.ResultValidatorsOutput
	usedUtxos     map[string]bool
	// Set of all chains whose gas price should be recalculated.
	calGasChains *sync.Map
}

func NewGlobalData(cfg config.Config) GlobalData {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	cdc := codec.NewLegacyAmino()
	cryptoCdc.RegisterCrypto(cdc)

	return &GlobalDataDefault{
		httpClient:    httpClient,
		isCatchingUp:  true,
		lock:          &sync.RWMutex{},
		cdc:           cdc,
		validatorSets: new(rpc.ResultValidatorsOutput),
		cfg:           cfg,
		usedUtxos:     make(map[string]bool),
		isDataInit:    atomic.NewBool(false),
		calGasChains:  &sync.Map{},
		txInQueue:     make([]*types.TxInDetails, 0),
	}
}

// Initialize common variables that could be used throughout this app.
func (a *GlobalDataDefault) Init() {
	sisuConfig := a.cfg.Sisu

	defaultConfigTomlFile := sisuConfig.Dir + "/config/config.toml"
	data, err := os.ReadFile(defaultConfigTomlFile)
	if err != nil {
		panic(err)
	}

	var configToml struct {
		PrivValidatorKeyFile   string `toml:"priv_validator_key_file"`
		PrivValidatorStateFile string `toml:"priv_validator_state_file"`
	}

	if _, err := toml.Decode(string(data), &configToml); err != nil {
		panic(err)
	}

	privValidator := pvm.LoadFilePV(
		sisuConfig.Dir+"/"+configToml.PrivValidatorKeyFile,
		sisuConfig.Dir+"/"+configToml.PrivValidatorStateFile,
	)
	// Get the tendermint address of this node.
	a.myPubkey, _ = privValidator.GetPubKey()
}

func (a *GlobalDataDefault) UpdateCatchingUp() bool {
	oldValue := a.IsCatchingUp()

	url := "http://127.0.0.1:26657/status"
	body, _, err := utils.HttpGet(a.httpClient, url)
	if err != nil {
		log.Error(fmt.Errorf("Cannot get status data: %w", err))
		return oldValue
	}

	var resp struct {
		Result struct {
			SyncInfo struct {
				CatchingUp bool `json:"catching_up"`
			} `json:"sync_info"`
		} `json:"result"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
		} `json:"validator_info"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		log.Error(fmt.Errorf("Cannot parse tendermint status: %w", err))
		return oldValue
	}

	a.lock.Lock()
	a.isCatchingUp = resp.Result.SyncInfo.CatchingUp
	a.lock.Unlock()

	return resp.Result.SyncInfo.CatchingUp
}

func (a *GlobalDataDefault) UpdateValidatorSets() {
	url := "http://127.0.0.1:1317/validatorsets/latest"
	body, _, err := utils.HttpGet(a.httpClient, url)
	if err != nil {
		log.Error(fmt.Errorf("Cannot get status data: %w", err))
		return
	}

	responseWithHeight := new(rest.ResponseWithHeight)
	err = a.cdc.UnmarshalJSON(body, responseWithHeight)
	if err != nil {
		return
	}

	response := new(rpc.ResultValidatorsOutput)
	err = a.cdc.UnmarshalJSON([]byte(responseWithHeight.Result), response)
	if err != nil {
		return
	}

	// TODO: make this atomic
	a.validatorSets = response
}

// Returns the latest validator set.
func (a *GlobalDataDefault) GetValidatorSet() []rpc.ValidatorOutput {
	a.lock.RLock()
	defer a.lock.RUnlock()

	copy := a.validatorSets.Validators
	return copy
}

func (a *GlobalDataDefault) IsCatchingUp() bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.isCatchingUp
}

func (a *GlobalDataDefault) ValidatorSize() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return len(a.validatorSets.Validators)
}

func (a *GlobalDataDefault) SetReadOnlyContext(ctx sdk.Context) {
	a.readOnlyContext.Store(ctx)
}

func (a *GlobalDataDefault) GetReadOnlyContext() sdk.Context {
	val := a.readOnlyContext.Load()
	if val == nil {
		log.Error(("Read only context is not set"))
		return sdk.Context{}
	}

	return val.(sdk.Context)
}

func (a *GlobalDataDefault) AppInitialized() bool {
	return a.isDataInit.Load()
}

func (a *GlobalDataDefault) SetAppInitialized() {
	a.isDataInit.Store(true)
}

func (a *GlobalDataDefault) RecalculateGas(chain string) {
	a.calGasChains.Store(chain, true)
}

func (a *GlobalDataDefault) GetRecalculateGas() []string {
	ret := make([]string, 0)
	a.calGasChains.Range(func(key, value any) bool {
		ret = append(ret, key.(string))
		return true
	})

	return ret
}

func (a *GlobalDataDefault) ResetGasCalculation() {
	a.calGasChains.Range(func(key interface{}, value interface{}) bool {
		a.calGasChains.Delete(key)
		return true
	})
}

func (a *GlobalDataDefault) GetMyPubkey() tcrypto.PubKey {
	return a.myPubkey
}

func (a *GlobalDataDefault) ConfirmTxIn(txIn *types.TxInDetails) {
	a.txInQueue = append(a.txInQueue, txIn)
}

func (a *GlobalDataDefault) GetTxInQueue() []*types.TxInDetails {
	copy := a.txInQueue
	return copy
}

func (a *GlobalDataDefault) ResetTxInQueue() {
	a.txInQueue = make([]*types.TxInDetails, 0)
}

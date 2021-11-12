package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/sisu-network/cosmos-sdk/client/rpc"
	"github.com/sisu-network/cosmos-sdk/codec"
	cryptoCdc "github.com/sisu-network/cosmos-sdk/crypto/codec"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/rest"
	"github.com/sisu-network/lib/log"
	pvm "github.com/sisu-network/tendermint/privval"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type GlobalData interface {
	Init()
	UpdateCatchingUp()
	UpdateValidatorSets()
	IsCatchingUp() bool
	GetValidatorSet() []rpc.ValidatorOutput
	GetMyValidatorAddr() string
}

type GlobalDataDefault struct {
	isCatchingUp  bool
	catchUpLock   *sync.RWMutex
	httpClient    *retryablehttp.Client
	cfg           config.Config
	myTmtConsAddr sdk.ConsAddress
	cdc           *codec.LegacyAmino

	validatorSets *rpc.ResultValidatorsOutput
}

func NewGlobalData(cfg config.Config) GlobalData {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	cdc := codec.NewLegacyAmino()
	cryptoCdc.RegisterCrypto(cdc)

	return &GlobalDataDefault{
		httpClient:    httpClient,
		isCatchingUp:  true,
		catchUpLock:   &sync.RWMutex{},
		cdc:           cdc,
		validatorSets: new(rpc.ResultValidatorsOutput),
		cfg:           cfg,
	}
}

// Initialize common variables that could be used throughout this app.
func (a *GlobalDataDefault) Init() {
	sisuConfig := a.cfg.Sisu

	defaultConfigTomlFile := sisuConfig.Dir + "/config/config.toml"
	data, err := ioutil.ReadFile(defaultConfigTomlFile)
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
	a.myTmtConsAddr = (sdk.ConsAddress)(privValidator.GetAddress())

	log.Info("My tendermint address = ", a.myTmtConsAddr.String())
}

func (a *GlobalDataDefault) UpdateCatchingUp() {
	url := "http://127.0.0.1:26657/status"

	body, _, err := utils.HttpGet(a.httpClient, url)
	if err != nil {
		log.Error(fmt.Errorf("Cannot get status data: %w", err))
		return
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
		return
	}

	a.catchUpLock.Lock()
	a.isCatchingUp = resp.Result.SyncInfo.CatchingUp
	a.catchUpLock.Unlock()
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
	copy := a.validatorSets.Validators
	return copy
}

func (a *GlobalDataDefault) IsCatchingUp() bool {
	a.catchUpLock.RLock()
	defer a.catchUpLock.RUnlock()

	return a.isCatchingUp
}

func (a *GlobalDataDefault) ValidatorSize() int {
	// TODO: make this atomic
	return len(a.validatorSets.Validators)
}

func (a *GlobalDataDefault) GetMyValidatorAddr() string {
	return a.myTmtConsAddr.String()
}

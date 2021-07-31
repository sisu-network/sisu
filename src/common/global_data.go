package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	pvm "github.com/tendermint/tendermint/privval"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type GlobalData struct {
	isCatchingUp  bool
	catchUpLock   *sync.RWMutex
	httpClient    *retryablehttp.Client
	cfg           config.Config
	myTmtConsAddr sdk.ConsAddress

	validatorSets []*Validator
}

func NewGlobalData(cfg config.Config) *GlobalData {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &GlobalData{
		httpClient:    httpClient,
		isCatchingUp:  true,
		catchUpLock:   &sync.RWMutex{},
		validatorSets: make([]*Validator, 0),
		cfg:           cfg,
	}
}

// Initialize common variables that could be used throughout this app.
func (a *GlobalData) Init() {
	sisuConfig := a.cfg.GetSisuConfig()

	defaultConfigTomlFile := sisuConfig.Home + "/config/config.toml"
	fmt.Println("defaultConfigTomlFile = ", defaultConfigTomlFile)

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
		sisuConfig.Home+"/"+configToml.PrivValidatorKeyFile,
		sisuConfig.Home+"/"+configToml.PrivValidatorStateFile,
	)
	// Get the tendermint address of this node.
	a.myTmtConsAddr = (sdk.ConsAddress)(privValidator.GetAddress())

	utils.LogInfo("My tendermint address = ", a.myTmtConsAddr.String())
}

func (a *GlobalData) UpdateCatchingUp() {
	url := "http://127.0.0.1:26657/status"

	body, _, err := utils.HttpGet(a.httpClient, url)
	if err != nil {
		utils.LogError(fmt.Errorf("Cannot get status data: %w", err))
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
		utils.LogError(fmt.Errorf("Cannot parse tendermint status: %w", err))
		return
	}

	a.catchUpLock.Lock()
	a.isCatchingUp = resp.Result.SyncInfo.CatchingUp
	a.catchUpLock.Unlock()
}

// TODO: Consider do this less often
func (a *GlobalData) UpdateValidatorSets() {
	url := "http://127.0.0.1:1317/validatorsets/latest"
	body, _, err := utils.HttpGet(a.httpClient, url)
	if err != nil {
		utils.LogError(fmt.Errorf("Cannot get status data: %w", err))
		return
	}

	var resp struct {
		Result struct {
			Validators []*Validator `json:"validators"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		utils.LogError(fmt.Errorf("Cannot parse validator set body: %w", err))
		return
	}

	a.validatorSets = resp.Result.Validators
}

// Returns the latest validator set.
func (a *GlobalData) GetValidatorSet() []*Validator {
	copy := a.validatorSets
	return copy
}

func (a *GlobalData) IsCatchingUp() bool {
	a.catchUpLock.RLock()
	defer a.catchUpLock.RUnlock()

	return a.isCatchingUp
}

func (a *GlobalData) ValidatorSize() int {
	return len(a.validatorSets)
}

package common

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/utils"
)

type GlobalData struct {
	isCatchingUp bool
	catchUpLock  *sync.RWMutex
	httpClient   *retryablehttp.Client

	validatorSize int
	validatorSets []*Validator
}

func NewGlobalData() *GlobalData {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &GlobalData{
		httpClient:    httpClient,
		isCatchingUp:  true,
		catchUpLock:   &sync.RWMutex{},
		validatorSets: make([]*Validator, 0),
		validatorSize: 1, // TODO: Get the real validator size instead of hardcoding it.
	}
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
	return a.validatorSize
}

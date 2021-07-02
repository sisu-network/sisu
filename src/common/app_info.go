package common

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/utils"
)

type AppInfo struct {
	isCatchingUp bool
	httpClient   *retryablehttp.Client
	catchUpLock  *sync.RWMutex
}

func NewAppInfo() *AppInfo {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &AppInfo{
		httpClient:   httpClient,
		isCatchingUp: true,
		catchUpLock:  &sync.RWMutex{},
	}
}

func (a *AppInfo) UpdateCatchingUp() {
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

func (a *AppInfo) IsCatchingUp() bool {
	a.catchUpLock.RLock()
	defer a.catchUpLock.RUnlock()

	return a.isCatchingUp
}

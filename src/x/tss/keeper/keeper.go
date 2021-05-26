package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/tuktukclient"
)

type Keeper struct {
	tssConfig  *config.TssConfig
	client     *tuktukclient.Client
	httpClient *retryablehttp.Client
}

func NewKeeper(tssConfig *config.TssConfig) *Keeper {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &Keeper{
		tssConfig:  tssConfig,
		httpClient: httpClient,
	}
}

func (k *Keeper) Initialize() {
	var err error
	k.client, err = tuktukclient.Dial(fmt.Sprintf("http://%s:%d", k.tssConfig.Host, k.tssConfig.Port))
	if err != nil {
		panic(err)
	}

	go func() {
		k.waitForAppReady()
	}()
}

// Wait until this app catch up with the rest of the its tendermint peers so that we do not sign
// old transactions.
func (k *Keeper) waitForAppReady() {
	utils.LogInfo("Waiting to for app to catch up with its peers....")
	validatorCount := -1

	for {
		var err error
		validatorCount, err = k.getValidatorCount()
		if err == nil && validatorCount == 1 {
			utils.LogInfo("We are the only validator.")
			return
		}

		catchup, err := k.isCatchingUp()
		if err != nil {
			utils.LogError(err)
		}

		if catchup {
			utils.LogInfo("App is catched up with its peers!")
			break
		}

		time.Sleep(time.Second * 5)
	}
}

func (k *Keeper) isCatchingUp() (bool, error) {
	url := "http://127.0.0.1:26657/status"

	body, _, err := utils.HttpGet(k.httpClient, url)
	if err != nil {
		return false, fmt.Errorf("Cannot get status data: %w", err)
	}

	var resp struct {
		Result struct {
			SyncInfo struct {
				CatchingUp bool `json:"catching_up"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return false, fmt.Errorf("Cannot parse tendermint status: %w", err)
	}

	return resp.Result.SyncInfo.CatchingUp, nil
}

func (k *Keeper) getValidatorCount() (int, error) {
	// TODO: put this into a config file.
	url := "http://127.0.0.1:1317/validatorsets/latest"
	body, _, err := utils.HttpGet(k.httpClient, url)
	if err != nil {
		return 0, fmt.Errorf("Cannot get status data: %w", err)
	}

	var resp struct {
		Result struct {
			Validators []struct {
			} `json:"validators"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("Cannot parse validator count response: %w", err)
	}

	return len(resp.Result.Validators), nil
}

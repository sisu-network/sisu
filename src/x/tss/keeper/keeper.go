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
	go func() {
		k.connectTutTuk()
		k.waitForAppReady()
	}()
}

func (k *Keeper) connectTutTuk() {
	var err error
	url := fmt.Sprintf("http://%s:%d", k.tssConfig.Host, k.tssConfig.Port)
	k.client, err = tuktukclient.Dial(url)
	if err != nil {
		panic(err)
	}

	for {
		_, err := k.client.GetVersion()
		if err == nil {
			utils.LogInfo("Connected to TukTuk")
			return
		}

		fmt.Println("Sleeping for tuktuk")
		time.Sleep(time.Second * 3)
	}
}

// Wait until this app catch up with the rest of the its tendermint peers so that we do not sign
// old transactions.
func (k *Keeper) waitForAppReady() {
	utils.LogInfo("Waiting to for app to catch up with its peers....")

	for {

		catchup, err := k.isCatchingUp()
		if err != nil {
			utils.LogError(err)
			continue
		}

		if !catchup {
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

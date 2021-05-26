package tss

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/tuktukclient"
)

type Bridge struct {
	tssConfig  *config.TssConfig
	client     *tuktukclient.Client
	httpClient *retryablehttp.Client
}

func NewBridge(tssConfig *config.TssConfig) *Bridge {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	return &Bridge{
		tssConfig:  tssConfig,
		httpClient: httpClient,
	}
}

func (b *Bridge) Initialize() {
	go func() {
		b.connectTutTuk()
		b.waitForAppReady()
	}()
}

func (b *Bridge) connectTutTuk() {
	var err error
	url := fmt.Sprintf("http://%s:%d", b.tssConfig.Host, b.tssConfig.Port)
	b.client, err = tuktukclient.Dial(url)
	if err != nil {
		panic(err)
	}

	for {
		_, err := b.client.GetVersion()
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
func (b *Bridge) waitForAppReady() {
	utils.LogInfo("Waiting to for app to catch up with its peers....")

	for {

		catchup, err := b.isCatchingUp()
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

func (b *Bridge) isCatchingUp() (bool, error) {
	url := "http://127.0.0.1:26657/status"

	body, _, err := utils.HttpGet(b.httpClient, url)
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

package tssclients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	cardano "github.com/sisu-network/cardano-be/src/handler"
	"github.com/sisu-network/lib/log"
)

type CardanoClient interface {
	Ping(ctx context.Context) error
	MintMultiAsset(ctx context.Context, req *cardano.MintRequest) (map[string]interface{}, error)
}

type DefaultCardanoClient struct {
	url    string
	client *retryablehttp.Client
}

func (c *DefaultCardanoClient) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/ping", c.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	r, err := retryablehttp.FromRequest(req)
	if err != nil {
		return err
	}

	if _, err := c.client.Do(r); err != nil {
		return err
	}

	return nil
}

func NewDefaultCardanoClient(url string) *DefaultCardanoClient {
	return &DefaultCardanoClient{
		url:    url,
		client: retryablehttp.NewClient(),
	}
}

func (c *DefaultCardanoClient) MintMultiAsset(ctx context.Context, req *cardano.MintRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/mint", c.url)

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("error when marshal mint multi asset body: ", err)
		return nil, err
	}

	mintReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r, err := retryablehttp.FromRequest(mintReq)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	defer resp.Body.Close()
	if err != nil {
		log.Error("error when calling mint asset", err)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error when read response body: ", err)
		return nil, err
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &m); err != nil {
		log.Error("error when unmarshalling response body: ", err)
		return nil, err
	}

	return m, nil
}

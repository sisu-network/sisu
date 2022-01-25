package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/lib/log"
)

func HttpGet(httpClient *retryablehttp.Client, url string) ([]byte, int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Get request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("Failed to close repsonse body")
		}
	}()

	buf, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return buf, resp.StatusCode, errors.New("Http get status is not ok: " + resp.Status)
	}

	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Failed to read response body: %w", err)
	}
	return buf, resp.StatusCode, nil
}

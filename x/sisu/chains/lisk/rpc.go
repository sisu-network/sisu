package lisk

import (
	"encoding/json"
	"fmt"
	"github.com/sisu-network/lib/log"
	"io/ioutil"
	"net/http"
	"reflect"

	ltype "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/deyes/config"
)

type APIErr struct {
	message string
}

func NewApiErr(message string) error {
	return &APIErr{message: message}
}

func (e *APIErr) Error() string {
	return fmt.Sprintf(e.message)
}

// LiskRPC  A wrapper around lisk.RPC so that we can mock in watcher tests.
type LiskRPC interface {
	GetAccount(address string) (*ltype.Account, error)
	CreateTransaction(txHash string) (string, error)
}

type defaultLiskRPC struct {
	chain string
	rpc   string
}

func NewLiskRPC(cfg config.Chain) LiskRPC {
	c := &defaultLiskRPC{
		chain: cfg.Chain,
		rpc:   cfg.Rpcs[0],
	}
	return c
}

func (c *defaultLiskRPC) execute(endpoint string, method string, params map[string]string) ([]byte, error) {
	keys := reflect.ValueOf(params).MapKeys()
	req, err := http.NewRequest(method, c.rpc+endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for _, key := range keys {
		q.Add(key.Interface().(string), params[key.Interface().(string)])
	}
	req.URL.RawQuery = q.Encode()
	response, err := http.Get(req.URL.String())
	if response == nil {
		return nil, NewApiErr("cannot fetch data " + endpoint)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, err
}

func (c *defaultLiskRPC) CreateTransaction(txHash string) (string, error) {
	params := map[string]string{
		"tx": txHash,
	}
	response, err := c.execute("/transaction", "POST", params)
	if err != nil {
		return "", err
	}

	var responseObject ltype.TransactionResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return "", err
	}
	message := responseObject.Message

	return message, nil
}

func (c *defaultLiskRPC) GetAccount(address string) (*ltype.Account, error) {
	params := map[string]string{
		"address": address,
	}
	response, err := c.execute("/accounts", "GET", params)
	if err != nil {
		return nil, err
	}

	var responseObject ltype.ResponseAccount
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Errorf("GetAccount: Failed to marshal response, err = %s", err)
		return nil, err
	}

	accounts := responseObject.Data
	if len(accounts) == 0 {
		return nil, NewApiErr("lisk account block is not found")
	}
	validAccount := accounts[0]

	return validAccount, nil
}

package lisk

import (
	"encoding/json"
	"fmt"
	ltype "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/deyes/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/chains/lisk/utils"
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

func (c *defaultLiskRPC) CreateTransaction(txHash string) (string, error) {
	params := map[string]string{
		"transaction": txHash,
	}
	http := utils.NewHttpRPC(c.rpc)
	response, err := http.Post("/transactions", params)
	if err != nil {
		return "", err
	}

	var responseObject ltype.TransactionResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return "", err
	}
	message := responseObject.TransactionId

	return message, nil
}

func (c *defaultLiskRPC) GetAccount(address string) (*ltype.Account, error) {
	params := map[string]string{
		"address": address,
	}
	http := utils.NewHttpRPC(c.rpc)
	response, err := http.Get("/accounts", params)
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

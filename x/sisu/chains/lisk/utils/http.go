package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

type HTTP interface {
	Get(endpoint string, params map[string]string) ([]byte, error)
	Post(endpoint string, body map[string]string) ([]byte, error)
}

type defaultHttpRPC struct {
	rpc string
}

func NewHttpRPC(rpc string) HTTP {
	c := &defaultHttpRPC{
		rpc: rpc,
	}
	return c
}

func (c *defaultHttpRPC) Get(endpoint string, params map[string]string) ([]byte, error) {
	keys := reflect.ValueOf(params).MapKeys()
	req, err := http.NewRequest("Get", c.rpc+endpoint, nil)
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

func (c *defaultHttpRPC) Post(endpoint string, body map[string]string) ([]byte, error) {
	jsonValue, _ := json.Marshal(body)

	response, err := http.Post(c.rpc+endpoint, "application/json", bytes.NewBuffer(jsonValue))
	if response == nil {
		return nil, NewApiErr("cannot post data " + endpoint)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, err
}

package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sisu-network/lib/log"
)

func ReadRESTReq(r *http.Request, outData interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

	if err := decoder.Decode(outData); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

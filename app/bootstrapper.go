package app

import (
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	tss "github.com/sisu-network/sisu/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

const (
	RETRY_TIMEOUT = time.Second * 3
)

type Bootstrapper interface {
	BootstrapInternalNetwork(
		tssConfig config.TssConfig,
		apiHandler *tss.ApiEndPoint,
		encryptedAes []byte,
		tendermintKeyType string,
	) (tssclients.DheartClient, tssclients.DeyesClient)
}

type DefaultBootstrapper struct {
	dheartConnected atomic.Value
	deyesConnected  atomic.Value
}

func NewBootstrapper() Bootstrapper {
	return &DefaultBootstrapper{}
}

func (b *DefaultBootstrapper) BootstrapInternalNetwork(
	tssConfig config.TssConfig,
	apiHandler *tss.ApiEndPoint,
	encryptedAes []byte,
	tendermintKeyType string,
) (tssclients.DheartClient, tssclients.DeyesClient) {
	// Dheart
	var dheartClient tssclients.DheartClient
	var err error
	for {
		url := fmt.Sprintf("http://%s:%d", tssConfig.DheartHost, tssConfig.DheartPort)
		log.Info("Connecting to Dheart server at", url)

		dheartClient, err = tssclients.DialDheart(url)
		if err != nil {
			log.Infof("cannot dial dheart, err = %v, sleeping before retry...", err)
			time.Sleep(RETRY_TIMEOUT)
		} else {
			break
		}
	}

	for {
		err := dheartClient.Ping("sisu")
		if err != nil {
			log.Infof("cannot check dheart health, err = %v, sleeping before retry...", err)
			time.Sleep(RETRY_TIMEOUT)
		} else {
			break
		}
	}

	// Pass encrypted private key to dheart
	if err := dheartClient.SetPrivKey(hex.EncodeToString(encryptedAes), tendermintKeyType); err != nil {
		panic(err)
	}

	// Deyes
	var deyesClient tssclients.DeyesClient
	for {
		deyesClient, err = tssclients.DialDeyes(tssConfig.DeyesUrl)
		if err != nil {
			log.Infof("cannot dial deyes, err = %v, sleeping before retry...", err)
			time.Sleep(RETRY_TIMEOUT)
		} else {
			break
		}
	}

	for {
		err := deyesClient.Ping("sisu")
		if err != nil {
			log.Infof("cannot ping deyes, err = %v, sleeping before retry...", err)
			time.Sleep(RETRY_TIMEOUT)
		} else {
			break
		}
	}

	return dheartClient, deyesClient
}

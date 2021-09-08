package tssclients

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	tTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/utils"
)

type DheartClient struct {
	client *rpc.Client
}

// DialDheart connects a client to the given URL.
func DialDheart(rawurl string) (*DheartClient, error) {
	return dialDheartContext(context.Background(), rawurl)
}

func dialDheartContext(ctx context.Context, rawurl string) (*DheartClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDheartClient(c), nil
}

func newDheartClient(c *rpc.Client) *DheartClient {
	return &DheartClient{c}
}

func (c *DheartClient) SetPrivKey(encodedKey string, keyType string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_setPrivKey", encodedKey, keyType)
	if err != nil {
		utils.LogError("Cannot do handshare with dheart, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) CheckHealth() error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "tss_checkHealth")
	if err != nil {
		utils.LogError("Cannot check Dheart health, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) KeyGen(chainSymbol string) error {
	utils.LogInfo("Broadcasting keygen to Dheart")

	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_keyGen", chainSymbol)
	if err != nil {
		utils.LogError("Cannot send keygen request, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) KeySign(req *tTypes.KeysignRequest) error {
	utils.LogVerbose("Broadcasting key signing to Dheart")

	fmt.Println("Len(serialized) = ", len(req.OutBytes))

	var r interface{}
	err := c.client.CallContext(context.Background(), &r, "tss_keySign", req)
	if err != nil {
		utils.LogError("Cannot send KeySign request, err = ", err)
		return err
	}

	return nil
}

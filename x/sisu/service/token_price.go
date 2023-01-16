package service

import (
	"math/big"

	"github.com/sisu-network/sisu/x/sisu/external"
)

type tokenPriceCache struct {
	updateTime int64
}

type TokenPrice interface {
}

// TODO: Add caching for the token price.
type defaultTokenPrice struct {
	deyesClient external.DeyesClient
}

func NewTokenPrice(deyesClient external.DeyesClient) TokenPrice {
	return &defaultTokenPrice{
		deyesClient: deyesClient,
	}
}

func (d *defaultTokenPrice) getTokenPrice(tokenId string) *big.Int {
	return nil
}

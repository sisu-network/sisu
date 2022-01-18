package dev

import (
	"crypto/ecdsa"
	"math/big"
	"time"
)

const (
	defaultMnemonic = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
	Blocktime       = time.Second * 3
)

var (
	privateKey0 *ecdsa.PrivateKey
	nonceMap    map[string]*big.Int
)

func init() {
}

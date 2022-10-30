package utils

import (
	"encoding/hex"

	"github.com/echovl/cardano-go"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"golang.org/x/crypto/sha3"
)

func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

func GetEthSender(tx *etypes.Transaction) (common.Address, error) {
	msg, err := tx.AsMessage(etypes.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		return common.Address{}, err
	}

	return msg.From(), nil
}

func GetAddressFromCardanoPubkey(pubkey []byte) cardano.Address {
	keyHash, err := cardano.Blake224Hash(pubkey)
	if err != nil {
		panic(err)
	}

	payment := cardano.StakeCredential{Type: cardano.KeyCredential, KeyHash: keyHash}
	enterpriseAddr, err := cardano.NewEnterpriseAddress(0, payment)
	if err != nil {
		panic(err)
	}

	return enterpriseAddr
}

func GetKeyTypeForChain(chain string) string {
	if libchain.IsETHBasedChain(chain) {
		return libchain.KEY_TYPE_ECDSA
	}

	if libchain.IsCardanoChain(chain) || libchain.IsSolanaChain(chain) {
		return libchain.KEY_TYPE_EDDSA
	}

	return "" // unknown
}

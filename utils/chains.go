package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"golang.org/x/crypto/sha3"
)

func GetTxHash(chain string, serialized []byte) (string, error) {
	if libchain.IsETHBasedChain(chain) {
		tx := &etypes.Transaction{}
		err := tx.UnmarshalBinary(serialized)
		if err != nil {
			return "", err
		}

		return KeccakHash32(string(serialized)), nil
	}

	// TODO: Support more chain other than ETH family.
	return "", fmt.Errorf("Unknwon chain: %s", chain)
}

func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

func GetEthChainSigners() map[string]etypes.Signer {
	m := make(map[string]etypes.Signer)

	// TODO: Add correct signer for each chain. For now, use NewEIP2930Signer for all chains.
	m["eth"] = etypes.NewEIP2930Signer(libchain.GetChainIntFromId("eth"))
	m["sisu-eth"] = etypes.NewEIP2930Signer(libchain.GetChainIntFromId("sisu-eth"))

	return m
}

func GetEthSender(tx *etypes.Transaction) (common.Address, error) {
	msg, err := tx.AsMessage(etypes.NewEIP2930Signer(tx.ChainId()))
	if err != nil {
		return common.Address{}, err
	}

	return msg.From(), nil
}

package main

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sisu-network/sisu/utils"
)

func main() {
	const abiString = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"getName","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"newName","type":"string"}],"name":"setName","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	myAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		panic(err)
	}

	input := `0xc47f0027000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000044141414100000000000000000000000000000000000000000000000000000000`
	payload := input[2:]

	decodeData, err := hex.DecodeString(payload)
	if err != nil {
		panic(err)
	}

	var funcName string
	for name, method := range myAbi.Methods {
		hash := hex.EncodeToString(crypto.Keccak256([]byte(method.Sig)))
		if hash[:8] == payload[:8] {
			utils.LogInfo("Found function name: ", name)
			funcName = name
			break
		}
	}

	if funcName == "" {
		panic("cannot find function name")
	}

	params, err := myAbi.Methods[funcName].Inputs.Unpack(decodeData[4:])
	if err != nil {
		panic(err)
	}

	utils.LogInfo("params = ", params)
}
